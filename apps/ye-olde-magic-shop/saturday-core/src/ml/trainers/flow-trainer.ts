import { BaseSite } from '../../base/base-site';
import { TrainingResult } from '../interfaces/ml-types';
import { BaseFlow } from '../../base/base-flow';

export class FlowTrainer {
  constructor(private site: BaseSite) {}
  async trainFlow(flow: BaseFlow, params: any, label?: string): Promise<TrainingResult> {
     return { success: true, modelLabel: label || 'flow', version: '0.0.1', metrics: { epochs: 1, trainingTime: 0 }, timestamp: new Date(), dataCount: 0 };
  }
}
