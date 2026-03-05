
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import * as Setup from '../src/hooks/setup';
import { browserOptions } from '../src/utils/browser-options';
import { SaturdayWorld } from '../src/world/saturday-world';

    vi.mock('playwright', () => {
    const mockBrowser = {
        close: vi.fn(),
        newContext: vi.fn().mockResolvedValue({
            newPage: vi.fn().mockResolvedValue({
                video: vi.fn().mockReturnValue({ path: vi.fn().mockResolvedValue('video.mp4') }),
                screenshot: vi.fn().mockResolvedValue(Buffer.from('img')),
                close: vi.fn(),
                on: vi.fn() // ConsoleLogger calls on
            }),
            close: vi.fn(),
            request: {},
            on: vi.fn() // TabManager calls on
        }),
    };
    return {
        chromium: { launch: vi.fn().mockResolvedValue(mockBrowser) },
        firefox: { launch: vi.fn().mockResolvedValue(mockBrowser) },
        webkit: { launch: vi.fn().mockResolvedValue(mockBrowser) }
    };
});

vi.mock('node:fs', () => ({
    default: {
        promises: {
            access: vi.fn().mockResolvedValue(undefined),
            readFile: vi.fn().mockResolvedValue(Buffer.from('video'))
        }
    }
}));

describe('Saturday Cucumber Coverage', () => {
    describe('Browser Options', () => {
        it('should have default options', () => {
            expect(browserOptions.headless).toBe(false);
            expect(browserOptions.args).toContain('--use-fake-ui-for-media-stream');
        });
    });

    describe('Hooks Setup', () => {
        beforeEach(async () => {
            process.env.BROWSER = 'chromium';
            await Setup.createBrowser()();
        });

        afterEach(async () => {
            await Setup.closeBrowser()();
        });

        it('should create browser for different engines', async () => {
            const { firefox, webkit, chromium } = await import('playwright');
            
            process.env.BROWSER = 'firefox';
            await Setup.createBrowser()();
            expect(firefox.launch).toHaveBeenCalled();

            process.env.BROWSER = 'webkit';
            await Setup.createBrowser()();
            expect(webkit.launch).toHaveBeenCalled();
            
            process.env.BROWSER = 'chrome';
            await Setup.createBrowser()();
            expect(chromium.launch).toHaveBeenCalled();
        });

        it('should default to chrome when BROWSER env is not set', async () => {
            const { chromium } = await import('playwright');
            delete process.env.BROWSER;
            
            await Setup.createBrowser()();
            expect(chromium.launch).toHaveBeenCalled();
        });

        it('should create and close context', async () => {
             const world = new SaturdayWorld({ attach: vi.fn(), log: vi.fn(), parameters: {}, link: vi.fn() });
             (world as any).attach = vi.fn();
             
             const pickle = { name: 'test scenario' };
             await Setup.createContext().call(world as any, { pickle } as any);
             
             expect(world.page).toBeDefined();
             expect(world.consoleLogger).toBeDefined();

             // Simulate test result
             const result = { status: 'PASSED', duration: { seconds: 1 } };
             await Setup.closeContext().call(world as any, { result } as any);
             
             expect(world.page?.close).toHaveBeenCalled();
             expect(world.attach).toHaveBeenCalled(); // Report attachment
        });

        it('should handle undefined result in closeContext', async () => {
            const world = new SaturdayWorld({ attach: vi.fn(), log: vi.fn(), parameters: {}, link: vi.fn() });
            (world as any).attach = vi.fn();
            
            const pickle = { name: 'test' };
            await Setup.createContext().call(world as any, { pickle } as any);

            // Call with undefined result
            await Setup.closeContext().call(world as any, { result: undefined } as any);
            
            expect(world.page?.close).toHaveBeenCalled();
        });
        
         it('should attach video on failure', async () => {
            const world = new SaturdayWorld({ attach: vi.fn(), log: vi.fn(), parameters: {}, link: vi.fn() });
            (world as any).attach = vi.fn();
            
            const pickle = { name: 'fail' };
            await Setup.createContext().call(world as any, { pickle } as any);

            const result = { status: 'FAILED', duration: { seconds: 1 } };
            // Simulate console logs to verify attachment
            (world.consoleLogger as any).getFormattedLogs = vi.fn().mockReturnValue('Logs');
            
            await Setup.closeContext().call(world as any, { result } as any);
            // Check video attachment which happens in closeContext -> attachVideo
            // attach is called 4 times: Status, screenshot, logs, video
            expect(world.attach).toHaveBeenCalledTimes(4); 
        });

        it('should handle video attachment errors', async () => {
            const world = new SaturdayWorld({ attach: vi.fn(), log: vi.fn(), parameters: {}, link: vi.fn() });
            (world as any).attach = vi.fn();
            (world as any).page = { 
                video: () => ({ path: async () => 'error.mp4' }), 
                close: async () => {},
                screenshot: vi.fn(), 
            };
            (world as any).consoleLogger = { getFormattedLogs: () => '' };
            (world as any).context = { close: async () => {}, on: vi.fn() };
            
            // Force error
             vi.mocked((await import('node:fs')).default.promises.access).mockRejectedValueOnce(new Error('fail'));
             const spy = vi.spyOn(console, 'error').mockImplementation(() => {});
             
             await Setup.closeContext().call(world as any, { result: { status: 'FAILED' } } as any);
             expect(spy).toHaveBeenCalledWith('Error accessing or reading video file:', expect.any(Error));
        });

        it('should handle missing video path', async () => {
            const world = new SaturdayWorld({ attach: vi.fn(), log: vi.fn(), parameters: {}, link: vi.fn() });
            (world as any).attach = vi.fn();
            (world as any).page = { 
                video: () => null, // No video
                close: async () => {},
                screenshot: vi.fn(),
            };
            (world as any).consoleLogger = { getFormattedLogs: () => '' };
            (world as any).context = { close: async () => {}, on: vi.fn() };
            
            // Should not throw
            await Setup.closeContext().call(world as any, { result: { status: 'PASSED' } } as any);
            expect(world.attach).toHaveBeenCalled();
        });

        it('should attach tabs info', async () => {
            const world = new SaturdayWorld({ attach: vi.fn(), log: vi.fn(), parameters: {}, link: vi.fn() });
            (world as any).attach = vi.fn();
            (world as any).context = { close: async () => {}, on: vi.fn() };
            world.page = { screenshot: vi.fn(), close: async () => {}, video: () => null } as any;
            (world as any).consoleLogger = { getFormattedLogs: () => '' };
            
            // mock tabManager
            world.initializeManagers();
            vi.spyOn(world.tabManager, 'count').mockReturnValue(2);
            vi.spyOn(world.tabManager, 'getActiveName').mockReturnValue('tab1');
            vi.spyOn(world.tabManager, 'forEach').mockImplementation(async (cb) => {
                await cb({ url: () => 'http://tab1' } as any, 'tab1', { purpose: 'testing', openedFrom: 'main' });
                await cb({ url: () => { throw new Error('no url'); } } as any, 'tab2', {});
            });

            await Setup.closeContext().call(world as any, { result: { status: 'FAILED', duration: { seconds: 1 } } } as any);
            expect(world.attach).toHaveBeenCalledWith(expect.stringContaining('=== Active Tabs ===\n→ tab1: http://tab1 [testing] (from: main)\n  tab2: unknown'), 'text/plain');
        });

        it('should attach sites info', async () => {
            const world = new SaturdayWorld({ attach: vi.fn(), log: vi.fn(), parameters: {}, link: vi.fn() });
            (world as any).attach = vi.fn();
            (world as any).context = { close: async () => {}, on: vi.fn() };
            world.page = { screenshot: vi.fn(), close: async () => {}, video: () => null } as any;
            (world as any).consoleLogger = { getFormattedLogs: () => '' };
            
            // mock siteManager
            world.initializeManagers();
            vi.spyOn(world.siteManager, 'count').mockReturnValue(1);
            vi.spyOn(world.siteManager, 'listSites').mockReturnValue(['site1', 'site2']);
            vi.spyOn(world.siteManager, 'getActiveName').mockReturnValue('site1');
            vi.spyOn(world.siteManager, 'get').mockImplementation((name: string) => {
                if (name === 'site1') return { constructor: { name: 'MySite' }, Object, getBaseUrl: () => 'http://site1' } as any;
                return { constructor: { name: 'OtherSite' }, getBaseUrl: () => { throw new Error('no base'); } } as any;
            });

            await Setup.closeContext().call(world as any, { result: { status: 'FAILED', duration: { seconds: 1 } } } as any);
            expect(world.attach).toHaveBeenCalledWith(expect.stringContaining('=== Registered Sites ===\n→ site1: http://site1 (MySite)\n  site2: unknown (OtherSite)'), 'text/plain');
        });

        it('should handle screenshot error', async () => {
            const world = new SaturdayWorld({ attach: vi.fn(), log: vi.fn(), parameters: {}, link: vi.fn() });
            (world as any).attach = vi.fn();
            (world as any).context = { close: async () => {}, on: vi.fn() };
            world.page = { screenshot: vi.fn().mockRejectedValue(new Error('no screenshot')), close: async () => {}, video: () => null } as any;
            (world as any).consoleLogger = { getFormattedLogs: () => '' };
            
            const consoleSpy = vi.spyOn(console, 'warn').mockImplementation(() => {});

            await Setup.closeContext().call(world as any, { result: { status: 'PASSED', duration: { seconds: 1 } } } as any);
            expect(consoleSpy).toHaveBeenCalledWith('Failed to take screenshot', expect.any(Error));
        });
    });
});
