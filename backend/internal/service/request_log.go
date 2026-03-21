package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

var (
	ErrRequestLogNotFound = infraerrors.NotFound("REQUEST_LOG_NOT_FOUND", "request log not found")
)

type RequestLog struct {
	ID int64

	UserID   int64
	APIKeyID int64

	RequestID        *string
	Model            string
	InboundEndpoint  *string
	UpstreamEndpoint *string
	Method           *string
	StatusCode       *int
	ErrorCode        *string
	ErrorMessage     *string

	InputTokens  int
	OutputTokens int
	TotalCost    float64

	Stream       bool
	DurationMs   *int
	FirstTokenMs *int
	UserAgent    *string
	IPAddress    *string

	CreatedAt time.Time
}

type RequestLogFilters struct {
	UserID     int64
	APIKeyID   int64
	Model      string
	StatusCode *int
	StartTime  *time.Time
	EndTime    *time.Time
}

type RequestLogRepository interface {
	Create(ctx context.Context, log *RequestLog) error
	CreateBestEffort(ctx context.Context, log *RequestLog) error
	GetByID(ctx context.Context, id int64) (*RequestLog, error)
	ListWithFilters(ctx context.Context, params pagination.PaginationParams, filters RequestLogFilters) ([]RequestLog, *pagination.PaginationResult, error)
}

type RequestLogService struct {
	repo RequestLogRepository
}

func NewRequestLogService(repo RequestLogRepository) *RequestLogService {
	return &RequestLogService{repo: repo}
}

func (s *RequestLogService) Create(ctx context.Context, log *RequestLog) error {
	if s == nil || s.repo == nil || log == nil {
		return nil
	}
	return s.repo.Create(ctx, log)
}

func (s *RequestLogService) GetByID(ctx context.Context, id int64) (*RequestLog, error) {
	if s == nil || s.repo == nil {
		return nil, fmt.Errorf("request log service not configured")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *RequestLogService) ListWithFilters(ctx context.Context, params pagination.PaginationParams, filters RequestLogFilters) ([]RequestLog, *pagination.PaginationResult, error) {
	if s == nil || s.repo == nil {
		return nil, nil, fmt.Errorf("request log service not configured")
	}
	return s.repo.ListWithFilters(ctx, params, filters)
}

func writeRequestLogBestEffort(ctx context.Context, repo RequestLogRepository, requestLog *RequestLog, logKey string) {
	if repo == nil || requestLog == nil {
		return
	}
	requestCtx, cancel := detachedBillingContext(ctx)
	defer cancel()

	if err := repo.CreateBestEffort(requestCtx, requestLog); err != nil {
		logger.LegacyPrintf(logKey, "Create request log failed: %v", err)
		if syncErr := repo.Create(requestCtx, requestLog); syncErr != nil {
			logger.LegacyPrintf(logKey, "Create request log sync fallback failed: %v", syncErr)
		}
	}
}

func newSuccessRequestLog(
	userID int64,
	apiKeyID int64,
	requestID string,
	model string,
	inboundEndpoint string,
	upstreamEndpoint string,
	userAgent string,
	ipAddress string,
	stream bool,
	durationMs *int,
	firstTokenMs *int,
	inputTokens int,
	outputTokens int,
	totalCost float64,
) *RequestLog {
	method := http.MethodPost
	statusCode := http.StatusOK

	return &RequestLog{
		UserID:           userID,
		APIKeyID:         apiKeyID,
		RequestID:        optionalTrimmedStringPtr(requestID),
		Model:            model,
		InboundEndpoint:  optionalTrimmedStringPtr(inboundEndpoint),
		UpstreamEndpoint: optionalTrimmedStringPtr(upstreamEndpoint),
		Method:           &method,
		StatusCode:       &statusCode,
		InputTokens:      inputTokens,
		OutputTokens:     outputTokens,
		TotalCost:        totalCost,
		Stream:           stream,
		DurationMs:       durationMs,
		FirstTokenMs:     firstTokenMs,
		UserAgent:        optionalTrimmedStringPtr(userAgent),
		IPAddress:        optionalTrimmedStringPtr(ipAddress),
		CreatedAt:        time.Now(),
	}
}
