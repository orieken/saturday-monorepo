package runs

// Run represents a single test execution.
type Run struct {
    ID         string `json:"id"`
    Status     string `json:"status"`
    StartedAt  string `json:"startedAt"`
    FinishedAt string `json:"finishedAt,omitempty"`
    ReportURL  string `json:"reportUrl,omitempty"`

    Framework  string `json:"framework"`
    SuiteId    string `json:"suiteId"`
    ScenarioId string `json:"scenarioId"`
    // Executor used to run this run (e.g. "docker" or "k8s").
    Executor   string `json:"executor,omitempty"`
}
