import { defineConfig, devices } from '@playwright/test';

// Same env-driven resolution pattern as apps/ye-olde-magic-shop/playwright.config.ts.
// Unlike saturday-playwright-leaflet, Cesium needs a real webgl2 context, so this
// config carries the same RUNNER_IMAGE_SOFTWARE_GL_ARGS launch args that app uses
// for its own 'chrome' project (see docs/runner-image.md).
const RUNNER_IMAGE_SOFTWARE_GL_ARGS = ['--use-gl=angle', '--use-angle=swiftshader'];

export function resolveRunnerSlowMo(): number {
  const parsed = Number(process.env.RUNNER_SLOW_MO);
  return Number.isFinite(parsed) ? parsed : 0;
}

const FIXTURE_SERVER_PORT = 4174;

export default defineConfig({
  testDir: './tests',
  testMatch: /fixture\.spec\.ts/,
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: [['list']],
  use: {
    baseURL: `http://127.0.0.1:${FIXTURE_SERVER_PORT}`,
    trace: 'on-first-retry',
  },
  projects: [
    {
      // Runs in the Saturday runner image (see docs/runner-image.md) in real,
      // headed Chrome with software WebGL — Cesium needs a webgl2 context
      // that Playwright's bundled Chromium project does not provide here.
      name: 'chrome',
      use: {
        ...devices['Desktop Chrome'],
        channel: 'chrome',
        headless: false,
        launchOptions: {
          slowMo: resolveRunnerSlowMo(),
          args: RUNNER_IMAGE_SOFTWARE_GL_ARGS,
        },
      },
    },
  ],
  webServer: {
    command: `node fixtures/serve.mjs`,
    url: `http://127.0.0.1:${FIXTURE_SERVER_PORT}`,
    reuseExistingServer: !process.env.CI,
    env: {
      FIXTURE_SERVER_PORT: String(FIXTURE_SERVER_PORT),
    },
  },
});
