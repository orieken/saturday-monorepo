import { BaseSite } from '../../base/base-site';
import { ValidationOptions, ValidationResult, Anomaly } from '../types/ml-types';

export abstract class BaseDetector {
  protected site: BaseSite;
  protected detectorType: string;

  constructor(site: BaseSite, detectorType: string) {
    this.site = site;
    this.detectorType = detectorType;
  }

  // Abstract methods that must be implemented by concrete detectors
  abstract detect(imageData: Buffer, modelLabel: string, options?: ValidationOptions): Promise<ValidationResult>;
  abstract validateImage(imageData: Buffer, modelLabel: string, options?: ValidationOptions): Promise<ValidationResult>;

  // Common detection workflow
  async executeDetection(
    imageCapture: () => Promise<Buffer>,
    modelLabel: string,
    options?: ValidationOptions
  ): Promise<ValidationResult> {
    try {
      // Capture current state
      const imageData = await imageCapture();

      // Validate model exists
      const modelExists = await this.validateModelExists(modelLabel);
      if (!modelExists) {
        throw new Error(`Model ${modelLabel} not found`);
      }

      // Preprocess image if needed
      const processedImage = await this.preprocessImage(imageData, options);

      // Execute detection
      const result = await this.detect(processedImage, modelLabel, options);

      // Post-process results
      const finalResult = await this.postprocessResults(result, options);

      // Log detection result
      this.logDetectionResult(modelLabel, finalResult);

      return finalResult;

    } catch (error) {
      return this.createErrorResult(modelLabel, error as Error);
    }
  }

  // Current page detection
  async validateCurrentPage(modelLabel: string, options?: ValidationOptions): Promise<ValidationResult> {
    return await this.executeDetection(
      () => this.captureCurrentPage(options),
      modelLabel,
      options
    );
  }

  // Element-specific detection
  async validateElement(selector: string, modelLabel: string, options?: ValidationOptions): Promise<ValidationResult> {
    return await this.executeDetection(
      () => this.captureElement(selector, options),
      modelLabel,
      options
    );
  }

  // Image capture methods
  protected async captureCurrentPage(options?: ValidationOptions): Promise<Buffer> {
    const screenshotOptions: any = {
      fullPage: options?.fullPage ?? true
    };

    if (options?.clip) {
      screenshotOptions.clip = options.clip;
    }

    return await this.site.page.screenshot(screenshotOptions);
  }

  protected async captureElement(selector: string, options?: ValidationOptions): Promise<Buffer> {
    const element = this.site.page.locator(selector);
    await element.waitFor({ state: 'visible' });

    const boundingBox = await element.boundingBox();
    if (!boundingBox) {
      throw new Error(`Element ${selector} not found or not visible`);
    }

    return await this.site.page.screenshot({
      clip: options?.clip || boundingBox
    });
  }

  // Image preprocessing
  protected async preprocessImage(imageData: Buffer, _options?: ValidationOptions): Promise<Buffer> {
    // Basic preprocessing - resize, normalize, etc.
    // In a real implementation, this would use image processing libraries
    return imageData;
  }

  // Result processing
  protected async postprocessResults(result: ValidationResult, options?: ValidationOptions): Promise<ValidationResult> {
    // Apply confidence thresholds
    if (options?.threshold && result.confidence < options.threshold) {
      result.isValid = false;
      result.anomalies.push({
        type: 'visual',
        severity: 'medium',
        description: `Confidence ${result.confidence} below threshold ${options.threshold}`,
        location: { x: 0, y: 0, width: 0, height: 0 },
        confidence: result.confidence
      });
    }

    // Filter anomalies by severity if needed
    result.anomalies = this.filterAnomaliesBySeverity(result.anomalies, options);

    // Recalculate overall score
    result.score = this.calculateOverallScore(result);

    return result;
  }

  protected filterAnomaliesBySeverity(anomalies: Anomaly[], options?: ValidationOptions): Anomaly[] {
    // Filter out low-severity anomalies if tolerance is high
    const tolerance = options?.tolerance || 0;

    if (tolerance > 0.8) {
      return anomalies.filter(a => a.severity !== 'low');
    } else if (tolerance > 0.5) {
      return anomalies.filter(a => a.severity !== 'low' || a.confidence > 0.9);
    }

    return anomalies;
  }

  protected calculateOverallScore(result: ValidationResult): number {
    if (result.anomalies.length === 0) {
      return result.confidence;
    }

    // Reduce score based on anomaly severity
    const severityPenalties = {
      'low': 0.05,
      'medium': 0.15,
      'high': 0.3,
      'critical': 0.5
    };

    const totalPenalty = result.anomalies.reduce((penalty, anomaly) => {
      return penalty + (severityPenalties[anomaly.severity] * anomaly.confidence);
    }, 0);

    return Math.max(0, result.confidence - totalPenalty);
  }

