# Testing Strategy & CI

## Unit Testing
- Table-driven tests, mocks, aim 85% coverage per package.

## Integration Testing
- Temp directories, verify generated TypeScript compiles, MCP handshake tests.

## Makefile (targets)
- `build`, `test`, `coverage`, `lint`, `clean`, `install`

## CI / Release
- GitHub Actions build matrix for linux/darwin/windows and archs; produce release artifacts.

