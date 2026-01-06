
export interface DetectorConfig {
  type: 'anomaly' | 'regression' | 'layout' | 'element';
  enabled: boolean;
  priority: number;
  thresholds: {
    confidence: number;
    similarity: number;
    anomaly: number;
    layout: number;
  };
  tolerances: {
    pixel: number;
    layout: number;
    color: number;
    position: number;
  };
  preprocessing: {
    resize: boolean;
    normalize: boolean;
    denoise: boolean;
    enhance: boolean;
  };
  postprocessing: {
    filterLowConfidence: boolean;
    mergeSimilarAnomalies: boolean;
    spatialFiltering: boolean;
    temporalFiltering: boolean;
  };
  performance: {
    maxConcurrent: number;
    timeout: number;
    cacheResults: boolean;
    batchSize: number;
  };
}

export interface AnomalyDetectionConfig extends DetectorConfig {
  type: 'anomaly';
  anomalyTypes: {
    visual: {
      enabled: boolean;
      sensitivity: number;
      regions: string[];
      excludeSelectors: string[];
    };
    content: {
      enabled: boolean;
      textChanges: boolean;
      imageChanges: boolean;
      linkChanges: boolean;
    };
    structure: {
      enabled: boolean;
      elementCount: boolean;
      hierarchy: boolean;
      attributes: boolean;
    };
    behavioral: {
      enabled: boolean;
      interactions: boolean;
      animations: boolean;
      loadTimes: boolean;
    };
  };
  learningParameters: {
    adaptiveThreshold: boolean;
    contextualLearning: boolean;
    temporalPatterns: boolean;
    seasonalAdjustment: boolean;
  };
}

export interface RegressionDetectionConfig extends DetectorConfig {
  type: 'regression';
  comparisonMethods: {
    pixelLevel: {
      enabled: boolean;
      algorithm: 'mse' | 'ssim' | 'psnr' | 'perceptual';
      weightMask: boolean;
      ignoreAntialiasing: boolean;
    };
    structural: {
      enabled: boolean;
      domComparison: boolean;
      layoutComparison: boolean;
      styleComparison: boolean;
    };
    semantic: {
      enabled: boolean;
      contentComparison: boolean;
      functionalComparison: boolean;
      accessibilityComparison: boolean;
    };
  };
  baselineManagement: {
    autoUpdate: boolean;
    versionControl: boolean;
    approvalWorkflow: boolean;
    rollbackCapability: boolean;
  };
}

export interface LayoutDetectionConfig extends DetectorConfig {
  type: 'layout';
  layoutChecks: {
    positioning: {
      enabled: boolean;
      absolutePositions: boolean;
      relativePositions: boolean;
      zIndex: boolean;
      overflow: boolean;
    };
    sizing: {
      enabled: boolean;
      elementDimensions: boolean;
      aspectRatios: boolean;
      responsive: boolean;
      minMaxConstraints: boolean;
    };
    spacing: {
      enabled: boolean;
      margins: boolean;
      padding: boolean;
      gaps: boolean;
      alignment: boolean;
    };
    typography: {
      enabled: boolean;
      fontSizes: boolean;
      lineHeights: boolean;
      letterSpacing: boolean;
      textAlignment: boolean;
    };
  };
  responsiveValidation: {
    enabled: boolean;
    breakpoints: Array<{
      name: string;
      width: number;
      height: number;
      tolerance: number;
    }>;
    orientations: ('portrait' | 'landscape')[];
    deviceTypes: ('mobile' | 'tablet' | 'desktop')[];
  };
}

export interface ElementDetectionConfig extends DetectorConfig {
  type: 'element';
  elementTypes: {
    navigation: {
      enabled: boolean;
      consistency: boolean;
      accessibility: boolean;
      interactions: boolean;
    };
    forms: {
      enabled: boolean;
      validation: boolean;
      styling: boolean;
      behavior: boolean;
    };
    media: {
      enabled: boolean;
      images: boolean;
      videos: boolean;
      loading: boolean;
    };
    interactive: {
      enabled: boolean;
      buttons: boolean;
      links: boolean;
      dropdowns: boolean;
    };
  };
  validationRules: {
    presence: boolean;
    visibility: boolean;
    accessibility: boolean;
    functionality: boolean;
    styling: boolean;
  };
}

