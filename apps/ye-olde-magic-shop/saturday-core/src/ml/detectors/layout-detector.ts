import { BaseSite } from '../../base/base-site';
import { LayoutDetectionOptions, LayoutAnomalyResult } from '../interfaces/ml-types';

export class LayoutDetector {
  constructor(private site: BaseSite) {}
  async detectAnomalies(options?: LayoutDetectionOptions): Promise<LayoutAnomalyResult[]> {
    return [];
  }
}
