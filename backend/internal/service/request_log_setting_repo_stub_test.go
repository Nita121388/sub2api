package service

import "context"

type requestLogSettingRepoStub struct {
	SettingRepository
	values map[string]string
	err    error
}

func (s *requestLogSettingRepoStub) GetValue(ctx context.Context, key string) (string, error) {
	if s.err != nil {
		return "", s.err
	}
	if value, ok := s.values[key]; ok {
		return value, nil
	}
	return "", ErrSettingNotFound
}
