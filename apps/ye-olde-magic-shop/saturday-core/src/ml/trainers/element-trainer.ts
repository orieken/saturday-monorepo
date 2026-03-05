import { BaseSite } from '../../base/base-site';
import { ElementTrainingOptions, TrainingResult } from '../interfaces/ml-types';

export class ElementTrainer {
  constructor(private site: BaseSite) {}
  async trainElement(selector: string, label: string, options?: ElementTrainingOptions): Promise<TrainingResult> {
     return { success: true, modelLabel: label, version: '0.0.1', metrics: { epochs: 1, trainingTime: 0 }, timestamp: new Date(), dataCount: 0 };
  }
}
