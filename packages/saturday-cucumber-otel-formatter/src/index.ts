import {
  Formatter,
  IFormatterOptions,
  formatterHelpers,
} from '@cucumber/cucumber'
import * as messages from '@cucumber/messages'
import { EOL as n } from 'os'
import { trace } from '@opentelemetry/api';
import * as fs from 'fs';
import { TracerSetup } from './tracer-setup';
import { SpanManager } from './span-manager';
import { StatusMapper } from './status-mapper';
import { CurrentScenario, CurrentStep } from './types';
import { loadConfig } from './config-loader';
import { OtelCucumberConfig, defaultConfig } from './config';

const { GherkinDocumentParser, PickleParser } = formatterHelpers
const { getGherkinScenarioMap, getGherkinStepMap } = GherkinDocumentParser
const { getPickleStepMap } = PickleParser

export default class OtelFormatter extends Formatter {
  private tracerSetup?: TracerSetup;
  private spanManager?: SpanManager;
  private currentScenario?: CurrentScenario;
  private uri?: string;
  private config: OtelCucumberConfig = defaultConfig;
  private isReady: boolean = false;
  private envelopeQueue: messages.Envelope[] = [];
  private initPromise: Promise<void>;

  constructor(options: IFormatterOptions) {
    super(options)
    
    // Start loading config asynchronously
    this.initPromise = this.init();

    options.eventBroadcaster.on('envelope', this.handleEnvelope.bind(this))
  }

  async finished(): Promise<void> {
    await this.initPromise;
    console.log('OtelFormatter: Report processing completed.');
    
    if (this.tracerSetup) {
      console.log('OtelFormatter: Shutting down OTel...');
      try {
        await this.tracerSetup.shutdown();
        console.log('OtelFormatter: OTel shutdown completed.');
        
        // Debugging active handles to diagnose hangs
        if (process.env.OTEL_DEBUG_LOGGING === 'true') {
           console.log('OtelFormatter: Inspecting active handles...');
           // @ts-ignore
           const handles = process._getActiveHandles();
           console.log(`OtelFormatter: Found ${handles.length} active handles.`);
           handles.forEach((h: any) => {
               // Try to identify the handle
               const type = h.constructor.name;
               // Extract some info if possible
               let info = '';
               if (type === 'Socket') {
                   info = `address=${JSON.stringify(h.address())} remote=${h.remoteAddress}:${h.remotePort}`;
               } else if (type === 'Timer') {
                   info = `hasRef=${h.hasRef()}`;
               }
               console.log(` - Handle: ${type} ${info}`);
           });
        }

      } catch (e) {
        console.error('OtelFormatter: OTel shutdown failed', e);
      }
    }
  }

  private async init() {
    try {
        if (process.env.OTEL_DEBUG_LOGGING === 'true') {
            fs.appendFileSync('otel-debug-init.log', `[${new Date().toISOString()}] Init start. ENABLE_OTEL=${process.env.ENABLE_OTEL} CONFIG=${process.env.OTEL_CUSTOM_CONFIG}\n`);
        }
        
        this.config = await loadConfig(process.env.OTEL_CUSTOM_CONFIG);
        this.tracerSetup = new TracerSetup('cucumber-tests', this.config.resourceAttributes);
        this.spanManager = new SpanManager(this.tracerSetup.getTracer());
        this.isReady = true;
        
        // Process queued envelopes
        for (const envelope of this.envelopeQueue) {
            this.parseEnvelope(envelope);
        }
        this.envelopeQueue = [];
    } catch (e) {
        console.error('OtelFormatter init failed', e);
    }
  }

  private handleEnvelope(envelope: messages.Envelope): void {
    if (!this.isReady) {
      this.envelopeQueue.push(envelope);
    } else {
      this.parseEnvelope(envelope);
    }
  }

  private parseEnvelope(envelope: messages.Envelope): void {
    if (!this.spanManager) return;

    if (envelope.testRunStarted) this.onTestRunStarted(envelope.testRunStarted)
    if (envelope.testCaseStarted) this.onTestCaseStarted(envelope.testCaseStarted)
    if (envelope.testStepStarted) this.onTestStepStarted(envelope.testStepStarted)
    if (envelope.testStepFinished) this.onTestStepFinished(envelope.testStepFinished)
    if (envelope.testCaseFinished) this.onTestCaseFinished(envelope.testCaseFinished)
    if (envelope.testRunFinished) this.onTestRunFinished(envelope.testRunFinished)
  }

  private onTestRunStarted(testRunStarted: messages.TestRunStarted): void {
    const startTime = (testRunStarted.timestamp.seconds as unknown as number) * 1000 + (testRunStarted.timestamp.nanos / 1000000);
    this.spanManager!.startTestRun(startTime);
  }

