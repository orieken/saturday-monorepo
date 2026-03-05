
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { BaseTrainer } from '../../src/ml/base/base-trainer';
import { BaseSite } from '../../src/base/base-site';

class TestTrainer extends BaseTrainer {
    constructor(site: BaseSite) { super(site, 'test'); }
    async train(data: any) { return { success: true, modelLabel: 'l', version: '1', metrics: {}, timestamp: new Date(), dataCount: 1 } as any; }
    async validateTrainingData(data: any) { return true; }
    async preprocessData(data: any) { return data; }
    
    // Check validation override
    async testValidateDataConsistency(d: any[]) { return this.validateDataConsistency(d); }
    async testValidateImageData(d: any[]) { return this.validateImageData(d); }
    async testCaptureElementState(sel: string, lbl: string) { return this.captureElementState(sel, lbl); }
}

describe('BaseTrainer Branch Coverage', () => {
    let mockSite: any;
    let trainer: TestTrainer;

    beforeEach(() => {
        mockSite = {
            page: {
                screenshot: vi.fn(),
                locator: vi.fn().mockReturnValue({
                    waitFor: vi.fn(),
                    boundingBox: vi.fn().mockResolvedValue({ x:0, y:0, width:10, height:10 }),
                    evaluate: vi.fn()
                }),
                url: vi.fn(),
                evaluate: vi.fn().mockResolvedValue({})
            },
            models: { saveModel: vi.fn() }
        };
        trainer = new TestTrainer(mockSite);
    });

    it('should handle execution errors in executeTraining', async () => {
        // No data
        const res = await trainer.executeTraining(async () => [], 'lbl');
        expect(res.success).toBe(false);

        // Validation failed
        const validateSpy = vi.spyOn(trainer, 'validateTrainingData').mockResolvedValue(false);
        const res2 = await trainer.executeTraining(async () => [{} as any], 'lbl');
        expect(res2.success).toBe(false);
    });

    it('should validate data consistency branches', async () => {
        expect(await trainer.testValidateDataConsistency([])).toBe(false);
        
        // Missing fields
        const invalidData = [{ id: '1' } as any];
        expect(await trainer.testValidateDataConsistency(invalidData)).toBe(false);

        // Missing metadata
        const invalidMeta = [{ id: '1', label: 'l', imageBuffer: Buffer.from('a'), metadata: {} } as any];
        expect(await trainer.testValidateDataConsistency(invalidMeta)).toBe(false); // missing pageUrl etc
        
         // Valid
         const valid = [{ 
             id: '1', label: 'l', imageBuffer: Buffer.from('a'), 
             metadata: { pageUrl: 'u', browserInfo: {}, viewportSize: {} } 
         } as any];
         expect(await trainer.testValidateDataConsistency(valid)).toBe(true);
    });

    it('should validate image data branches', async () => {
         const invalid = [{ imageBuffer: 'string' } as any];
         expect(await trainer.testValidateImageData(invalid)).toBe(false);
         
         const valid = [{ imageBuffer: Buffer.from('a') } as any];
         expect(await trainer.testValidateImageData(valid)).toBe(true);
    });

    it('should handle element capture failure', async () => {
        mockSite.page.locator().boundingBox.mockResolvedValue(null);
        await expect(trainer.testCaptureElementState('sel', 'lbl')).rejects.toThrow('not visible');
    });
});
