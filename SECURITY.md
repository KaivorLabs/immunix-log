# Security Policy

## Reporting

Don't open a public issue.

- [Private vulnerability reporting](https://github.com/KaivorLabs/immunix-log/security/advisories/new)
- `security@immunix.ai` (subject: `[immunix-log]`)

Include the version or commit, a minimal repro, and the observed
impact.

## Response

- Ack within 3 business days.
- Initial triage within 7 days of ack.
- Fix or mitigation plan within 30 days for high/critical.
- Disclosure coordinated with the reporter after a fix ships.

## Supported versions

Pre-1.0, only the latest tagged release receives fixes.

## In scope

- Flaws in the `Logger` interface that make secret leakage easy
  in routine use.
- Bugs in `Nop` or `FromSlog` (retained state, panics, lost `With`
  attributes, dropped records).

## Out of scope

- A consumer library passing secret material as a log value.
  Report to that library.
- `*slog.Logger` handler misconfiguration.
- Bugs in third-party log backends (hclog, zap, logrus, zerolog).
