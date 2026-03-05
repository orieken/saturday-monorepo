import { test as base } from '@playwright/test';
import * as fs from 'fs';
import * as path from 'path';
import { trackerScript } from './tracker';
import { scannerScript } from './scanner';

export const test = base.extend<{ heatmap: void }>({
  heatmap: [async ({ page }, use, testInfo) => {
    // 1. Inject Tracker Script
    await page.addInitScript(trackerScript);

    await use();

    // 2. Scan Interactables at the end of the test
    // Note: We might want to do this at multiple points, but for now, let's do it at the end
    // or maybe we don't scan at all and just rely on the reporter?
    // Actually, the plan said "Scanner Script ... at the end of each test".
    
    // Check if we care about this test (maybe skip if failed?)
    
    // 3. Collect Data
    const events: any[] = await page.evaluate(() => {
        return (window as any).__SATURDAY_HEATMAP_EVENTS__ || [];
    });

    const interactables: any[] = await page.evaluate(scannerScript);

    const testId = testInfo.titlePath.join('__').replace(/[^a-zA-Z0-9_]/g, '-');
    const outputDir = path.join(process.cwd(), 'heatmap-data');
    if (!fs.existsSync(outputDir)) {
        fs.mkdirSync(outputDir, { recursive: true });
    }

    const data = {
        testId,
        testTitle: testInfo.title,
        status: testInfo.status,
        events,
        // For simplicity, we capture the final state of interactables. 
        // Ideally we'd capture them throughout, but that's expensive.
        interactables,
        snapshot: await page.screenshot({ fullPage: true }).then(b => b.toString('base64')).catch(() => null)
    };

    fs.writeFileSync(path.join(outputDir, `${testId}.json`), JSON.stringify(data, null, 2));

  }, { auto: true }]
});
