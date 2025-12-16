import { Reporter, FullConfig, Suite, TestCase, TestResult, TestStep, FullResult } from '@playwright/test/reporter';
import * as fs from 'fs';
import { pathToFileURL } from 'url';
import { trace, Context, Counter } from '@opentelemetry/api';
import { TracerSetup } from './tracer-setup';
import { MetricsSetup } from './metrics-setup';
import { SpanManager } from './span-manager';
import { loadConfig } from './config-loader';
import { defaultConfig, OtelPlaywrightConfig } from './config';

export default class OtelReporter implements Reporter {
  private tracerSetup?: TracerSetup;
  private metricsSetup?: MetricsSetup;
  private spanManager?: SpanManager;
  private testCounter?: Counter;
  private config: OtelPlaywrightConfig = defaultConfig;
  private testContexts: Map<string, Context> = new Map();
  private stepContexts: WeakMap<TestStep, Context> = new WeakMap();

  async onBegin(config: FullConfig, suite: Suite) {
    if (process.env.OTEL_DEBUG_LOGGING === 'true') {
        console.log('OtelReporter: onBegin started');
    }
    
    // Load config
    this.config = await loadConfig(process.env.OTEL_CUSTOM_CONFIG);


    this.tracerSetup = new TracerSetup('playwright-tests', this.config.resourceAttributes);
    this.spanManager = new SpanManager(this.tracerSetup.getTracer());
    
    // Initialize Metrics
    this.metricsSetup = new MetricsSetup('playwright-tests', this.tracerSetup.getResource());
    const meter = this.metricsSetup.getMeter();
    this.testCounter = meter.createCounter('playwright.test.cases', {
      description: 'Counts the number of test cases executed',
    });
    
    this.spanManager.startTestRun(Date.now());
  }

  onTestBegin(test: TestCase, result: TestResult) {
    if (process.env.OTEL_DEBUG_LOGGING === 'true') {
        console.log(`OtelReporter: onTestBegin: ${test.title}`);
    }
    if (!this.spanManager) return;
    
    let customAttributes: Record<string, string | number | boolean> = {};
    if (this.config.testAttributes) {
        if (typeof this.config.testAttributes === 'function') {
            customAttributes = this.config.testAttributes(test);
        } else {
            customAttributes = this.config.testAttributes;
        }
    }

    const { ctx } = this.spanManager.startTest(test.title, {
        'test.type': 'test-case',
        'test.case.title': test.title,
        'test.case.file': test.location.file,
        'test.case.line': test.location.line,
        ...customAttributes
    });
    
    this.testContexts.set(test.id, ctx);
  }

  onStepBegin(test: TestCase, result: TestResult, step: TestStep) {
    if (process.env.OTEL_DEBUG_LOGGING === 'true') {
        console.log(`OtelReporter: onStepBegin: ${step.title}`);
    }
    if (!this.spanManager) return;
    
    // Determine parent context: either parent step's context or the test case context
    const parentCtx = step.parent ? this.stepContexts.get(step.parent) : this.testContexts.get(test.id);
    if (!parentCtx) {
        if (process.env.OTEL_DEBUG_LOGGING === 'true') {
            console.log(`OtelReporter: onStepBegin: No parent context found for step ${step.title}`);
        }
        return;
    }

    let customAttributes: Record<string, string | number | boolean> = {};
    if (this.config.stepAttributes) {
        if (typeof this.config.stepAttributes === 'function') {
            customAttributes = this.config.stepAttributes(step);
        } else {
            customAttributes = this.config.stepAttributes;
        }
    }
    
    const { ctx } = this.spanManager.startStep(step.title, {
        'test.type': 'step',
        'test.step.title': step.title,
        'test.step.category': step.category,
        ...customAttributes
    }, parentCtx);
    
    this.stepContexts.set(step, ctx);
  }

  onStepEnd(test: TestCase, result: TestResult, step: TestStep) {
      if (process.env.OTEL_DEBUG_LOGGING === 'true') {
          console.log(`OtelReporter: onStepEnd: ${step.title}`);
      }
      if (!this.spanManager) return;
      
      const ctx = this.stepContexts.get(step);
      if (!ctx) return;
      
      const span = trace.getSpan(ctx);
      if (span) {
          // Playwright doesn't always have explicit status on step objects until the end?
          // step.error might be populated.
          const status = step.error ? 'failed' : 'passed'; 
          this.spanManager.endSpan(span, status, step.error?.message, Date.now());
      }
  }

  onTestEnd(test: TestCase, result: TestResult) {
    if (process.env.OTEL_DEBUG_LOGGING === 'true') {
        console.log(`OtelReporter: onTestEnd: ${test.title}`);
    }
    if (!this.spanManager) return;
    const ctx = this.testContexts.get(test.id);
    if (!ctx) return;
    
    const span = trace.getSpan(ctx);
    if (span) {
        this.spanManager.endSpan(span, result.status, result.error?.message, Date.now());

        // Record Metric
        if (this.testCounter) {
            this.testCounter.add(1, {
                'test.status': result.status, // passed, failed, timedOut, skipped
                'test.browser': test.parent.project()?.use.browserName || 'unknown',
                'test.file': test.location.file || 'unknown'
            });
        }
    }
  }

  async onEnd(result: FullResult) {
    if (!this.spanManager || !this.tracerSetup) return;
    
    const success = result.status === 'passed';
    this.spanManager.endTestRun(success, Date.now());
    
    if (process.env.OTEL_DEBUG_LOGGING === 'true') {
        console.log('OtelReporter: Shutting down...');
    }

    if (this.metricsSetup) {
        await this.metricsSetup.shutdown();
    }
    await this.tracerSetup.shutdown();
  }
}
