package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type RequestLogSettings struct {
	CaptureRequestBody   bool `json:"capture_request_body"`
	CaptureResponseBody  bool `json:"capture_response_body"`
	MaxRequestBodyBytes  int  `json:"max_request_body_bytes"`
	MaxResponseBodyBytes int  `json:"max_response_body_bytes"`
	RetentionDays        int  `json:"retention_days"`
}

func DefaultRequestLogSettings() *RequestLogSettings {
	return &RequestLogSettings{
		CaptureRequestBody:   true,
		CaptureResponseBody:  true,
		MaxRequestBodyBytes:  16 * 1024,
		MaxResponseBodyBytes: 32 * 1024,
		RetentionDays:        30,
	}
}

func normalizeRequestLogSettings(settings *RequestLogSettings) *RequestLogSettings {
	effective := DefaultRequestLogSettings()
	if settings == nil {
		return effective
	}

	effective.CaptureRequestBody = settings.CaptureRequestBody
	effective.CaptureResponseBody = settings.CaptureResponseBody
	if settings.MaxRequestBodyBytes > 0 {
		effective.MaxRequestBodyBytes = settings.MaxRequestBodyBytes
	}
	if settings.MaxResponseBodyBytes > 0 {
		effective.MaxResponseBodyBytes = settings.MaxResponseBodyBytes
	}
	effective.RetentionDays = settings.RetentionDays
	if effective.RetentionDays < 0 {
		effective.RetentionDays = 0
	}

	return effective
}

func validateRequestLogSettings(settings *RequestLogSettings) error {
	if settings == nil {
		return fmt.Errorf("settings cannot be nil")
	}
	if settings.MaxRequestBodyBytes < 1024 || settings.MaxRequestBodyBytes > 256*1024 {
		return fmt.Errorf("max_request_body_bytes must be between 1024 and 262144")
	}
	if settings.MaxResponseBodyBytes < 1024 || settings.MaxResponseBodyBytes > 512*1024 {
		return fmt.Errorf("max_response_body_bytes must be between 1024 and 524288")
	}
	if settings.RetentionDays < 0 || settings.RetentionDays > 365 {
		return fmt.Errorf("retention_days must be between 0 and 365")
	}
	return nil
}

func (s *SettingService) GetRequestLogSettings(ctx context.Context) (*RequestLogSettings, error) {
	value, err := s.settingRepo.GetValue(ctx, SettingKeyRequestLogSettings)
	if err != nil {
		if errors.Is(err, ErrSettingNotFound) {
			return DefaultRequestLogSettings(), nil
		}
		return nil, fmt.Errorf("get request log settings: %w", err)
	}
	if value == "" {
		return DefaultRequestLogSettings(), nil
	}

	var settings RequestLogSettings
	if err := json.Unmarshal([]byte(value), &settings); err != nil {
		return DefaultRequestLogSettings(), nil
	}

	return normalizeRequestLogSettings(&settings), nil
}

func (s *SettingService) SetRequestLogSettings(ctx context.Context, settings *RequestLogSettings) error {
	effective := normalizeRequestLogSettings(settings)
	if err := validateRequestLogSettings(effective); err != nil {
		return err
	}

	data, err := json.Marshal(effective)
	if err != nil {
		return fmt.Errorf("marshal request log settings: %w", err)
	}

	return s.settingRepo.Set(ctx, SettingKeyRequestLogSettings, string(data))
}
