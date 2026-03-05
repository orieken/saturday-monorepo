import { LaunchOptions, Browser } from 'playwright';

export interface SaturdayConfig {
  browser: {
    headless?: boolean;
    slowMo?: number;
    args?: string[];
    channel?: string;
    [key: string]: any;
  };
  context: {
    recordVideo?: {
      dir: string;
      size: { width: number; height: number };
    };
    viewport?: { width: number; height: number };
    [key: string]: any;
  };
  customBrowsers?: {
    [key: string]: (options: LaunchOptions) => Promise<Browser>;
  };
}

export const defaultConfig: SaturdayConfig = {
  browser: {
    headless: process.env.HEADLESS === 'true',
    slowMo: 0,
    args: ['--use-fake-ui-for-media-stream', '--use-fake-device-for-media-stream'],
  },
  context: {
    recordVideo: {
      dir: 'reports/videos/',
      size: { width: 1280, height: 720 },
    },
    viewport: { width: 1280, height: 720 },
  },
};

let currentConfig: SaturdayConfig = { ...defaultConfig };

export function configure(config: Partial<SaturdayConfig>): void {
  currentConfig = {
    ...currentConfig,
    ...config,
    browser: { ...currentConfig.browser, ...config.browser },
    context: { ...currentConfig.context, ...config.context },
  };
}

export function getConfig(): SaturdayConfig {
  return currentConfig;
}
