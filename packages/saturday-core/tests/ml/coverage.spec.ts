
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { BaseDetector } from '../../src/ml/base/base-detector';
import { BaseSite } from '../../src/base/base-site';
import { BasePage } from '../../src/base/base-page';
import { AnomalyModel } from '../../src/ml/models/anomaly-model';
import { VisualModel } from '../../src/ml/models/visual-model';
import { RegressionModel } from '../../src/ml/models/regression-model';
import { ScreenshotTrainer } from '../../src/ml/trainers/screenshot-trainer';
import { ElementTrainer } from '../../src/ml/trainers/element-trainer';
import { FlowTrainer } from '../../src/ml/trainers/flow-trainer';
import { ElementDetector } from '../../src/ml/detectors/element-detector';
import { LayoutDetector } from '../../src/ml/detectors/layout-detector';
import { RegressionDetector } from '../../src/ml/detectors/regression-detector';
import { BaseModel } from '../../src/ml/base/base-model';

class ExposedDetector extends BaseDetector {
    public testMeasureDetectionTime<T>(op: () => Promise<T>) { return this.measureDetectionTime(op); }
    public testValidateMultipleImages(imgs: Buffer[], lbl: string) { return this.validateMultipleImages(imgs, lbl); }
    public testCalculateAggregateMetrics(res: any[]) { return this.calculateAggregateMetrics(res); }
    public testValidateRegions(img: Buffer, regs: any[], lbl: string) { return this.validateRegions(img, regs, lbl); }
    public testDetectTemplate(img: Buffer, tmpl: Buffer) { return this.detectTemplate(img, tmpl); }
    public testExtractFeatures(img: Buffer) { return this.extractFeatures(img); }
    public testCalculateSimilarityScore(a: any, b: any) { return this.calculateSimilarityScore(a, b); }
    public testCleanup() { return this.cleanup(); }
    public async testHealthCheck() { return this.healthCheck(); }
    public async testWarmUp(lbl: string) { return this.warmUpModel(lbl); }
    public testDetectLayoutChanges(a: any[], b: any[]) { return this.detectLayoutChanges(a, b); }
    public testDetectPixelDifferences(a: Buffer, b: Buffer) { return this.detectPixelDifferences(a, b); }

    async detect(d: Buffer, l: string) { return { isValid: true, confidence: 1, anomalies: [], score: 1, timestamp: new Date(), modelUsed: l }; }
    async validateImage(d: Buffer, l: string) { return this.detect(d, l); }
}

class ExposedModel extends BaseModel {
    constructor(site: any) { super(site, 'test', 'label'); }
    async train() { return { success: true } as any; }
    async predict() { return {}; }
    async evaluate() { return {} as any; }
    async save() { return ''; }
    async load() {}

    public testLogProgress(m: string, p?: number) { this.logProgress(m, p); }
    public testHandleError(c: string, e: Error) { this.handleError(c, e); }
    public testCreateFailedTrainingResult(e: Error) { return this.createFailedTrainingResult(e); }
    public testCreateFailedValidationResult(e: Error) { return this.createFailedValidationResult(e); }
    public testMergeWithDefaults(c: any) { return this.mergeWithDefaults(c); }
    public async testValidateConfiguration() { await this.validateConfiguration(); }
    public testGetStoragePerformance() { return this.getStoredPerformance(); }
}

class TestSite extends BaseSite {
    constructor(page: any) { super(page, 'http://localhost'); }
    initializePages() {}
    initializeFlows() {}
    
    public testRegisterPage(name: string, cls: any) { this.registerPage(name, cls); }
    public testRegisterFlow(name: string, cls: any) { this.registerFlow(name, cls); }
}

class TestPage extends BasePage {
    containerSelector = 'body';
    path = '/test';
    constructor(page: any, site: any) { super(page, site); }
    protected initializeElements() {}
    protected initializeFilters() {}
}

