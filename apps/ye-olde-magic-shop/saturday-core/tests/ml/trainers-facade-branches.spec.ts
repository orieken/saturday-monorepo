
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { TrainersFacade } from '../../src/ml/facades/trainers.facade';

describe('TrainersFacade Branch Coverage', () => {
    let mockSite: any;
    let facade: TrainersFacade;

    beforeEach(() => {
        mockSite = {
            getPage: vi.fn().mockReturnValue({ visit: vi.fn(), setViewportSize: vi.fn() as any }),
            getFlow: vi.fn().mockReturnValue({}),
            pageObject: { setViewportSize: vi.fn().mockResolvedValue(undefined) }
        };
        facade = new TrainersFacade(mockSite);
        // Mock sub-trainers
        (facade as any)._pageTrainer = { trainCurrentPage: vi.fn().mockResolvedValue({ success: true }) };
        (facade as any)._elementTrainer = { trainElement: vi.fn().mockResolvedValue({ success: true }) };
        (facade as any)._flowTrainer = { trainFlow: vi.fn().mockResolvedValue({ success: true }) };
    });

    it('should train complete user journey branches', async () => {
        const steps = [
            { type: 'page', target: 'p' },
            { type: 'element', target: 'e' },
            { type: 'flow', target: 'f' }
        ];

        const results = await facade.trainCompleteUserJourney('j', steps as any);
        expect(results).toHaveLength(3);
        
        // Invalid type
        await expect(facade.trainCompleteUserJourney('j', [{ type: 'unknown' } as any])).rejects.toThrow('Unknown training step type');
    });

    it('should train multiple pages', async () => {
        const results = await facade.trainMultiplePages([{ pageName: 'p', label: 'l' }]);
        expect(results).toHaveLength(1);
    });

    it('should train across browsers', async () => {
        const results = await facade.trainAcrossBrowsers('l', ['chrome', 'firefox']);
        expect(results).toHaveLength(2);
        // Verify metadata was passed
        const spy = (facade.pageTrainer.trainCurrentPage as any);
        expect(spy.mock.calls[0][1].metadata.browser).toBe('chrome');
    });

    it('should train responsive design', async () => {
        const viewports = [{ name: 'mobile', width: 300, height: 600 }];
        const results = await facade.trainResponsiveDesign('l', viewports);
        expect(results).toHaveLength(1);
        expect(mockSite.pageObject.setViewportSize).toHaveBeenCalledWith({ width: 300, height: 600 });
    });

    it('should get training stats', async () => {
        const stats = await facade.getTrainingStats();
        expect(stats.totalSamples).toBe(0);
    });

    it('should compare model versions', async () => {
        const res = await facade.compareModelVersions('l', ['1']);
        expect(res).toHaveLength(1);
    });
});
