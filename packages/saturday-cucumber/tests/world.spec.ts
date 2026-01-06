
import { describe, it, expect, vi } from 'vitest';
import { CustomWorld, installSaturdayWorld } from '../src/world/custom-world';
import * as cucumber from '@cucumber/cucumber';

vi.mock('@cucumber/cucumber', () => ({
    setWorldConstructor: vi.fn(),
    World: class { constructor(public options: any) {} }
}));

describe('CustomWorld', () => {
    it('should be instantiable', () => {
        const world = new CustomWorld({} as any);
        expect(world).toBeDefined();
    });

    it('should install world constructor', () => {
        installSaturdayWorld();
        expect(cucumber.setWorldConstructor).toHaveBeenCalledWith(CustomWorld);
    });
});
