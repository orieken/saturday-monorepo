import { BaseSite } from '../../base/base-site';
import {
  ModelConfig,
  TrainingData,
  TrainingResult,
  ValidationResult,
  ModelInfo,
  ModelPerformanceMetrics,
  TrainingOptions
} from '../interfaces/ml-types';

export abstract class BaseModel {
  protected site: BaseSite;
  protected modelType: string;
  protected modelLabel: string;
  protected version: string;
  protected config: ModelConfig;
  protected isLoaded: boolean = false;
  protected modelData: any;
  protected metadata: Record<string, any> = {};

  constructor(site: BaseSite, modelType: string, modelLabel: string, config?: Partial<ModelConfig>) {
    this.site = site;
    this.modelType = modelType;
    this.modelLabel = modelLabel;
    this.version = '1.0';
    this.config = this.mergeWithDefaults(config);
  }

  // Abstract methods that must be implemented by concrete models
  abstract train(data: TrainingData[], options?: TrainingOptions): Promise<TrainingResult>;
  abstract predict(imageData: Buffer, options?: any): Promise<any>;
  abstract evaluate(testData: TrainingData[]): Promise<ModelPerformanceMetrics>;
  abstract save(path?: string): Promise<string>;
  abstract load(path: string, version?: string): Promise<void>;

  // Model lifecycle management
  async initialize(): Promise<void> {
    this.logProgress('Initializing model', 0);

    try {
      await this.setupModel();
      await this.validateConfiguration();
      this.isLoaded = true;
      this.logProgress('Model initialized successfully', 1);
    } catch (error) {
      this.handleError('initialization', error as Error);
    }
  }

  async dispose(): Promise<void> {
    try {
      await this.cleanup();
      this.isLoaded = false;
      this.modelData = null;
      this.logProgress('Model disposed successfully');
    } catch (error) {
      this.handleError('disposal', error as Error);
    }
  }

  // Training workflow
  async executeTraining(
    data: TrainingData[],
    options?: TrainingOptions
  ): Promise<TrainingResult> {
    if (!this.isLoaded) {
      await this.initialize();
    }

    try {
      this.logProgress('Starting training process', 0);

      // Validate training data
      const isValid = await this.validateTrainingData(data);
      if (!isValid) {
        throw new Error('Training data validation failed');
      }

      // Preprocess data
      const processedData = await this.preprocessTrainingData(data, options);
      this.logProgress('Data preprocessing complete', 0.2);

      // Split data for training and validation
      const { trainData, validationData } = await this.splitData(processedData, options);
      this.logProgress('Data split complete', 0.3);

      // Execute actual training
      const result = await this.train(trainData, options);
      this.logProgress('Training complete', 0.8);

      // Validate model performance
      if (validationData.length > 0) {
        const performance = await this.evaluate(validationData);
        result.metrics = { ...result.metrics, ...performance };
      }
      this.logProgress('Validation complete', 0.9);

      // Save model if training was successful
      if (result.success) {
        const savedPath = await this.save();
        this.logProgress(`Model saved to ${savedPath}`, 1);
      }

      return result;

    } catch (error) {
      return this.createFailedTrainingResult(error as Error);
    }
  }

  // Prediction workflow
  async executePrediction(
    imageData: Buffer,
    options?: any
  ): Promise<ValidationResult> {
    if (!this.isLoaded) {
      await this.load(this.modelLabel);
    }

    try {
      // Preprocess input image
      const processedImage = await this.preprocessImage(imageData);

      // Execute prediction
      const predictions = await this.predict(processedImage, options);

      // Postprocess predictions to validation result
      const result = await this.postprocessPredictions(predictions, options);

      return result;

    } catch (error) {
      return this.createFailedValidationResult(error as Error);
    }
  }

