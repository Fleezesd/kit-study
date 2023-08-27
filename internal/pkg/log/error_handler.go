package log

import (
	"context"
	"fmt"
)

type LogErrorHandler struct {
	logger *ZapLogger
}

func NewLogErrorHandler(logger *ZapLogger) *LogErrorHandler {
	return &LogErrorHandler{
		logger: logger,
	}
}

func (h *LogErrorHandler) Handle(ctx context.Context, err error) {
	h.logger.Warnw(fmt.Sprint(ctx.Value("req_uuid")), "err", err)
}
