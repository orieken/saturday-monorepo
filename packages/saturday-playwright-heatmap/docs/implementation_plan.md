# ML Analysis Implementation Plan

## Goal Description
Integrate Machine Learning (specifically Clustering algorithms) to analyze the test execution heatmap data. The goal is to mathematically identify "Hotspots" (areas of heavy testing) and "Coldspots" (interactable areas with zero or low testing coverage) to provide actionable insights.

## User Review Required
> [!NOTE]
> We will use `ml-kmeans` (Node.js) to keep the stack consistent (TypeScript) rather than introducing Python.

## Proposed Changes

### 1. Update `saturday-playwright-heatmap`
*   **Modify `src/fixture.ts`**: Capture the test result (pass/fail) in the JSON output. This allows future analysis to correlate interaction patterns with failures.

### 2. New Package: `packages/saturday-ml-analyzer`
A new TypeScript package dedicated to data analysis.

#### Components:
*   **Data Loader**: Reads `heatmap-data` JSON files.
*   **Cluster Engine**:
    *   Uses **K-Means Clustering** to group interaction points `(x, y)` into clusters.
    *   Calculates centroids of these clusters.
*   **Coverage Logic**:
    *   Compares **Interactable Elements** against **Interaction Clusters**.
    *   If an interactable element's position is > Threshold distance from nearest cluster centroid, mark as **"Cold Spot" (Risk Area)**.
*   **Reporter**:
    *   Generates a text summary or updates the HTML report with "Risk Scores".

### 3. Dependencies
*   `ml-kmeans`: For clustering algorithms.
*   `simple-statistics`: For basic statistical functions.

## Verification Plan

### Automated
*   **Unit Tests**: Verify K-Means logic works on a simple set of 2D points.
*   **Integration**: Run `analyzer` on the existing `ye-olde-magic-shop` heatmap data and verify it identifies non-interacted buttons as "Cold Spots".

### Manual
*   Run the analyzer command.
*   Check the console output for a list of "High Risk / Untested Elements".
