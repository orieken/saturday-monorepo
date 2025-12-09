
# @saturday/k6-redaction-basic
Basic redaction policy plugin for `@saturday/playwright-k6-exporter`.

- Detects common secret headers and sensitive body fields.
- Replaces values with `${ENV}` placeholders and returns findings.

## Usage
```ts
import { createDefaultRedactionPolicy } from '@saturday/k6-redaction-basic';
const policy = createDefaultRedactionPolicy({ envPrefix: 'K6_' });
```
