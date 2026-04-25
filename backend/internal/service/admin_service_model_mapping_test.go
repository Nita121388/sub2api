//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type accountRepoStubForModelMapping struct {
	accountRepoStub
	accountsByID map[int64]*Account
	updates      map[int64]map[string]any
}

func (s *accountRepoStubForModelMapping) GetByIDs(_ context.Context, ids []int64) ([]*Account, error) {
	out := make([]*Account, 0, len(ids))
	for _, id := range ids {
		if account, ok := s.accountsByID[id]; ok {
			out = append(out, account)
		}
	}
	return out, nil
}

func (s *accountRepoStubForModelMapping) UpdateCredentials(_ context.Context, id int64, credentials map[string]any) error {
	if s.updates == nil {
		s.updates = make(map[int64]map[string]any)
	}
	s.updates[id] = cloneCredentials(credentials)
	return nil
}

func TestAdminService_AppendAccountModelMappings_AppendsAndPreservesCredentials(t *testing.T) {
	repo := &accountRepoStubForModelMapping{
		accountsByID: map[int64]*Account{
			1: {
				ID:       1,
				Name:     "openai-1",
				Platform: PlatformOpenAI,
				Status:   StatusActive,
				Credentials: map[string]any{
					"access_token": "secret",
					"model_mapping": map[string]any{
						"gpt-5.4": "gpt-5.4",
					},
				},
			},
		},
	}
	svc := &adminServiceImpl{accountRepo: repo}

	result, err := svc.AppendAccountModelMappings(context.Background(), &AppendAccountModelMappingsInput{
		Platform:                PlatformOpenAI,
		AccountIDs:              []int64{1},
		Mappings:                map[string]string{"gpt-5.5": "gpt-5.5"},
		OnlyWithExistingMapping: true,
	})

	require.NoError(t, err)
	require.Equal(t, 1, result.Changed)
	require.Equal(t, "changed", result.Results[0].Action)
	require.Equal(t, "secret", repo.updates[1]["access_token"])
	mapping := repo.updates[1]["model_mapping"].(map[string]any)
	require.Equal(t, "gpt-5.4", mapping["gpt-5.4"])
	require.Equal(t, "gpt-5.5", mapping["gpt-5.5"])
}

func TestAdminService_AppendAccountModelMappings_SkipsAccountsWithoutMapping(t *testing.T) {
	repo := &accountRepoStubForModelMapping{
		accountsByID: map[int64]*Account{
			1: {
				ID:          1,
				Name:        "openai-1",
				Platform:    PlatformOpenAI,
				Status:      StatusActive,
				Credentials: map[string]any{"access_token": "secret"},
			},
		},
	}
	svc := &adminServiceImpl{accountRepo: repo}

	result, err := svc.AppendAccountModelMappings(context.Background(), &AppendAccountModelMappingsInput{
		Platform:                PlatformOpenAI,
		AccountIDs:              []int64{1},
		Mappings:                map[string]string{"gpt-5.5": "gpt-5.5"},
		OnlyWithExistingMapping: true,
	})

	require.NoError(t, err)
	require.Equal(t, 0, result.Changed)
	require.Equal(t, 1, result.Skipped)
	require.Equal(t, "skipped", result.Results[0].Action)
	require.Equal(t, "no_model_mapping", result.Results[0].Reason)
	require.Empty(t, repo.updates)
}

func TestAdminService_AppendAccountModelMappings_ConflictDoesNotPartiallyUpdate(t *testing.T) {
	repo := &accountRepoStubForModelMapping{
		accountsByID: map[int64]*Account{
			1: {
				ID:       1,
				Name:     "openai-1",
				Platform: PlatformOpenAI,
				Status:   StatusActive,
				Credentials: map[string]any{
					"model_mapping": map[string]any{
						"gpt-5.5": "other-model",
					},
				},
			},
		},
	}
	svc := &adminServiceImpl{accountRepo: repo}

	result, err := svc.AppendAccountModelMappings(context.Background(), &AppendAccountModelMappingsInput{
		Platform:                PlatformOpenAI,
		AccountIDs:              []int64{1},
		Mappings:                map[string]string{"gpt-5.5": "gpt-5.5", "gpt-5.6": "gpt-5.6"},
		OnlyWithExistingMapping: true,
	})

	require.NoError(t, err)
	require.Equal(t, 1, result.Conflicted)
	require.Equal(t, "conflict", result.Results[0].Action)
	require.Len(t, result.Results[0].Conflicts, 1)
	require.Empty(t, repo.updates)
}

func TestAdminService_AppendAccountModelMappings_DryRunDoesNotPersist(t *testing.T) {
	repo := &accountRepoStubForModelMapping{
		accountsByID: map[int64]*Account{
			1: {
				ID:       1,
				Name:     "openai-1",
				Platform: PlatformOpenAI,
				Status:   StatusActive,
				Credentials: map[string]any{
					"model_mapping": map[string]any{
						"gpt-5.4": "gpt-5.4",
					},
				},
			},
		},
	}
	svc := &adminServiceImpl{accountRepo: repo}

	result, err := svc.AppendAccountModelMappings(context.Background(), &AppendAccountModelMappingsInput{
		Platform:                PlatformOpenAI,
		AccountIDs:              []int64{1},
		Mappings:                map[string]string{"gpt-5.5": "gpt-5.5"},
		DryRun:                  true,
		OnlyWithExistingMapping: true,
	})

	require.NoError(t, err)
	require.True(t, result.DryRun)
	require.Equal(t, 1, result.Changed)
	require.Empty(t, repo.updates)
}
