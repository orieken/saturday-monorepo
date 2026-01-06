
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { BaseModel } from '../../src/ml/base/base-model';
import { BaseSite } from '../../src/base/base-site';
import { TrainingData, TrainingResult, ValidationResult, ModelPerformanceMetrics } from '../../src/ml/interfaces/ml-types';

class TestModel extends BaseModel {
  constructor(site: BaseSite, label: string) {
      super(site, 'test', label);
  }

  async train(data: TrainingData[], options: any): Promise<TrainingResult> {
      return {
          success: true,
          metrics: { accuracy: 1.0, trainingTime: 100, epochs: 1 },
          modelLabel: this.modelLabel,
          version: this.version,
          timestamp: new Date(),
          dataCount: data.length
      };
  }

  async predict(imageData: Buffer, options: any): Promise<any> {
      return { confidence: 0.95, anomalies: [] };
  }

  async evaluate(testData: TrainingData[]): Promise<ModelPerformanceMetrics> {
      return { accuracy: 0.9, precision: 0.9, recall: 0.9, f1Score: 0.9, validationLoss: 0.1, trainingTime: 0 };
  }

  async save(path?: string): Promise<string> {
      return '1.0';
  }

  async load(path: string, version?: string): Promise<void> {
      this.isLoaded = true;
  }
}

describe('BaseModel', () => {
    let mockSite: any;
    let model: TestModel;

    beforeEach(() => {
        mockSite = {
            models: { 
                saveModel: vi.fn().mockResolvedValue('1.0'), 
                loadModel: vi.fn(), 
                getModelInfo: vi.fn().mockResolvedValue({}),
                getModelPerformance: vi.fn().mockResolvedValue({ accuracy: 0.9 })
            }
        };
        model = new TestModel(mockSite as BaseSite, 'test-model');
    });

    it('should initialize', async () => {
        await model.initialize();
        expect(model.isModelLoaded).toBe(true);
    });

    it('should dispose', async () => {
        await model.dispose();
        expect(model.isModelLoaded).toBe(false);
    });

    it('should execute training', async () => {
        const dummyData: TrainingData[] = Array(20).fill({
            imageBuffer: Buffer.from('data'),
            label: 'test',
            id: '1',
            metadata: {}
        });

        const result = await model.executeTraining(dummyData);
        expect(result.success).toBe(true);
    });

    it('should fail training with insufficient data', async () => {
        const result = await model.executeTraining([]);
        expect(result.success).toBe(false);
    });

    it('should execute prediction', async () => {
        await model.initialize();
        const result = await model.executePrediction(Buffer.from('image'));
        expect(result.isValid).toBe(true);
    });

    it('should validate training data', async () => {
         const validation = await (model as any).validateTrainingData([]);
         expect(validation).toBe(false);
    });

    it('should split data', async () => {
         const data = Array(10).fill({ id: '1' });
         const split = await (model as any).splitData(data, { metadata: { validationSplit: 0.2 } });
         expect(split.trainData.length).toBe(8);
         expect(split.validationData.length).toBe(2);
    });

    it('should augment data', async () => {
        const data = [{ imageBuffer: Buffer.from('img'), id: '1', metadata: {}, label: 'test' }];
        const augmented = await (model as any).augmentData(data);
        expect(augmented.length).toBe(4); // 4 augmentations in loop
    });

    it('should update config and metadata', () => {
        model.updateConfig({ batchSize: 32 });
        expect(model.modelConfiguration.batchSize).toBe(32);

        model.setMetadata('key', 'value');
        expect(model.getMetadata('key')).toBe('value');
    });

    it('should calculate metrics', async () => {
        const metrics = await model.calculateMetrics([0.9, 0.2], [1.0, 0.0]); // Matches
        expect(metrics.accuracy).toBe(1.0);
    });
    
    it('should generate report', async () => {
        const report = await model.generateModelReport();
        expect(report.summary).toContain('test-model');
    });
});
