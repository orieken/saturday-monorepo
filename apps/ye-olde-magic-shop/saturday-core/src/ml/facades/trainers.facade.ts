import { BaseSite } from '../../base/base-site';
import { ScreenshotTrainer } from '../trainers/screenshot-trainer';
import { ElementTrainer } from '../trainers/element-trainer';
import { PageTrainer } from '../trainers/page-trainer';
import { FlowTrainer } from '../trainers/flow-trainer';
import { TrainingOptions, TrainingResult, ElementTrainingOptions, PageTrainingConfig } from '../interfaces/ml-types';

export class TrainersFacade {
  private _screenshotTrainer?: ScreenshotTrainer;
  private _elementTrainer?: ElementTrainer;
  private _pageTrainer?: PageTrainer;
  private _flowTrainer?: FlowTrainer;

  constructor(private site: BaseSite) {}

  get screenshotTrainer(): ScreenshotTrainer {
    if (!this._screenshotTrainer) {
      this._screenshotTrainer = new ScreenshotTrainer(this.site);
    }
    return this._screenshotTrainer;
  }

  get elementTrainer(): ElementTrainer {
    if (!this._elementTrainer) {
      this._elementTrainer = new ElementTrainer(this.site);
    }
    return this._elementTrainer;
  }

  get pageTrainer(): PageTrainer {
    if (!this._pageTrainer) {
      this._pageTrainer = new PageTrainer(this.site);
    }
    return this._pageTrainer;
  }

  get flowTrainer(): FlowTrainer {
    if (!this._flowTrainer) {
      this._flowTrainer = new FlowTrainer(this.site);
    }
    return this._flowTrainer;
  }

  // Convenience methods with fluent interface
  async trainCurrentPage(label: string, options?: TrainingOptions): Promise<TrainingResult> {
    return await this.pageTrainer.trainCurrentPage(label, options);
  }

  async trainCurrentElement(selector: string, label: string, options?: ElementTrainingOptions): Promise<TrainingResult> {
    return await this.elementTrainer.trainElement(selector, label, options);
  }

  async trainFlow(flowName: string, params: any, label?: string): Promise<TrainingResult> {
    const flow = this.site.getFlow(flowName);
    return await this.flowTrainer.trainFlow(flow, params, label);
  }

  // Batch training operations
  async trainMultiplePages(pageConfigs: PageTrainingConfig[]): Promise<TrainingResult[]> {
    const results: TrainingResult[] = [];
    for (const config of pageConfigs) {
      const page = this.site.getPage(config.pageName);
      await page.visit();
      const result = await this.trainCurrentPage(config.label, config.options);
      results.push(result);
    }
    return results;
  }

  // High-level training workflows
  async trainCompleteUserJourney(journeyName: string, steps: Array<{
    type: 'page' | 'element' | 'flow';
    target: string;
    params?: any;
    label?: string;
  }>): Promise<TrainingResult[]> {
    const results: TrainingResult[] = [];

    for (const step of steps) {
      let result: TrainingResult;
      const stepLabel = step.label || `${journeyName}_${step.type}_${step.target}`;

      switch (step.type) {
        case 'page':
          const page = this.site.getPage(step.target);
          await page.visit();
          result = await this.trainCurrentPage(stepLabel);
          break;

        case 'element':
          result = await this.trainCurrentElement(step.target, stepLabel);
          break;

        case 'flow':
          result = await this.trainFlow(step.target, step.params || {}, stepLabel);
          break;

        default:
          throw new Error(`Unknown training step type: ${step.type}`);
      }

      results.push(result);
    }

    return results;
  }

  // Cross-browser training
  async trainAcrossBrowsers(label: string, browsers: string[], options?: TrainingOptions): Promise<TrainingResult[]> {
    const results: TrainingResult[] = [];

    for (const browser of browsers) {
      // This would require browser switching capability
      // Implementation depends on test runner setup
      const browserLabel = `${label}_${browser}`;
      const result = await this.trainCurrentPage(browserLabel, {
        ...options,
        metadata: { ...options?.metadata, browser }
      });
      results.push(result);
    }

    return results;
  }

  // Responsive training
  async trainResponsiveDesign(label: string, viewports: Array<{
    name: string;
    width: number;
    height: number;
  }>, options?: TrainingOptions): Promise<TrainingResult[]> {
    const results: TrainingResult[] = [];

    for (const viewport of viewports) {
      await this.site.pageObject.setViewportSize({
        width: viewport.width,
        height: viewport.height
      });

      const viewportLabel = `${label}_${viewport.name}`;
      const result = await this.trainCurrentPage(viewportLabel, {
        ...options,
        metadata: {
          ...options?.metadata,
          viewport: viewport.name,
          dimensions: { width: viewport.width, height: viewport.height }
        }
      });
      results.push(result);
    }

    return results;
  }

  // Training data management
  async getTrainingStats(label?: string): Promise<{
    totalSamples: number;
    modelCount: number;
    lastTrainingDate: Date;
    averageAccuracy: number;
  }> {
    // This would integrate with the training data manager
    return {
      totalSamples: 0,
      modelCount: 0,
      lastTrainingDate: new Date(),
      averageAccuracy: 0
    };
  }

  // Model versioning and comparison
  async compareModelVersions(modelLabel: string, versions: string[]): Promise<{
    version: string;
    accuracy: number;
    trainingTime: number;
    dataCount: number;
  }[]> {
    // Implementation would depend on model storage system
    return versions.map(version => ({
      version,
      accuracy: 0.95,
      trainingTime: 1000,
      dataCount: 100
    }));
  }
}
