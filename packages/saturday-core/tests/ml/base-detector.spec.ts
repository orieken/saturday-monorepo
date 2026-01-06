
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { BaseDetector } from '../../src/ml/base/base-detector';
import { BaseSite } from '../../src/base/base-site';
import { ValidationResult } from '../../src/ml/interfaces/ml-types';
import { Page } from 'playwright';

class TestDetector extends BaseDetector {
  async detect(imageData: Buffer, modelLabel: string): Promise<ValidationResult> {
    if (modelLabel === 'fail') throw new Error('Simulated failure');
    return {
      isValid: true,
      confidence: 0.95,
      anomalies: [],
      score: 0.95,
      timestamp: new Date(),
      modelUsed: modelLabel
    };
  }

  async validateImage(imageData: Buffer, modelLabel: string): Promise<ValidationResult> {
      return this.detect(imageData, modelLabel);
  }
}

describe('BaseDetector', () => {
    let mockPage: any;
    let mockSite: any;
    let detector: TestDetector;

    beforeEach(() => {
        mockPage = {
            screenshot: vi.fn().mockResolvedValue(Buffer.from('fake-screenshot')),
            locator: vi.fn().mockReturnValue({
                waitFor: vi.fn(),
                boundingBox: vi.fn().mockResolvedValue({ x: 0, y: 0, width: 100, height: 100 })
            })
        };
        mockSite = {
            page: mockPage,
            models: { disconnect: vi.fn(), getModelInfo: vi.fn().mockResolvedValue({}) },
        };
        detector = new TestDetector(mockSite as BaseSite, 'test-type');
    });

    it('should validate current page', async () => {
        const result = await detector.validateCurrentPage('model-label');
        expect(result.isValid).toBe(true);
        expect(mockPage.screenshot).toHaveBeenCalled();
    });

    it('should validate element', async () => {
        const result = await detector.validateElement('.selector', 'model-label');
        expect(result.isValid).toBe(true);
        expect(mockPage.locator).toHaveBeenCalledWith('.selector');
        expect(mockPage.screenshot).toHaveBeenCalled(); // Should check clip?
    });

    it('should handle model errors', async () => {
        const result = await detector.executeDetection(
            () => Promise.resolve(Buffer.from('')),
            'fail'
        );
        expect(result.isValid).toBe(false);
        expect(result.anomalies[0].severity).toBe('critical');
    });

    it('should log detection result', async () => {
        const spy = vi.spyOn(console, 'log');
        await detector.validateCurrentPage('test');
        expect(spy).toHaveBeenCalled();
    });

    it('should calculate confidence', async () => {
        // Access protected method via any or subclass
        const score = (detector as any).calculateConfidence([0.9, 0.8, 0.7]);
        expect(score).toBeGreaterThan(0.8);
    });

    it('should filter anomalies by severity', async () => {
        const anomalies = [
            { severity: 'low', confidence: 0.9 },
            { severity: 'medium', confidence: 0.8 }
        ];
        const filtered = (detector as any).filterAnomaliesBySeverity(anomalies, { tolerance: 0.9 });
        expect(filtered).toHaveLength(1); // 'low' should be removed
    });
    
    it('should handle missing model', async () => {
         mockSite.models.getModelInfo.mockResolvedValueOnce(null);
         const result = await detector.executeDetection(
             () => Promise.resolve(Buffer.from('')),
             'missing'
         );
         expect(result.isValid).toBe(false);
         expect(result.anomalies[0].description).toContain('not found');
    });
});