// Default configurations for different environments
const DevelopmentConfig = {
    anomaly: {
      type: 'anomaly',
      enabled: true,
      priority: 1,
      thresholds: {
        confidence: 0.7,
        similarity: 0.8,
        anomaly: 0.3,
        layout: 0.1
      },
      tolerances: {
        pixel: 0.05,
        layout: 3,
        color: 0.1,
        position: 2
      },
      preprocessing: {
        resize: true,
        normalize: true,
        denoise: false,
        enhance: false
      },
      postprocessing: {
        filterLowConfidence: false,
        mergeSimilarAnomalies: true,
        spatialFiltering: false,
        temporalFiltering: false
      },
      performance: {
        maxConcurrent: 2,
        timeout: 30000,
        cacheResults: true,
        batchSize: 4
      },
      anomalyTypes: {
        visual: {
          enabled: true,
          sensitivity: 0.7,
          regions: ['header', 'main', 'footer'],
          excludeSelectors: ['.dynamic-content', '.timestamp']
        },
        content: {
          enabled: true,
          textChanges: true,
          imageChanges: true,
          linkChanges: false
        },
        structure: {
          enabled: true,
          elementCount: true,
          hierarchy: true,
          attributes: false
        },
        behavioral: {
          enabled: false,
          interactions: false,
          animations: false,
          loadTimes: false
        }
      },
      learningParameters: {
        adaptiveThreshold: true,
        contextualLearning: false,
        temporalPatterns: false,
        seasonalAdjustment: false
      }
    } as AnomalyDetectionConfig,

    regression: {
      type: 'regression',
      enabled: true,
      priority: 2,
      thresholds: {
        confidence: 0.85,
        similarity: 0.95,
        anomaly: 0.15,
        layout: 0.05
      },
      tolerances: {
        pixel: 0.02,
        layout: 1,
        color: 0.05,
        position: 1
      },
      preprocessing: {
        resize: true,
        normalize: true,
        denoise: true,
        enhance: false
      },
      postprocessing: {
        filterLowConfidence: true,
        mergeSimilarAnomalies: true,
        spatialFiltering: true,
        temporalFiltering: false
      },
      performance: {
        maxConcurrent: 3,
        timeout: 20000,
        cacheResults: true,
        batchSize: 8
      },
      comparisonMethods: {
        pixelLevel: {
          enabled: true,
          algorithm: 'ssim',
          weightMask: false,
          ignoreAntialiasing: true
        },
        structural: {
          enabled: true,
          domComparison: true,
          layoutComparison: true,
          styleComparison: false
        },
        semantic: {
          enabled: false,
          contentComparison: false,
          functionalComparison: false,
          accessibilityComparison: false
        }
      },
      baselineManagement: {
        autoUpdate: false,
        versionControl: true,
        approvalWorkflow: false,
        rollbackCapability: true
      }
    } as RegressionDetectionConfig,

    layout: {
      type: 'layout',
      enabled: true,
      priority: 3,
      thresholds: {
        confidence: 0.8,
        similarity: 0.9,
        anomaly: 0.2,
        layout: 0.08
      },
      tolerances: {
        pixel: 0.03,
        layout: 2,
        color: 0.1,
        position: 2
      },
      preprocessing: {
        resize: false,
        normalize: false,
        denoise: false,
        enhance: false
      },
      postprocessing: {
        filterLowConfidence: true,
        mergeSimilarAnomalies: true,
        spatialFiltering: false,
        temporalFiltering: false
      },
      performance: {
        maxConcurrent: 4,
        timeout: 15000,
        cacheResults: false,
        batchSize: 1
      },
      layoutChecks: {
        positioning: {
          enabled: true,
          absolutePositions: true,
          relativePositions: true,
          zIndex: false,
          overflow: true
        },
        sizing: {
          enabled: true,
          elementDimensions: true,
          aspectRatios: false,
          responsive: true,
          minMaxConstraints: false
        },
        spacing: {
          enabled: true,
          margins: true,
          padding: true,
          gaps: false,
          alignment: true
        },
        typography: {
          enabled: false,
          fontSizes: false,
          lineHeights: false,
          letterSpacing: false,
          textAlignment: false
        }
      },
      responsiveValidation: {
        enabled: true,
        breakpoints: [
          { name: 'mobile', width: 375, height: 667, tolerance: 0.1 },
          { name: 'tablet', width: 768, height: 1024, tolerance: 0.08 },
          { name: 'desktop', width: 1920, height: 1080, tolerance: 0.05 }
        ],
        orientations: ['portrait', 'landscape'],
        deviceTypes: ['mobile', 'tablet', 'desktop']
      }
    } as LayoutDetectionConfig,

    element: {
      type: 'element',
      enabled: true,
      priority: 4,
      thresholds: {
        confidence: 0.9,
        similarity: 0.95,
        anomaly: 0.1,
        layout: 0.05
      },
      tolerances: {
        pixel: 0.01,
        layout: 1,
        color: 0.03,
        position: 1
      },
      preprocessing: {
        resize: true,
        normalize: true,
        denoise: false,
        enhance: true
      },
      postprocessing: {
        filterLowConfidence: true,
        mergeSimilarAnomalies: false,
        spatialFiltering: false,
        temporalFiltering: false
      },
      performance: {
        maxConcurrent: 5,
        timeout: 10000,
        cacheResults: true,
        batchSize: 10
      },
      elementTypes: {
        navigation: {
          enabled: true,
          consistency: true,
          accessibility: false,
          interactions: false
        },
        forms: {
          enabled: true,
          validation: true,
          styling: true,
          behavior: false
        },
        media: {
          enabled: true,
          images: true,
          videos: false,
          loading: false
        },
        interactive: {
          enabled: true,
          buttons: true,
          links: true,
          dropdowns: false
        }
      },
      validationRules: {
        presence: true,
        visibility: true,
        accessibility: false,
        functionality: false,
        styling: true
      }
    } as ElementDetectionConfig
};

