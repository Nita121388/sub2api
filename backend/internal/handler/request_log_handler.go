package handler

import (
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/pkg/timezone"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

type RequestLogHandler struct {
	requestLogService *service.RequestLogService
	apiKeyService     *service.APIKeyService
}

func NewRequestLogHandler(requestLogService *service.RequestLogService, apiKeyService *service.APIKeyService) *RequestLogHandler {
	return &RequestLogHandler{
		requestLogService: requestLogService,
		apiKeyService:     apiKeyService,
	}
}

func (h *RequestLogHandler) List(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	page, pageSize := response.ParsePagination(c)

	var apiKeyID int64
	if apiKeyIDStr := strings.TrimSpace(c.Query("api_key_id")); apiKeyIDStr != "" {
		id, err := strconv.ParseInt(apiKeyIDStr, 10, 64)
		if err != nil {
			response.BadRequest(c, "Invalid api_key_id")
			return
		}
		apiKey, err := h.apiKeyService.GetByID(c.Request.Context(), id)
		if err != nil {
			response.ErrorFrom(c, err)
			return
		}
		if apiKey.UserID != subject.UserID {
			response.Forbidden(c, "Not authorized to access this API key's request logs")
			return
		}
		apiKeyID = id
	}

	var statusCode *int
	if rawStatusCode := strings.TrimSpace(c.Query("status_code")); rawStatusCode != "" {
		parsed, err := strconv.Atoi(rawStatusCode)
		if err != nil {
			response.BadRequest(c, "Invalid status_code")
			return
		}
		statusCode = &parsed
	}

	var startTime, endTime *time.Time
	userTZ := c.Query("timezone")
	if startDateStr := strings.TrimSpace(c.Query("start_date")); startDateStr != "" {
		t, err := timezone.ParseInUserLocation("2006-01-02", startDateStr, userTZ)
		if err != nil {
			response.BadRequest(c, "Invalid start_date format, use YYYY-MM-DD")
			return
		}
		startTime = &t
	}
	if endDateStr := strings.TrimSpace(c.Query("end_date")); endDateStr != "" {
		t, err := timezone.ParseInUserLocation("2006-01-02", endDateStr, userTZ)
		if err != nil {
			response.BadRequest(c, "Invalid end_date format, use YYYY-MM-DD")
			return
		}
		t = t.AddDate(0, 0, 1)
		endTime = &t
	}

	logs, result, err := h.requestLogService.ListWithFilters(c.Request.Context(), pagination.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}, service.RequestLogFilters{
		UserID:     subject.UserID,
		APIKeyID:   apiKeyID,
		Model:      strings.TrimSpace(c.Query("model")),
		StatusCode: statusCode,
		StartTime:  startTime,
		EndTime:    endTime,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Paginated(c, logs, result.Total, result.Page, result.PageSize)
}

func (h *RequestLogHandler) GetByID(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid request log ID")
		return
	}

	log, err := h.requestLogService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	if log.UserID != subject.UserID {
		response.Forbidden(c, "Not authorized to access this request log")
		return
	}

	response.Success(c, log)
}