  // Data preprocessing
  protected async preprocessTrainingData(
    data: TrainingData[],
    options?: TrainingOptions
  ): Promise<TrainingData[]> {
    const processedData: TrainingData[] = [];

    for (let i = 0; i < data.length; i++) {
      const item = data[i];

      // Resize image to model input size
      const resizedImage = await this.resizeImage(item.imageBuffer, this.config.inputSize);

      // Normalize image data
      const normalizedImage = await this.normalizeImage(resizedImage);

      const processedItem: TrainingData = {
        ...item,
        imageBuffer: normalizedImage,
        metadata: {
          ...item.metadata,
          preprocessed: true,
          originalSize: await this.getImageDimensions(item.imageBuffer),
          processedSize: this.config.inputSize
        }
      };

      processedData.push(processedItem);

      // Report progress
      if (i % 10 === 0) {
        this.logProgress('Preprocessing data', i / data.length);
      }
    }

    // Apply data augmentation if enabled
    if (options?.metadata?.enableAugmentation) {
      const augmentedData = await this.augmentData(processedData);
      processedData.push(...augmentedData);
    }

    return processedData;
  }

  protected async preprocessImage(imageData: Buffer): Promise<Buffer> {
    // Resize to model input size
    const resized = await this.resizeImage(imageData, this.config.inputSize);

    // Normalize
    const normalized = await this.normalizeImage(resized);

    return normalized;
  }

  // Data validation
  protected async validateTrainingData(data: TrainingData[]): Promise<boolean> {
    if (data.length === 0) {
      this.logProgress('Error: No training data provided');
      return false;
    }

    if (data.length < 10) {
      this.logProgress('Warning: Very small training dataset');
    }

    // Check data consistency
    const labelCounts = data.reduce((counts, item) => {
      counts[item.label] = (counts[item.label] || 0) + 1;
      return counts;
    }, {} as Record<string, number>);

    const labels = Object.keys(labelCounts);
    if (labels.length < 2) {
      this.logProgress('Warning: Training data contains only one label');
    }

    // Check for class imbalance
    const maxCount = Math.max(...Object.values(labelCounts));
    const minCount = Math.min(...Object.values(labelCounts));
    const imbalanceRatio = maxCount / minCount;

    if (imbalanceRatio > 10) {
      this.logProgress(`Warning: High class imbalance detected (ratio: ${imbalanceRatio.toFixed(2)})`);
    }

    // Validate image data
    for (const item of data) {
      if (!Buffer.isBuffer(item.imageBuffer) || item.imageBuffer.length === 0) {
        this.logProgress(`Error: Invalid image data for item ${item.id}`);
        return false;
      }
    }

    return true;
  }

  // Data splitting
  protected async splitData(
    data: TrainingData[],
    options?: TrainingOptions
  ): Promise<{
    trainData: TrainingData[];
    validationData: TrainingData[];
  }> {
    const validationSplit = options?.metadata?.validationSplit || this.config.validationSplit || 0.2;

    // Shuffle data
    const shuffled = [...data].sort(() => Math.random() - 0.5);

    const splitIndex = Math.floor(shuffled.length * (1 - validationSplit));

    return {
      trainData: shuffled.slice(0, splitIndex),
      validationData: shuffled.slice(splitIndex)
    };
  }

  // Data augmentation
  protected async augmentData(data: TrainingData[]): Promise<TrainingData[]> {
    const augmentedData: TrainingData[] = [];

    for (const item of data) {
      // Apply various augmentations
      const augmentations = [
        { name: 'brightness', params: { factor: 1.1 } },
        { name: 'contrast', params: { factor: 1.1 } },
        { name: 'rotation', params: { degrees: 2 } },
        { name: 'noise', params: { intensity: 0.05 } }
      ];

      for (const aug of augmentations) {
        const augmentedImage = await this.applyAugmentation(item.imageBuffer, aug);

        const augmentedItem: TrainingData = {
          ...item,
          id: `${item.id}_${aug.name}`,
          imageBuffer: augmentedImage,
          metadata: {
            ...item.metadata,
            augmented: true,
            augmentationType: aug.name,
            augmentationParams: aug.params
          }
        };

        augmentedData.push(augmentedItem);
      }
    }

    return augmentedData;
  }

  protected async applyAugmentation(
    imageData: Buffer,
    augmentation: { name: string; params: any }
  ): Promise<Buffer> {
    // Placeholder for actual image augmentation
    // Would use libraries like Sharp, Canvas, or TensorFlow.js image ops
    switch (augmentation.name) {
      case 'brightness':
      case 'contrast':
      case 'rotation':
      case 'noise':
        // Apply specific augmentation
        return imageData; // Placeholder
      default:
        return imageData;
    }
  }