const StagingConfig = {
    anomaly: {
      ...DevelopmentConfig.anomaly,
      thresholds: {
        confidence: 0.8,
        similarity: 0.9,
        anomaly: 0.2,
        layout: 0.05
      },
      performance: {
        maxConcurrent: 3,
        timeout: 25000,
        cacheResults: true,
        batchSize: 6
      }
    } as AnomalyDetectionConfig,

    regression: {
      ...DevelopmentConfig.regression,
      thresholds: {
        confidence: 0.9,
        similarity: 0.97,
        anomaly: 0.1,
        layout: 0.03
      },
      baselineManagement: {
        autoUpdate: false,
        versionControl: true,
        approvalWorkflow: true,
        rollbackCapability: true
      }
    } as RegressionDetectionConfig,

    layout: {
      ...DevelopmentConfig.layout,
      thresholds: {
        confidence: 0.85,
        similarity: 0.95,
        anomaly: 0.15,
        layout: 0.05
      }
    } as LayoutDetectionConfig,

    element: {
      ...DevelopmentConfig.element,
      thresholds: {
        confidence: 0.92,
        similarity: 0.97,
        anomaly: 0.08,
        layout: 0.03
      },
      elementTypes: {
        ...DevelopmentConfig.element.elementTypes,
        navigation: {
          enabled: true,
          consistency: true,
          accessibility: true,
          interactions: false
        }
      }
    } as ElementDetectionConfig
};

const ProductionConfig = {
    anomaly: {
      ...StagingConfig.anomaly,
      thresholds: {
        confidence: 0.85,
        similarity: 0.95,
        anomaly: 0.15,
        layout: 0.03
      },
      performance: {
        maxConcurrent: 4,
        timeout: 20000,
        cacheResults: true,
        batchSize: 8
      },
      learningParameters: {
        adaptiveThreshold: true,
        contextualLearning: true,
        temporalPatterns: true,
        seasonalAdjustment: false
      }
    } as AnomalyDetectionConfig,

    regression: {
      ...StagingConfig.regression,
      thresholds: {
        confidence: 0.92,
        similarity: 0.98,
        anomaly: 0.08,
        layout: 0.02
      },
      comparisonMethods: {
        ...StagingConfig.regression.comparisonMethods,
        semantic: {
          enabled: true,
          contentComparison: true,
          functionalComparison: false,
          accessibilityComparison: true
        }
      }
    } as RegressionDetectionConfig,

    layout: {
      ...StagingConfig.layout,
      thresholds: {
        confidence: 0.9,
        similarity: 0.97,
        anomaly: 0.1,
        layout: 0.03
      },
      layoutChecks: {
        ...StagingConfig.layout.layoutChecks,
        typography: {
          enabled: true,
          fontSizes: true,
          lineHeights: true,
          letterSpacing: false,
          textAlignment: true
        }
      }
    } as LayoutDetectionConfig,

    element: {
      ...StagingConfig.element,
      thresholds: {
        confidence: 0.95,
        similarity: 0.98,
        anomaly: 0.05,
        layout: 0.02
      },
      elementTypes: {
        navigation: {
          enabled: true,
          consistency: true,
          accessibility: true,
          interactions: true
        },
        forms: {
          enabled: true,
          validation: true,
          styling: true,
          behavior: true
        },
        media: {
          enabled: true,
          images: true,
          videos: true,
          loading: true
        },
        interactive: {
          enabled: true,
          buttons: true,
          links: true,
          dropdowns: true
        }
      },
      validationRules: {
        presence: true,
        visibility: true,
        accessibility: true,
        functionality: true,
        styling: true
      }
    } as ElementDetectionConfig
};

