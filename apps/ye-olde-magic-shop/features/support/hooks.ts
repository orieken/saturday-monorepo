import { Before, After, BeforeAll, AfterAll, setDefaultTimeout, Status } from '@cucumber/cucumber';
import { TestWorld, WorldImpl } from './world';
import * as fs from 'fs';
import { installSaturdayHooks } from '@orieken/saturday-cucumber';

installSaturdayHooks();

// Increase default step timeout to accommodate browser startup on CI/containers
setDefaultTimeout(60_000);

BeforeAll(function() {
    console.log('[Cucumber Hooks] Test Suite Starting...');
});

AfterAll(function() {
    console.log('[Cucumber Hooks] Test Suite Finished.');
});

Before(async function (this: WorldImpl) {  
  if (this.page) {
      // @ts-ignore
      const { injectHeatmap } = await import('@orieken/saturday-playwright-heatmap');
      await injectHeatmap(this.page);
  }
});

After(async function (this: WorldImpl, scenario) {
  if (this.page) {
       // @ts-ignore
      const { saveHeatmapData } = await import('@orieken/saturday-playwright-heatmap');
      await saveHeatmapData(this.page, scenario.pickle.name, scenario.result?.status?.toString());
  }
});
