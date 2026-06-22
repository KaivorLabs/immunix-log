// SPDX-License-Identifier: Apache-2.0

package immunixlog

import (
	"context"
	"log/slog"
)

// Logger is the logging interface the Immunix Go libraries accept.
// Methods take ctx first and slog-style key/value args. Wrap a
// *slog.Logger via FromSlog. Implementations must be safe for
// concurrent use; callers must not pass secret material as a value.
type Logger interface {
	DebugContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)

	// Debug/Info/Warn/Error are the convenience variants for callers
	// with no request context (boot, background loops). They mirror the
	// same-named slog.Logger methods and behave as the *Context form
	// with context.Background(). slogAdapter satisfies them via the
	// embedded *slog.Logger (method promotion), so slog's caller-PC skip
	// stays correct and source= still points at the real call site —
	// implementations must NOT hand-write these as wrappers around the
	// *Context form.
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)

	// With returns a Logger that adds args to every subsequent
	// record.
	With(args ...any) Logger
}

// Nop is a no-op Logger. Its zero value is usable; With returns
// the receiver.
type Nop struct{}

func (Nop) DebugContext(context.Context, string, ...any) {}
func (Nop) InfoContext(context.Context, string, ...any)  {}
func (Nop) WarnContext(context.Context, string, ...any)  {}
func (Nop) ErrorContext(context.Context, string, ...any) {}
func (Nop) Debug(string, ...any)                         {}
func (Nop) Info(string, ...any)                          {}
func (Nop) Warn(string, ...any)                          {}
func (Nop) Error(string, ...any)                         {}
func (n Nop) With(...any) Logger                         { return n }

// NewLogger returns a Logger that writes via h. A nil h returns Nop.
func NewLogger(h slog.Handler) Logger {
	if h == nil {
		return Nop{}
	}
	return slogAdapter{slog.New(h)}
}

// FromSlog wraps a *slog.Logger so it satisfies Logger. A nil l
// returns Nop.
func FromSlog(l *slog.Logger) Logger {
	if l == nil {
		return Nop{}
	}
	return slogAdapter{l}
}

type slogAdapter struct{ *slog.Logger }

func (s slogAdapter) With(args ...any) Logger {
	return slogAdapter{s.Logger.With(args...)}
}

var (
	_ Logger = Nop{}
	_ Logger = slogAdapter{}
)
