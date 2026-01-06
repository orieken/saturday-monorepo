
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { DetectorsFacade } from '../../src/ml/facades/detectors.facade';

describe('DetectorsFacade Branch Coverage', () => {
    let mockSite: any;
    let facade: DetectorsFacade;

    beforeEach(() => {
        vi.useFakeTimers();
        mockSite = {
            page: {
                screenshot: vi.fn().mockResolvedValue(Buffer.from('test')),
                goto: vi.fn().mockResolvedValue(null)
            },
            getPage: vi.fn().mockReturnValue({ visit: vi.fn() }),
            createDetectorsFacade: vi.fn()
        };
        // Explicitly set pageObject for the mock since it's a getter in the real class
        Object.defineProperty(mockSite, 'pageObject', {
            get: () => mockSite.page
        });
        facade = new DetectorsFacade(mockSite);
        // Mock sub-detectors
        (facade as any)._anomalyDetector = {
            validateImage: vi.fn().mockResolvedValue({ isValid: true, confidence: 1, anomalies: [], score: 1 }),
            validateCurrentPage: vi.fn().mockResolvedValue({ isValid: true, confidence: 1, anomalies: [], score: 1 }),
            detectLayoutChanges: vi.fn().mockResolvedValue([]),
            validateRegions: vi.fn().mockResolvedValue([])
        };
        (facade as any)._elementDetector = { validateElement: vi.fn() };
        (facade as any)._layoutDetector = { detectAnomalies: vi.fn().mockResolvedValue([]) };
        (facade as any)._regressionDetector = { 
            compareWithBaseline: vi.fn().mockResolvedValue({ similarity: 1, overallMatch: true }),  
            compareScreenshots: vi.fn().mockResolvedValue({ similarity: 1, overallMatch: true })
        };
    });

    afterEach(() => {
        vi.useRealTimers();
    });

    it('should cover compareVariants branches', async () => {
        // We need to verify what is actually called. 
        // compareVariants calls validateCurrentPage.
        // It also calls regressionDetector.compareScreenshots. which is mocked in beforeEach.
        (facade.regressionDetector.compareScreenshots as any).mockResolvedValue({ similarity: 1, overallMatch: true });

        const spy = vi.spyOn(facade, 'validateCurrentPage');
        
        // B significantly better
        spy.mockResolvedValueOnce({ confidence: 0.5 } as any)
           .mockResolvedValueOnce({ confidence: 0.8 } as any);
        let res = await facade.compareVariants('A', 'B', 'm');
        expect(res.recommendation).toBe('B');

        // A significantly better
        spy.mockReset();
        spy.mockResolvedValueOnce({ confidence: 0.8 } as any)
           .mockResolvedValueOnce({ confidence: 0.5 } as any);
        res = await facade.compareVariants('A', 'B', 'm');
        expect(res.recommendation).toBe('A');

        // Close enough (inconclusive)
        spy.mockReset();
        spy.mockResolvedValueOnce({ confidence: 0.8 } as any)
           .mockResolvedValueOnce({ confidence: 0.79 } as any);
        res = await facade.compareVariants('A', 'B', 'm');
        expect(res.recommendation).toBe('inconclusive');
    });

    it('should cover monitorPagePerformance loop', async () => {
        const spyDetect = vi.spyOn(facade, 'validateCurrentPage').mockResolvedValue({ confidence: 1, anomalies: [] } as any);
        
        const promise = facade.monitorPagePerformance('p', 3);
        
        await vi.runAllTimersAsync();
        
        const res = await promise;
        expect(res.averageConfidence).toBe(1);
        expect(spyDetect).toHaveBeenCalledTimes(3);
    });

    it('should cover runFullValidation aggregation', async () => {
        const res = await facade.runFullValidation('m');
        expect(res.overall.isValid).toBe(true);
        // Maybe mock partial failures to hit branches in aggregation?
        // if aggregation logic exists
    });

     it('should cover validateMultipleElements empty', async () => {
        const res = await facade.validateMultipleElements([]);
        expect(res).toHaveLength(0);
    });
});
