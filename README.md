# immunix-log

[![Go Reference](https://pkg.go.dev/badge/github.com/KaivorLabs/immunix-log.svg)](https://pkg.go.dev/github.com/KaivorLabs/immunix-log)
[![License: Apache-2.0](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](./LICENSE)

The `Logger` interface used by the Immunix Go libraries.

The module ships the interface, a no-op default, and a
`*slog.Logger` adapter.

## Install

```bash
go get github.com/KaivorLabs/immunix-log@latest
```

Requires Go 1.26+.

## Interface

```go
type Logger interface {
    DebugContext(ctx context.Context, msg string, args ...any)
    InfoContext(ctx context.Context, msg string, args ...any)
    WarnContext(ctx context.Context, msg string, args ...any)
    ErrorContext(ctx context.Context, msg string, args ...any)
    With(args ...any) Logger
}
```

`args` are key/value pairs (slog convention).

## With slog

```go
import (
    "log/slog"
    "os"

    immunixlog "github.com/KaivorLabs/immunix-log"
    vault "github.com/KaivorLabs/immunix-broker-vault"
)

logger := immunixlog.NewLogger(slog.NewJSONHandler(os.Stderr, nil))
v, err := vault.New(ctx, conf, vault.WithLogger(logger))
```

Already have a `*slog.Logger`? Wrap it with `FromSlog(l)`. Both
return `Nop` on a nil argument.

## With another backend

Implement the five methods on an adapter type:

```go
type hclogAdapter struct{ inner hclog.Logger }

func (a hclogAdapter) InfoContext(_ context.Context, msg string, args ...any) {
    a.inner.Info(msg, args...)
}
// ... three more, plus:
func (a hclogAdapter) With(args ...any) immunixlog.Logger {
    return hclogAdapter{a.inner.With(args...)}
}
```

## Nop

```go
var logger immunixlog.Logger = immunixlog.Nop{}
```

Zero value usable. `Nop.With(...)` returns the receiver.

## Security

Don't pass keys, sealed blobs, credentials, or tokens as a log
value. Report vulnerabilities via [private advisory](https://github.com/KaivorLabs/immunix-log/security/advisories/new)
or `security@immunix.ai`; see [`SECURITY.md`](./SECURITY.md).

## License

[Apache 2.0](./LICENSE).
