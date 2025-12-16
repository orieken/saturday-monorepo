import { Pickle, GherkinDocument, PickleStep, Step } from '@cucumber/messages';

export type Attributes = Record<string, string | number | boolean>;

export interface OtelCucumberConfig {
  resourceAttributes?: Attributes;
  scenarioAttributes?: Attributes | ((pickle: Pickle, gherkinDocument: GherkinDocument) => Attributes);
  stepAttributes?: Attributes | ((pickleStep: PickleStep, gherkinStep: Step) => Attributes);
}

export const defaultConfig: OtelCucumberConfig = {
  resourceAttributes: {},
  scenarioAttributes: (pickle: Pickle, gherkinDocument: GherkinDocument) => {
    return {
       'custom.feature.name': gherkinDocument.feature?.name || 'unknown',
       'custom.scenario.tags': pickle.tags.map(t => t.name).join(','),
    };
  },
  stepAttributes: {}
};
