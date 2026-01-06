
import { describe, it, expect, vi } from 'vitest';
import { BaseFlow } from '../../src/base/base-flow';
import { BaseSite } from '../../src/base/base-site';

class TestFlow extends BaseFlow {
  async execute() {}
}

describe('BaseFlow', () => {
    it('should be instantiable and delegate to site', () => {
        const mockSite = {
            getPage: vi.fn().mockReturnValue('page'),
            getFlow: vi.fn().mockReturnValue('flow')
        };
        const flow = new TestFlow(mockSite as any as BaseSite);
        
        expect((flow as any).getPage('p')).toBe('page');
        expect((flow as any).getFlow('f')).toBe('flow');
        expect(mockSite.getPage).toHaveBeenCalledWith('p');
        expect(mockSite.getFlow).toHaveBeenCalledWith('f');
    });
});
