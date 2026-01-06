
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ConsoleLogger } from '../../src/utils/console-logger';

describe('ConsoleLogger Coverage', () => {
    let mockPage: any;
    let logger: ConsoleLogger;
    let consoleListener: (msg: any) => void;

    beforeEach(() => {
        mockPage = {
            on: vi.fn((event, callback) => {
                if (event === 'console') consoleListener = callback;
            })
        };
        logger = new ConsoleLogger(mockPage);
    });

    it('should capture console logs', () => {
        const msg = {
            type: () => 'log',
            text: () => 'hello',
            location: () => ({ url: 'url' })
        };
        consoleListener(msg);

        expect(logger.getLogs()).toHaveLength(1);
        expect(logger.getLogs()[0].text).toBe('hello');
    });

    it('should formatting logs', () => {
        const msg = {
            type: () => 'error',
            text: () => 'fail',
            location: () => ({ url: 'loc' })
        };
        consoleListener(msg);

        const formatted = logger.getFormattedLogs();
        expect(formatted).toContain('ERROR: fail (loc)');
    });

    it('should format logs without location', () => {
         const msg = {
            type: () => 'warn',
            text: () => 'warn',
            location: () => ({ url: undefined })
        };
        consoleListener(msg);
        expect(logger.getFormattedLogs()).not.toContain('()'); 
    });

    it('should clear logs', () => {
        const msg = { type: () => 'log', text: () => 't', location: () => ({}) };
        consoleListener(msg);
        expect(logger.getLogs()).toHaveLength(1);
        
        logger.clear();
        expect(logger.getLogs()).toHaveLength(0);
    });
});
