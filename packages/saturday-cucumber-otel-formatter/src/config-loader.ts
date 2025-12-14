import * as path from 'path';
import * as fs from 'fs';
import { pathToFileURL } from 'url';
import { OtelCucumberConfig, defaultConfig } from './config';

export async function loadConfig(configPath?: string): Promise<OtelCucumberConfig> {
  const searchPath = configPath || path.join(process.cwd(), 'cucumber-otel.config.js');
  
  if (fs.existsSync(searchPath)) {
    try {
      // Use dynamic import for ESM support
      const fileUrl = pathToFileURL(searchPath).href;
      const module = await import(fileUrl);
      return { ...defaultConfig, ...(module.default || module) };
    } catch (error) {
      console.error(`[OtelFormatter] Failed to load config from ${searchPath}:`, error);
      return defaultConfig;
    }
  }

  return defaultConfig;
}
