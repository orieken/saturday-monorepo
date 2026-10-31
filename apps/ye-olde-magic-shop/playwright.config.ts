import { defineConfig, devices } from '@playwright/test';
require('dotenv').config();

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
      // Runs in the Saturday runner image (see docs/runner-image.md) in
      // real Chrome under Xvfb, instead of Playwright's bundled Chromium —
      // Cesium's upcoming WebGL work needs a real browser's GL stack.
      name: 'chrome',
      use: {
        ...devices['Desktop Chrome'],
        channel: 'chrome',
        headless: false,
        launchOptions: {
          slowMo: Number(process.env.RUNNER_SLOW_MO ?? 0),
          args: ['--use-gl=angle', '--use-angle=swiftshader'],
        },
        video: process.env.RUNNER_VIDEO === '1' ? 'on' : 'off',
      },
    },
  ],
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:8000',
    reuseExistingServer: !process.env.CI,
  },
});
