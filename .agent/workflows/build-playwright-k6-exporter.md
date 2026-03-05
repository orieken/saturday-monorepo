---
description: Build and verify the playwright-k6-exporter package
---

1. Navigate to the package directory
cd packages/playwright-k6-exporter

2. Install dependencies
npm install

3. Build the package
npm run build

4. Verify the build output
ls -l dist/index.cjs dist/index.mjs
