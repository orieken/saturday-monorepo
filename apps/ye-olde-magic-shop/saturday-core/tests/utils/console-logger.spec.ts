
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ConsoleLogger } from '../../src/utils/console-logger';
import { Page, ConsoleMessage } from 'playwright';

describe('ConsoleLogger', () => {
  let mockPage: any;
  let consoleLogger: ConsoleLogger;
  let consoleListener: (msg: ConsoleMessage) => void;

  beforeEach(() => {
    mockPage = {
      on: vi.fn((event, callback) => {
        if (event === 'console') {
          consoleListener = callback;
        }
      }),
    };
    consoleLogger = new ConsoleLogger(mockPage as Page);
  });

  it('should register console listener on construction', () => {
    expect(mockPage.on).toHaveBeenCalledWith('console', expect.any(Function));
  });

  it('should collect logs from console messages', async () => {
    const mockMsg = {
      type: () => 'log',
      text: () => 'Test log message',
      location: () => ({ url: 'http://test.com' }),
    } as unknown as ConsoleMessage;

    consoleListener(mockMsg);

    const logs = consoleLogger.getLogs();
    expect(logs).toHaveLength(1);
    expect(logs[0].text).toBe('Test log message');
    expect(logs[0].type).toBe('log');
    expect(logs[0].location).toBe('http://test.com');
  });

  it('should format logs correctly', () => {
    const mockMsg = {
      type: () => 'error',
      text: () => 'Test error',
      location: () => ({ url: 'http://test.com' }),
    } as unknown as ConsoleMessage;

    consoleListener(mockMsg);

    const formatted = consoleLogger.getFormattedLogs();
    expect(formatted).toContain('ERROR: Test error (http://test.com)');
  });

  it('should format logs correctly without location', () => {
    const mockMsg = {
        type: () => 'info',
        text: () => 'Test info',
        location: () => ({ url: undefined }),
    } as unknown as ConsoleMessage;

    consoleListener(mockMsg);

    const formatted = consoleLogger.getFormattedLogs();
    expect(formatted).toContain('INFO: Test info');
    expect(formatted).not.toContain('()');
  });

  it('should clear logs', () => {
    const mockMsg = {
        type: () => 'log',
        text: () => 'Test',
        location: () => ({ url: '' }),
      } as unknown as ConsoleMessage;
  
      consoleListener(mockMsg);
      expect(consoleLogger.getLogs()).toHaveLength(1);
  
      consoleLogger.clear();
      expect(consoleLogger.getLogs()).toHaveLength(0);
  });
});
