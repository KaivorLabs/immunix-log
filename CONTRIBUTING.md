# Contributing

## Issues

File a GitHub issue with the version or commit, what you expected,
what you got. Security issues: see [SECURITY.md](./SECURITY.md).

## Pull requests

- Open an issue first for anything larger than a typo. Saves a
  finished-but-wrong PR.
- Branch off `main`. One concern per branch.
- Imperative commit subjects (`add Logger.With` not `added`).
- No `Co-Authored-By` trailers.
- `make check` before pushing.

## Scope

In: `Logger`, `Nop`, `FromSlog`, and stdlib adapters added in
future Go releases.

Out: concrete loggers, handlers, formatters, level constants, log
shippers, redaction, sampling, third-party adapters (hclog, zap,
logrus, zerolog). Those belong in the consumer's logging package.

## Tests

If you change the interface or `Nop`:

- Exercise every affected method.
- Add a `var _ Logger = ...` compile-time assertion if you add a
  new reference implementation.
- For `With`, confirm the receiver is unaffected and bound
  attributes appear in subsequent records.

A breaking change to the interface breaks every Immunix Go library
that consumes it. Pre-1.0 gets a CHANGELOG migration note, post-1.0
a major version bump.

## License

By contributing, you agree your changes are licensed under
[Apache License 2.0](./LICENSE).