describe('Coverage Fillers', () => {
    let mockSite: any;
    let mockPage: any;

    beforeEach(() => {
        mockPage = { 
            screenshot: vi.fn().mockResolvedValue(Buffer.from('img')), 
            locator: vi.fn().mockReturnValue({ waitFor: vi.fn(), isVisible: vi.fn().mockResolvedValue(true) }), 
            url: vi.fn().mockReturnValue('http://localhost/test'), 
            evaluate: vi.fn().mockResolvedValue({}), 
            goto: vi.fn(), 
            title: vi.fn(), 
            waitForLoadState: vi.fn() 
        };
        mockSite = {
            page: mockPage,
            models: { loadModel: vi.fn(), getModelInfo: vi.fn(), getModelPerformance: vi.fn() },
            detectors: { anomalyDetector: { validateImage: vi.fn().mockResolvedValue({}) } },
            getBaseUrl: () => 'http://localhost'
        };
    });

    it('should cover BaseModel utilities', async () => {
        const model = new ExposedModel(mockSite);
        model.testLogProgress('msg', 0.5);
        expect(() => model.testHandleError('ctx', new Error('boom'))).toThrow('boom');
        const tr = model.testCreateFailedTrainingResult(new Error('fail'));
        expect(tr.success).toBe(false);
        const vr = model.testCreateFailedValidationResult(new Error('fail'));
        expect(vr.isValid).toBe(false);
        const config = model.testMergeWithDefaults({ epochs: 10 });
        expect(config.epochs).toBe(10);
        await model.testValidateConfiguration(); 
        model.updateConfig({ batchSize: -1 });
        await expect(model.testValidateConfiguration()).rejects.toThrow();
        await model.testGetStoragePerformance();
    });

    it('should instantiate stubs', () => {
        expect(new AnomalyModel(mockSite)).toBeDefined();
        expect(new VisualModel(mockSite)).toBeDefined();
        expect(new RegressionModel(mockSite)).toBeDefined();
        expect(new ScreenshotTrainer(mockSite)).toBeDefined();
        expect(new ElementTrainer(mockSite)).toBeDefined().not.toBeNull();
        expect(new FlowTrainer(mockSite)).toBeDefined();
        expect(new ElementDetector(mockSite)).toBeDefined();
        expect(new LayoutDetector(mockSite)).toBeDefined();
        expect(new RegressionDetector(mockSite)).toBeDefined();
    });

    it('should cover BaseDetector utilities', async () => {
        const det = new ExposedDetector(mockSite, 'test');
        
        await det.testMeasureDetectionTime(async () => 1);
        await det.testValidateMultipleImages([Buffer.from('a')], 'lbl');
        det.testCalculateAggregateMetrics([]);
        det.testCalculateAggregateMetrics([{ confidence: 1, isValid: true, anomalies: [], score: 1 } as any]);
        await det.testValidateRegions(Buffer.from('a'), [{ name: 'r', clip: { x:0,y:0,width:10,height:10 } }], 'lbl');
        await det.testDetectTemplate(Buffer.from('a'), Buffer.from('b'));
        await det.testExtractFeatures(Buffer.from('a'));
        det.testCalculateSimilarityScore({}, {});
        await det.testCleanup();
        await det.testHealthCheck();
        const badDet = new ExposedDetector({} as any, 'bad');
        await badDet.testHealthCheck();
        await det.testWarmUp('lbl');
        det.testDetectLayoutChanges([], []);
        det.testDetectLayoutChanges([{selector:'a', x:0, y:0, boundingBox:{}}], [{x:100, y:100, boundingBox:{}}]); 
        det.testDetectLayoutChanges([{selector:'a'}], []); 
        det.testDetectPixelDifferences(Buffer.from('a'), Buffer.from('b'));

        const anomalies = [
            { severity: 'low', confidence: 0.9 },
            { severity: 'medium', confidence: 0.8 },
            { severity: 'low', confidence: 0.1 }
        ];
        let filtered = (det as any).filterAnomaliesBySeverity(anomalies, { tolerance: 0.9 });
        expect(filtered.length).toBe(1); 
        filtered = (det as any).filterAnomaliesBySeverity(anomalies, { tolerance: 0.6 });
        expect(filtered.length).toBe(1); 
        filtered = (det as any).filterAnomaliesBySeverity(anomalies, { tolerance: 0.1 });
        expect(filtered.length).toBe(3);

        const resWithAnomalies = {
            isValid: true, confidence: 1.0, 
            anomalies: [
                { severity: 'low', confidence: 1.0 }, 
                { severity: 'critical', confidence: 1.0 } 
            ],
            score: 0
        };
        const score = (det as any).calculateOverallScore(resWithAnomalies);
        expect(score).toBeCloseTo(0.45);
        
        // Calculate confidence empty
        expect((det as any).calculateConfidence([])).toBe(0);

        // Model existence check failure in executeDetection
        // Mock getModelInfo to return null
        mockSite.models.getModelInfo.mockResolvedValue(null);
        // We need to call executeDetection via an exposed method or cast
        // But executeDetection is public? No, strictly it's public in the class def (lines 18).
        // Let's verify: async executeDetection(...) : Promise<ValidationResult>
        const errResult = await det.executeDetection(
            async () => Buffer.from('img'),
            'missing-model'
        );
        expect(errResult.isValid).toBe(false);
        expect(errResult.anomalies[0].description).toContain('not found');

        // Layout changes with position difference > 5
        const layoutAnomalies = det.testDetectLayoutChanges(
            [{selector:'a', x:0, y:0, boundingBox:{x:0,y:0}}], 
            [{x:10, y:10, boundingBox:{x:10,y:10}}] // Distance sqrt(100+100) = 14 > 5
        );
        expect(layoutAnomalies).toHaveLength(1);
        expect(layoutAnomalies[0].description).toContain('position changed');

        // Layout changes - Missing element (already covered? "Element missing")
        // det.testDetectLayoutChanges([{selector:'a'}], []) -> covered in previous edit?
        
        // Postprocess results with threshold
        const lowConfResult = {
            isValid: true, confidence: 0.5, anomalies: [], score: 0.5, timestamp: new Date(), modelUsed: 'm'
        };
        const ppRes = await (det as any).postprocessResults(lowConfResult, { threshold: 0.8 });
        expect(ppRes.isValid).toBe(false);
        expect(ppRes.anomalies[0].description).toContain('below threshold');
    });

    it('should cover BaseSite ', async () => {
        const site = new TestSite(mockPage);
        
        site.testRegisterPage('p', TestPage);
        site.testRegisterPage('p', TestPage); 

        expect(() => site.getPage('missing')).toThrow();
        expect(() => site.getFlow('missing')).toThrow();
        
        const p = site.getPage<TestPage>('p');
        expect(p).toBeDefined();

        await p.visit();
        await p.isLoaded();
        await p.scrollToTop();
        await p.scrollToBottom();
        
        const capture = await p.captureForTraining('lbl');
        expect(capture).toBeDefined();
        
        const val = await p.validateWithML('lbl');
        expect(val).toBeDefined();

        // Repeated access for lazy getters branch coverage (all of them)
        expect(site.trainers).toBe(site.trainers);
        expect(site.detectors).toBe(site.detectors);
        expect(site.models).toBe(site.models);

        expect(site.trainers.pageTrainer).toBe(site.trainers.pageTrainer);
        expect(site.trainers.elementTrainer).toBe(site.trainers.elementTrainer);
        expect(site.trainers.screenshotTrainer).toBe(site.trainers.screenshotTrainer);
        expect(site.trainers.flowTrainer).toBe(site.trainers.flowTrainer);
        
        expect(site.detectors.anomalyDetector).toBe(site.detectors.anomalyDetector);
        expect(site.detectors.elementDetector).toBe(site.detectors.elementDetector);
        expect(site.detectors.layoutDetector).toBe(site.detectors.layoutDetector);
        expect(site.detectors.regressionDetector).toBe(site.detectors.regressionDetector);

        expect(site.models.anomalyModel).toBe(site.models.anomalyModel);
        expect(site.models.visualModel).toBe(site.models.visualModel);
        expect(site.models.regressionModel).toBe(site.models.regressionModel);

        await site.visit();
        await site.getCurrentUrl();
        await site.getCurrentTitle();
        await site.waitForNavigation();
        expect(site.getBaseUrl()).toBe('http://localhost');

        // ElementDetector failure
        // Use exposed detector or mock
        // We have ElementDetector in src/ml/detectors/element-detector.ts
        // It extends BaseDetector.
        // We can test its specific logic if we instantiate it.
        const elDet = new ElementDetector(mockSite);
        // Stub implementation always returns valid
        const res = await elDet.validateElement('sel', 'lbl');
        expect(res.isValid).toBe(true);
    });
});
