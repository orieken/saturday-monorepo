
# Claude / LLM Agent Readme

This repository contains `@saturday/playwright-k6-exporter`, a Playwright plugin that records API calls and emits k6 scripts.

## What to Do
- Follow `docs/ARCHITECTURE.md` for data flow and boundaries.
- Obey `COPILOT_INSTRUCTIONS.md` quality gates and coding rules.
- Prefer small, pure functions and extract helpers to keep complexity low.

## Implementation Norms
- Write/modify tests first (TDD).
- Keep domain types (`RecordedCall`) free of Playwright specifics.
- Proxy only HTTP verbs/fetch; do not alter unrelated APIs.
- Normalize headers and body; include `k6Name` if provided.
- Ensure **idempotent** k6 file names from test title slug.

## Output Norms
- Never write secrets. If in doubt, redact or mask.
- Provide clear comments in generated k6 files indicating call purpose.
- Maintain deterministic ordering of calls.
