Cucumber Indexer (Final Project)

This is a tiny TypeScript CLI that scans `.feature` files and produces a `CucumberIndex` JSON compatible with the Final backend. It can also POST the generated index to the backend.

Install / Build

- cd final/tools/cucumber-indexer
- npm install
- npm run build

Usage

- Generate an index JSON from final/tests/features into final/backend/data:
  - node dist/index.js --features ../../tests/features --suiteId demo-ecommerce --out ../../backend/data/cucumber_index.json

- Optionally POST directly to a running backend:
  - node dist/index.js --features ../../tests/features --suiteId demo-ecommerce --out ../../backend/data/cucumber_index.json --backend http://localhost:8080

CLI Flags

- --features <dir>: Root folder to scan for .feature files. Defaults to final/tests/features based on repo layout.
- --suiteId <string>: Suite identifier to embed in the index. Default: demo-ecommerce.
- --out <path>: Output JSON file. Default: final/backend/data/cucumber_index.json based on repo layout.
- --backend <url>: When provided, the tool will POST the generated index to <url>/api/cucumber/index.

Notes

- The parser is deliberately simple (keyword-based). For robustness, consider swapping in a real Gherkin parser later.
- When running with the local dev cluster (final/local-cluster), start the backend first, then run the indexer with --backend http://localhost:8080 to register the suite.
