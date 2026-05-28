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
