# Tool Flow Architecture

## generate_site
Sequence:
- Client → Server `generate_site(params)`
- Validate input
- Generate site structure via generator and template processor
- Validate generated code
- Write files and return manifest

Returns: `{ files: [], structure: {}, imports: [] }`

## generate_step_definitions
Sequence:
- Server reads feature file, parses Gherkin
- Analyze existing steps, generate missing steps
- Write step definitions

Returns: `{ newSteps: [], existingSteps: [], suggestions: [] }`

## analyze_framework_usage
Sequence:
- Server triggers framework analyzer which parses TypeScript, extracts sites/pages/flows, validates patterns and calculates metrics

Returns: `{ sites: [], pages: [], flows: [], metrics: {}, violations: [], suggestions: [] }`

## Resource Access Flow
- `saturday://` URIs handled by Resource Provider with caching, content processing and discovery.

