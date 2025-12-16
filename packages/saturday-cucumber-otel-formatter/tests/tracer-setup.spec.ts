// Mock OTel dependencies
jest.mock('@opentelemetry/sdk-trace-node', () => ({
  NodeTracerProvider: jest.fn().mockImplementation(() => ({
    register: jest.fn(),
  })),
}));

jest.mock('@opentelemetry/resources', () => ({
  Resource: jest.fn(),
  resourceFromAttributes: jest.fn(),
}));

jest.mock('@opentelemetry/semantic-conventions', () => ({
  SemanticResourceAttributes: {
    SERVICE_NAME: 'service.name',
    SERVICE_VERSION: 'service.version',
    SERVICE_NAMESPACE: 'service.namespace',
    DEPLOYMENT_ENVIRONMENT: 'deployment.environment',
  },
}));

jest.mock('@opentelemetry/sdk-trace-base', () => ({
  SimpleSpanProcessor: jest.fn(),
}));

jest.mock('@opentelemetry/exporter-trace-otlp-http', () => ({
  OTLPTraceExporter: jest.fn(),
}));

jest.mock('@opentelemetry/api', () => ({
  trace: {
    getTracer: jest.fn().mockReturnValue({}),
  },
}));

import { TracerSetup } from '../src/tracer-setup';
import { NodeTracerProvider } from '@opentelemetry/sdk-trace-node';
import { resourceFromAttributes } from '@opentelemetry/resources';
import { SimpleSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http';
import { trace } from '@opentelemetry/api';

describe('TracerSetup', () => {
  const originalEnv = process.env;

  beforeEach(() => {
    jest.clearAllMocks();
    process.env = { ...originalEnv };
    process.env.ENABLE_OTEL = 'true';
  });

  afterAll(() => {
    process.env = originalEnv;
  });

  it('should initialize with default values', () => {
    const tracerSetup = new TracerSetup();
    
    // Verify NodeTracerProvider created
    expect(NodeTracerProvider).toHaveBeenCalledTimes(1);
    expect(resourceFromAttributes).toHaveBeenCalledWith(expect.objectContaining({
      'service.name': 'cucumber-tests',
    }));

    // Verify Exporter created
    expect(OTLPTraceExporter).toHaveBeenCalledWith(expect.objectContaining({
        timeoutMillis: 15000
    }));

    // Verify SpanProcessor
    expect(SimpleSpanProcessor).toHaveBeenCalledTimes(1);

    // Verify Trace getTracer
    expect(trace.getTracer).toHaveBeenCalledWith('cucumber-tests');
    
    expect(tracerSetup.getTracer()).toBeDefined();
  });

  it('should use environment variables for configuration', () => {
    process.env.OTEL_SERVICE_NAME = 'env-service-name';
    process.env.OTEL_SERVICE_VERSION = '1.0.0';
    process.env.OTEL_SERVICE_NAMESPACE = 'my-namespace';
    process.env.NODE_ENV = 'production';
    process.env.OTEL_EXPORTER_OTLP_ENDPOINT = 'http://localhost:4318';
    process.env.OTEL_EXPORTER_OTLP_TIMEOUT = '5000';

    new TracerSetup();

    expect(resourceFromAttributes).toHaveBeenCalledWith(expect.objectContaining({
      'service.name': 'env-service-name',
      'service.version': '1.0.0',
      'service.namespace': 'my-namespace',
      'deployment.environment': 'production',
    }));

    expect(OTLPTraceExporter).toHaveBeenCalledWith(expect.objectContaining({
      url: 'http://localhost:4318',
      timeoutMillis: 5000,
    }));
  });

  it('should accept service name in constructor but prefer env var', () => {
    process.env.OTEL_SERVICE_NAME = 'env-overrides-ctor';
    new TracerSetup('ctor-service-name');
    
    expect(resourceFromAttributes).toHaveBeenCalledWith(expect.objectContaining({
      'service.name': 'env-overrides-ctor',
    }));
  });

  it('should use constructor argument if env var is missing', () => {
    delete process.env.OTEL_SERVICE_NAME;
    new TracerSetup('ctor-service-name');
    
    expect(resourceFromAttributes).toHaveBeenCalledWith(expect.objectContaining({
      'service.name': 'ctor-service-name',
    }));
  });
  
  it('should use OTEL_DEFAULT_ENVIRONMENT if NODE_ENV is missing', () => {
    delete process.env.NODE_ENV;
    process.env.OTEL_DEFAULT_ENVIRONMENT = 'staging';
    
    new TracerSetup();
    
    expect(resourceFromAttributes).toHaveBeenCalledWith(expect.objectContaining({
      'deployment.environment': 'staging',
    }));
  });

  it('should return no-op tracer if ENABLE_OTEL is not true', () => {
      process.env.ENABLE_OTEL = 'false';
      const setup = new TracerSetup();
      const tracer = setup.getTracer();
      expect(tracer).toBeDefined();
      expect(NodeTracerProvider).not.toHaveBeenCalled();
  });

  it('shutdown should be safe if provider not initialized', async () => {
      process.env.ENABLE_OTEL = 'false';
      const setup = new TracerSetup();
      await expect(setup.shutdown()).resolves.not.toThrow();
  });
});
