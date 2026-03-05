
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { BaseDetector } from '../../src/ml/base/base-detector';
import { BaseSite } from '../../src/base/base-site';

class TestDetector extends BaseDetector {
    constructor(site: BaseSite) { super(site, 'test'); }
    async detect(data: any) { return { isValid: true, confidence: 1, anomalies: [], score: 1 } as any; }
    async validateImage(data: any) { return { isValid: true, confidence: 1, anomalies: [], score: 1 } as any; }

    // Expose protected
    async testValidateModelExists(l: string) { return this.validateModelExists(l); }
    async testPostprocessResults(r: any, o: any) { return this.postprocessResults(r, o); }
    async testFilterAnomalies(a: any, o: any) { return this.filterAnomaliesBySeverity(a, o); }
    async testCalculateOverallScore(r: any) { return this.calculateOverallScore(r); }
    async testDetectLayoutChanges(b: any, a: any) { return this.detectLayoutChanges(b, a); }
    async testDetectPixelDifferences(a: any, b: any) { return this.detectPixelDifferences(a, b); }
    async testHealthCheck() { return this.healthCheck(); }
}

describe('BaseDetector Branch Coverage', () => {
    let mockSite: any;
    let detector: TestDetector;

    beforeEach(() => {
        mockSite = {
            page: {
                screenshot: vi.fn(),
                locator: vi.fn()
            },
            models: {
                getModelInfo: vi.fn().mockResolvedValue({}),
                loadModel: vi.fn()
            }
        };
        detector = new TestDetector(mockSite);
    });

    it('should validate model exists', async () => {
        expect(await detector.testValidateModelExists('l')).toBe(true);
        mockSite.models.getModelInfo.mockResolvedValue(null);
        expect(await detector.testValidateModelExists('l')).toBe(false);
        mockSite.models.getModelInfo.mockRejectedValue(new Error('fail'));
        expect(await detector.testValidateModelExists('l')).toBe(false);
    });

    it('should filter anomalies by severity', async () => {
        const anomalies = [
            { severity: 'low', confidence: 0.8 },
            { severity: 'medium', confidence: 0.8 },
            { severity: 'low', confidence: 0.95 }
        ];

        // Tolerance > 0.8 (filters low)
        let filtered = await detector.testFilterAnomalies(anomalies as any, { tolerance: 0.9 });
        expect(filtered).toHaveLength(1); // ONLY medium kept.

        // Tolerance > 0.5 (filters low UNLESS conf > 0.9)
        filtered = await detector.testFilterAnomalies(anomalies as any, { tolerance: 0.6 });
        expect(filtered).toHaveLength(2);
    });

    it('should calculate overall score', async () => {
        // No anomalies
        expect(await detector.testCalculateOverallScore({ confidence: 1, anomalies: [] })).toBe(1);

        // With anomalies
        const res = {
            confidence: 1,
            anomalies: [
                { severity: 'low', confidence: 1 }, // penalty 0.05
                { severity: 'critical', confidence: 1 } // penalty 0.5
            ]
        };
        // 1 - (0.05 + 0.5) = 0.45
        expect(await detector.testCalculateOverallScore(res as any)).toBeCloseTo(0.45);
    });

    it('should detect layout changes', async () => {
        const before = [{ selector: 'a', x:0, y:0, boundingBox: {} }];
        
        // Missing
        const afterMissing = []; 
        let anomalies = await detector.testDetectLayoutChanges(before, afterMissing);
        expect(anomalies[0].description).toContain('Element missing');

        // Moved
        const afterMoved = [{ x: 10, y: 10, boundingBox: {} }]; // dist > 5
        anomalies = await detector.testDetectLayoutChanges(before, afterMoved);
        expect(anomalies[0].description).toContain('position changed');
        
        // Same
        const afterSame = [{ x: 1, y: 1 }]; // dist < 5
        anomalies = await detector.testDetectLayoutChanges(before, afterSame);
        expect(anomalies).toHaveLength(0);
    });

    it('should handle healthCheck branches', async () => {
        // Healthy (fast)
        let hc = await detector.testHealthCheck();
        expect(hc.status).toBe('healthy');

        // Degraded (slow)
        // Stub detect implementation to be slow
        // We need to cast to any to overwrite abstract method with concrete compatible signature
        (detector as any).detect = async () => { await new Promise(r => setTimeout(r, 1100)); return {} as any; };
        hc = await detector.testHealthCheck();
        expect(hc.status).toBe('degraded');

        (detector as any).detect = () => { throw new Error('fail'); };
        hc = await detector.testHealthCheck();
        expect(hc.status).toBe('unhealthy');
    });

    it('should calculate aggregate metrics', async () => {
        // Empty
        let metrics = (detector as any).calculateAggregateMetrics([]);
        expect(metrics.totalAnomalies).toBe(0);

        // Populated
        const results = [
            { confidence: 1, isValid: true, anomalies: [], score: 1 },
            { confidence: 0.5, isValid: false, anomalies: [{}], score: 0.5 }
        ];
        metrics = (detector as any).calculateAggregateMetrics(results);
        expect(metrics.averageConfidence).toBe(0.75);
        expect(metrics.successRate).toBe(0.5);
    });

    it('should validate regions', async () => {
        const regions = [{ name: 'r1', clip: { x:0, y:0, width:10, height:10 } }];
        const res = await (detector as any).validateRegions(Buffer.from('img'), regions, 'l');
        expect(res).toHaveLength(1);
        expect(res[0].region).toBe('r1');
    });

    it('should cover stubs and utilities', async () => {
        // detectTemplate
        expect(await (detector as any).detectTemplate()).toEqual([]);
        // extractFeatures
        const f = await (detector as any).extractFeatures();
        expect(f.edges).toBe(0);
        // calculateSimilarityScore
        const sim = (detector as any).calculateSimilarityScore();
        expect(sim).toBeGreaterThanOrEqual(0.7);
        // cleanup
        await detector.cleanup();
    });

    it('should warm up model', async () => {
        // Success
        await detector.warmUpModel('l');
        
        // Fail
        mockSite.models.loadModel.mockRejectedValue(new Error('fail'));
        await detector.warmUpModel('l'); // Should not throw, logs error
    });

    it('should postprocess results thresholds', async () => {
        const res = { isValid: true, confidence: 0.5, anomalies: [], score: 0.5 };
        const pp = await detector.testPostprocessResults(res, { threshold: 0.8 });
        expect(pp.isValid).toBe(false);
        expect(pp.anomalies[0].description).toContain('below threshold');
    });
});
