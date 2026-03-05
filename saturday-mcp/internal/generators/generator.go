package generators

// Generator is the main facade for all code generators
type Generator struct {
	SiteGenerator    *SiteGenerator
	PageGenerator    *PageGenerator
	FlowGenerator    *FlowGenerator
	StepGenerator    *StepGenerator
	ElementGenerator   *ElementGenerator
	ServiceGenerator       *ServiceGenerator
	MigrationGenerator     *MigrationGenerator
	DocumentationGenerator *DocumentationGenerator
}

// NewGenerator creates a new generator facade
func NewGenerator(siteGen *SiteGenerator, pageGen *PageGenerator, flowGen *FlowGenerator, stepGen *StepGenerator, elementGen *ElementGenerator, serviceGen *ServiceGenerator, migrationGen *MigrationGenerator, docGen *DocumentationGenerator) *Generator {
	return &Generator{
		SiteGenerator:          siteGen,
		PageGenerator:          pageGen,
		FlowGenerator:          flowGen,
		StepGenerator:          stepGen,
		ElementGenerator:       elementGen,
		ServiceGenerator:       serviceGen,
		MigrationGenerator:     migrationGen,
		DocumentationGenerator: docGen,
	}
}
