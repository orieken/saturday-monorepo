export class TrainingDataManager {
  // Core CRUD operations
  async saveTrainingData(data: TrainingData[], label: string, metadata?: Record<string, any>): Promise<string>;
  async loadTrainingData(label: string, filters?: DataFilter): Promise<TrainingData[]>;
  async deleteTrainingData(label: string, conditions?: DeleteCondition): Promise<boolean>;

  // Batch operations
  async batchProcessData(data: TrainingData[], operations: ProcessingOperation[]): Promise<TrainingData[]>;
  async mergeDatasetsData(datasets: string[], newLabel: string): Promise<TrainingData[]>;
  async splitDataset(label: string, ratio: number, stratify?: boolean): Promise<{ train: TrainingData[]; test: TrainingData[] }>;

  // Data validation and quality
  async validateDataQuality(data: TrainingData[]): Promise<DataQualityReport>;
  async detectDuplicates(data: TrainingData[]): Promise<DuplicateReport>;
  async balanceDataset(data: TrainingData[], method: 'oversample' | 'undersample' | 'smote'): Promise<TrainingData[]>;

  // Data augmentation
  async augmentDataset(data: TrainingData[], config: AugmentationConfig): Promise<TrainingData[]>;
  async generateSyntheticData(baseData: TrainingData[], count: number): Promise<TrainingData[]>;

  // Analytics and reporting
  async getDatasetStatistics(label: string): Promise<DatasetStatistics>;
  async generateDataReport(label: string): Promise<DataReport>;
  async analyzeDataDistribution(data: TrainingData[]): Promise<DistributionAnalysis>;

  // Data lifecycle management
  async archiveOldData(criteria: ArchiveCriteria): Promise<ArchiveResult>;
  async cleanupExpiredData(): Promise<CleanupResult>;
  async exportDataset(label: string, format: 'json' | 'csv' | 'parquet'): Promise<string>;
  async importDataset(path: string, format: string, label: string): Promise<TrainingData[]>;

  // Version control
  async createDataVersion(label: string, description?: string): Promise<string>;
  async listDataVersions(label: string): Promise<DataVersion[]>;
  async rollbackToVersion(label: string, version: string): Promise<boolean>;
}