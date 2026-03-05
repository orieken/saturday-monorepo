/**
 * Cucumber JSON Report Types
 * 
 * These types match the Cucumber JSON format output from cucumber-js
 * See: https://github.com/cucumber/cucumber-js/blob/main/docs/formatters.md
 */

export type StepStatus = 'passed' | 'failed' | 'skipped' | 'pending' | 'undefined' | 'ambiguous';

export interface CucumberStepMatch {
  location: string;
}

export interface CucumberStepResult {
  status: StepStatus;
  duration?: number; // nanoseconds
  error_message?: string;
}

export interface CucumberStep {
  keyword: string;
  name: string;
  line: number;
  match?: CucumberStepMatch;
  result: CucumberStepResult;
}

export interface CucumberTag {
  name: string;
  line: number;
}

export interface CucumberScenario {
  id: string;
  keyword: string;
  name: string;
  description?: string;
  line: number;
  type: string;
  tags?: CucumberTag[];
  steps: CucumberStep[];
}

export interface CucumberFeature {
  uri: string;
  id: string;
  keyword: string;
  name: string;
  description?: string;
  line: number;
  tags?: CucumberTag[];
  elements: CucumberScenario[];
}

export type CucumberReport = CucumberFeature[];

/**
 * Helper type for scenario status (derived from steps)
 */
export type ScenarioStatus = 'passed' | 'failed' | 'skipped' | 'pending' | 'unknown';

/**
 * Report statistics
 */
export interface ReportStats {
  totalScenarios: number;
  passedScenarios: number;
  failedScenarios: number;
  skippedScenarios: number;
  pendingScenarios: number;
  totalDuration: number; // milliseconds
  successRate: number; // percentage
}
