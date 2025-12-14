# @saturday/playwright-k6-exporter

Record Playwright API calls during tests and automatically generate **k6** scripts you can use for performance testing.

## Install

```bash
npm i -D @saturday/playwright-k6-exporter
```

## Quick Start

1) **Scaffold an example test** (optional helper):
```bash
npx k6-exporter-init
```

2) **Write/Use a test with the fixture**:

```ts
// e2e/api/my-flow.spec.ts
import { expect } from '@playwright/test';
import { test } from '@saturday/playwright-k6-exporter/fixture';
import { createK6Recorder } from '@saturday/playwright-k6-exporter';

test('users flow @k6', async ({}, testInfo) => {
  const setup = await createK6Recorder(testInfo.title);
  if (!setup) test.skip(true, 'Set K6_EXPORT=1 to export k6 script');
  const { ctx, recorder } = setup!;

  const res = await ctx.get('https://api.example.com/users', { k6Name: 'list users' });
  expect(res.status()).toBe(200);

  await recorder.flushToK6();
  await ctx.dispose();
});
```

3) **Export k6**:
```bash
# only tests tagged with @k6
K6_EXPORT=1 npm run k6:export

# export from entire suite
K6_EXPORT=1 npm run k6:export:all
```

k6 files are written to `perf/k6/<test-title-slug>.k6.js`.

## API

```ts
import { createK6Recorder, K6Recorder } from '@saturday/playwright-k6-exporter';

// create a wrapped APIRequestContext that records calls
const setup = await createK6Recorder(testInfo.title, 'perf/k6');
const { ctx, recorder } = setup!;

// use ctx.get/ctx.post/... as usual
await ctx.get('https://httpbin.org/get', { k6Name: 'Get hello world' });

// write file
await recorder.flushToK6();
await ctx.dispose();
```

### `k6Name` Option
Provide `k6Name` in the request options to give each call a friendly comment in the generated file.

## Notes
- This package does **not** bundle Playwright. You must install `@playwright/test` in your project.
- Re-running a test overwrites the same k6 file for that test title (idempotent updates).
- Consider redacting secrets (tokens/PII) before committing generated scripts.
- You can customize the generated k6 `options` by post-processing the file, or by extending the generator.

## License
MIT

---

## Documentation & Agent Guides
- Architecture: `docs/ARCHITECTURE.md`
- Copilot/Agent Rules: `COPILOT_INSTRUCTIONS.md`
- Claude Agent Readme: `CLAUDE_README.md`
- Contributing: `CONTRIBUTING.md`

**Quality Gates:** We enforce clean code with cyclomatic complexity < 7, short functions, and TDD/BDD as standard practice.


---

## Redaction Policy & Env Artifacts

This package supports a **pluggable redaction policy** to keep secrets out of generated files.

- Default policy: detects common secret headers (`authorization`, `x-api-key`, `cookie`, etc.) and body fields (`*.token`, `*.password`, `*.apiKey`).
- On detection, values are replaced with `${ENV_VAR}` placeholders and:
  - written to `.env.apis` with their discovered values (local only),
  - appended to `.env.example` (as empty keys) so teams know which env vars are required.

### Usage (default policy via fixture)
The provided fixture already enables the default redaction policy:
```ts
import { test } from '@saturday/playwright-k6-exporter/fixture';
// ... your @k6 tests
```

### Custom policy
You can provide your own policy:
```ts
import { createK6Recorder } from '@saturday/playwright-k6-exporter';
import type { RedactionPolicy } from '@saturday/playwright-k6-exporter/dist/policies/defaultRedaction';

const policy: RedactionPolicy = {
  redactHeader: (name, value) => name.toLowerCase()==='authorization'
    ? { value: '${K6_AUTHORIZATION}', finding: { envName: 'K6_AUTHORIZATION', value, source: 'header', path: `headers.${name}` } }
    : null,
  redactBody: (path, val) => /password$/i.test(path)
    ? { value: '${K6_PASSWORD}', finding: { envName: 'K6_PASSWORD', value: String(val), source: 'body', path } }
    : null
};

const setup = await createK6Recorder(testInfo.title, 'perf/k6', { policy });
```

> **Security Note:** `.env.apis` is written in the project root by default and should be **git-ignored**. It contains raw values discovered during test execution. `.env.example` only lists empty keys for teammates.

## Development Workflows

For contributors and agents working on this package, we have defined workflows to streamline common tasks.

### Agent Workflows
Location: `.agent/workflows/`

- **Build and Verify**: `/build-playwright-k6-exporter`
  - Installs dependencies
  - Builds the package
  - Verifies the output

### Manual Development Steps

1. **Install Dependencies**:
   ```bash
   npm install
   ```

2. **Build**:
   ```bash
   npm run build
   ```

3. **Link for Local Testing**:
   ```bash
   npm link
   ```

4. **Run Tests** (if available):
   ```bash
   npm test
   ```
