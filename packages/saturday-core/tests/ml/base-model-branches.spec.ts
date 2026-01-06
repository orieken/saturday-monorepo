
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { BaseModel } from '../../src/ml/base/base-model';
import { BaseSite } from '../../src/base/base-site';

class TestModel extends BaseModel {
    constructor(site: BaseSite, config?: any) { super(site, 'testType', 'testLabel', config); }
    async train(data: any) { return { success: true, metrics: {} } as any; }
    async predict(data: any) { return [0.9]; }
    async evaluate(data: any) { return { accuracy: 0.9, precision: 0.9, recall: 0.9, f1Score: 0.9 } as any; }
    async save() { return 'path'; }
    async load(p: string) { this.isLoaded = true; }
    
    // Expose protected methods
    async testValidateConfiguration() { return this.validateConfiguration(); }
    async testValidateTrainingData(d: any) { return this.validateTrainingData(d); }
    async testCalculateMetrics(p: any[], g: any[]) { return this.calculateMetrics(p, g); }
    async testGenerateModelReport() { return this.generateModelReport(); }
    async testExecuteTraining(d: any) { return this.executeTraining(d); }
    async testExecutePrediction(d: any) { return this.executePrediction(d); }
}

describe('BaseModel Branch Coverage', () => {
    let mockSite: any;
    let model: TestModel;

    beforeEach(() => {
        mockSite = {
            models: {
                saveModel: vi.fn().mockResolvedValue('1.0'),
                loadModel: vi.fn(),
                getModelInfo: vi.fn().mockResolvedValue({ metadata: { config: {} } }),
                getModelPerformance: vi.fn().mockResolvedValue({ accuracy: 0.9, f1Score: 0.9 })
            }
        };
        model = new TestModel(mockSite);
    });

    it('should validate configuration branches', async () => {
        const invalidConfig = [
            { batchSize: 0 },
            { epochs: 0 },
            { learningRate: 0 },
            { learningRate: 2 },
            { inputSize: { width: 0, height: 10 } }
        ];

        for (const cfg of invalidConfig) {
            const m = new TestModel(mockSite, cfg);
            await expect(m.testValidateConfiguration()).rejects.toThrow();
        }
        
        const validM = new TestModel(mockSite, { batchSize: 1, epochs: 1, learningRate: 0.5, inputSize: { width: 1, height: 1 } });
        await expect(validM.testValidateConfiguration()).resolves.not.toThrow();
    });

    it('should validate training data branches', async () => {
        // Empty
        expect(await model.testValidateTrainingData([])).toBe(false);
        
        // Small dataset (<10) - logs warning but likely returns true or false depending on other checks
        // Implementation: if length < 10 logs warning.
        // It returns true unless consistency checks fail.
        const smallValid = Array(5).fill({ 
            id: '1', label: 'A', imageBuffer: Buffer.from('a'), metadata: {} 
        });
        // Needs labels >= 2
        
        // Single label
        const singleLabel = Array(10).fill({ id: '1', label: 'A', imageBuffer: Buffer.from('a') });
        await model.testValidateTrainingData(singleLabel); // Should return true but log valid warnings
        
        // Imbalance > 10
        const imbalance = [
            ...Array(20).fill({ label: 'A', imageBuffer: Buffer.from('a') }),
            { label: 'B', imageBuffer: Buffer.from('a') }
        ];
        await model.testValidateTrainingData(imbalance); // Should cover imbalance warning

         // Invalid buffer
        const invalidBuffer = [{ label: 'A', imageBuffer: 'string' }];
        expect(await model.testValidateTrainingData(invalidBuffer as any)).toBe(false);
    });

    it('should calculate metrics branches', async () => {
        // Length mismatch
        await expect(model.testCalculateMetrics([1], [1, 2])).rejects.toThrow();
        
        // All TP, FP, FN, TN
        // Predictions > 0.5 is Positive
        // TP: Pred 0.6, Actual 0.6
        // FP: Pred 0.6, Actual 0.4
        // FN: Pred 0.4, Actual 0.6
        // TN: Pred 0.4, Actual 0.4
        const preds = [0.6, 0.6, 0.4, 0.4];
        const truth = [0.6, 0.4, 0.6, 0.4];
        const metrics = await model.testCalculateMetrics(preds, truth);
        expect(metrics.accuracy).toBe(0.5); // 2/4
        expect(metrics.precision).toBe(0.5); // 1 TP / 2 Pred pos (1TP+1FP) = 0.5
        expect(metrics.recall).toBe(0.5); // 1 TP / 2 Actual pos (1TP+1FN) = 0.5
        
        // Division by zero safe checks (all zero)
        const metricsZero = await model.testCalculateMetrics([], []);
        expect(metricsZero.accuracy).toBeNaN(); // 0/0
        // wait, if length is 0 loop doesn't run.
        // accuracy = 0/0 = NaN
    });

    it('should generate model report recommendations', async () => {
        // Mock low accuracy
        mockSite.models.getModelPerformance.mockResolvedValue({ accuracy: 0.5, f1Score: 0.9 });
        let report = await model.testGenerateModelReport();
        expect(report.recommendations[0]).toContain('increasing training data');
        
        // Mock low f1
        mockSite.models.getModelPerformance.mockResolvedValue({ accuracy: 0.9, f1Score: 0.5 });
        report = await model.testGenerateModelReport();
        expect(report.recommendations[0]).toContain('class imbalance');
        
        // Mock high complexity low accuracy
        // Since analyzeModel is hardcoded to 'medium', we can't easily trigger the complexity branch 
        // unless we overwrite analyzeModel in TestModel or mock it if it was external.
        // analyzeModel is protected async on BaseModel.
        // We can override it in TestModel.
    });

    it('should cover executeTraining failure branches', async () => {
        // Validation data > 0
        // splitData defaults to 0.2 split.
        // We pass enough data to have validation set.
        const data = Array(10).fill({ id: '1', label: 'A', imageBuffer: Buffer.from('a'), metadata: {} });
        await model.testExecuteTraining(data);

        // handle error
        vi.spyOn(model, 'validateTrainingData').mockRejectedValue(new Error('fail'));
        const res = await model.testExecuteTraining([]);
        expect(res.success).toBe(false);
    });
    
    it('should cover executePrediction failure branches', async () => {
         vi.spyOn(model, 'predict').mockRejectedValue(new Error('fail'));
         const res = await model.testExecutePrediction(Buffer.from('a'));
         expect(res.isValid).toBe(false);
         expect(res.anomalies[0].description).toContain('prediction failed');
    });

    it('should validate getters/setters', () => {
        expect(model.isModelLoaded).toBe(false);
        expect(model.modelConfiguration).toBeDefined();
        model.setMetadata('k', 'v');
        expect(model.getMetadata('k')).toBe('v');
        expect(model.modelMetadata).toEqual(expect.objectContaining({k:'v'}));
        model.updateConfig({ epochs: 100 });
        expect(model.modelConfiguration.epochs).toBe(100);
    });
});
