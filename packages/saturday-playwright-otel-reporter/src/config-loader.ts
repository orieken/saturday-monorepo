import * as path from 'path';
import * as fs from 'fs';
import { pathToFileURL } from 'url';
import { OtelPlaywrightConfig, defaultConfig } from './config';

export async function loadConfig(configPath?: string): Promise<OtelPlaywrightConfig> {
  let searchPath = configPath;

  if (!searchPath) {
      const defaultMjs = path.join(process.cwd(), 'playwright-otel.config.mjs');
      const defaultJs = path.join(process.cwd(), 'playwright-otel.config.js');
      
      if (fs.existsSync(defaultMjs)) {
          searchPath = defaultMjs;
      } else {
          searchPath = defaultJs;
      }
  }
  
  if (fs.existsSync(searchPath)) {
    try {
      // Use dynamic import for ESM support
      const fileUrl = pathToFileURL(searchPath).href;
      const module = await import(fileUrl);
      return { ...defaultConfig, ...(module.default || module) };
    } catch (error) {
      console.error(`[OtelReporter] Failed to load config from ${searchPath}:`, error);
      return defaultConfig;
    }
  }

  return defaultConfig;
}
