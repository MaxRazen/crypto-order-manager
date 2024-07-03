package logger

import (
	"context"
	"log/slog"
)

type logHandler struct {
	handler slog.Handler
}

func (h *logHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *logHandler) Handle(ctx context.Context, r slog.Record) error {
	if reqId, ok := ctx.Value("req-id").(string); ok {
		r.AddAttrs(slog.String("req-id", reqId))
	}
	return h.handler.Handle(ctx, r)
}
func (h *logHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &logHandler{handler: h.handler.WithAttrs(attrs)}
}
func (h *logHandler) WithGroup(name string) slog.Handler {
	return &logHandler{handler: h.handler.WithGroup(name)}
}
