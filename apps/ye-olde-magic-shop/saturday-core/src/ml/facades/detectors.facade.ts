import { BaseSite } from '../../base/base-site';
import { AnomalyDetector } from '../detectors/anomaly-detector';
import { RegressionDetector } from '../detectors/regression-detector';
import { LayoutDetector } from '../detectors/layout-detector';
import { ElementDetector } from '../detectors/element-detector';
import {
  ValidationOptions,
  ValidationResult,
  ComparisonOptions,
  ComparisonResult,
  LayoutDetectionOptions,
  LayoutAnomalyResult,
  FullValidationResult,
  ValidationScore
} from '../interfaces/ml-types';

export class DetectorsFacade {
  private _anomalyDetector?: AnomalyDetector;
  private _regressionDetector?: RegressionDetector;
  private _layoutDetector?: LayoutDetector;
  private _elementDetector?: ElementDetector;

  constructor(private site: BaseSite) {
  }

  get anomalyDetector(): AnomalyDetector {
    if (!this._anomalyDetector) {
      this._anomalyDetector = new AnomalyDetector(this.site);
    }
    return this._anomalyDetector;
  }

  get regressionDetector(): RegressionDetector {
    if (!this._regressionDetector) {
      this._regressionDetector = new RegressionDetector(this.site);
    }
    return this._regressionDetector;
  }

  get layoutDetector(): LayoutDetector {
    if (!this._layoutDetector) {
      this._layoutDetector = new LayoutDetector(this.site);
    }
    return this._layoutDetector;
  }

  get elementDetector(): ElementDetector {
    if (!this._elementDetector) {
      this._elementDetector = new ElementDetector(this.site);
    }
    return this._elementDetector;
  }

  // Convenience methods
  async validateCurrentPage(modelLabel: string, options?: ValidationOptions): Promise<ValidationResult> {
    return await this.anomalyDetector.validateCurrentPage(modelLabel, options);
  }

  async compareWithBaseline(baselineLabel: string, options?: ComparisonOptions): Promise<ComparisonResult> {
    return await this.regressionDetector.compareWithBaseline(baselineLabel, options);
  }

  async detectLayoutAnomalies(options?: LayoutDetectionOptions): Promise<LayoutAnomalyResult[]> {
    return await this.layoutDetector.detectAnomalies(options);
  }

  // Comprehensive validation
  async runFullValidation(modelLabel: string, options?: {
    validationOptions?: ValidationOptions;
    comparisonOptions?: ComparisonOptions;
    layoutOptions?: LayoutDetectionOptions;
  }): Promise<FullValidationResult> {
    const [anomalies, regression, layout] = await Promise.all([
      this.anomalyDetector.validateCurrentPage(modelLabel, options?.validationOptions),
      this.regressionDetector.compareWithBaseline(modelLabel, options?.comparisonOptions),
      this.layoutDetector.detectAnomalies(options?.layoutOptions)
    ]);

    return {
      anomalies,
      regression,
      layout,
      overall: this.calculateOverallScore([anomalies, regression, layout])
    };
  }

  // Element-specific validation
  async validateElement(selector: string, modelLabel: string, options?: ValidationOptions): Promise<ValidationResult> {
    return await this.elementDetector.validateElement(selector, modelLabel, options);
  }

  // Batch validation
  async validateMultipleElements(elements: Array<{
    selector: string;
    modelLabel: string;
    options?: ValidationOptions;
  }>): Promise<ValidationResult[]> {
    const results: ValidationResult[] = [];

    for (const element of elements) {
      const result = await this.validateElement(element.selector, element.modelLabel, element.options);
      results.push(result);
    }

    return results;
  }

  // Cross-page validation
  async validateUserJourney(journeySteps: Array<{
    pageName: string;
    modelLabel: string;
    elements?: Array<{ selector: string; modelLabel: string }>;
    options?: ValidationOptions;
  }>): Promise<Array<{
    pageName: string;
    pageValidation: ValidationResult;
    elementValidations: ValidationResult[];
    overallValid: boolean;
  }>> {
    const results: Array<{
      pageName: string;
      pageValidation: ValidationResult;
      elementValidations: ValidationResult[];
      overallValid: boolean;
    }> = [];

    for (const step of journeySteps) {
      const page = this.site.getPage(step.pageName);
      await page.visit();

      const pageValidation = await this.validateCurrentPage(step.modelLabel, step.options);

      const elementValidations: ValidationResult[] = [];
      if (step.elements) {
        for (const element of step.elements) {
          const elementValidation = await this.validateElement(element.selector, element.modelLabel);
          elementValidations.push(elementValidation);
        }
      }

      const overallValid = pageValidation.isValid &&
        elementValidations.every(ev => ev.isValid);

      results.push({
        pageName: step.pageName,
        pageValidation,
        elementValidations,
        overallValid
      });
    }

    return results;
  }

