package service

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

const requestLogPayloadEncodingGzip = "gzip"

type RequestLogPayload struct {
	RequestBody          *string `json:"request_body,omitempty"`
	RequestBodyBytes     *int    `json:"request_body_bytes,omitempty"`
	RequestBodyTruncated bool    `json:"request_body_truncated"`

	ResponseBody          *string `json:"response_body,omitempty"`
	ResponseBodyBytes     *int    `json:"response_body_bytes,omitempty"`
	ResponseBodyTruncated bool    `json:"response_body_truncated"`
}

func BuildRequestLogPayload(settings *RequestLogSettings, requestBody []byte, responseBody []byte) *RequestLogPayload {
	effective := normalizeRequestLogSettings(settings)
	if effective == nil {
		return nil
	}

	payload := &RequestLogPayload{}

	if effective.CaptureRequestBody && len(requestBody) > 0 {
		if body, truncated, bytesLen := sanitizeAndTrimRequestBody(requestBody, effective.MaxRequestBodyBytes); body != "" {
			payload.RequestBody = &body
			payload.RequestBodyTruncated = truncated
			payload.RequestBodyBytes = intPtr(bytesLen)
		}
	}

	if effective.CaptureResponseBody && len(responseBody) > 0 {
		body, truncated := sanitizeResponseBodyForRequestLog(responseBody, effective.MaxResponseBodyBytes)
		if body != "" {
			payload.ResponseBody = &body
			payload.ResponseBodyTruncated = truncated
			payload.ResponseBodyBytes = intPtr(len(responseBody))
		}
	}

	if payload.RequestBody == nil && payload.ResponseBody == nil {
		return nil
	}

	return payload
}

func EncodeRequestLogPayload(value *string) ([]byte, *string, error) {
	if value == nil {
		return nil, nil, nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil, nil, nil
	}

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if _, err := zw.Write([]byte(trimmed)); err != nil {
		_ = zw.Close()
		return nil, nil, fmt.Errorf("gzip request log payload: %w", err)
	}
	if err := zw.Close(); err != nil {
		return nil, nil, fmt.Errorf("close gzip request log payload: %w", err)
	}

	encoding := requestLogPayloadEncodingGzip
	return buf.Bytes(), &encoding, nil
}

func DecodeRequestLogPayload(data []byte, encoding *string) (*string, error) {
	if len(data) == 0 {
		return nil, nil
	}

	codec := requestLogPayloadEncodingGzip
	if encoding != nil && strings.TrimSpace(*encoding) != "" {
		codec = strings.TrimSpace(*encoding)
	}

	switch codec {
	case requestLogPayloadEncodingGzip:
		reader, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("open gzip request log payload: %w", err)
		}
		defer func() {
			_ = reader.Close()
		}()

		raw, err := io.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("read gzip request log payload: %w", err)
		}
		out := string(raw)
		return &out, nil
	case "":
		out := string(data)
		return &out, nil
	default:
		return nil, fmt.Errorf("unsupported request log payload encoding: %s", codec)
	}
}

func sanitizeResponseBodyForRequestLog(raw []byte, maxBytes int) (string, bool) {
	if len(raw) == 0 {
		return "", false
	}

	if json.Valid(raw) {
		if out, truncated, _ := sanitizeAndTrimRequestBody(raw, maxBytes); out != "" {
			return out, truncated
		}
	}

	return sanitizeErrorBodyForStorage(string(raw), maxBytes)
}
