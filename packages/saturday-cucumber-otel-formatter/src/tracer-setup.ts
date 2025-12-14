import * as os from 'os';
import { NodeTracerProvider } from '@opentelemetry/sdk-trace-node';
import { Resource } from '@opentelemetry/resources';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';
import { SimpleSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http';
import { Tracer, trace } from '@opentelemetry/api';

import { Attributes } from './config';

export class TracerSetup {
  private provider?: NodeTracerProvider;
  private tracer: Tracer;

  constructor(serviceName: string = 'cucumber-tests', customAttributes: Attributes = {}) {
    if (process.env.ENABLE_OTEL !== 'true') {
      this.tracer = trace.getTracer(serviceName);
      return;
    }

    const envServiceName = process.env.OTEL_SERVICE_NAME;
    const finalServiceName = envServiceName || serviceName;

    this.provider = new NodeTracerProvider({
      resource: new Resource({
        [SemanticResourceAttributes.SERVICE_NAME]: finalServiceName,
        [SemanticResourceAttributes.SERVICE_VERSION]: process.env.OTEL_SERVICE_VERSION,
        [SemanticResourceAttributes.SERVICE_NAMESPACE]: process.env.OTEL_SERVICE_NAMESPACE,
        [SemanticResourceAttributes.DEPLOYMENT_ENVIRONMENT]: process.env.NODE_ENV || process.env.OTEL_DEFAULT_ENVIRONMENT,
        'test.environment': process.env.NODE_ENV || process.env.OTEL_DEFAULT_ENVIRONMENT || 'development',
        'test.platform': process.platform,
        'test.os': os.type(),
        'test.browser': process.env.BROWSER || 'unknown',
        ...customAttributes
      }),
    });

    const otlpExporter = new OTLPTraceExporter({
      url: process.env.OTEL_EXPORTER_OTLP_ENDPOINT,
      headers: {},
      timeoutMillis: parseInt(process.env.OTEL_EXPORTER_OTLP_TIMEOUT || '15000'),
    });

    this.provider.addSpanProcessor(new SimpleSpanProcessor(otlpExporter));
    this.provider.register();

    this.tracer = trace.getTracer(finalServiceName);
  }

  getTracer(): Tracer {
    return this.tracer;
  }

  async shutdown(): Promise<void> {
    if (this.provider) {
      await this.provider.shutdown();
    }
  }
}
