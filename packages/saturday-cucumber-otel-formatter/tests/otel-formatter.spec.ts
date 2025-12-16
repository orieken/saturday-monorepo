// Mock OTel dependencies aggressively to prevent import crashes
// ... (keep exisitng mocks)

jest.mock('@cucumber/cucumber', () => ({
  SummaryFormatter: class {
    eventDataCollector: any;
    constructor(options: any) {
      if (options.eventBroadcaster) {
         // Minimal mock of base class subscription
      }
      this.eventDataCollector = options.eventDataCollector;
    }
  },
  Formatter: class {
    eventDataCollector: any;
    constructor(options: any) {
       this.eventDataCollector = options.eventDataCollector;
    }
  },
  formatterHelpers: {
    GherkinDocumentParser: {
      getGherkinScenarioMap: jest.fn().mockReturnValue({
        '1': { location: { line: 10 } }
      }),
      getGherkinStepMap: jest.fn().mockReturnValue({
        '1': { keyword: 'Given ', location: { line: 11 } }
      })
    },
    PickleParser: {
      getPickleStepMap: jest.fn().mockReturnValue({
        'pickle-step-1': { astNodeIds: ['1'], text: 'step text' }
      })
    }
  }
}));

jest.mock('@opentelemetry/sdk-trace-node', () => ({
  NodeTracerProvider: class {
    register() {}
    addSpanProcessor() {}
  },
}));
jest.mock('@opentelemetry/resources', () => ({
  Resource: class {},
}));
jest.mock('@opentelemetry/semantic-conventions', () => ({
  SemanticResourceAttributes: {},
}));
jest.mock('@opentelemetry/sdk-trace-base', () => ({
  SimpleSpanProcessor: class {},
}));
jest.mock('@opentelemetry/exporter-trace-otlp-http', () => ({
  OTLPTraceExporter: class {},
}));
jest.mock('@opentelemetry/api', () => ({
  trace: {
    getSpan: jest.fn().mockReturnValue({}),
    getTracer: jest.fn().mockReturnValue({}),
  },
  context: {
    active: jest.fn(),
  },
  SpanStatusCode: {
    OK: 1,
    ERROR: 2
  },
  createContextKey: jest.fn(),
}));

jest.mock('@opentelemetry/sdk-metrics', () => ({
  MeterProvider: class {
    getMeter() {
      return {
        createCounter: jest.fn().mockReturnValue({
          add: jest.fn(),
        }),
      };
    }
    shutdown() {}
  },
  PeriodicExportingMetricReader: class {},
}));

jest.mock('@opentelemetry/exporter-metrics-otlp-http', () => ({
  OTLPMetricExporter: class {},
}));

import OtelFormatter from '../src/index';
import { EventEmitter } from 'events';
import { SpanManager } from '../src/span-manager';
import { TracerSetup } from '../src/tracer-setup';
import { MetricsSetup } from '../src/metrics-setup';

jest.mock('../src/tracer-setup');
jest.mock('../src/metrics-setup');
jest.mock('../src/span-manager');