  // Performance monitoring
  async monitorPagePerformance(modelLabel: string, iterations: number = 5): Promise<{
    averageConfidence: number;
    consistencyScore: number;
    anomaliesDetected: number;
    performanceMetrics: {
      validationTime: number;
      modelLoadTime: number;
      screenshotTime: number;
    };
  }> {
    const results: ValidationResult[] = [];
    let totalValidationTime = 0;
    let totalModelLoadTime = 0;
    let totalScreenshotTime = 0;

    for (let i = 0; i < iterations; i++) {
      const startTime = Date.now();

      // Simulate model load time
      const modelLoadStart = Date.now();
      await new Promise(resolve => setTimeout(resolve, 10)); // Placeholder
      const modelLoadTime = Date.now() - modelLoadStart;

      // Simulate screenshot time
      const screenshotStart = Date.now();
      await this.site.pageObject.screenshot();
      const screenshotTime = Date.now() - screenshotStart;

      const result = await this.validateCurrentPage(modelLabel);
      const validationTime = Date.now() - startTime;

      results.push(result);
      totalValidationTime += validationTime;
      totalModelLoadTime += modelLoadTime;
      totalScreenshotTime += screenshotTime;
    }

    const averageConfidence = results.reduce((sum, r) => sum + r.confidence, 0) / results.length;
    const anomaliesDetected = results.reduce((sum, r) => sum + r.anomalies.length, 0);

    // Calculate consistency score based on confidence variance
    const confidenceVariance = results.reduce((sum, r) =>
      sum + Math.pow(r.confidence - averageConfidence, 2), 0) / results.length;
    const consistencyScore = Math.max(0, 1 - Math.sqrt(confidenceVariance));

    return {
      averageConfidence,
      consistencyScore,
      anomaliesDetected,
      performanceMetrics: {
        validationTime: totalValidationTime / iterations,
        modelLoadTime: totalModelLoadTime / iterations,
        screenshotTime: totalScreenshotTime / iterations
      }
    };
  }

  // A/B testing support
  async compareVariants(variantA: string, variantB: string, modelLabel: string): Promise<{
    variantA: ValidationResult;
    variantB: ValidationResult;
    comparison: ComparisonResult;
    recommendation: 'A' | 'B' | 'inconclusive';
  }> {
    // Navigate to variant A
    await this.site.pageObject.goto(variantA);
    const resultA = await this.validateCurrentPage(modelLabel);

    // Navigate to variant B
    await this.site.pageObject.goto(variantB);
    const resultB = await this.validateCurrentPage(modelLabel);

    // Compare the two variants
    const comparison = await this.regressionDetector.compareScreenshots(
      await this.site.pageObject.screenshot(),
      await this.site.pageObject.screenshot() // This would be variant A screenshot in real implementation
    );

    let recommendation: 'A' | 'B' | 'inconclusive' = 'inconclusive';
    if (resultA.confidence > resultB.confidence + 0.05) {
      recommendation = 'A';
    } else if (resultB.confidence > resultA.confidence + 0.05) {
      recommendation = 'B';
    }

    return {
      variantA: resultA,
      variantB: resultB,
      comparison,
      recommendation
    };
  }

  private calculateOverallScore(results: any[]): ValidationScore {
    // Calculate weighted score based on different validation types
    const anomalyWeight = 0.4;
    const regressionWeight = 0.4;
    const layoutWeight = 0.2;

    const [anomalies, regression, layout] = results;

    const anomalyScore = anomalies.isValid ? anomalies.confidence : 0;
    const regressionScore = regression.overallMatch ? regression.similarity : 0;
    const layoutScore = layout.length === 0 ? 1.0 : Math.max(0, 1 - (layout.length * 0.1));

    const overallScore = (anomalyScore * anomalyWeight) +
      (regressionScore * regressionWeight) +
      (layoutScore * layoutWeight);

    const confidence = Math.min(anomalies.confidence, regression.similarity || 1.0);

    return {
      score: overallScore,
      confidence,
      isValid: overallScore > 0.8 && confidence > 0.85
    };
  }

  // Debugging and diagnostics
  async generateDiagnosticReport(modelLabel: string): Promise<{
    modelInfo: any;
    pageMetrics: any;
    validationHistory: ValidationResult[];
    recommendations: string[];
  }> {
    // This would integrate with model management and history tracking
    return {
      modelInfo: { version: '1.0', accuracy: 0.95, trainingDate: new Date() },
      pageMetrics: { loadTime: 1500, elementCount: 25, complexity: 'medium' },
      validationHistory: [],
      recommendations: [
        'Consider retraining model with more diverse data',
        'Increase confidence threshold for critical validations',
        'Add element-specific validation for dynamic content'
      ]
    };
  }
}

