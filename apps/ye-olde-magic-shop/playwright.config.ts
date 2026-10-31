import { defineConfig, devices } from '@playwright/test';
require('dotenv').config();

// Same literal args as SOFTWARE_GL_ARGS in
// containers/saturday-runner/runner-check/index.js — that script is the
// verified contract for webgl2 working in the runner image. Keep the two in
// sync by inspection; there's no shared module boundary between an app
// Playwright config and a container-internal Node script.
const RUNNER_IMAGE_SOFTWARE_GL_ARGS = ['--use-gl=angle', '--use-angle=swiftshader'];

// Exported (not just local to this config) so tests/playwright-config.spec.ts
// can verify the env-var resolution directly, per each call reading
// process.env fresh rather than once at config-load time.
export function resolveRunnerSlowMo(): number {
  const parsed = Number(process.env.RUNNER_SLOW_MO);
  return Number.isFinite(parsed) ? parsed : 0;
}

export function resolveRunnerVideo(): 'on' | 'off' {
  return process.env.RUNNER_VIDEO === '1' ? 'on' : 'off';
}

export default defineConfig({
  testDir: './tests',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: process.env.ENABLE_OTEL === 'true'
    ? [['list'], ['html'], ['@orieken/saturday-playwright-otel-reporter']]
    : [['list'], ['html']],
  use: {
    baseURL: 'http://localhost:8000',
    trace: 'on-first-retry',
    video: 'on',
    launchOptions: {
      slowMo: 1000,
    },
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      // Runs in the Saturday runner image (see docs/runner-image.md) in real
      // Chrome under Xvfb, instead of Playwright's bundled Chromium — Cesium
      // needs a real browser's GL stack for a webgl2 context. Only runs via
      // --project=chrome or scripts/ui-suite-in-chrome.sh; the default
      // (no --project) invocation used by test:e2e/test:pw stays chromium-only.
      name: 'chrome',
      use: {
        ...devices['Desktop Chrome'],
        channel: 'chrome',
        headless: false,
        launchOptions: {
          slowMo: resolveRunnerSlowMo(),
          args: RUNNER_IMAGE_SOFTWARE_GL_ARGS,
        },
        video: resolveRunnerVideo(),
      },
    },
  ],
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:8000',
    reuseExistingServer: !process.env.CI,
  },
});
