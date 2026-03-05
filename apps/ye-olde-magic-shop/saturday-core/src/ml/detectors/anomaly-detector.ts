import { BaseSite } from '../../base/base-site';
import { ValidationOptions, ValidationResult } from '../interfaces/ml-types';

export class AnomalyDetector {
  constructor(private site: BaseSite) {}
  async validateCurrentPage(modelLabel: string, options?: ValidationOptions): Promise<ValidationResult> {
    return { isValid: true, confidence: 1.0, anomalies: [], score: 1.0, timestamp: new Date(), modelUsed: modelLabel };
  }
  async validateImage(image: Buffer, modelLabel: string, options?: any): Promise<ValidationResult> {
    return { isValid: true, confidence: 1.0, anomalies: [], score: 1.0, timestamp: new Date(), modelUsed: modelLabel };
  }
}
