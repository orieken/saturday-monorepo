package models

// AnalysisResult represents the result of a framework analysis
type AnalysisResult struct {
	Stats     ProjectStats     `json:"stats"`
	Structure ProjectStructure `json:"structure"`
	Inventory ProjectInventory `json:"inventory"`
	Issues    []ValidationIssue `json:"issues,omitempty"`
}

// ProjectStats contains high-level metrics about the project
type ProjectStats struct {
	TotalFiles       int `json:"totalFiles"`
	PageCount        int `json:"pageCount"`
	FlowCount        int `json:"flowCount"`
	StepDefCount     int `json:"stepDefCount"`
	TestCount        int `json:"testCount"`
	HelperCount      int `json:"helperCount"`
}

// ProjectStructure represents the directory structure locations
type ProjectStructure struct {
	RootPath      string `json:"rootPath"`
	PagesPath     string `json:"pagesPath,omitempty"`
	FlowsPath     string `json:"flowsPath,omitempty"`
	StepsPath     string `json:"stepsPath,omitempty"`
	GenericUtils  string `json:"genericUtils,omitempty"`
}

// ProjectInventory contains lists of identified components
type ProjectInventory struct {
	Pages     []PageInfo     `json:"pages,omitempty"`
	Flows     []FlowInfo     `json:"flows,omitempty"`
	StepDefs  []StepDefInfo  `json:"stepDefs,omitempty"`
}

// PageInfo represents a discovered Page Object
type PageInfo struct {
	Name        string   `json:"name"`
	FilePath    string   `json:"filePath"`
	Methods     []string `json:"methods,omitempty"`
	Selectors   []string `json:"selectors,omitempty"`
}

// FlowInfo represents a discovered Flow
type FlowInfo struct {
	Name     string   `json:"name"`
	FilePath string   `json:"filePath"`
	Steps    []string `json:"steps,omitempty"`
}

// StepDefInfo represents a discovered Step Definition file
type StepDefInfo struct {
	FilePath string `json:"filePath"`
	Steps    []Step `json:"steps,omitempty"`
}

// Step represents a single cucumber step definition
type Step struct {
	Type    string `json:"type"` // Given, When, Then
	Pattern string `json:"pattern"`
	Line    int    `json:"line"`
}

