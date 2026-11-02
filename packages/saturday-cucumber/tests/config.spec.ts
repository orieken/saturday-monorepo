
import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';

describe('defaultConfig.browser.headless', () => {
  const originalHeadless = process.env.HEADLESS;

  beforeEach(() => {
    vi.resetModules();
  });

  afterEach(() => {
    if (originalHeadless === undefined) {
      delete process.env.HEADLESS;
    } else {
      process.env.HEADLESS = originalHeadless;
    }
  });

  it('is true when HEADLESS is the literal string "true" at import time', async () => {
    process.env.HEADLESS = 'true';
    const { defaultConfig } = await import('../src/config/index');
    expect(defaultConfig.browser.headless).toBe(true);
  });

  it('is false when HEADLESS is unset', async () => {
    delete process.env.HEADLESS;
    const { defaultConfig } = await import('../src/config/index');
    expect(defaultConfig.browser.headless).toBe(false);
  });

  it('is false when HEADLESS is set to any value other than "true"', async () => {
    process.env.HEADLESS = 'false';
    const { defaultConfig } = await import('../src/config/index');
    expect(defaultConfig.browser.headless).toBe(false);
  });
});

describe('configure() and getConfig()', () => {
  beforeEach(() => {
    vi.resetModules();
  });

  it('merges browser and context overrides without discarding untouched fields', async () => {
    const { configure, getConfig, defaultConfig } = await import('../src/config/index');

    configure({ browser: { slowMo: 50 }, context: { viewport: { width: 800, height: 600 } } });
    const config = getConfig();

    expect(config.browser.slowMo).toBe(50);
    expect(config.browser.headless).toBe(defaultConfig.browser.headless);
    expect(config.context.viewport).toEqual({ width: 800, height: 600 });
    expect(config.context.recordVideo).toEqual(defaultConfig.context.recordVideo);
  });
});
