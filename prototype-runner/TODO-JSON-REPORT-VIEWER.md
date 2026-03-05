# TODO: Implement Cartridge JSON Report Viewer

**Goal**: Replace the raw HTML iframe/link with a native, integrated JSON report viewer in Cartridge.
**Status**: In Progress

## 1. Backend: Serve JSON Reports (Console) ✅
- [x] Ensure `test-runner-service` generates `cucumber.json`.
- [x] Create API endpoint `GET /api/runs/{runId}/report`.
  - Should read `cucumber.json` from the reports directory.
  - specific path: `./reports/<suiteId>/<runId>/cucumber.json`.
- [x] Ensure CORS allows Cartridge to fetch this JSON.

## 2. Frontend: API Client (Cartridge) ✅
- [x] Add types for Cucumber JSON (`CucumberFeature`, `CucumberScenario`, `CucumberStep`, etc.) in `src/types/cucumber.ts`.
- [x] Add `fetchRunReport(runId)` to `src/api/runs.ts` to call the Console API.

## 3. Frontend: Report Viewer Component ✅
- [x] Create `src/components/CucumberReportViewer.vue`.
- [x] **Design Implementation**:
  - [x] **Header**: Test Summary Stats (Total, Passed, Failed, Skipped, Success Rate).
  - [x] **Feature Section**: Title, Description, Tags.
  - [x] **Scenario Cards**: 
    - Collapsible/Expandable.
    - Status borders (Green/Red).
    - Status badges.
  - [x] **Steps**: 
    - Icons (✓, ✗, -).
    - Durations in ms/s.
    - Error messages for failures.
- [x] Styling: Use Tailwind CSS dark mode (slate-800/900).

## 4. Frontend: Integration (RunModal) ✅
- [x] Import `CucumberReportViewer` in `RunModal.vue`.
- [x] Add "Report" tab to the modal interface.
- [x] **Logic**:
  - Fetch report when tab is active.
  - Show loading spinner.
  - Handle error states.
- [x] **Replace HTML View**: Ensure this is the primary way to view results.

## 5. Cleanup & Polish (Next Steps)
- [ ] **Verify styling**: Ensure it matches the requested "Dark Dashboard" mockup.
- [ ] **Data Validation**: Ensure large reports render performantly.
- [ ] **HTML Report Link**: Decide whether to keep "Open HTML report" as a fallback or remove it.
