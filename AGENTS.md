# Mercury Client Go – Agent Notes

## Context
- This repository is a trimmed extraction of the broader [sprucelabsai-community/mercury-workspace](https://github.com/sprucelabsai-community/mercury-workspace) that focuses exclusively on the Go SDK.
- Aim for functional parity with the TypeScript client in the workspace; reference it whenever behaviour or payload structure is unclear.
- Keep references to shared schemas up to date (`spruce-core-schemas` and friends) so responses stay compatible with the Mercury event bus.

## Environment Setup
- Mirror the workspace convention: tests expect a `.env` file at the repo root with `TEST_HOST` pointing to the running Mercury Theatre (e.g. `http://127.0.0.1:8081`).
- `factory_test.go` loads that file via `godotenv`, so missing or stale values lead to connection failures and nil pointer panics.
- The Mercury Theatre is the local mono-repo environment that boots the `spruce-mercury-skill` (the Mercury API) along with supporting skills; keep it running on localhost while developing or exercising integration tests here.
- This client talks directly to that Mercury API (`spruce-mercury-skill`) at whatever `TEST_HOST` points to—normally the theatre’s default `http://127.0.0.1:8081`.

## Preferred Workflow
- Follow the three laws of TDD:
  1. Write no production code until you have a failing test that describes the desired behaviour.
  2. Write only enough test to see that failure (compilation counts as failure).
  3. Write only the minimal production code required to make that test pass.
- Keep the red → green → refactor cadence tight; resist bulk changes without a guiding test.

## Critical Protocol
- This document is authoritative; treat every instruction here as mandatory unless the user explicitly overrides it.
- No AI agent will make any filesystem changes without first presenting the plan to the user and obtaining explicit sign-off.
