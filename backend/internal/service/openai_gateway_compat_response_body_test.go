package service

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestHandleChatBufferedStreamingResponse_ReturnsResponseBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	svc := &OpenAIGatewayService{}
	resp := &http.Response{
		Header: http.Header{
			"Content-Type": []string{"text/event-stream"},
		},
		Body: io.NopCloser(strings.NewReader(strings.Join([]string{
			`data: {"type":"response.completed","response":{"id":"resp_chat_buffered","model":"gpt-4.1-mini","output":[{"type":"message","role":"assistant","content":[{"type":"output_text","text":"hello from buffered"}]}],"usage":{"input_tokens":11,"output_tokens":22}}}`,
			`data: [DONE]`,
			``,
		}, "\n"))),
	}

	result, err := svc.handleChatBufferedStreamingResponse(resp, c, "gpt-4.1-mini", "gpt-4.1-mini", time.Now())
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotEmpty(t, result.ResponseBody)
	require.Contains(t, string(result.ResponseBody), "hello from buffered")
}

func TestHandleChatBufferedStreamingResponse_ResponseDoneReturnsResponseBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	svc := &OpenAIGatewayService{}
	resp := &http.Response{
		Header: http.Header{
			"Content-Type": []string{"text/event-stream"},
		},
		Body: io.NopCloser(strings.NewReader(strings.Join([]string{
			`data: {"type":"response.done","response":{"id":"resp_chat_done","model":"gpt-4.1-mini","output":[{"type":"message","role":"assistant","content":[{"type":"output_text","text":"hello from done"}]}],"usage":{"input_tokens":3,"output_tokens":4}}}`,
			``,
		}, "\n"))),
	}

	result, err := svc.handleChatBufferedStreamingResponse(resp, c, "gpt-4.1-mini", "gpt-4.1-mini", time.Now())
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotEmpty(t, result.ResponseBody)
	require.Contains(t, string(result.ResponseBody), "hello from done")
}

func TestHandleChatStreamingResponse_ReturnsResponseBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	svc := &OpenAIGatewayService{}
	resp := &http.Response{
		Header: http.Header{
			"Content-Type": []string{"text/event-stream"},
		},
		Body: io.NopCloser(strings.NewReader(strings.Join([]string{
			`data: {"type":"response.completed","response":{"id":"resp_chat_stream","model":"gpt-4.1-mini","output":[{"type":"message","role":"assistant","content":[{"type":"output_text","text":"hello from stream"}]}],"usage":{"input_tokens":7,"output_tokens":9}}}`,
			`data: [DONE]`,
			``,
		}, "\n"))),
	}

	result, err := svc.handleChatStreamingResponse(resp, c, "gpt-4.1-mini", "gpt-4.1-mini", false, time.Now())
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotEmpty(t, result.ResponseBody)
	require.Contains(t, string(result.ResponseBody), "hello from stream")
}

func TestHandleAnthropicBufferedStreamingResponse_ReturnsResponseBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	svc := &OpenAIGatewayService{}
	resp := &http.Response{
		Header: http.Header{
			"Content-Type": []string{"text/event-stream"},
		},
		Body: io.NopCloser(strings.NewReader(strings.Join([]string{
			`data: {"type":"response.completed","response":{"id":"resp_messages_buffered","model":"gpt-4.1-mini","output":[{"type":"message","role":"assistant","content":[{"type":"output_text","text":"hello buffered anthropic"}]}],"usage":{"input_tokens":6,"output_tokens":10}}}`,
			`data: [DONE]`,
			``,
		}, "\n"))),
	}

	result, err := svc.handleAnthropicBufferedStreamingResponse(resp, c, "claude-3.5-sonnet", "gpt-4.1-mini", time.Now())
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotEmpty(t, result.ResponseBody)
	require.Contains(t, string(result.ResponseBody), "hello buffered anthropic")
}

func TestHandleAnthropicStreamingResponse_ReturnsResponseBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	svc := &OpenAIGatewayService{}
	resp := &http.Response{
		Header: http.Header{
			"Content-Type": []string{"text/event-stream"},
		},
		Body: io.NopCloser(strings.NewReader(strings.Join([]string{
			`data: {"type":"response.completed","response":{"id":"resp_messages_stream","model":"gpt-4.1-mini","output":[{"type":"message","role":"assistant","content":[{"type":"output_text","text":"hello anthropic"}]}],"usage":{"input_tokens":5,"output_tokens":8}}}`,
			`data: [DONE]`,
			``,
		}, "\n"))),
	}

	result, err := svc.handleAnthropicStreamingResponse(resp, c, "claude-3.5-sonnet", "gpt-4.1-mini", time.Now())
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotEmpty(t, result.ResponseBody)
	require.Contains(t, string(result.ResponseBody), "hello anthropic")
}

func TestHandleStreamingResponsePassthrough_ReturnsResponseBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	svc := &OpenAIGatewayService{}
	resp := &http.Response{
		Header: http.Header{
			"Content-Type": []string{"text/event-stream"},
		},
		Body: io.NopCloser(strings.NewReader(strings.Join([]string{
			`data: {"type":"response.completed","response":{"id":"resp_passthrough","output":[{"type":"message","role":"assistant","content":[{"type":"output_text","text":"passthrough body"}]}],"usage":{"input_tokens":1,"output_tokens":2}}}`,
			`data: [DONE]`,
			``,
		}, "\n"))),
	}

	result, err := svc.handleStreamingResponsePassthrough(c.Request.Context(), resp, c, &Account{ID: 1}, time.Now())
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotEmpty(t, result.responseBody)
	require.Contains(t, string(result.responseBody), `"id":"resp_passthrough"`)
}
