import { BaseSite } from '../../base/base-site';
import { ComparisonOptions, ComparisonResult } from '../interfaces/ml-types';

export class RegressionDetector {
  constructor(private site: BaseSite) {}
  async compareWithBaseline(baselineLabel: string, options?: ComparisonOptions): Promise<ComparisonResult> {
    return { similarity: 1.0, differences: [], overallMatch: true, timestamp: new Date() };
  }
  async compareScreenshots(buffer1: Buffer, buffer2: Buffer): Promise<ComparisonResult> {
    return { similarity: 1.0, differences: [], overallMatch: true, timestamp: new Date() };
  }
}
