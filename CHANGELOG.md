# Changelog

Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
[Semantic Versioning](https://semver.org/) applies from v1.0.0.

## [Unreleased]

## [0.1.0] - 2026-05-28

- `Logger` interface: `DebugContext`, `InfoContext`, `WarnContext`,
  `ErrorContext` (each takes `context.Context` first), plus
  `With(args ...any) Logger`.
- `Nop` no-op implementation. Zero value usable; `Nop.With`
  returns the receiver.
- `NewLogger(slog.Handler) Logger` builds a Logger from any
  `slog.Handler`. `NewLogger(nil)` returns `Nop`.
- `FromSlog(*slog.Logger) Logger` wraps a pre-built `*slog.Logger`.
  `FromSlog(nil)` returns `Nop`.
