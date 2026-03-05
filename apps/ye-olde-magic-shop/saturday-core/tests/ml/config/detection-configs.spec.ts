
import { describe, it, expect } from 'vitest';
import { 
    getDetectionConfig, 
    createCustomDetectionConfig, 
    validateDetectionConfig, 
    optimizeConfigForPerformance, 
    optimizeConfigForAccuracy, 
    DetectionConfigs 
} from '../../../src/ml/config/detection-configs';

describe('Detection Configs Coverage', () => {
    it('should get default detection configs', () => {
        const config = getDetectionConfig('development', 'anomaly');
        expect(config.enabled).toBe(true);
        expect(config.type).toBe('anomaly');
    });

    it('should apply site overrides', () => {
        const config = getDetectionConfig('development', 'anomaly', 'foo');
        // Check override from SiteDetectionConfigs
        // 'foo' sensitivity is 0.8 vs dev default 0.7
        expect((config as any).anomalyTypes.visual.sensitivity).toBe(0.8);
    });

    it('should create custom detection config', () => {
         const base = getDetectionConfig('development', 'anomaly');
         const custom = createCustomDetectionConfig(base, { priority: 99 });
         expect(custom.priority).toBe(99);
    });

    it('should validate detection config', () => {
        const base = getDetectionConfig('development', 'anomaly');
        expect(validateDetectionConfig(base)).toBe(true);

        expect(validateDetectionConfig({ ...base, thresholds: { ...base.thresholds, confidence: 1.1 } })).toBe(false);
        expect(validateDetectionConfig({ ...base, thresholds: { ...base.thresholds, confidence: -1 } })).toBe(false);
        expect(validateDetectionConfig({ ...base, thresholds: { ...base.thresholds, similarity: 1.1 } })).toBe(false);
        expect(validateDetectionConfig({ ...base, thresholds: { ...base.thresholds, similarity: -1 } })).toBe(false);
        expect(validateDetectionConfig({ ...base, thresholds: { ...base.thresholds, anomaly: 1.1 } })).toBe(false);
        expect(validateDetectionConfig({ ...base, thresholds: { ...base.thresholds, anomaly: -1 } })).toBe(false);
        expect(validateDetectionConfig({ ...base, thresholds: { ...base.thresholds, layout: -1 } })).toBe(false);
        
        expect(validateDetectionConfig({ ...base, performance: { ...base.performance, maxConcurrent: 0 } })).toBe(false);
        expect(validateDetectionConfig({ ...base, performance: { ...base.performance, timeout: 0 } })).toBe(false);
        expect(validateDetectionConfig({ ...base, performance: { ...base.performance, batchSize: 0 } })).toBe(false);
    });

    it('should optimize for performance', () => {
        const base = getDetectionConfig('development', 'anomaly');
        const opt = optimizeConfigForPerformance(base);
        expect(opt.preprocessing.denoise).toBe(false);
        expect(opt.performance.cacheResults).toBe(true);
    });

    it('should optimize for accuracy', () => {
        const base = getDetectionConfig('development', 'anomaly');
        const opt = optimizeConfigForAccuracy(base);
        expect(opt.preprocessing.denoise).toBe(true);
        expect(opt.thresholds.confidence).toBeGreaterThan(base.thresholds.confidence);
    });

    it('should handle deep merge recursion', () => {
        // deepMerge is internal but tested via getDetectionConfig with site overrides
        // We can double check nested objects merge
        const c = getDetectionConfig('development', 'anomaly', 'foo');
        // 'anomalyTypes' is nested, 'visual' is nested.
        // visual enabled should remain true (merged from base)
        expect((c as any).anomalyTypes.visual.enabled).toBe(true); 
    });
});