  // Model validation
  protected async validateModelExists(modelLabel: string): Promise<boolean> {
    try {
      const modelInfo = await this.site.models.getModelInfo(modelLabel);
      return modelInfo !== null;
    } catch (error) {
      return false;
    }
  }

  protected async loadModel(modelLabel: string, version?: string): Promise<any> {
    return await this.site.models.loadModel(modelLabel, version);
  }

  // Anomaly detection helpers
  protected createAnomaly(
    type: Anomaly['type'],
    severity: Anomaly['severity'],
    description: string,
    location: { x: number; y: number; width: number; height: number },
    confidence: number
  ): Anomaly {
    return {
      type,
      severity,
      description,
      location,
      confidence
    };
  }

  protected detectPixelDifferences(
    imageA: Buffer,
    imageB: Buffer,
    threshold: number = 0.1
  ): Anomaly[] {
    // Placeholder for pixel-level comparison
    // Real implementation would use libraries like pixelmatch
    const anomalies: Anomaly[] = [];

    // Simulate finding differences
    if (Math.random() > 0.8) {
      anomalies.push(this.createAnomaly(
        'visual',
        'medium',
        'Pixel differences detected in content area',
        { x: 100, y: 200, width: 50, height: 30 },
        0.85
      ));
    }

    return anomalies;
  }

  protected detectLayoutChanges(
    beforeElements: any[],
    afterElements: any[]
  ): Anomaly[] {
    const anomalies: Anomaly[] = [];

    // Compare element positions and sizes
    beforeElements.forEach((beforeEl, index) => {
      const afterEl = afterElements[index];
      if (!afterEl) {
        anomalies.push(this.createAnomaly(
          'layout',
          'high',
          `Element missing: ${beforeEl.selector}`,
          beforeEl.boundingBox,
          0.95
        ));
        return;
      }

      // Check position changes
      const positionDiff = Math.sqrt(
        Math.pow(afterEl.x - beforeEl.x, 2) +
        Math.pow(afterEl.y - beforeEl.y, 2)
      );

      if (positionDiff > 5) {
        anomalies.push(this.createAnomaly(
          'layout',
          'medium',
          `Element position changed: ${beforeEl.selector}`,
          afterEl.boundingBox,
          0.9
        ));
      }
    });

    return anomalies;
  }

  // Confidence scoring
  protected calculateConfidence(
    predictions: number[],
    threshold: number = 0.5
  ): number {
    if (predictions.length === 0) return 0;

    const maxPrediction = Math.max(...predictions);
    const avgPrediction = predictions.reduce((sum, p) => sum + p, 0) / predictions.length;

    // Combine max and average for confidence score
    return (maxPrediction * 0.7) + (avgPrediction * 0.3);
  }

  // Error handling
  protected createErrorResult(modelLabel: string, error: Error): ValidationResult {
    return {
      isValid: false,
      confidence: 0,
      anomalies: [{
        type: 'visual',
        severity: 'critical',
        description: `Detection failed: ${error.message}`,
        location: { x: 0, y: 0, width: 0, height: 0 },
        confidence: 1.0
      }],
      score: 0,
      timestamp: new Date(),
      modelUsed: modelLabel
    };
  }

  // Logging and monitoring
  protected logDetectionResult(modelLabel: string, result: ValidationResult): void {
    const timestamp = new Date().toISOString();
    const status = result.isValid ? 'PASS' : 'FAIL';
    const confidence = (result.confidence * 100).toFixed(1);
    const anomalyCount = result.anomalies.length;

    console.log(
      `[${timestamp}] ${this.constructor.name}: ${status} - ` +
      `Model: ${modelLabel}, Confidence: ${confidence}%, ` +
      `Anomalies: ${anomalyCount}, Score: ${result.score.toFixed(3)}`
    );

    // Log individual anomalies if any
    if (anomalyCount > 0) {
      result.anomalies.forEach((anomaly, index) => {
        console.log(
          `  Anomaly ${index + 1}: ${anomaly.type}/${anomaly.severity} - ` +
          `${anomaly.description} (${(anomaly.confidence * 100).toFixed(1)}%)`
        );
      });
    }
  }

  // Performance monitoring
  protected async measureDetectionTime<T>(operation: () => Promise<T>): Promise<{
    result: T;
    duration: number;
  }> {
    const startTime = Date.now();
    const result = await operation();
    const duration = Date.now() - startTime;

    return { result, duration };
  }

  // Batch detection
  async validateMultipleImages(
    imageData: Buffer[],
    modelLabel: string,
    options?: ValidationOptions
  ): Promise<ValidationResult[]> {
    const results: ValidationResult[] = [];

    for (const image of imageData) {
      const result = await this.validateImage(image, modelLabel, options);
      results.push(result);
    }

    return results;
  }