  // Image processing utilities
  protected async resizeImage(
    imageData: Buffer,
    targetSize: { width: number; height: number }
  ): Promise<Buffer> {
    // Placeholder for image resizing
    // Would use Sharp or Canvas API
    return imageData;
  }

  protected async normalizeImage(imageData: Buffer): Promise<Buffer> {
    // Placeholder for image normalization
    // Would convert to normalized tensor values [0-1] or [-1,1]
    return imageData;
  }

  protected async getImageDimensions(imageData: Buffer): Promise<{ width: number; height: number }> {
    // Placeholder for getting image dimensions
    // Would use Sharp or image-size library
    return { width: 0, height: 0 };
  }

  // Model configuration
  protected mergeWithDefaults(userConfig?: Partial<ModelConfig>): ModelConfig {
    const defaults: ModelConfig = {
      architecture: 'cnn',
      inputSize: { width: 224, height: 224 },
      batchSize: 16,
      epochs: 50,
      learningRate: 0.001,
      optimizerType: 'adam',
      validationSplit: 0.2
    };

    return { ...defaults, ...userConfig };
  }

  protected async validateConfiguration(): Promise<void> {
    if (this.config.batchSize <= 0) {
      throw new Error('Batch size must be positive');
    }

    if (this.config.epochs <= 0) {
      throw new Error('Epochs must be positive');
    }

    if (this.config.learningRate <= 0 || this.config.learningRate > 1) {
      throw new Error('Learning rate must be between 0 and 1');
    }

    if (this.config.inputSize.width <= 0 || this.config.inputSize.height <= 0) {
      throw new Error('Input size dimensions must be positive');
    }
  }

  // Model persistence
  async saveModel(metadata?: Record<string, any>): Promise<string> {
    const modelInfo: ModelInfo = {
      label: this.modelLabel,
      version: this.version,
      type: this.modelType,
      createdAt: new Date(),
      lastUsed: new Date(),
      dataCount: 0, // Would be set during training
      metadata: {
        config: this.config,
        ...this.metadata,
        ...metadata
      }
    };

    // Save through the site's model facade
    const savedVersion = await this.site.models.saveModel(this.modelLabel, this.modelData, modelInfo.metadata);
    this.version = savedVersion;

    return `${this.modelLabel}@${this.version}`;
  }

  async loadModel(label: string, version?: string): Promise<void> {
    try {
      this.modelData = await this.site.models.loadModel(label, version);
      this.modelLabel = label;
      this.isLoaded = true;

      // Load metadata
      const modelInfo = await this.site.models.getModelInfo(label, version);
      if (modelInfo?.metadata) {
        this.metadata = modelInfo.metadata;
        if (modelInfo.metadata.config) {
          this.config = { ...this.config, ...modelInfo.metadata.config };
        }
      }

    } catch (error) {
      throw new Error(`Failed to load model ${label}: ${(error as Error).message}`);
    }
  }

  // Performance metrics
  async calculateMetrics(
    predictions: any[],
    groundTruth: any[]
  ): Promise<ModelPerformanceMetrics> {
    if (predictions.length !== groundTruth.length) {
      throw new Error('Predictions and ground truth must have same length');
    }

    // Calculate basic metrics
    let truePositives = 0;
    let falsePositives = 0;
    let falseNegatives = 0;
    let trueNegatives = 0;

    for (let i = 0; i < predictions.length; i++) {
      const predicted = predictions[i] > 0.5; // Binary threshold
      const actual = groundTruth[i] > 0.5;

      if (predicted && actual) truePositives++;
      else if (predicted && !actual) falsePositives++;
      else if (!predicted && actual) falseNegatives++;
      else trueNegatives++;
    }

    const accuracy = (truePositives + trueNegatives) / predictions.length;
    const precision = truePositives / (truePositives + falsePositives) || 0;
    const recall = truePositives / (truePositives + falseNegatives) || 0;
    const f1Score = 2 * (precision * recall) / (precision + recall) || 0;

    return {
      accuracy,
      precision,
      recall,
      f1Score,
      validationLoss: 0, // Would be calculated during training
      trainingTime: 0    // Would be tracked during training
    };
  }