  private onTestCaseStarted(testCaseStarted: messages.TestCaseStarted): void {
    const { gherkinDocument, pickle } = this.eventDataCollector.getTestCaseAttempt(testCaseStarted.id)

    if (!gherkinDocument.feature) return;

    const gherkinScenarioMap = getGherkinScenarioMap(gherkinDocument)
    const scenario = gherkinScenarioMap[pickle.astNodeIds[0]]
    
    this.currentScenario = {
      name: pickle.name,
      line: scenario.location.line,
      file: this.uri || gherkinDocument.uri || '',
      steps: [],
      hasFailedSteps: false
    };

    let customAttributes: Record<string, string | number | boolean> = {};
    if (this.config.scenarioAttributes) {
      if (typeof this.config.scenarioAttributes === 'function') {
        customAttributes = this.config.scenarioAttributes(pickle, gherkinDocument);
      } else {
        customAttributes = this.config.scenarioAttributes;
      }
    }

    const { ctx } = this.spanManager!.startScenario(this.currentScenario.name, {
      'test.type': 'scenario',
      'test.scenario.name': this.currentScenario.name,
      'test.scenario.line': this.currentScenario.line,
      'test.scenario.file': this.currentScenario.file,
      ...customAttributes
    });
    this.currentScenario.context = ctx;
  }

  private onTestStepStarted(testStepStarted: messages.TestStepStarted): void {
    if (!this.currentScenario) return;

    const { gherkinDocument, pickle, testCase } = this.eventDataCollector.getTestCaseAttempt(testStepStarted.testCaseStartedId)
    const pickleStepMap = getPickleStepMap(pickle)
    const gherkinStepMap = getGherkinStepMap(gherkinDocument)
    const testStep = testCase.testSteps.find((item) => item.id === testStepStarted.testStepId)

    if (testStep && testStep.pickleStepId) {
      const pickleStep = pickleStepMap[testStep.pickleStepId]
      const astNodeId = pickleStep.astNodeIds[0]
      const gherkinStep = gherkinStepMap[astNodeId]
      
      const step: CurrentStep = {
        name: `${gherkinStep.keyword}${pickleStep.text}`,
        keyword: gherkinStep.keyword,
        text: pickleStep.text,
        line: gherkinStep.location.line,
        file: this.uri || gherkinDocument.uri || ''
      };

      if (this.currentScenario.context) {
        let customAttributes: Record<string, string | number | boolean> = {};
        if (this.config.stepAttributes) {
          if (typeof this.config.stepAttributes === 'function') {
            customAttributes = this.config.stepAttributes(pickleStep, gherkinStep);
          } else {
            customAttributes = this.config.stepAttributes;
          }
        }

        const { ctx } = this.spanManager!.startStep(step.name, {
            'test.type': 'step',
            'test.step.name': step.name,
            'test.step.keyword': step.keyword,
            'test.step.text': step.text,
            'test.step.line': step.line,
            'test.step.file': step.file,
            ...customAttributes
        }, this.currentScenario.context);
        
        step.context = ctx;
      }
      
      this.currentScenario.steps.push(step);
    }
  }

  private onTestStepFinished(testStepFinished: messages.TestStepFinished): void {
    if (!this.currentScenario || this.currentScenario.steps.length === 0) return;

    const currentStep = this.currentScenario.steps[this.currentScenario.steps.length - 1];
    
    const status = StatusMapper.map(testStepFinished.testStepResult.status);
    const error = testStepFinished.testStepResult.message;
    const duration = testStepFinished.testStepResult.duration 
      ? (testStepFinished.testStepResult.duration.seconds as unknown as number) * 1000 + (testStepFinished.testStepResult.duration.nanos / 1000000)
      : undefined;

    if (status === 'error') {
      this.currentScenario.hasFailedSteps = true;
      this.currentScenario.error = error;
    }

    if (currentStep.context) {
      const span = trace.getSpan(currentStep.context);
      if (span) {
        this.spanManager!.endSpan(span, status, error, duration ? Date.now() : undefined); 
      }
    }
  }

  private onTestCaseFinished(testCaseFinished: messages.TestCaseFinished): void {
    if (this.currentScenario && this.currentScenario.context) {
      const span = trace.getSpan(this.currentScenario.context);
      if (span) {
        const status = this.currentScenario.hasFailedSteps ? 'error' : 'ok';
        this.spanManager!.endSpan(span, status, this.currentScenario.error);
      }
      this.currentScenario = undefined;
    }
  }

  private async onTestRunFinished(testRunFinished: messages.TestRunFinished): Promise<void> {
    const success = testRunFinished.success;
    const endTime = (testRunFinished.timestamp.seconds as unknown as number) * 1000 + (testRunFinished.timestamp.nanos / 1000000);
    this.spanManager!.endTestRun(success, endTime);
  }
}
