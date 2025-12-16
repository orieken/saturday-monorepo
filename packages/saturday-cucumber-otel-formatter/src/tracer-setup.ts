import * as os from 'os';
import { NodeTracerProvider } from '@opentelemetry/sdk-trace-node';
import { Resource, resourceFromAttributes } from '@opentelemetry/resources';
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

    const spanProcessors: SimpleSpanProcessor[] = [];

    const otlpExporter = new OTLPTraceExporter({
      url: process.env.OTEL_EXPORTER_OTLP_ENDPOINT,
      headers: {},
      timeoutMillis: parseInt(process.env.OTEL_EXPORTER_OTLP_TIMEOUT || '15000'),
    });
    spanProcessors.push(new SimpleSpanProcessor(otlpExporter));

    if (process.env.OTEL_SAVE_PAYLOADS === 'true') {
      const { FileSpanExporter } = require('./file-exporter');
      console.log('OTEL: Enabling FileSpanExporter');
      const timestamp = new Date().toISOString().replace(/:/g, '-').replace(/\..+/, '');
      spanProcessors.push(new SimpleSpanProcessor(new FileSpanExporter('./reports', `otel-cucumber-spans-${timestamp}.json`)));
    }

    this.provider = new NodeTracerProvider({
      resource: resourceFromAttributes({
        [SemanticResourceAttributes.SERVICE_NAME]: finalServiceName,
        [SemanticResourceAttributes.SERVICE_VERSION]: process.env.OTEL_SERVICE_VERSION,
        [SemanticResourceAttributes.SERVICE_NAMESPACE]: process.env.OTEL_SERVICE_NAMESPACE,
        [SemanticResourceAttributes.DEPLOYMENT_ENVIRONMENT]: process.env.NODE_ENV || process.env.OTEL_DEFAULT_ENVIRONMENT,
        'test.environment': process.env.NODE_ENV || process.env.OTEL_DEFAULT_ENVIRONMENT || 'development',
        'test.platform': process.platform,
        'test.os': os.type(),
        'test.browser': process.env.BROWSER || 'unknown',
        'test.reporter': 'cucumber',
        ...customAttributes
      }),
      spanProcessors: spanProcessors,
    });

    this.provider.register();

    this.tracer = trace.getTracer(finalServiceName);
  }

  getTracer(): Tracer {
    return this.tracer;
  }

  getResource(): Resource {
     return (this.provider as any)?.resource || resourceFromAttributes({});
  }

  async shutdown(): Promise<void> {
    if (this.provider) {
      await this.provider.shutdown();
    }
  }
}
