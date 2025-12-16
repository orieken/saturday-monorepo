import { TestCase, TestStep } from '@playwright/test/reporter';

export type Attributes = Record<string, string | number | boolean>;

export interface OtelPlaywrightConfig {
  resourceAttributes?: Attributes;
  testAttributes?: Attributes | ((test: TestCase) => Attributes);
  stepAttributes?: Attributes | ((step: TestStep) => Attributes);
}

export const defaultConfig: OtelPlaywrightConfig = {
  resourceAttributes: {},
  testAttributes: {},
  stepAttributes: {}
};