// Default configurations for different environments
export const DetectionConfigs = {
  development: DevelopmentConfig,
  staging: StagingConfig,
  production: ProductionConfig
};

// Site-specific configuration overrides
export const SiteDetectionConfigs = {
  foo: {
    // E-commerce specific configurations
    anomaly: {
      anomalyTypes: {
        visual: {
          enabled: true,
          sensitivity: 0.8,
          regions: ['header', 'product-gallery', 'pricing', 'cart', 'footer'],
          excludeSelectors: ['.dynamic-price', '.stock-counter', '.timestamp']
        },
        content: {
          enabled: true,
          textChanges: true,
          imageChanges: true,
          linkChanges: true // Important for e-commerce
        }
      }
    },

    element: {
      elementTypes: {
        navigation: {
          enabled: true,
          consistency: true,
          accessibility: true,
          interactions: true
        },
        forms: {
          enabled: true,
          validation: true,
          styling: true,
          behavior: true // Critical for checkout forms
        },
        media: {
          enabled: true,
          images: true,
          videos: false,
          loading: true // Product image loading is critical
        }
      }
    }
  }
};

// Configuration utility functions
export function getDetectionConfig(
  environment: 'development' | 'staging' | 'production',
  detectorType: 'anomaly' | 'regression' | 'layout' | 'element',
  siteName?: string
): DetectorConfig {
  let config = DetectionConfigs[environment][detectorType];

  // Apply site-specific overrides
  if (siteName && SiteDetectionConfigs[siteName]) {
    const siteConfig = SiteDetectionConfigs[siteName][detectorType];
    if (siteConfig) {
      config = deepMerge(config, siteConfig);
    }
  }

  return config;
}

export function createCustomDetectionConfig(
  baseConfig: DetectorConfig,
  overrides: Partial<DetectorConfig>
): DetectorConfig {
  return deepMerge(baseConfig, overrides);
}

export function validateDetectionConfig(config: DetectorConfig): boolean {
  // Validate threshold ranges
  const thresholds = config.thresholds;
  if (thresholds.confidence < 0 || thresholds.confidence > 1) return false;
  if (thresholds.similarity < 0 || thresholds.similarity > 1) return false;
  if (thresholds.anomaly < 0 || thresholds.anomaly > 1) return false;
  if (thresholds.layout < 0) return false;

  // Validate performance settings
  if (config.performance.maxConcurrent <= 0) return false;
  if (config.performance.timeout <= 0) return false;
  if (config.performance.batchSize <= 0) return false;

  return true;
}

export function optimizeConfigForPerformance(config: DetectorConfig): DetectorConfig {
  return {
    ...config,
    preprocessing: {
      ...config.preprocessing,
      denoise: false, // Disable expensive operations
      enhance: false
    },
    postprocessing: {
      ...config.postprocessing,
      spatialFiltering: false,
      temporalFiltering: false
    },
    performance: {
      ...config.performance,
      maxConcurrent: Math.min(config.performance.maxConcurrent * 2, 8),
      cacheResults: true,
      batchSize: Math.min(config.performance.batchSize * 2, 16)
    }
  };
}

export function optimizeConfigForAccuracy(config: DetectorConfig): DetectorConfig {
  return {
    ...config,
    thresholds: {
      confidence: Math.min(config.thresholds.confidence + 0.05, 0.98),
      similarity: Math.min(config.thresholds.similarity + 0.02, 0.99),
      anomaly: Math.max(config.thresholds.anomaly - 0.05, 0.01),
      layout: Math.max(config.thresholds.layout - 0.02, 0.01)
    },
    preprocessing: {
      ...config.preprocessing,
      denoise: true,
      enhance: true
    },
    postprocessing: {
      ...config.postprocessing,
      filterLowConfidence: true,
      mergeSimilarAnomalies: true,
      spatialFiltering: true
    }
  };
}

// Helper function for deep merging configurations
function deepMerge(target: any, source: any): any {
  const result = { ...target };

  for (const key in source) {
    if (source.hasOwnProperty(key)) {
      if (typeof source[key] === 'object' && source[key] !== null && !Array.isArray(source[key])) {
        result[key] = deepMerge(target[key] || {}, source[key]);
      } else {
        result[key] = source[key];
      }
    }
  }

  return result;
}