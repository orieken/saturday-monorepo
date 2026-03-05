import { Page } from '@playwright/test';
import * as fs from 'fs';
import * as path from 'path';
import { trackerScript } from './tracker';
import { scannerScript } from './scanner';

export async function injectHeatmap(page: Page) {
    await page.addInitScript(trackerScript);
}

export async function saveHeatmapData(page: Page, testTitle: string, status: string = 'unknown') {
    // 1. Collect Events
    const events: any[] = await page.evaluate(() => {
        return (window as any).__SATURDAY_HEATMAP_EVENTS__ || [];
    });

    // 2. Scan Interactables
    const interactables: any[] = await page.evaluate(scannerScript);

    const testId = testTitle.replace(/[^a-zA-Z0-9_]/g, '-');
    const outputDir = path.join(process.cwd(), 'heatmap-data');
    if (!fs.existsSync(outputDir)) {
        fs.mkdirSync(outputDir, { recursive: true });
    }

    const data = {
        testId,
        testTitle,
        status,
        events,
        interactables,
        snapshot: await page.screenshot({ fullPage: true }).then(b => b.toString('base64')).catch(() => null)
    };

    fs.writeFileSync(path.join(outputDir, `${testId}.json`), JSON.stringify(data, null, 2));
}
