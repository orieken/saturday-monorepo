
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { BaseTrainer } from '../../src/ml/base/base-trainer';
import { BaseSite } from '../../src/base/base-site';
import { TrainingData, TrainingResult } from '../../src/ml/interfaces/ml-types';

class TestTrainer extends BaseTrainer {
  constructor(site: BaseSite) {
      super(site, 'test-trainer');
  }

  async train(data: TrainingData[], label: string, options: any): Promise<TrainingResult> {
      return {
          success: true,
          metrics: { accuracy: 1.0, trainingTime: 100, epochs: 1 },
          modelLabel: label,
          version: '1.0',
          timestamp: new Date(),
          dataCount: data.length
      };
  }

  async validateTrainingData(data: TrainingData[]): Promise<boolean> {
      return data.length > 0;
  }

  async preprocessData(data: TrainingData[]): Promise<TrainingData[]> {
      return data;
  }
}

describe('BaseTrainer', () => {
    let mockSite: any;
    let trainer: TestTrainer;
    let mockPage: any;

    beforeEach(() => {
        mockPage = {
            screenshot: vi.fn().mockResolvedValue(Buffer.from('img')),
            url: vi.fn().mockReturnValue('http://localhost'),
            evaluate: vi.fn().mockReturnValue({}),
            locator: vi.fn().mockReturnValue({
                waitFor: vi.fn(),
                boundingBox: vi.fn().mockResolvedValue({ x:0, y:0, width:100, height:100 }),
                evaluate: vi.fn().mockReturnValue('div')
            })
        };
        mockSite = {
            page: mockPage,
            models: { saveModel: vi.fn().mockResolvedValue('1.0') }
        };
        trainer = new TestTrainer(mockSite as BaseSite);
        
        vi.stubGlobal('navigator', { userAgent:'ua', platform:'plat', cookieEnabled:true });
        vi.stubGlobal('window', { innerWidth:1024, innerHeight:768 });
    });

    it('should execute training workflow', async () => {
        const dataCollector = vi.fn().mockResolvedValue([{ id:'1', label:'l', imageBuffer:Buffer.from('b'), metadata:{} }]);
        const result = await trainer.executeTraining(dataCollector, 'test-label');
        expect(result.success).toBe(true);
        expect(mockSite.models.saveModel).toHaveBeenCalled();
    });

    it('should handle training errors', async () => {
        const dataCollector = vi.fn().mockResolvedValue([]);
        const result = await trainer.executeTraining(dataCollector, 'test-label');
        expect(result.success).toBe(false); // No data throws error, caught same way
    });

    it('should capture current state', async () => {
        const data = await (trainer as any).captureCurrentState('label');
        expect(data.imageBuffer).toBeDefined();
        expect(data.metadata.pageUrl).toBe('http://localhost');
    });

    it('should capture element state', async () => {
        const data = await (trainer as any).captureElementState('.sel', 'label');
        expect(data.metadata.elementSelector).toBe('.sel');
    });

    it('should augment data', async () => {
        const d = [{ id:'1', label:'l', imageBuffer:Buffer.from('b'), metadata:{} }];
        const augmented = await (trainer as any).augmentData(d);
        expect(augmented.length).toBe(2);
    });

    it('should validate data consistency', () => {
         const good = [{ id:'1', label:'l', imageBuffer:Buffer.from('b'), metadata:{ pageUrl:'u', browserInfo:{}, viewportSize:{} } }];
         expect((trainer as any).validateDataConsistency(good)).toBe(true);
         
         const bad = [{ id:'1' }];
         expect((trainer as any).validateDataConsistency(bad)).toBe(false);
    });
});
