import { defineStore } from 'pinia';
import { demoEcommerceCucumberIndex, type CucumberFeatureRef, type CucumberScenarioRef } from '../canned/demoEcommerceCucumberIndex';

export type RunStatus = 'idle' | 'running' | 'passed' | 'failed';

export interface ScenarioRunInfo {
  scenarioId: string;
  line: number;
  status: RunStatus;
  lastRunAt?: string;
  fakeReportUrl?: string;
}

export interface ActiveRun {
  feature: CucumberFeatureRef;
  scenario: CucumberScenarioRef;
  status: RunStatus;
  logs: string[];
}

interface CucumberDemoState {
  suiteId: string;
  features: CucumberFeatureRef[];
  scenarioRuns: Record<string, ScenarioRunInfo>;
  activeRun: ActiveRun | null;
}

function scenarioKey(s: CucumberScenarioRef): string {
  return `${s.file}:${s.line}`;
}

export const useCucumberDemoStore = defineStore('cucumberDemo', {
  state: (): CucumberDemoState => ({
    suiteId: demoEcommerceCucumberIndex.suiteId,
    features: demoEcommerceCucumberIndex.features,
    scenarioRuns: {},
    activeRun: null
  }),

  getters: {
    getScenarioRun:
      (state) =>
      (scenario: CucumberScenarioRef): ScenarioRunInfo => {
        const key = scenarioKey(scenario);
        return (
          state.scenarioRuns[key] ?? {
            scenarioId: scenario.id,
            line: scenario.line,
            status: 'idle'
          }
        );
      }
  },

  actions: {
    runScenario(feature: CucumberFeatureRef, scenario: CucumberScenarioRef) {
      const key = scenarioKey(scenario);

      const now = new Date().toISOString();

      this.scenarioRuns[key] = {
        scenarioId: scenario.id,
        line: scenario.line,
        status: 'running',
        lastRunAt: now,
        fakeReportUrl: '#'
      };

      const logs: string[] = [
        `$ test-runner --framework cucumber --suite ${this.suiteId} --scenario "${scenario.name}"`,
        '',
        `[${now}] INFO  Starting scenario "${scenario.name}"`,
        `[${now}] INFO  Using feature file: ${scenario.file}:${scenario.line}`,
        `[${now}] INFO  Step 1/3: Given I am on the relevant page`,
        `[${now}] INFO  Step 2/3: When I perform the primary user action`,
        `[${now}] INFO  Step 3/3: Then I should see the expected outcome`,
        '',
        `[${now}] PASS  Scenario "${scenario.name}" finished successfully`,
        '',
        `Generated HTML report: /reports/cucumber/${this.suiteId}/${scenario.id}.html`
      ];

      this.activeRun = {
        feature,
        scenario,
        status: 'running',
        logs
      };

      setTimeout(() => {
        this.scenarioRuns[key] = {
          ...this.scenarioRuns[key],
          status: 'passed',
          lastRunAt: new Date().toISOString(),
          fakeReportUrl: `/reports/cucumber/${this.suiteId}/${scenario.id}.html`
        };

        if (this.activeRun && this.activeRun.scenario.id === scenario.id) {
          this.activeRun.status = 'passed';
        }
      }, 2500);
    },

    closeActiveRun() {
      this.activeRun = null;
    }
  }
});
