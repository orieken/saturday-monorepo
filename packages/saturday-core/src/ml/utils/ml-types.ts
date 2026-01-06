export interface DataFilter {
  dateRange?: { start: Date; end: Date };
  labels?: string[];
  metadata?: Record<string, any>;
  minConfidence?: number;
  pageUrls?: string[];
  browsers?: string[];
  viewports?: string[];
}

export interface DataQualityReport {
  totalSamples: number;
  validSamples: number;
  invalidSamples: number;
  duplicates: number;
  missingLabels: number;
  corruptedImages: number;
  inconsistentMetadata: number;
  recommendations: string[];
}

export interface AugmentationConfig {
  rotation: { enabled: boolean; range: number };
  brightness: { enabled: boolean; range: number };
  contrast: { enabled: boolean; range: number };
  noise: { enabled: boolean; intensity: number };
  flip: { enabled: boolean; horizontal: boolean; vertical: boolean };
  crop: { enabled: boolean; ratio: number };
  multiplier: number; // How many augmented samples per original
}

export interface DatasetStatistics {
  totalSamples: number;
  labelDistribution: Record<string, number>;
  averageImageSize: { width: number; height: number };
  totalStorageSize: string;
  createdDate: Date;
  lastModified: Date;
  samplesByDate: Array<{ date: string; count: number }>;
  metadataFields: string[];
}

export interface ModelMetadata {
  name: string;
  description: string;
  author: string;
  tags: string[];
  framework: string;
  architecture: string;
  inputShape: number[];
  outputShape: number[];
  trainingConfig: ModelConfig;
  performance: ModelPerformanceMetrics;
  dependencies: string[];
  environment: Record<string, string>;
}

export interface OptimizationConfig {
  type: 'quantization' | 'pruning' | 'distillation' | 'tensorrt';
  target: 'speed' | 'size' | 'accuracy';
  parameters: Record<string, any>;
  preserveAccuracy: boolean;
  maxAccuracyLoss: number;
}

export interface DeploymentEnvironment {
  name: string;
  type: 'local' | 'cloud' | 'edge';
  resources: { cpu: string; memory: string; gpu?: string };
  scalingConfig: { minReplicas: number; maxReplicas: number };
  healthCheck: { enabled: boolean; endpoint: string; interval: number };
}

export interface ImageInfo {
  width: number;
  height: number;
  channels: number;
  format: string;
  colorSpace: ColorSpace;
  hasAlpha: boolean;
  bitDepth: number;
  fileSize: number;
  compressionRatio?: number;
}

export interface SimilarityMetric {
  type: 'mse' | 'ssim' | 'psnr' | 'perceptual' | 'histogram';
  parameters?: Record<string, any>;
}

export interface DifferenceResult {
  overallSimilarity: number;
  differences: Array<{
    region: { x: number; y: number; width: number; height: number };
    severity: 'low' | 'medium' | 'high';
    type: 'color' | 'structure' | 'content';
    description: string;
  }>;
  differenceMap?: Buffer;
  pixelDifferenceCount: number;
  percentageDifferent: number;
}

export interface QualityMetrics {
  sharpness: number;
  brightness: number;
  contrast: number;
  colorfulness: number;
  noiseLevel: number;
  overallQuality: number;
  recommendations: string[];
}

export interface ImageOperation {
  type: 'resize' | 'crop' | 'rotate' | 'enhance' | 'convert';
  parameters: Record<string, any>;
  order: number;
}

export interface ImageAugmentation {
  type: 'rotation' | 'flip' | 'brightness' | 'contrast' | 'noise' | 'crop' | 'scale';
  probability: number;
  parameters: Record<string, any>;
}