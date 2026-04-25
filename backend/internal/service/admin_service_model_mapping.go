package service

import (
	"context"
	"strings"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

type AppendAccountModelMappingsInput struct {
	Platform                string
	AccountIDs              []int64
	Mappings                map[string]string
	DryRun                  bool
	OverwriteExisting       bool
	OnlyWithExistingMapping bool
}

type AccountModelMappingConflict struct {
	Model     string `json:"model"`
	Existing  string `json:"existing"`
	Requested string `json:"requested"`
}

type AppendAccountModelMappingAccountResult struct {
	AccountID     int64                         `json:"account_id"`
	AccountName   string                        `json:"account_name,omitempty"`
	Platform      string                        `json:"platform,omitempty"`
	Status        string                        `json:"status,omitempty"`
	Action        string                        `json:"action"`
	Reason        string                        `json:"reason,omitempty"`
	ExistingCount int                           `json:"existing_count"`
	FinalCount    int                           `json:"final_count"`
	Added         map[string]string             `json:"added,omitempty"`
	Overwritten   map[string]string             `json:"overwritten,omitempty"`
	Conflicts     []AccountModelMappingConflict `json:"conflicts,omitempty"`
}

type AppendAccountModelMappingsResult struct {
	DryRun     bool                                     `json:"dry_run"`
	Platform   string                                   `json:"platform"`
	Total      int                                      `json:"total"`
	Changed    int                                      `json:"changed"`
	Skipped    int                                      `json:"skipped"`
	Conflicted int                                      `json:"conflicted"`
	Results    []AppendAccountModelMappingAccountResult `json:"results"`
}

func (s *adminServiceImpl) AppendAccountModelMappings(ctx context.Context, input *AppendAccountModelMappingsInput) (*AppendAccountModelMappingsResult, error) {
	if input == nil {
		return nil, infraerrors.BadRequest("ACCOUNT_MODEL_MAPPING_NIL_INPUT", "input cannot be nil")
	}

	platform := strings.TrimSpace(input.Platform)
	if platform == "" {
		platform = PlatformOpenAI
	}

	mappings := sanitizeModelMappingPatch(input.Mappings)
	if len(mappings) == 0 {
		return nil, infraerrors.BadRequest("ACCOUNT_MODEL_MAPPING_EMPTY", "model mappings cannot be empty")
	}

	result := &AppendAccountModelMappingsResult{
		DryRun:   input.DryRun,
		Platform: platform,
		Results:  []AppendAccountModelMappingAccountResult{},
	}

	accounts, missingIDs, err := s.loadAccountsForModelMappingAppend(ctx, platform, input.AccountIDs)
	if err != nil {
		return nil, err
	}

	for _, id := range missingIDs {
		result.Results = append(result.Results, AppendAccountModelMappingAccountResult{
			AccountID: id,
			Action:    "skipped",
			Reason:    "not_found",
		})
		result.Skipped++
	}

	for _, account := range accounts {
		accountResult, err := s.appendModelMappingsToAccount(ctx, account, platform, mappings, input)
		if err != nil {
			return nil, err
		}
		result.Results = append(result.Results, accountResult)
		switch accountResult.Action {
		case "changed":
			result.Changed++
		case "conflict":
			result.Conflicted++
		default:
			result.Skipped++
		}
	}

	result.Total = len(result.Results)
	return result, nil
}

func sanitizeModelMappingPatch(in map[string]string) map[string]string {
	out := make(map[string]string)
	for from, to := range in {
		from = strings.TrimSpace(from)
		to = strings.TrimSpace(to)
		if from == "" {
			continue
		}
		if to == "" {
			to = from
		}
		out[from] = to
	}
	return out
}

func (s *adminServiceImpl) loadAccountsForModelMappingAppend(ctx context.Context, platform string, accountIDs []int64) ([]*Account, []int64, error) {
	if len(accountIDs) > 0 {
		accounts, err := s.accountRepo.GetByIDs(ctx, accountIDs)
		if err != nil {
			return nil, nil, err
		}

		byID := make(map[int64]*Account, len(accounts))
		for _, account := range accounts {
			if account != nil {
				byID[account.ID] = account
			}
		}

		ordered := make([]*Account, 0, len(accountIDs))
		missing := make([]int64, 0)
		seen := make(map[int64]struct{}, len(accountIDs))
		for _, id := range accountIDs {
			if id <= 0 {
				continue
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			if account, ok := byID[id]; ok {
				ordered = append(ordered, account)
			} else {
				missing = append(missing, id)
			}
		}
		return ordered, missing, nil
	}

	const pageSize = 500
	all := make([]*Account, 0)
	for page := 1; ; page++ {
		rows, total, err := s.ListAccounts(ctx, page, pageSize, platform, "", "", "", 0, "", "id", "asc")
		if err != nil {
			return nil, nil, err
		}
		for i := range rows {
			row := rows[i]
			all = append(all, &row)
		}
		if len(rows) == 0 || len(all) >= int(total) {
			break
		}
	}
	return all, nil, nil
}

func (s *adminServiceImpl) appendModelMappingsToAccount(
	ctx context.Context,
	account *Account,
	platform string,
	mappings map[string]string,
	input *AppendAccountModelMappingsInput,
) (AppendAccountModelMappingAccountResult, error) {
	if account == nil {
		return AppendAccountModelMappingAccountResult{Action: "skipped", Reason: "not_found"}, nil
	}

	result := AppendAccountModelMappingAccountResult{
		AccountID:   account.ID,
		AccountName: account.Name,
		Platform:    account.Platform,
		Status:      account.Status,
		Action:      "skipped",
	}

	if platform != "" && account.Platform != platform {
		result.Reason = "platform_mismatch"
		return result, nil
	}

	credentials := cloneCredentials(account.Credentials)
	currentMapping, hasMapping := credentialsModelMapping(credentials)
	result.ExistingCount = len(currentMapping)
	result.FinalCount = len(currentMapping)

	if input.OnlyWithExistingMapping && (!hasMapping || len(currentMapping) == 0) {
		result.Reason = "no_model_mapping"
		return result, nil
	}

	conflicts := make([]AccountModelMappingConflict, 0)
	added := make(map[string]string)
	overwritten := make(map[string]string)
	merged := make(map[string]string, len(currentMapping)+len(mappings))
	for k, v := range currentMapping {
		merged[k] = v
	}

	for from, to := range mappings {
		if existing, ok := merged[from]; ok {
			if existing == to {
				continue
			}
			if !input.OverwriteExisting {
				conflicts = append(conflicts, AccountModelMappingConflict{
					Model:     from,
					Existing:  existing,
					Requested: to,
				})
				continue
			}
			merged[from] = to
			overwritten[from] = to
			continue
		}
		merged[from] = to
		added[from] = to
	}

	if len(conflicts) > 0 {
		result.Action = "conflict"
		result.Conflicts = conflicts
		return result, nil
	}

	if len(added) == 0 && len(overwritten) == 0 {
		result.Action = "unchanged"
		result.Reason = "already_present"
		return result, nil
	}

	result.Action = "changed"
	result.Added = added
	result.Overwritten = overwritten
	result.FinalCount = len(merged)

	if input.DryRun {
		return result, nil
	}

	credentials["model_mapping"] = stringMapToAnyMap(merged)
	if err := persistAccountCredentials(ctx, s.accountRepo, account, credentials); err != nil {
		return result, err
	}
	return result, nil
}

func credentialsModelMapping(credentials map[string]any) (map[string]string, bool) {
	if credentials == nil {
		return map[string]string{}, false
	}
	raw, ok := credentials["model_mapping"]
	if !ok || raw == nil {
		return map[string]string{}, false
	}

	out := make(map[string]string)
	switch typed := raw.(type) {
	case map[string]any:
		for k, v := range typed {
			if s, ok := v.(string); ok {
				out[k] = s
			}
		}
	case map[string]string:
		for k, v := range typed {
			out[k] = v
		}
	default:
		return out, false
	}
	return out, true
}

func stringMapToAnyMap(in map[string]string) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
