package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
)

const (
	requestLogCleanupWorkerName = "request_log_cleanup_worker"
)

type RequestLogCleanupService struct {
	db             *sql.DB
	timingWheel    *TimingWheelService
	settingService *SettingService

	running   int32
	startOnce sync.Once
	stopOnce  sync.Once

	workerCtx    context.Context
	workerCancel context.CancelFunc
}

func NewRequestLogCleanupService(db *sql.DB, timingWheel *TimingWheelService, settingService *SettingService) *RequestLogCleanupService {
	workerCtx, workerCancel := context.WithCancel(context.Background())
	return &RequestLogCleanupService{
		db:             db,
		timingWheel:    timingWheel,
		settingService: settingService,
		workerCtx:      workerCtx,
		workerCancel:   workerCancel,
	}
}

func (s *RequestLogCleanupService) Start() {
	if s == nil || s.db == nil || s.timingWheel == nil || s.settingService == nil {
		logger.LegacyPrintf("service.request_log_cleanup", "[RequestLogCleanup] not started (missing deps)")
		return
	}

	s.startOnce.Do(func() {
		interval := 6 * time.Hour
		s.timingWheel.ScheduleRecurring(requestLogCleanupWorkerName, interval, s.runOnce)
		logger.LegacyPrintf("service.request_log_cleanup", "[RequestLogCleanup] started (interval=%s)", interval)
	})
}

func (s *RequestLogCleanupService) Stop() {
	if s == nil {
		return
	}
	s.stopOnce.Do(func() {
		if s.workerCancel != nil {
			s.workerCancel()
		}
		if s.timingWheel != nil {
			s.timingWheel.Cancel(requestLogCleanupWorkerName)
		}
		logger.LegacyPrintf("service.request_log_cleanup", "[RequestLogCleanup] stopped")
	})
}

func (s *RequestLogCleanupService) runOnce() {
	if s == nil || s.db == nil || s.settingService == nil {
		return
	}
	if !atomic.CompareAndSwapInt32(&s.running, 0, 1) {
		return
	}
	defer atomic.StoreInt32(&s.running, 0)

	parent := context.Background()
	if s.workerCtx != nil {
		parent = s.workerCtx
	}
	ctx, cancel := context.WithTimeout(parent, 30*time.Minute)
	defer cancel()

	settings, err := s.settingService.GetRequestLogSettings(ctx)
	if err != nil {
		logger.LegacyPrintf("service.request_log_cleanup", "[RequestLogCleanup] load settings failed: %v", err)
		return
	}
	if settings.RetentionDays <= 0 {
		return
	}

	cutoff := time.Now().UTC().AddDate(0, 0, -settings.RetentionDays)
	const batchSize = 5000

	payloadRows, err := deleteOldRowsByID(ctx, s.db, "request_log_payloads", "created_at", cutoff, batchSize, false)
	if err != nil {
		logger.LegacyPrintf("service.request_log_cleanup", "[RequestLogCleanup] delete payload rows failed: %v", err)
		return
	}

	logRows, err := deleteOldRowsByID(ctx, s.db, "request_logs", "created_at", cutoff, batchSize, false)
	if err != nil {
		logger.LegacyPrintf("service.request_log_cleanup", "[RequestLogCleanup] delete request rows failed: %v", err)
		return
	}

	if payloadRows > 0 || logRows > 0 {
		logger.LegacyPrintf(
			"service.request_log_cleanup",
			"[RequestLogCleanup] cleanup complete: %s",
			strings.TrimSpace(fmt.Sprintf("request_log_payloads=%d request_logs=%d", payloadRows, logRows)),
		)
	}
}
