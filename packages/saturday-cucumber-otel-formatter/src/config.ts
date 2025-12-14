import { Pickle, GherkinDocument, PickleStep, Step } from '@cucumber/messages';

export type Attributes = Record<string, string | number | boolean>;

export interface OtelCucumberConfig {
  resourceAttributes?: Attributes;
  scenarioAttributes?: Attributes | ((pickle: Pickle, gherkinDocument: GherkinDocument) => Attributes);
  stepAttributes?: Attributes | ((pickleStep: PickleStep, gherkinStep: Step) => Attributes);
}

export const defaultConfig: OtelCucumberConfig = {
  resourceAttributes: {},
  scenarioAttributes: {},
  stepAttributes: {}
};
