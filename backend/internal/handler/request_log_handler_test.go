package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type userRequestLogRepoCapture struct {
	service.RequestLogRepository

	listFilters service.RequestLogFilters
	getLog      *service.RequestLog
}

func (s *userRequestLogRepoCapture) ListWithFilters(ctx context.Context, params pagination.PaginationParams, filters service.RequestLogFilters) ([]service.RequestLog, *pagination.PaginationResult, error) {
	s.listFilters = filters
	return []service.RequestLog{}, &pagination.PaginationResult{
		Total:    0,
		Page:     params.Page,
		PageSize: params.PageSize,
		Pages:    0,
	}, nil
}

func (s *userRequestLogRepoCapture) GetByID(ctx context.Context, id int64) (*service.RequestLog, error) {
	return s.getLog, nil
}

type apiKeyRepoForRequestLogTest struct {
	service.APIKeyRepository
	apiKey *service.APIKey
}

func (s *apiKeyRepoForRequestLogTest) GetByID(ctx context.Context, id int64) (*service.APIKey, error) {
	return s.apiKey, nil
}

func newUserRequestLogTestRouter(repo *userRequestLogRepoCapture, apiKey *service.APIKey) *gin.Engine {
	gin.SetMode(gin.TestMode)
	requestLogSvc := service.NewRequestLogService(repo)
	apiKeySvc := service.NewAPIKeyService(&apiKeyRepoForRequestLogTest{apiKey: apiKey}, nil, nil, nil, nil, nil, nil)
	handler := NewRequestLogHandler(requestLogSvc, apiKeySvc)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(string(middleware2.ContextKeyUser), middleware2.AuthSubject{UserID: 42})
		c.Next()
	})
	router.GET("/request-logs", handler.List)
	router.GET("/request-logs/:id", handler.GetByID)
	return router
}

func TestUserRequestLogListUsesCurrentUserScope(t *testing.T) {
	repo := &userRequestLogRepoCapture{}
	router := newUserRequestLogTestRouter(repo, &service.APIKey{ID: 9, UserID: 42})

	req := httptest.NewRequest(http.MethodGet, "/request-logs?api_key_id=9&model=gpt-5.4", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, int64(42), repo.listFilters.UserID)
	require.Equal(t, int64(9), repo.listFilters.APIKeyID)
	require.Equal(t, "gpt-5.4", repo.listFilters.Model)
}

func TestUserRequestLogGetByIDRejectsForeignRecord(t *testing.T) {
	repo := &userRequestLogRepoCapture{
		getLog: &service.RequestLog{
			ID:     1,
			UserID: 99,
		},
	}
	router := newUserRequestLogTestRouter(repo, &service.APIKey{ID: 9, UserID: 42})

	req := httptest.NewRequest(http.MethodGet, "/request-logs/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusForbidden, rec.Code)
}
