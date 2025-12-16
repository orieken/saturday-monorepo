import { Context, Tracer, Span, trace, context, SpanStatusCode } from '@opentelemetry/api';

export interface SpanAttributes {
  [key: string]: string | number | boolean;
}

export class SpanManager {
  private tracer: Tracer;
  private testRunSpan?: Span;
  private testRunContext?: Context;
  
  constructor(tracer: Tracer) {
    this.tracer = tracer;
  }

  startTestRun(startTime?: number): Context {
    const rootSpan = this.tracer.startSpan('test-run', {
      startTime,
      attributes: {
        'test.type': 'test-run',
        'test.status': 'started'
      }
    });
    this.testRunSpan = rootSpan;
    this.testRunContext = trace.setSpan(context.active(), rootSpan);
    return this.testRunContext;
  }

  endTestRun(success: boolean, endTime?: number) {
    if (this.testRunSpan) {
      this.testRunSpan.setAttributes({
        'test.status': success ? 'passed' : 'failed',
        'test.success': success
      });
      this.testRunSpan.end(endTime);
    }
  }

  startTest(name: string, attributes: SpanAttributes, parentContext?: Context): { span: Span, ctx: Context } {
    const ctx = parentContext || this.testRunContext || context.active();
    const span = this.tracer.startSpan(`test: ${name}`, { attributes }, ctx);
    const newCtx = trace.setSpan(ctx, span);
    return { span, ctx: newCtx };
  }

  startStep(name: string, attributes: SpanAttributes, parentContext: Context): { span: Span, ctx: Context } {
    const span = this.tracer.startSpan(`step: ${name}`, { attributes }, parentContext);
    const newCtx = trace.setSpan(parentContext, span);
    return { span, ctx: newCtx };
  }

  endSpan(span: Span, status: string, error?: string, endTime?: number) {
     span.setAttributes({
       'test.status': status,
       'test.error': error || ''
     });
     
     if (status === 'error' || status === 'failed' || status === 'timedOut') {
       span.setStatus({
         code: SpanStatusCode.ERROR,
         message: error
       });
     } else {
       span.setStatus({ code: SpanStatusCode.OK });
     }
     
     span.end(endTime);
  }
}
