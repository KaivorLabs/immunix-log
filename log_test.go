// SPDX-License-Identifier: Apache-2.0

package immunixlog_test

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	immunixlog "github.com/KaivorLabs/immunix-log"
)

func TestNop(t *testing.T) {
	var l immunixlog.Logger = immunixlog.Nop{}
	ctx := context.Background()
	l.DebugContext(ctx, "d", "k", "v")
	l.InfoContext(ctx, "i")
	l.WarnContext(ctx, "w")
	l.ErrorContext(ctx, "e")
	l.With("k", "v").InfoContext(ctx, "after with")
}

func TestNopZeroValue(t *testing.T) {
	var n immunixlog.Nop
	n.InfoContext(context.Background(), "zero value")
}

func TestFromSlogNil(t *testing.T) {
	l := immunixlog.FromSlog(nil)
	if l == nil {
		t.Fatal("FromSlog(nil) must return Nop, not a nil interface")
	}
	l.InfoContext(context.Background(), "nil slog accepted")
}

func TestNewLogger(t *testing.T) {
	var buf bytes.Buffer
	l := immunixlog.NewLogger(slog.NewTextHandler(&buf, nil))
	l.InfoContext(context.Background(), "from-newlogger", "k", "v")
	for _, want := range []string{"from-newlogger", "k=v"} {
		if !strings.Contains(buf.String(), want) {
			t.Errorf("output missing %q; got %q", want, buf.String())
		}
	}
}

func TestNewLoggerNil(t *testing.T) {
	l := immunixlog.NewLogger(nil)
	if l == nil {
		t.Fatal("NewLogger(nil) must return Nop, not a nil interface")
	}
	l.InfoContext(context.Background(), "nil handler accepted")
}

func TestFromSlogRoutesAllLevels(t *testing.T) {
	var buf bytes.Buffer
	l := immunixlog.FromSlog(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))

	ctx := context.Background()
	l.DebugContext(ctx, "debug-line")
	l.InfoContext(ctx, "info-line")
	l.WarnContext(ctx, "warn-line")
	l.ErrorContext(ctx, "error-line")

	for _, want := range []string{"debug-line", "info-line", "warn-line", "error-line"} {
		if !strings.Contains(buf.String(), want) {
			t.Errorf("output missing %q; got %q", want, buf.String())
		}
	}
}

func TestFromSlogWith(t *testing.T) {
	var buf bytes.Buffer
	l := immunixlog.FromSlog(slog.New(slog.NewTextHandler(&buf, nil))).
		With("backend", "aws_kms", "key_id", "arn:aws:kms:...")
	l.InfoContext(context.Background(), "kms: ready")

	for _, want := range []string{"backend=aws_kms", "key_id=arn:aws:kms:...", "kms: ready"} {
		if !strings.Contains(buf.String(), want) {
			t.Errorf("output missing %q; got %q", want, buf.String())
		}
	}
}

// TestConvenienceMethodsRouteAllLevels exercises the non-context
// Debug/Info/Warn/Error convenience methods so a regression that drops
// one from slogAdapter's promoted set (or stubs it on Nop) fails here.
func TestConvenienceMethodsRouteAllLevels(t *testing.T) {
	var buf bytes.Buffer
	l := immunixlog.NewLogger(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	l.Debug("debug-plain")
	l.Info("info-plain")
	l.Warn("warn-plain")
	l.Error("error-plain")
	l.With("k", "v").Info("with-plain")

	for _, want := range []string{"debug-plain", "info-plain", "warn-plain", "error-plain", "with-plain", "k=v"} {
		if !strings.Contains(buf.String(), want) {
			t.Errorf("output missing %q; got %q", want, buf.String())
		}
	}
}

// TestConvenienceMethodsSourceAttribution locks the load-bearing
// guarantee: the convenience methods are satisfied via slogAdapter's
// embedded *slog.Logger (method promotion), NOT hand-written wrappers
// that call the *Context form. A wrapper would add a stack frame and
// make slog's hardcoded caller skip resolve source= to log.go instead
// of the real call site. With AddSource on, the emitted source.file
// must be THIS test file, never the library's log.go.
func TestConvenienceMethodsSourceAttribution(t *testing.T) {
	var buf bytes.Buffer
	l := immunixlog.NewLogger(slog.NewTextHandler(&buf, &slog.HandlerOptions{AddSource: true}))
	l.Info("attrib-check")

	out := buf.String()
	if !strings.Contains(out, "attrib-check") {
		t.Fatalf("message missing: %q", out)
	}
	if strings.Contains(out, "log.go") {
		t.Errorf("source attribution leaked into the library (promotion broken — convenience method is wrapping, not promoted): %q", out)
	}
	if !strings.Contains(out, "log_test.go") {
		t.Errorf("source should point at the real call site (log_test.go); got %q", out)
	}
}

// With must not mutate the receiver: a derived logger's attributes
// don't leak back to the root.
func TestFromSlogWithIsImmutable(t *testing.T) {
	var buf bytes.Buffer
	root := immunixlog.FromSlog(slog.New(slog.NewTextHandler(&buf, nil)))
	_ = root.With("scope", "branch")

	root.InfoContext(context.Background(), "root")
	if strings.Contains(buf.String(), "scope=branch") {
		t.Errorf("root logger leaked With attribute: %q", buf.String())
	}
}
