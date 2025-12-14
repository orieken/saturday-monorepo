import { Context } from '@opentelemetry/api';
import { TestStepResultStatus, Duration } from '@cucumber/messages';

export interface CurrentStep {
  name: string;
  keyword: string;
  text: string;
  line: number;
  file: string;
  status?: TestStepResultStatus;
  duration?: Duration;
  error?: string;
  context?: Context;
}

export interface CurrentScenario {
  name: string;
  line: number;
  file: string;
  steps: CurrentStep[];
  status?: TestStepResultStatus;
  duration?: Duration;
  context?: Context;
  hasFailedSteps: boolean;
  error?: string;
}

export interface CurrentFeature {
  name: string;
  line: number;
  file: string;
  scenarios: CurrentScenario[];
  context?: Context;
}
