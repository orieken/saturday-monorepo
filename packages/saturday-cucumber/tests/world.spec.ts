
import { describe, it, expect, vi } from 'vitest';
import { SaturdayWorld, installSaturdayWorld } from '../src/world/saturday-world';
import * as cucumber from '@cucumber/cucumber';

vi.mock('@cucumber/cucumber', () => ({
    setWorldConstructor: vi.fn(),
    World: class { constructor(public options: any) {} }
}));

describe('SaturdayWorld', () => {
    it('should be instantiable', () => {
        const world = new SaturdayWorld({} as any);
        expect(world).toBeDefined();
    });

    it('should install world constructor', () => {
        installSaturdayWorld();
        expect(cucumber.setWorldConstructor).toHaveBeenCalledWith(SaturdayWorld);
    });

    it('should throw when accessing siteManager without page', () => {
        const world = new SaturdayWorld({} as any);
        expect(() => world.siteManager).toThrow('Page must be initialized before utilizing SiteManager defaults');
    });

    it('should initialize siteManager when page is present', () => {
        const world = new SaturdayWorld({} as any);
        world.page = {} as any;
        expect(world.siteManager).toBeDefined();
    });

    it('should throw when accessing tabManager without context', () => {
        const world = new SaturdayWorld({} as any);
        expect(() => world.tabManager).toThrow('Browser context must be initialized before using TabManager');
    });

    it('should initialize tabManager when context and page are present', () => {
        const world = new SaturdayWorld({} as any);
        world.context = { on: vi.fn() } as any;
        world.page = {} as any;
        expect(world.tabManager).toBeDefined();
        world.initializeManagers(); // cover initializeManagers
    });

    it('should cleanup managers and handle errors', async () => {
        const world = new SaturdayWorld({} as any);
        world.context = { on: vi.fn() } as any;
        world.page = {} as any;
        
        const tabManager = world.tabManager;
        const siteManager = world.siteManager;
        
        vi.spyOn(tabManager, 'closeAll').mockRejectedValue(new Error('close error'));
        const consoleSpy = vi.spyOn(console, 'warn').mockImplementation(() => {});
        vi.spyOn(siteManager, 'clear');

        await world.cleanupManagers();

        expect(consoleSpy).toHaveBeenCalledWith('Error closing tabs:', expect.any(Error));
        expect(siteManager.clear).toHaveBeenCalled();
    });
});
