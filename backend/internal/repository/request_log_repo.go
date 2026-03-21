package repository

import (
	"context"
	"database/sql"
	"strings"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	dbrequestlog "github.com/Wei-Shaw/sub2api/ent/requestlog"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type requestLogRepository struct {
	client *dbent.Client
}

func NewRequestLogRepository(client *dbent.Client, _ *sql.DB) service.RequestLogRepository {
	return &requestLogRepository{client: client}
}

func (r *requestLogRepository) Create(ctx context.Context, log *service.RequestLog) error {
	if log == nil {
		return nil
	}

	tx, err := r.client.Tx(ctx)
	if err != nil {
		return err
	}

	builder := tx.RequestLog.Create().
		SetUserID(log.UserID).
		SetAPIKeyID(log.APIKeyID).
		SetModel(log.Model).
		SetInputTokens(log.InputTokens).
		SetOutputTokens(log.OutputTokens).
		SetTotalCost(log.TotalCost).
		SetStream(log.Stream).
		SetCreatedAt(log.CreatedAt)

	if log.RequestID != nil {
		builder.SetRequestID(strings.TrimSpace(*log.RequestID))
	}
	if log.InboundEndpoint != nil {
		builder.SetInboundEndpoint(strings.TrimSpace(*log.InboundEndpoint))
	}
	if log.UpstreamEndpoint != nil {
		builder.SetUpstreamEndpoint(strings.TrimSpace(*log.UpstreamEndpoint))
	}
	if log.Method != nil {
		builder.SetMethod(strings.TrimSpace(*log.Method))
	}
	if log.StatusCode != nil {
		builder.SetStatusCode(*log.StatusCode)
	}
	if log.ErrorCode != nil {
		builder.SetErrorCode(strings.TrimSpace(*log.ErrorCode))
	}
	if log.ErrorMessage != nil {
		builder.SetErrorMessage(*log.ErrorMessage)
	}
	if log.DurationMs != nil {
		builder.SetDurationMs(*log.DurationMs)
	}
	if log.FirstTokenMs != nil {
		builder.SetFirstTokenMs(*log.FirstTokenMs)
	}
	if log.UserAgent != nil {
		builder.SetUserAgent(*log.UserAgent)
	}
	if log.IPAddress != nil {
		builder.SetIPAddress(*log.IPAddress)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if log.Payload != nil {
		payloadBuilder := tx.RequestLogPayload.Create().
			SetRequestLogID(created.ID).
			SetRequestBodyTruncated(log.Payload.RequestBodyTruncated).
			SetResponseBodyTruncated(log.Payload.ResponseBodyTruncated).
			SetCreatedAt(log.CreatedAt)

		if encoded, encoding, encodeErr := service.EncodeRequestLogPayload(log.Payload.RequestBody); encodeErr != nil {
			_ = tx.Rollback()
			return encodeErr
		} else if len(encoded) > 0 {
			payloadBuilder.SetRequestBody(encoded)
			if encoding != nil {
				payloadBuilder.SetRequestBodyEncoding(*encoding)
			}
		}
		if log.Payload.RequestBodyBytes != nil {
			payloadBuilder.SetRequestBodyBytes(*log.Payload.RequestBodyBytes)
		}

		if encoded, encoding, encodeErr := service.EncodeRequestLogPayload(log.Payload.ResponseBody); encodeErr != nil {
			_ = tx.Rollback()
			return encodeErr
		} else if len(encoded) > 0 {
			payloadBuilder.SetResponseBody(encoded)
			if encoding != nil {
				payloadBuilder.SetResponseBodyEncoding(*encoding)
			}
		}
		if log.Payload.ResponseBodyBytes != nil {
			payloadBuilder.SetResponseBodyBytes(*log.Payload.ResponseBodyBytes)
		}

		if _, err := payloadBuilder.Save(ctx); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.ID = created.ID
	return nil
}

func (r *requestLogRepository) CreateBestEffort(ctx context.Context, log *service.RequestLog) error {
	return r.Create(ctx, log)
}

func (r *requestLogRepository) GetByID(ctx context.Context, id int64) (*service.RequestLog, error) {
	m, err := r.client.RequestLog.Query().
		Where(dbrequestlog.IDEQ(id)).
		WithPayload().
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrRequestLogNotFound, nil)
	}
	return requestLogEntityToService(m), nil
}

func (r *requestLogRepository) ListWithFilters(ctx context.Context, params pagination.PaginationParams, filters service.RequestLogFilters) ([]service.RequestLog, *pagination.PaginationResult, error) {
	q := r.client.RequestLog.Query()

	if filters.UserID > 0 {
		q = q.Where(dbrequestlog.UserIDEQ(filters.UserID))
	}
	if filters.APIKeyID > 0 {
		q = q.Where(dbrequestlog.APIKeyIDEQ(filters.APIKeyID))
	}
	if model := strings.TrimSpace(filters.Model); model != "" {
		q = q.Where(dbrequestlog.ModelContainsFold(model))
	}
	if filters.StatusCode != nil {
		q = q.Where(dbrequestlog.StatusCodeEQ(*filters.StatusCode))
	}
	if filters.StartTime != nil {
		q = q.Where(dbrequestlog.CreatedAtGTE(*filters.StartTime))
	}
	if filters.EndTime != nil {
		q = q.Where(dbrequestlog.CreatedAtLT(*filters.EndTime))
	}

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	rows, err := q.
		Offset(params.Offset()).
		Limit(params.Limit()).
		Order(dbent.Desc(dbrequestlog.FieldCreatedAt), dbent.Desc(dbrequestlog.FieldID)).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}

	out := make([]service.RequestLog, 0, len(rows))
	for _, row := range rows {
		out = append(out, *requestLogEntityToService(row))
	}
	return out, paginationResultFromTotal(int64(total), params), nil
}

