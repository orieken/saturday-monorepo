import { BaseSite } from '../../base/base-site';
import { TrainingOptions, TrainingResult } from '../interfaces/ml-types';

export class PageTrainer {
  constructor(private site: BaseSite) {}
  async trainCurrentPage(label: string, options?: TrainingOptions): Promise<TrainingResult> {
     return {
         success: true,
         modelLabel: label,
         version: '0.0.1',
         metrics: { epochs: 1, trainingTime: 0 },
         timestamp: new Date(),
         dataCount: 0
     };
  }
}
