import { test as base, request } from '@playwright/test';
import { createK6Recorder, K6Recorder } from './index';
import { createDefaultRedactionPolicy } from '@orieken/saturday-k6-redaction-basic';

type K6Fixtures = {
  k6Request: import('@playwright/test').APIRequestContext;
  k6Recorder: K6Recorder | null;
};

export const test = base.extend<K6Fixtures>({
  k6Recorder: [async ({}, use, testInfo) => {
    if (!process.env.K6_EXPORT || !testInfo.title.includes('@k6')) {
      await use(null);
      return;
    }
    const setup = await createK6Recorder(testInfo.title, undefined /* outDir */, { policy: createDefaultRedactionPolicy({ envPrefix: 'K6_' }) });
    if (!setup) { await use(null); return; }
    await use(setup.recorder);
    if (setup.recorder.hasCalls()) {
      await setup.recorder.flushToK6();
    }
    await setup.ctx.dispose();
  }, { scope: 'test' }],

  k6Request: async ({ k6Recorder }, use) => {
    const ctx = await request.newContext();
    await use(ctx);
    await ctx.dispose();
  }
});

export const expect = test.expect;