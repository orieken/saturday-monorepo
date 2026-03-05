import { BaseSite } from '../../base/base-site';
import { TrainingData, TrainingOptions, TrainingResult, ModelConfig } from '../types/ml-types';

export abstract class BaseTrainer {
  protected site: BaseSite;
  protected modelType: string;

  constructor(site: BaseSite, modelType: string) {
    this.site = site;
    this.modelType = modelType;
  }

  // Abstract methods that must be implemented by concrete trainers
  abstract train(data: TrainingData[], label: string, options?: TrainingOptions): Promise<TrainingResult>;
  abstract validateTrainingData(data: TrainingData[]): Promise<boolean>;
  abstract preprocessData(data: TrainingData[]): Promise<TrainingData[]>;

  // Common training workflow
  async executeTraining(
    dataCollector: () => Promise<TrainingData[]>,
    label: string,
    options?: TrainingOptions
  ): Promise<TrainingResult> {
    try {
      // Collect training data
      const rawData = await dataCollector();
      if (rawData.length === 0) {
        throw new Error('No training data collected');
      }

      // Validate data quality
      const isValid = await this.validateTrainingData(rawData);
      if (!isValid) {
        throw new Error('Training data validation failed');
      }

      // Preprocess data
      const processedData = await this.preprocessData(rawData);

      // Execute training
      const result = await this.train(processedData, label, options);

      // Save model if training was successful
      if (result.success) {
        await this.saveModel(label, result);
      }

      return result;

    } catch (error) {
      return {
        success: false,
        modelLabel: label,
        version: '0.0',
        metrics: {
          accuracy: 0,
          loss: 1.0,
          epochs: 0,
          trainingTime: 0
        },
        timestamp: new Date(),
        dataCount: 0
      };
    }
  }

  // Data collection helpers
  protected async captureCurrentState(label: string, metadata?: Record<string, any>): Promise<TrainingData> {
    const screenshot = await this.site.page.screenshot({ fullPage: true });

    return {
      id: this.generateId(),
      timestamp: new Date(),
      label,
      imageBuffer: screenshot,
      metadata: {
        pageUrl: this.site.page.url(),
        browserInfo: await this.getBrowserInfo(),
        viewportSize: await this.getViewportSize(),
        ...metadata
      }
    };
  }

  protected async captureElementState(
    selector: string,
    label: string,
    metadata?: Record<string, any>
  ): Promise<TrainingData> {
    const element = this.site.page.locator(selector);
    await element.waitFor({ state: 'visible' });

    const boundingBox = await element.boundingBox();
    if (!boundingBox) {
      throw new Error(`Element ${selector} not found or not visible`);
    }

    const screenshot = await this.site.page.screenshot({ clip: boundingBox });

    return {
      id: this.generateId(),
      timestamp: new Date(),
      label,
      imageBuffer: screenshot,
      metadata: {
        pageUrl: this.site.page.url(),
        elementSelector: selector,
        boundingBox,
        elementType: await element.evaluate(el => el.tagName.toLowerCase()),
        browserInfo: await this.getBrowserInfo(),
        viewportSize: await this.getViewportSize(),
        ...metadata
      }
    };
  }

  // Data augmentation
  protected async augmentData(data: TrainingData[]): Promise<TrainingData[]> {
    const augmentedData: TrainingData[] = [...data];

    for (const item of data) {
      // Add slight variations - this would use image processing libraries
      // For now, we'll just duplicate with different metadata
      const augmented: TrainingData = {
        ...item,
        id: this.generateId(),
        timestamp: new Date(),
        metadata: {
          ...item.metadata,
          augmented: true,
          augmentationType: 'baseline'
        }
      };
      augmentedData.push(augmented);
    }

    return augmentedData;
  }

  // Model management
  protected async saveModel(label: string, result: TrainingResult): Promise<void> {
    // Implementation would save to model storage
    console.log(`Saving model ${label} with version ${result.version}`);

    // This would integrate with the site's model facade
    await this.site.models.saveModel(label, {
      type: this.modelType,
      metrics: result.metrics,
      timestamp: result.timestamp,
      dataCount: result.dataCount
    }, {
      trainerType: this.constructor.name,
      trainingOptions: {}
    });
  }

  protected async loadModel(label: string, version?: string): Promise<any> {
    return await this.site.models.loadModel(label, version);
  }

  // Training data management
  protected async saveTrainingData(data: TrainingData[], label: string): Promise<void> {
    // Implementation would save training data for future use
    console.log(`Saving ${data.length} training samples for ${label}`);
  }

  protected async loadTrainingData(label: string): Promise<TrainingData[]> {
    // Implementation would load existing training data
    console.log(`Loading training data for ${label}`);
    return [];
  }

  // Model configuration
  protected getDefaultModelConfig(): ModelConfig {
    return {
      architecture: 'cnn',
      inputSize: { width: 224, height: 224 },
      batchSize: 16,
      epochs: 50,
      learningRate: 0.001,
      optimizerType: 'adam'
    };
  }

  protected mergeModelConfig(defaultConfig: ModelConfig, userConfig?: Partial<ModelConfig>): ModelConfig {
    return {
      ...defaultConfig,
      ...userConfig
    };
  }

  // Validation helpers
  protected validateDataConsistency(data: TrainingData[]): boolean {
    if (data.length === 0) return false;

    const firstItem = data[0];
    return data.every(item => {
      // Check for required fields
      if (!item.id || !item.label || !item.imageBuffer || !item.metadata) {
        return false;
      }

      // Check metadata consistency
      const requiredMetadata = ['pageUrl', 'browserInfo', 'viewportSize'];
      return requiredMetadata.every(field =>
        item.metadata.hasOwnProperty(field)
      );
    });
  }

  protected validateImageData(data: TrainingData[]): boolean {
    return data.every(item => {
      // Check if image buffer is valid
      return Buffer.isBuffer(item.imageBuffer) && item.imageBuffer.length > 0;
    });
  }

  // Utility methods
  protected generateId(): string {
    return Date.now().toString(36) + Math.random().toString(36).substr(2);
  }

  protected async getBrowserInfo() {
    return await this.site.page.evaluate(() => ({
      userAgent: navigator.userAgent,
      platform: navigator.platform,
      cookieEnabled: navigator.cookieEnabled
    }));
  }

  protected async getViewportSize() {
    return await this.site.page.evaluate(() => ({
      width: window.innerWidth,
      height: window.innerHeight
    }));
  }

  // Progress tracking
  protected logProgress(message: string, progress?: number): void {
    const timestamp = new Date().toISOString();
    const progressStr = progress !== undefined ? ` (${Math.round(progress * 100)}%)` : '';
    console.log(`[${timestamp}] ${this.constructor.name}: ${message}${progressStr}`);
  }

  // Error handling
  protected handleTrainingError(error: Error, context: string): never {
    const timestamp = new Date().toISOString();
    console.error(`[${timestamp}] ${this.constructor.name} error in ${context}:`, error);
    throw new Error(`Training failed in ${context}: ${error.message}`);
  }
}
