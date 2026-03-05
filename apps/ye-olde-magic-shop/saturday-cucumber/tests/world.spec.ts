
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
});
