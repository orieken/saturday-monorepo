/**
 * @issue ui-suite-real-chrome
 * Verifies the 'chrome' project's resolved Playwright configuration and the
 * env-var resolution it depends on, without needing a real browser or the
 * runner container. The behavioral proof that the chrome project actually
 * launches real Chrome with a working GL stack is tests/webgl.spec.ts,
 * exercised inside scripts/ui-suite-in-chrome.sh.
 */
import { test, expect } from '@playwright/test';
import playwrightConfig, { resolveRunnerSlowMo, resolveRunnerVideo } from '../playwright.config';

type ProjectConfig = {
  name?: string;
  use?: {
    channel?: string;
    headless?: boolean;
    launchOptions?: { args?: string[]; slowMo?: number };
    video?: string;
    defaultBrowserType?: string;
  };
};

function findProject(name: string): ProjectConfig {
  const projects = (playwrightConfig.projects ?? []) as ProjectConfig[];
  const found = projects.find((project) => project.name === name);
  if (!found) {
    throw new Error(`no '${name}' project defined in playwright.config.ts`);
  }
  return found;
}

test.describe('chrome project configuration (AC1)', () => {
  test('resolves channel "chrome" and headed mode', () => {
    const chrome = findProject('chrome');

    expect(chrome.use?.channel).toBe('chrome');
    expect(chrome.use?.headless).toBe(false);
  });

  test('resolves the ANGLE SwiftShader software-GL launch args', () => {
    const chrome = findProject('chrome');

    expect(chrome.use?.launchOptions?.args).toEqual(['--use-gl=angle', '--use-angle=swiftshader']);
  });

  test('the bundled-Chromium project remains defined for local runs', () => {
    const chromium = findProject('chromium');

    expect(chromium.use?.channel).toBeUndefined();
    expect(chromium.use?.defaultBrowserType).toBe('chromium');
  });
});

test.describe.serial('RUNNER_SLOW_MO / RUNNER_VIDEO resolution (AC1 edge cases)', () => {
  const ORIGINAL_SLOW_MO = process.env.RUNNER_SLOW_MO;
  const ORIGINAL_VIDEO = process.env.RUNNER_VIDEO;

  test.afterEach(() => {
    if (ORIGINAL_SLOW_MO === undefined) delete process.env.RUNNER_SLOW_MO;
    else process.env.RUNNER_SLOW_MO = ORIGINAL_SLOW_MO;
    if (ORIGINAL_VIDEO === undefined) delete process.env.RUNNER_VIDEO;
    else process.env.RUNNER_VIDEO = ORIGINAL_VIDEO;
  });

  test('slowMo defaults to 0 when RUNNER_SLOW_MO is unset', () => {
    delete process.env.RUNNER_SLOW_MO;

    expect(resolveRunnerSlowMo()).toBe(0);
  });

  test('slowMo resolves to the parsed number when RUNNER_SLOW_MO is numeric', () => {
    process.env.RUNNER_SLOW_MO = '500';

    expect(resolveRunnerSlowMo()).toBe(500);
  });

  test('slowMo falls back to 0 when RUNNER_SLOW_MO is not a number', () => {
    process.env.RUNNER_SLOW_MO = 'not-a-number';

    expect(resolveRunnerSlowMo()).toBe(0);
  });

  test('video defaults to "off" when RUNNER_VIDEO is unset', () => {
    delete process.env.RUNNER_VIDEO;

    expect(resolveRunnerVideo()).toBe('off');
  });

  test('video resolves to "on" when RUNNER_VIDEO is "1"', () => {
    process.env.RUNNER_VIDEO = '1';

    expect(resolveRunnerVideo()).toBe('on');
  });

  test('video stays "off" for any other RUNNER_VIDEO value', () => {
    process.env.RUNNER_VIDEO = 'true';

    expect(resolveRunnerVideo()).toBe('off');
  });
});
