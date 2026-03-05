export class ImageProcessor {
  // Basic image operations
  async resizeImage(imageData: Buffer, width: number, height: number, options?: ResizeOptions): Promise<Buffer>;
  async cropImage(imageData: Buffer, region: CropRegion): Promise<Buffer>;
  async rotateImage(imageData: Buffer, angle: number): Promise<Buffer>;
  async flipImage(imageData: Buffer, direction: 'horizontal' | 'vertical'): Promise<Buffer>;

  // Image enhancement
  async adjustBrightness(imageData: Buffer, factor: number): Promise<Buffer>;
  async adjustContrast(imageData: Buffer, factor: number): Promise<Buffer>;
  async adjustSaturation(imageData: Buffer, factor: number): Promise<Buffer>;
  async sharpenImage(imageData: Buffer, intensity: number): Promise<Buffer>;
  async blurImage(imageData: Buffer, radius: number): Promise<Buffer>;
  async denoiseImage(imageData: Buffer, method: 'gaussian' | 'median' | 'bilateral'): Promise<Buffer>;

  // Color operations
  async convertColorSpace(imageData: Buffer, from: ColorSpace, to: ColorSpace): Promise<Buffer>;
  async normalizeColors(imageData: Buffer): Promise<Buffer>;
  async extractColorPalette(imageData: Buffer, count: number): Promise<Color[]>;
  async replaceColors(imageData: Buffer, colorMap: Map<Color, Color>): Promise<Buffer>;

  // Image analysis
  async getImageDimensions(imageData: Buffer): Promise<{ width: number; height: number }>;
  async getImageInfo(imageData: Buffer): Promise<ImageInfo>;
  async calculateHistogram(imageData: Buffer, channel?: 'red' | 'green' | 'blue' | 'alpha'): Promise<number[]>;
  async detectEdges(imageData: Buffer, method: 'canny' | 'sobel' | 'laplacian'): Promise<Buffer>;
  async extractFeatures(imageData: Buffer, method: 'orb' | 'sift' | 'surf'): Promise<ImageFeature[]>;

  // Image comparison
  async calculateSimilarity(imageA: Buffer, imageB: Buffer, metric: SimilarityMetric): Promise<number>;
  async findDifferences(imageA: Buffer, imageB: Buffer, threshold?: number): Promise<DifferenceResult>;
  async alignImages(imageA: Buffer, imageB: Buffer): Promise<{ alignedA: Buffer; alignedB: Buffer; transform: Transform }>;
  async generateDifferenceMap(imageA: Buffer, imageB: Buffer): Promise<Buffer>;

  // Template matching
  async findTemplate(image: Buffer, template: Buffer, threshold: number): Promise<TemplateMatch[]>;
  async matchFeatures(imageA: Buffer, imageB: Buffer): Promise<FeatureMatch[]>;

  // Quality assessment
  async assessImageQuality(imageData: Buffer): Promise<QualityMetrics>;
  async detectBlur(imageData: Buffer): Promise<{ isBlurry: boolean; blurScore: number }>;
  async detectNoise(imageData: Buffer): Promise<{ hasNoise: boolean; noiseLevel: number }>;
  async validateImageIntegrity(imageData: Buffer): Promise<IntegrityReport>;

  // Format operations
  async convertFormat(imageData: Buffer, targetFormat: ImageFormat, quality?: number): Promise<Buffer>;
  async optimizeImage(imageData: Buffer, options: OptimizationOptions): Promise<{ data: Buffer; originalSize: number; newSize: number }>;
  async generateThumbnail(imageData: Buffer, size: number, crop?: boolean): Promise<Buffer>;

  // Batch operations
  async batchProcess(images: Buffer[], operations: ImageOperation[]): Promise<Buffer[]>;
  async batchResize(images: Buffer[], dimensions: { width: number; height: number }): Promise<Buffer[]>;
  async batchConvert(images: Buffer[], targetFormat: ImageFormat): Promise<Buffer[]>;

  // Advanced operations
  async removeBackground(imageData: Buffer, method: 'color' | 'ai' | 'mask'): Promise<Buffer>;
  async applyFilter(imageData: Buffer, filter: ImageFilter): Promise<Buffer>;
  async createMontage(images: Buffer[], layout: MontageLayout): Promise<Buffer>;
  async watermarkImage(imageData: Buffer, watermark: Buffer | string, position: WatermarkPosition): Promise<Buffer>;

  // ML preprocessing
  async normalizeForML(imageData: Buffer, targetSize: { width: number; height: number }): Promise<Float32Array>;
  async augmentImage(imageData: Buffer, augmentations: ImageAugmentation[]): Promise<Buffer[]>;
  async extractPatches(imageData: Buffer, patchSize: number, overlap: number): Promise<Buffer[]>;
  async preprocessForModel(imageData: Buffer, modelConfig: ModelPreprocessConfig): Promise<Buffer>;
}