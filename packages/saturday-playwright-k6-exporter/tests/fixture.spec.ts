
import { test as baseTest } from '@playwright/test';
import { createK6Recorder } from '../src/index';

// Verify that fixture.ts executes the extend
// We need to require it to trigger side effect
import '../src/fixture';

jest.mock('../src/index', () => ({
  createK6Recorder: jest.fn(),
  K6Recorder: class {
      hasCalls() { return false; }
      flushToK6() {}
  }
}));

describe('fixture', () => {
    it('should extend base test with k6Recorder', () => {
        expect(baseTest.extend).toHaveBeenCalled();
        const fixtures = (baseTest as any).fixtures;
        expect(fixtures).toBeDefined();
        expect(fixtures.k6Recorder).toBeDefined();
        expect(fixtures.k6Request).toBeDefined();
    });

    it('should initialize k6Recorder fixture', async () => {
         const fixtures = (baseTest as any).fixtures;
         const use = jest.fn();
         const testInfo = { title: 'my test @k6' };

         // Mock createK6Recorder
         const mockRecorder = { hasCalls: jest.fn().mockReturnValue(true), flushToK6: jest.fn(), ctx: { dispose: jest.fn() } };
         (createK6Recorder as jest.Mock).mockResolvedValue({ ctx: mockRecorder.ctx, recorder: mockRecorder });

         process.env.K6_EXPORT = '1';
         
         // Execute the fixture function (it's [fn, { scope }] or just fn)
         // In fixture.ts: k6Recorder: [async ({}, use, testInfo) => ..., { scope: 'test' }]
         const fixtureFn = fixtures.k6Recorder[0];
         
         await fixtureFn({}, use, testInfo);

         expect(use).toHaveBeenCalledWith(mockRecorder);
         expect(mockRecorder.flushToK6).toHaveBeenCalled();
         expect(mockRecorder.ctx.dispose).toHaveBeenCalled();
    });
    
    it('should skip recorder if K6_EXPORT not set', async () => {
        delete process.env.K6_EXPORT;
        const fixtures = (baseTest as any).fixtures;
        const fixtureFn = fixtures.k6Recorder[0];
        const use = jest.fn();
        
        await fixtureFn({}, use, { title: 'test @k6'});
        expect(use).toHaveBeenCalledWith(null);
    });
});
