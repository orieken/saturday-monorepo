// Core ML Types and Interfaces

export interface TrainingData {
  id: string;
  timestamp: Date;
  label: string;
  imageBuffer: Buffer;
  metadata: TrainingMetadata;
}

export interface TrainingMetadata {
  pageUrl: string;
  elementSelector?: string;
  boundingBox?: { x: number; y: number; width: number; height: number };
  browserInfo: BrowserInfo;
  viewportSize: { width: number; height: number };
  userType?: string;
  testEnvironment?: string;
  [key: string]: any;
}

export interface BrowserInfo {
  userAgent: string;
  platform: string;
  cookieEnabled: boolean;
}

export interface TrainingOptions {
  fullPage?: boolean;
  clip?: { x: number; y: number; width: number; height: number };
  metadata?: Record<string, any>;
  modelConfig?: ModelConfig;
}

export interface ElementCaptureOptions {
  screenshotOptions?: {
    quality?: number;
    type?: 'png' | 'jpeg';
  };
  metadata?: Record<string, any>;
}

export interface CaptureOptions {
  fullPage?: boolean;
  clip?: { x: number; y: number; width: number; height: number };
  metadata?: Record<string, any>;
}

export interface ValidationOptions {
  threshold?: number;
  fullPage?: boolean;
  clip?: { x: number; y: number; width: number; height: number };
  regions?: string[];
  tolerance?: number;
}

export interface ElementValidationOptions {
  threshold?: number;
  tolerance?: number;
  compareStructure?: boolean;
}

export interface TrainingResult {
  success: boolean;
  modelLabel: string;
  version: string;
  metrics: TrainingMetrics;
  timestamp: Date;
  dataCount: number;
}

export interface TrainingMetrics {
  accuracy?: number;
  loss?: number;
  epochs: number;
  trainingTime: number;
}

export interface ValidationResult {
  isValid: boolean;
  confidence: number;
  anomalies: Anomaly[];
  score: number;
  timestamp: Date;
  modelUsed: string;
}

export interface ElementAnomalyResult {
  isValid: boolean;
  confidence: number;
  anomalies: Anomaly[];
  elementSelector: string;
  timestamp: Date;
  modelUsed: string;
}

export interface Anomaly {
  type: 'visual' | 'layout' | 'content' | 'structure';
  severity: 'low' | 'medium' | 'high' | 'critical';
  description: string;
  location: { x: number; y: number; width: number; height: number };
  confidence: number;
}

export interface ComparisonResult {
  similarity: number;
  differences: Difference[];
  overallMatch: boolean;
  timestamp: Date;
}

export interface Difference {
  type: 'pixel' | 'structure' | 'layout';
  severity: 'low' | 'medium' | 'high';
  area: { x: number; y: number; width: number; height: number };
  description: string;
}

export interface ModelConfig {
  architecture: 'cnn' | 'autoencoder' | 'siamese';
  inputSize: { width: number; height: number };
  batchSize: number;
  epochs: number;
  learningRate: number;
  optimizerType: 'adam' | 'sgd' | 'rmsprop';
}

export interface LayoutDetectionOptions {
  tolerance?: number;
  excludeSelectors?: string[];
  includeHidden?: boolean;
}

export interface LayoutAnomalyResult {
  type: 'misalignment' | 'overflow' | 'missing' | 'unexpected';
  severity: 'low' | 'medium' | 'high' | 'critical';
  element: string;
  expected: any;
  actual: any;
  location: { x: number; y: number; width: number; height: number };
}

export interface ComparisonOptions {
  tolerance?: number;
  ignoreAntialiasing?: boolean;
  threshold?: number;
  includeAA?: boolean;
}

export interface FullValidationResult {
  anomalies: ValidationResult;
  regression: ComparisonResult;
  layout: LayoutAnomalyResult[];
  overall: ValidationScore;
}

export interface ValidationScore {
  score: number;
  confidence: number;
  isValid: boolean;
}

export interface PageTrainingConfig {
  pageName: string;
  label: string;
  options?: TrainingOptions;
}

export interface ElementTrainingOptions extends TrainingOptions {
  waitForStable?: boolean;
  captureStates?: string[];
}

export interface FlowTrainingConfig {
  flowName: string;
  params: Record<string, any>;
  label?: string;
  capturePoints?: string[];
}

// Training and Detection Configuration Types
export interface TrainingConfiguration {
  modelType: 'anomaly' | 'regression' | 'classification';
  dataAugmentation: boolean;
  validationSplit: number;
  earlyStoppingPatience: number;
  saveCheckpoints: boolean;
}

export interface DetectionConfiguration {
  anomalyThreshold: number;
  similarityThreshold: number;
  enableBoundingBoxes: boolean;
  maxAnomalies: number;
  confidenceThreshold: number;
}

// Model Management Types
export interface ModelInfo {
  label: string;
  version: string;
  type: string;
  createdAt: Date;
  lastUsed: Date;
  accuracy?: number;
  dataCount: number;
  metadata: Record<string, any>;
}

export interface ModelPerformanceMetrics {
  accuracy: number;
  precision: number;
  recall: number;
  f1Score: number;
  validationLoss: number;
  trainingTime: number;
}

// Storage and Persistence Types
export interface DataStorageOptions {
  path: string;
  compression: boolean;
  encryption?: boolean;
  retention: string;
  cleanup: boolean;
}

export interface ModelStorageOptions {
  path: string;
  versioning: boolean;
  compression: boolean;
  metadata: boolean;
  backup: boolean;
}