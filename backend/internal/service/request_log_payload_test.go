//go:build unit

package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildRequestLogPayload_RedactsAndCapturesBodies(t *testing.T) {
	settings := &RequestLogSettings{
		CaptureRequestBody:   true,
		CaptureResponseBody:  true,
		MaxRequestBodyBytes:  4096,
		MaxResponseBodyBytes: 4096,
		RetentionDays:        30,
	}

	payload := BuildRequestLogPayload(
		settings,
		[]byte(`{"model":"gpt-5.4","api_key":"secret-token","messages":[{"role":"user","content":"hello"}]}`),
		[]byte(`{"id":"resp_1","output":[{"type":"message","content":[{"type":"output_text","text":"ok"}]}]}`),
	)

	require.NotNil(t, payload)
	require.NotNil(t, payload.RequestBody)
	require.Contains(t, *payload.RequestBody, "[REDACTED]")
	require.NotContains(t, *payload.RequestBody, "secret-token")
	require.NotNil(t, payload.RequestBodyBytes)
	require.Greater(t, *payload.RequestBodyBytes, 0)
	require.NotNil(t, payload.ResponseBody)
	require.Contains(t, *payload.ResponseBody, `"id":"resp_1"`)
	require.False(t, payload.RequestBodyTruncated)
	require.False(t, payload.ResponseBodyTruncated)
}

func TestEncodeDecodeRequestLogPayload_RoundTrip(t *testing.T) {
	raw := `{"hello":"world"}`
	encoded, encoding, err := EncodeRequestLogPayload(&raw)
	require.NoError(t, err)
	require.NotEmpty(t, encoded)
	require.NotNil(t, encoding)
	require.Equal(t, requestLogPayloadEncodingGzip, *encoding)

	decoded, err := DecodeRequestLogPayload(encoded, encoding)
	require.NoError(t, err)
	require.NotNil(t, decoded)
	require.Equal(t, raw, *decoded)
}
