
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { DetectorsFacade } from '../../src/ml/facades/detectors.facade';
import { BaseSite } from '../../src/base/base-site';

describe('DetectorsFacade', () => {
    let mockSite: any;
    let facade: DetectorsFacade;

    beforeEach(() => {
       mockSite = {
           getPage: vi.fn().mockReturnValue({ visit: vi.fn(), captureForTraining: vi.fn() }),
       };
       facade = new DetectorsFacade(mockSite as BaseSite);
    });

    it('should initialize detectors lazy', () => {
        expect(facade.anomalyDetector).toBeDefined();
        expect(facade.elementDetector).toBeDefined();
        expect(facade.layoutDetector).toBeDefined();
        expect(facade.regressionDetector).toBeDefined();
    });

    it('should validate current page', async () => {
        const result = await facade.validateCurrentPage('label');
        expect(result.isValid).toBe(true);
    });

    it('should validate element', async () => {
        const result = await facade.validateElement('.selector', 'label');
        expect(result.isValid).toBe(true);
    });



    it('should validate user journey', async () => {
        const results = await facade.validateUserJourney([
            { pageName: 'home', modelLabel: 'label' }
        ]);
        expect(results).toHaveLength(1);
    });

    it('should validate multiple elements', async () => {
        const results = await facade.validateMultipleElements([
            { selector: '.a', modelLabel: 'm' }
        ]);
        expect(results).toHaveLength(1);
    });

    it('should run full validation', async () => {
        const result = await facade.runFullValidation('model');
        expect(result.overall.isValid).toBeDefined();
    });

    it('should monitor page performance', async () => {
        // Mock site.pageObject methods used in monitorPagePerformance
        mockSite.pageObject = { screenshot: vi.fn(), goto: vi.fn() };
        
        const result = await facade.monitorPagePerformance('model', 2);
        expect(result.averageConfidence).toBe(1);
        expect(result.performanceMetrics).toBeDefined();
    });

    it('should compare variants', async () => {
        mockSite.pageObject = { screenshot: vi.fn().mockResolvedValue(Buffer.from('a')), goto: vi.fn() };
        
        const result = await facade.compareVariants('urlA', 'urlB', 'model');
        expect(result.recommendation).toBeDefined();
    });
    
    it('should detect layout anomalies', async () => {
        const res = await facade.detectLayoutAnomalies();
        expect(res).toBeDefined();
    });

    it('should compare with baseline', async () => {
        const res = await facade.compareWithBaseline('base');
        expect(res).toBeDefined();
    });
    
    it('should generate diagnostic report', async () => {
        const res = await facade.generateDiagnosticReport('model');
        expect(res.recommendations.length).toBeGreaterThan(0);
    });

    it('should compare variants where B is better', async () => {
        mockSite.pageObject = { screenshot: vi.fn(), goto: vi.fn() };
        const spy = vi.spyOn(facade, 'validateCurrentPage');
        spy.mockResolvedValueOnce({ isValid: true, confidence: 0.1, anomalies: [], score: 0.1, timestamp: new Date(), modelUsed: 'm' } as any)
           .mockResolvedValueOnce({ isValid: true, confidence: 0.9, anomalies: [], score: 0.9, timestamp: new Date(), modelUsed: 'm' } as any);
           
        const result = await facade.compareVariants('A', 'B', 'model');
        expect(result.recommendation).toBe('B');
    });
});
