jest.mock('@opentelemetry/api', () => ({
  trace: {
    setSpan: jest.fn((ctx, span) => ctx),
    getTracer: jest.fn(),
  },
  context: {
    active: jest.fn(() => ({})),
  },
  SpanStatusCode: {
    UNSET: 0,
    OK: 1,
    ERROR: 2,
  },
}));

import { SpanManager } from '../src/span-manager';
import { SpanStatusCode } from '@opentelemetry/api';

describe('SpanManager', () => {
  let mockSpan: any;
  let mockTracer: any;
  let spanManager: SpanManager;

  beforeEach(() => {
    mockSpan = {
      setAttributes: jest.fn(),
      setStatus: jest.fn(),
      end: jest.fn(),
    };
    mockTracer = {
      startSpan: jest.fn().mockReturnValue(mockSpan),
    };
    spanManager = new SpanManager(mockTracer);
  });

  it('should start and end test run', () => {
    spanManager.startTestRun(123456);
    expect(mockTracer.startSpan).toHaveBeenCalledWith('test-run', expect.objectContaining({ startTime: 123456 }));
    
    spanManager.endTestRun(false, 123457);
    expect(mockSpan.setAttributes).toHaveBeenCalledWith(expect.objectContaining({ 'test.success': false }));
    expect(mockSpan.end).toHaveBeenCalledWith(123457);
  });

  it('should start scenario', () => {
    const parentCtx = {};
    const result = spanManager.startScenario('scenario 1', { foo: 'bar' }, parentCtx as any);
    expect(mockTracer.startSpan).toHaveBeenCalledWith('scenario: scenario 1', { attributes: { foo: 'bar' } }, parentCtx);
    expect(result.span).toBe(mockSpan);
    expect(result.ctx).toBeTruthy();
  });

  it('should start step', () => {
    const parentCtx = {};
    const result = spanManager.startStep('step 1', { baz: 'qux' }, parentCtx as any);
    expect(mockTracer.startSpan).toHaveBeenCalledWith('step: step 1', { attributes: { baz: 'qux' } }, parentCtx);
    expect(result.span).toBe(mockSpan);
    expect(result.ctx).toBeTruthy();
  });
  it('should end span correctly', () => {
    spanManager.endSpan(mockSpan, 'ok', undefined, 100);
    expect(mockSpan.setStatus).toHaveBeenCalledWith({ code: SpanStatusCode.OK });
    expect(mockSpan.end).toHaveBeenCalledWith(100);
  });

  it('should end span with error', () => {
    spanManager.endSpan(mockSpan, 'error', 'failed', 100);
    expect(mockSpan.setStatus).toHaveBeenCalledWith({ code: SpanStatusCode.ERROR, message: 'failed' });
  });
});
