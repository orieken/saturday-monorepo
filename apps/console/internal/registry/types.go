package registry

// Types representing a minimal Cucumber index the backend serves and uses to resolve runs.

type CucumberStepRef struct {
    Text string `json:"text"`
    Line int    `json:"line"`
}

type CucumberScenarioRef struct {
    Id    string            `json:"id"`
    Line  int               `json:"line"`
    Name  string            `json:"name"`
    File  string            `json:"file"`
    Tags  []string          `json:"tags"`
    Steps []CucumberStepRef `json:"steps"`
}

type CucumberFeatureRef struct {
    Id          string                `json:"id"`
    Name        string                `json:"name"`
    File        string                `json:"file"`
    Description string                `json:"description"`
    Scenarios   []CucumberScenarioRef `json:"scenarios"`
}

type CucumberIndex struct {
    Framework string                `json:"framework"`
    SuiteId   string                `json:"suiteId"`
    Features  []CucumberFeatureRef  `json:"features"`
}
