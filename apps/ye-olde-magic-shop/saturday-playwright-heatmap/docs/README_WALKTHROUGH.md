# Test Execution Heatmap Walkthrough

I have implemented the `saturday-playwright-heatmap` package to visualize test coverage by tracking interactions and interactable elements.

## Features
1.  **Interaction Tracking**: Records clicks, inputs, and changes during tests.
2.  **Element Scanning**: Identifies all interactable elements (buttons, links, inputs) at the end of each test to show what was *available* vs what was *touched*.
3.  **Visual Report**: Generates an HTML report overlaying interactions and interactables on a page screenshot.

## Usage

### 1. Installation
The package is installed in the monorepo.
```bash
pnpm install
```

### 2. Integration
In your test file, import the `test` fixture from the package:
```typescript
import { test } from '@orieken/saturday-playwright-heatmap';

test('my feature', async ({ page, heatmap }) => {
  // Your test code...
});
```

### 3. Running Tests
Run your Playwright tests as usual. The fixture will automatically capture data into `heatmap-data/`.
```bash
pnpm exec playwright test
```

### 4. Generating the Report
Run the reporter script to generate the HTML file:
```bash
node ../../packages/saturday-playwright-heatmap/dist/reporter.js heatmap-data heatmap-report.html
```

## Example Result
A verification test was created at `apps/ye-olde-magic-shop/tests/heatmap.spec.ts`.
Running it generated a report showing:
*   **Cyan Boxes**: Detectable interactable elements.
*   **Red Dots**: Places where the test performed an interaction.

This allows you to instantly see which buttons or links were missed by your test suite.

## Demonstration

### Heatmap Report Screenshot
![Heatmap Report](/Users/oscarrieken/.gemini/antigravity/brain/0515f9cb-b45f-4380-9e59-74a1d56a731c/heatmap_report_view_1766348506561.png)

### Report Navigation Recording
![Report Navigation](/Users/oscarrieken/.gemini/antigravity/brain/0515f9cb-b45f-4380-9e59-74a1d56a731c/heatmap_demo_1766348488145.webp)

## ML Analysis
We have added a Machine Learning Analyzer to identifying "Cold Spots" (interactable elements that are rarely or never clicked).

### Running Analysis
```bash
node ../../packages/saturday-ml-analyzer/dist/index.js heatmap-data
```

### Example Output
```text
Analyzing 1 test files in heatmap-data...

Test: heatmap verification (passed)
  Coverage Score: 1.6%
  Hotspots (Clusters): 1
  Cold Spots (Untested): 63
  Top 3 Cold Spots:
    - a "Shop" ([data-testid="shop-link"])
    - a "Cart" ([data-testid="cart-link"])
    - a "🔮 Login" ([data-testid="login-link"])
```
This output uses K-Means clustering to identify where you test and highlights what you missed.
