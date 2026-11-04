/**
 * @issue reagent-map-globe
 * AC1: "The /reagent-map route and its nav link exist only when
 * VITE_FEATURE_REAGENT_MAP is exactly 'true'" — isReagentMapEnabled is the
 * small pure function that decision depends on, and this file proves it is
 * testable for 'true', 'false', and unset without a browser.
 *
 * AC6: these specs are named tests/reagent-map-*.spec.ts so they run only
 * under the 'chrome' Playwright project (see playwright.config.ts's
 * testIgnore on the chromium project) — including this file itself, which
 * needs no browser at all. That is intentional: the naming convention, not
 * the test's own requirements, decides which project runs it.
 */
import { test, expect } from '@playwright/test';
import { isReagentMapEnabled } from '../src/featureFlags';
import playwrightConfig from '../playwright.config';

type ProjectConfig = { name?: string; testIgnore?: string[] | RegExp[] };

test.describe('isReagentMapEnabled (AC1)', () => {
  test('returns true for the literal string "true"', () => {
    expect(isReagentMapEnabled('true')).toBe(true);
  });

  test('returns false for the literal string "false"', () => {
    expect(isReagentMapEnabled('false')).toBe(false);
  });

  test('returns false when unset', () => {
    expect(isReagentMapEnabled(undefined)).toBe(false);
  });

  test('returns false for the empty string', () => {
    expect(isReagentMapEnabled('')).toBe(false);
  });

  test('returns false for near-true case and whitespace variants (AC1 edge case)', () => {
    expect(isReagentMapEnabled('TRUE')).toBe(false);
    expect(isReagentMapEnabled('True')).toBe(false);
    expect(isReagentMapEnabled(' true')).toBe(false);
    expect(isReagentMapEnabled('true ')).toBe(false);
    expect(isReagentMapEnabled('1')).toBe(false);
  });
});

test.describe('reagent-map-*.spec.ts scoping (AC6)', () => {
  test('the chromium project ignores reagent-map-*.spec.ts', () => {
    const projects = (playwrightConfig.projects ?? []) as ProjectConfig[];
    const chromium = projects.find((project) => project.name === 'chromium');
    if (!chromium) throw new Error("no 'chromium' project defined in playwright.config.ts");

    const ignoresReagentMapSpecs = (chromium.testIgnore ?? []).some((pattern) =>
      String(pattern).includes('reagent-map-')
    );
    expect(ignoresReagentMapSpecs).toBe(true);
  });

  test('the chrome project defines no such exclusion', () => {
    const projects = (playwrightConfig.projects ?? []) as ProjectConfig[];
    const chrome = projects.find((project) => project.name === 'chrome');
    if (!chrome) throw new Error("no 'chrome' project defined in playwright.config.ts");

    expect(chrome.testIgnore).toBeUndefined();
  });
});