func requestLogEntityToService(m *dbent.RequestLog) *service.RequestLog {
	if m == nil {
		return nil
	}
	return &service.RequestLog{
		ID:               m.ID,
		UserID:           m.UserID,
		APIKeyID:         m.APIKeyID,
		RequestID:        m.RequestID,
		Model:            m.Model,
		InboundEndpoint:  m.InboundEndpoint,
		UpstreamEndpoint: m.UpstreamEndpoint,
		Method:           m.Method,
		StatusCode:       m.StatusCode,
		ErrorCode:        m.ErrorCode,
		ErrorMessage:     m.ErrorMessage,
		InputTokens:      m.InputTokens,
		OutputTokens:     m.OutputTokens,
		TotalCost:        m.TotalCost,
		Stream:           m.Stream,
		DurationMs:       m.DurationMs,
		FirstTokenMs:     m.FirstTokenMs,
		UserAgent:        m.UserAgent,
		IPAddress:        m.IPAddress,
		Payload:          requestLogPayloadEntityToService(m.Edges.Payload),
		CreatedAt:        m.CreatedAt,
	}
}

func requestLogPayloadEntityToService(m *dbent.RequestLogPayload) *service.RequestLogPayload {
	if m == nil {
		return nil
	}

	requestBody, err := service.DecodeRequestLogPayload(m.RequestBody, m.RequestBodyEncoding)
	if err != nil {
		requestBody = nil
	}
	responseBody, err := service.DecodeRequestLogPayload(m.ResponseBody, m.ResponseBodyEncoding)
	if err != nil {
		responseBody = nil
	}

	return &service.RequestLogPayload{
		RequestBody:           requestBody,
		RequestBodyBytes:      m.RequestBodyBytes,
		RequestBodyTruncated:  m.RequestBodyTruncated,
		ResponseBody:          responseBody,
		ResponseBodyBytes:     m.ResponseBodyBytes,
		ResponseBodyTruncated: m.ResponseBodyTruncated,
	}
}
