import { defineConfig, devices } from '@playwright/test';

// Same env-driven resolution pattern as apps/ye-olde-magic-shop/playwright.config.ts;
// Leaflet needs no WebGL context, so this config omits that file's
// RUNNER_IMAGE_SOFTWARE_GL_ARGS launch args.
export function resolveRunnerSlowMo(): number {
  const parsed = Number(process.env.RUNNER_SLOW_MO);
  return Number.isFinite(parsed) ? parsed : 0;
}

const FIXTURE_SERVER_PORT = 4173;

export default defineConfig({
  testDir: './tests',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: [['list']],
  use: {
    baseURL: `http://localhost:${FIXTURE_SERVER_PORT}`,
    trace: 'on-first-retry',
  },
  projects: [
    {
      // Runs in the Saturday runner image (see docs/runner-image.md) in real
      // Chrome, per AC4 — no bundled-Chromium project exists for this
      // package since the fixture is only ever exercised this way.
      name: 'chrome',
      use: {
        ...devices['Desktop Chrome'],
        channel: 'chrome',
        launchOptions: {
          slowMo: resolveRunnerSlowMo(),
        },
      },
    },
  ],
  webServer: {
    command: `node fixtures/serve.mjs`,
    url: `http://localhost:${FIXTURE_SERVER_PORT}`,
    reuseExistingServer: !process.env.CI,
    env: {
      FIXTURE_SERVER_PORT: String(FIXTURE_SERVER_PORT),
    },
  },
});
