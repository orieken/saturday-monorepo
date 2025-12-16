import { SpanManager } from '../src/span-manager';
import { SpanStatusCode } from '@opentelemetry/api';

describe('SpanManager', () => {
  let tracerMock: any;
  let spanMock: any;
  let spanManager: SpanManager;

  beforeEach(() => {
    spanMock = {
      end: jest.fn(),
      setAttributes: jest.fn(),
      setStatus: jest.fn(),
    };
    tracerMock = {
      startSpan: jest.fn().mockReturnValue(spanMock),
    };
    spanManager = new SpanManager(tracerMock);
  });

  describe('startTestRun', () => {
    it('should start a test-run span', () => {
      const startTime = Date.now();
      spanManager.startTestRun(startTime);
      expect(tracerMock.startSpan).toHaveBeenCalledWith('test-run', {
        startTime,
        attributes: {
          'test.type': 'test-run',
          'test.status': 'started',
        },
      });
    });
  });

  describe('endTestRun', () => {
    it('should end the test-run span with success status', () => {
      spanManager.startTestRun();
      spanManager.endTestRun(true);
      expect(spanMock.setAttributes).toHaveBeenCalledWith({
        'test.status': 'passed',
        'test.success': true,
      });
      expect(spanMock.end).toHaveBeenCalled();
    });

    it('should end the test-run span with failure status', () => {
      spanManager.startTestRun();
      spanManager.endTestRun(false);
      expect(spanMock.setAttributes).toHaveBeenCalledWith({
        'test.status': 'failed',
        'test.success': false,
      });
      expect(spanMock.end).toHaveBeenCalled();
    });
  });

  describe('endSpan', () => {
    it('should set status to OK for passed steps', () => {
        spanManager.endSpan(spanMock, 'passed', undefined);
        expect(spanMock.setStatus).toHaveBeenCalledWith({ code: SpanStatusCode.OK });
    });

    it('should set status to ERROR for failed steps', () => {
        spanManager.endSpan(spanMock, 'failed', 'some error');
        expect(spanMock.setStatus).toHaveBeenCalledWith({ code: SpanStatusCode.ERROR, message: 'some error' });
    });
  });
});