  // Statistical analysis
  protected calculateAggregateMetrics(results: ValidationResult[]): {
    averageConfidence: number;
    successRate: number;
    totalAnomalies: number;
    averageScore: number;
  } {
    if (results.length === 0) {
      return {
        averageConfidence: 0,
        successRate: 0,
        totalAnomalies: 0,
        averageScore: 0
      };
    }

    const averageConfidence = results.reduce((sum, r) => sum + r.confidence, 0) / results.length;
    const successRate = results.filter(r => r.isValid).length / results.length;
    const totalAnomalies = results.reduce((sum, r) => sum + r.anomalies.length, 0);
    const averageScore = results.reduce((sum, r) => sum + r.score, 0) / results.length;

    return {
      averageConfidence,
      successRate,
      totalAnomalies,
      averageScore
    };
  }

  // Region-based detection
  protected async validateRegions(
    imageData: Buffer,
    regions: Array<{ name: string; clip: { x: number; y: number; width: number; height: number } }>,
    modelLabel: string,
    options?: ValidationOptions
  ): Promise<Array<{ region: string; result: ValidationResult }>> {
    const results: Array<{ region: string; result: ValidationResult }> = [];

    for (const region of regions) {
      // Extract region from image (would use image processing library)
      const regionImage = await this.extractImageRegion(imageData, region.clip);

      const result = await this.validateImage(regionImage, `${modelLabel}_${region.name}`, options);
      results.push({ region: region.name, result });
    }

    return results;
  }

  protected async extractImageRegion(
    imageData: Buffer,
    clip: { x: number; y: number; width: number; height: number }
  ): Promise<Buffer> {
    // Placeholder - would use Sharp or similar library to extract region
    // For now, return original image
    return imageData;
  }

  // Template matching
  protected async detectTemplate(
    imageData: Buffer,
    templateImage: Buffer,
    threshold: number = 0.8
  ): Promise<Array<{ x: number; y: number; confidence: number }>> {
    // Placeholder for template matching implementation
    // Would use OpenCV.js or similar library
    return [];
  }

  // Feature detection
  protected async extractFeatures(imageData: Buffer): Promise<{
    edges: number;
    corners: number;
    textRegions: Array<{ x: number; y: number; width: number; height: number }>;
    colorHistogram: number[];
  }> {
    // Placeholder for feature extraction
    // Would use computer vision libraries
    return {
      edges: 0,
      corners: 0,
      textRegions: [],
      colorHistogram: []
    };
  }

  // Similarity scoring
  protected calculateSimilarityScore(
    featuresA: any,
    featuresB: any
  ): number {
    // Placeholder for similarity calculation
    // Would implement actual similarity metrics (SSIM, PSNR, etc.)
    return Math.random() * 0.3 + 0.7; // Random score between 0.7-1.0
  }

  // Model warm-up
  async warmUpModel(modelLabel: string): Promise<void> {
    try {
      // Load model into memory for faster subsequent detections
      await this.loadModel(modelLabel);

      // Run a dummy detection to warm up the model
      const dummyImage = Buffer.alloc(1024); // Small dummy image
      await this.detect(dummyImage, modelLabel).catch(() => {
        // Ignore errors during warm-up
      });

      this.logProgress(`Model ${modelLabel} warmed up successfully`);
    } catch (error) {
      this.logProgress(`Failed to warm up model ${modelLabel}: ${error}`);
    }
  }

  // Progress tracking
  protected logProgress(message: string, progress?: number): void {
    const timestamp = new Date().toISOString();
    const progressStr = progress !== undefined ? ` (${Math.round(progress * 100)}%)` : '';
    console.log(`[${timestamp}] ${this.constructor.name}: ${message}${progressStr}`);
  }

  // Cleanup resources
  async cleanup(): Promise<void> {
    // Override in concrete implementations to cleanup model resources
    this.logProgress('Cleaning up detector resources');
  }

  // Health check
  async healthCheck(): Promise<{
    status: 'healthy' | 'degraded' | 'unhealthy';
    modelCount: number;
    averageResponseTime: number;
    lastError?: string;
  }> {
    try {
      const startTime = Date.now();

      // Test with a small dummy validation
      const dummyImage = Buffer.alloc(100);
      await this.detect(dummyImage, 'health_check_model').catch(() => {
        // Expected to fail, just testing response time
      });

      const responseTime = Date.now() - startTime;

      return {
        status: responseTime < 1000 ? 'healthy' : 'degraded',
        modelCount: 0, // Would query actual model count
        averageResponseTime: responseTime
      };
    } catch (error) {
      return {
        status: 'unhealthy',
        modelCount: 0,
        averageResponseTime: 0,
        lastError: (error as Error).message
      };
    }
  }
}