  // Prediction postprocessing
  protected async postprocessPredictions(
    predictions: any,
    options?: any
  ): Promise<ValidationResult> {
    // Convert model predictions to ValidationResult format
    const confidence = Array.isArray(predictions) ?
      Math.max(...predictions) :
      predictions.confidence || 0;

    const isValid = confidence > (options?.threshold || 0.8);

    return {
      isValid,
      confidence,
      anomalies: predictions.anomalies || [],
      score: confidence,
      timestamp: new Date(),
      modelUsed: this.modelLabel
    };
  }

  // Model analysis and debugging
  async analyzeModel(): Promise<{
    layerCount: number;
    parameterCount: number;
    modelSize: string;
    complexity: 'low' | 'medium' | 'high';
  }> {
    // Placeholder for model analysis
    // Would inspect actual model architecture
    return {
      layerCount: 0,
      parameterCount: 0,
      modelSize: '0 MB',
      complexity: 'medium'
    };
  }

  async generateModelReport(): Promise<{
    summary: string;
    performance: ModelPerformanceMetrics;
    config: ModelConfig;
    recommendations: string[];
  }> {
    const performance = await this.getStoredPerformance();
    const analysis = await this.analyzeModel();

    const recommendations: string[] = [];

    if (performance.accuracy < 0.8) {
      recommendations.push('Consider increasing training data or adjusting hyperparameters');
    }

    if (performance.f1Score < 0.7) {
      recommendations.push('Model may be suffering from class imbalance');
    }

    if (analysis.complexity === 'high' && performance.accuracy < 0.9) {
      recommendations.push('Consider simplifying model architecture to reduce overfitting');
    }

    return {
      summary: `Model ${this.modelLabel} v${this.version} - ${this.modelType}`,
      performance,
      config: this.config,
      recommendations
    };
  }

  // Utility methods
  protected async setupModel(): Promise<void> {
    // Override in concrete implementations for model-specific setup
  }

  protected async cleanup(): Promise<void> {
    // Override in concrete implementations for cleanup
  }

  protected async getStoredPerformance(): Promise<ModelPerformanceMetrics> {
    // Get performance metrics from storage
    return await this.site.models.getModelPerformance(this.modelLabel, this.version) || {
      accuracy: 0,
      precision: 0,
      recall: 0,
      f1Score: 0,
      validationLoss: 0,
      trainingTime: 0
    };
  }

  protected createFailedTrainingResult(error: Error): TrainingResult {
    return {
      success: false,
      modelLabel: this.modelLabel,
      version: this.version,
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

  protected createFailedValidationResult(error: Error): ValidationResult {
    return {
      isValid: false,
      confidence: 0,
      anomalies: [{
        type: 'visual',
        severity: 'critical',
        description: `Model prediction failed: ${error.message}`,
        location: { x: 0, y: 0, width: 0, height: 0 },
        confidence: 1.0
      }],
      score: 0,
      timestamp: new Date(),
      modelUsed: this.modelLabel
    };
  }

  // Logging and monitoring
  protected logProgress(message: string, progress?: number): void {
    const timestamp = new Date().toISOString();
    const progressStr = progress !== undefined ? ` (${Math.round(progress * 100)}%)` : '';
    console.log(`[${timestamp}] ${this.constructor.name} [${this.modelLabel}]: ${message}${progressStr}`);
  }

  protected handleError(context: string, error: Error): never {
    const timestamp = new Date().toISOString();
    console.error(`[${timestamp}] ${this.constructor.name} [${this.modelLabel}] error in ${context}:`, error);
    throw new Error(`${context} failed: ${error.message}`);
  }

  // Getters and setters
  get isModelLoaded(): boolean {
    return this.isLoaded;
  }

  get modelConfiguration(): ModelConfig {
    return { ...this.config };
  }

  get modelMetadata(): Record<string, any> {
    return { ...this.metadata };
  }

  setMetadata(key: string, value: any): void {
    this.metadata[key] = value;
  }

  getMetadata(key: string): any {
    return this.metadata[key];
  }

  updateConfig(updates: Partial<ModelConfig>): void {
    this.config = { ...this.config, ...updates };
  }
}
