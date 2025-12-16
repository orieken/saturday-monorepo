import { MeterProvider, PeriodicExportingMetricReader } from '@opentelemetry/sdk-metrics';
import { Resource } from '@opentelemetry/resources';
import { OTLPMetricExporter } from '@opentelemetry/exporter-metrics-otlp-http';
import { Meter } from '@opentelemetry/api';

export class MetricsSetup {
  private provider?: MeterProvider;
  private meter: Meter;

  constructor(serviceName: string, resource: Resource) {
    if (process.env.ENABLE_OTEL !== 'true') {
      // Return a no-op meter if OTEL is disabled
      this.provider = new MeterProvider();
      this.meter = this.provider.getMeter(serviceName);
      return;
    }

    const metricExporter = new OTLPMetricExporter({
      // Default to standard OTLP metrics endpoint if not specified
      url: process.env.OTEL_EXPORTER_OTLP_METRICS_ENDPOINT || 
           (process.env.OTEL_EXPORTER_OTLP_ENDPOINT ? process.env.OTEL_EXPORTER_OTLP_ENDPOINT.replace('/v1/traces', '/v1/metrics') : undefined),
    });

    this.provider = new MeterProvider({
      resource: resource,
      readers: [
        new PeriodicExportingMetricReader({
          exporter: metricExporter,
          exportIntervalMillis: 1000, // Export every second for faster feedback in tests
        }),
      ],
    });

    this.meter = this.provider.getMeter(serviceName);
  }

  getMeter(): Meter {
    return this.meter;
  }

  async shutdown(): Promise<void> {
    if (this.provider) {
      await this.provider.shutdown();
    }
  }
}