describe('OtelFormatter', () => {
  let formatter: OtelFormatter;
  let eventBroadcaster: EventEmitter;
  let mockTracerSetup: any;
  let mockMetricsSetup: any;
  let mockSpanManager: any;


  beforeEach(() => {
    eventBroadcaster = new EventEmitter();
    
    // Setup mocks
    mockTracerSetup = {
      getTracer: jest.fn().mockReturnValue({}),
      getResource: jest.fn().mockReturnValue({}),
    };
    (TracerSetup as jest.Mock).mockImplementation(() => mockTracerSetup);

    mockMetricsSetup = {
        getMeter: jest.fn().mockReturnValue({
            createCounter: jest.fn().mockReturnValue({
                add: jest.fn()
            })
        }),
        shutdown: jest.fn()
    };
    (MetricsSetup as jest.Mock).mockImplementation(() => mockMetricsSetup);

    mockSpanManager = {
      startTestRun: jest.fn(),
      endTestRun: jest.fn(),
      startScenario: jest.fn().mockReturnValue({ span: {}, ctx: {} }),
      startStep: jest.fn().mockReturnValue({ span: {}, ctx: {} }),
      endSpan: jest.fn(),
    };
    (SpanManager as jest.Mock).mockImplementation(() => mockSpanManager);

    formatter = new OtelFormatter({
      eventBroadcaster,
      eventDataCollector: {
        getTestCaseAttempt: jest.fn().mockReturnValue({
          gherkinDocument: { feature: { uri: 'feature.feature' }, uri: 'feature.feature' },
          pickle: { astNodeIds: ['1'], name: 'scenario', tags: [] },
          testCase: {
             testSteps: [{ id: 'step-1', pickleStepId: 'pickle-step-1' }]
          }
        }),
        getTestCaseAttempts: jest.fn().mockReturnValue([]),
      } as any,
      supportCodeLibrary: {} as any,
      colorFns: {} as any,
      cwd: '',
      log: jest.fn(),
      parsedArgvOptions: {},
      stream: {} as any,
      cleanup: jest.fn(),
      snippetBuilder: {} as any
    });

    // Prevent SummaryFormatter from crashing on logSummary
    (formatter as any).logSummary = jest.fn();
  });

  it('should initialize tracer and span manager', () => {
    expect(TracerSetup).toHaveBeenCalled();
    expect(SpanManager).toHaveBeenCalled();
  });

  // Basic event flow test
  it('should handle test run start', () => {
    eventBroadcaster.emit('envelope', {
      testRunStarted: { timestamp: { seconds: 1, nanos: 0 } }
    });
    expect(mockSpanManager.startTestRun).toHaveBeenCalled();
  });

  it('should handle scenario start', () => {
    eventBroadcaster.emit('envelope', {
      testCaseStarted: { 
        id: 'tc-1', 
        testCaseId: 'tcid-1',
        attempt: 0,
        timestamp: { seconds: 1, nanos: 0 }
      }
    });

    expect(mockSpanManager.startScenario).toHaveBeenCalledWith('scenario', expect.any(Object));
  });

  it('should handle step start', () => {
    eventBroadcaster.emit('envelope', {
      testCaseStarted: { id: 'tc-1', testCaseId: 'tcid-1', attempt: 0, timestamp: { seconds: 1, nanos: 0 } }
    });

    eventBroadcaster.emit('envelope', {
      testStepStarted: { 
        testCaseStartedId: 'tc-1',
        testStepId: 'step-1',
        timestamp: { seconds: 2, nanos: 0 }
      }
    });

    expect(mockSpanManager.startStep).toHaveBeenCalledWith('Given step text', expect.any(Object), expect.any(Object));
  });

  it('should handle step finish', () => {
    eventBroadcaster.emit('envelope', {
      testCaseStarted: { id: 'tc-1', testCaseId: 'tcid-1', attempt: 0, timestamp: { seconds: 1, nanos: 0 } }
    });
    eventBroadcaster.emit('envelope', {
      testStepStarted: { testCaseStartedId: 'tc-1', testStepId: 'step-1', timestamp: { seconds: 2, nanos: 0 } }
    });

    eventBroadcaster.emit('envelope', {
      testStepFinished: { 
        testCaseStartedId: 'tc-1', 
        testStepId: 'step-1', 
        testStepResult: { status: 'PASSED', duration: { seconds: 0, nanos: 1000 } },
        timestamp: { seconds: 3, nanos: 0 }
      }
    });

    expect(mockSpanManager.endSpan).toHaveBeenCalled();
  });

  it('should handle failed step', () => {
    // Setup state
    eventBroadcaster.emit('envelope', {
      testCaseStarted: { id: 'tc-fail', testCaseId: 'tcid-fail', attempt: 0, timestamp: { seconds: 1, nanos: 0 } }
    });
    eventBroadcaster.emit('envelope', {
      testStepStarted: { testCaseStartedId: 'tc-fail', testStepId: 'step-1', timestamp: { seconds: 2, nanos: 0 } }
    });

    eventBroadcaster.emit('envelope', {
      testStepFinished: { 
        testCaseStartedId: 'tc-fail', 
        testStepId: 'step-1', 
        testStepResult: { status: 'FAILED', message: 'Error msg', duration: { seconds: 0, nanos: 1000 } },
        timestamp: { seconds: 3, nanos: 0 }
      }
    });

    expect(mockSpanManager.endSpan).toHaveBeenCalledWith(
        expect.anything(), 
        'error', 
        'Error msg', 
        expect.any(Number)
    );
  });

  it('should handle scenario finish', () => {
     eventBroadcaster.emit('envelope', {
      testCaseStarted: { id: 'tc-1', testCaseId: 'tcid-1', attempt: 0, timestamp: { seconds: 1, nanos: 0 } }
    });
    eventBroadcaster.emit('envelope', {
        testCaseFinished: { testCaseStartedId: 'tc-1', timestamp: { seconds: 4, nanos: 0 } }
    });
    expect(mockSpanManager.endSpan).toHaveBeenCalled();
  });

  it('should handle test run finish', () => {
    eventBroadcaster.emit('envelope', {
      testRunStarted: { timestamp: { seconds: 1, nanos: 0 } }
    });
    eventBroadcaster.emit('envelope', {
      testRunFinished: { success: true, timestamp: { seconds: 5, nanos: 0 } }
    });
    expect(mockSpanManager.endTestRun).toHaveBeenCalledWith(true, 5000);
  });

  // More detailed tests would require mocking formatterHelpers which is complex as it's part of Cucumber internals
  // validating that the logic calls SpanManager is the main goal here.
});
