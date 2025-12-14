
import { K6Recorder, createK6Recorder, ExporterOptions } from '../src/index';
import * as fs from 'fs';
import * as path from 'path';
import { request } from '@playwright/test';

// Mock fs to avoid writing files
jest.mock('fs', () => ({
  ...jest.requireActual('fs'),
  writeFileSync: jest.fn(),
  mkdirSync: jest.fn(),
  existsSync: jest.fn().mockReturnValue(true),
  readFileSync: jest.fn().mockReturnValue('EXISTING_VAR=val'),
}));

jest.mock('@playwright/test', () => ({
    request: {
        newContext: jest.fn().mockResolvedValue({})
    }
}));

describe('K6Recorder Coverage', () => {
    const outDir = path.join(__dirname, 'output');

    beforeEach(() => {
        jest.clearAllMocks();
        process.env.K6_EXPORT = 'true';
    });
    
    afterEach(() => {
        delete process.env.K6_EXPORT;
        delete process.env.K6_REDACTION_POLICY;
    });

    test('manual record method adds calls', () => {
        const recorder = new K6Recorder(outDir, 'test-slug');
        recorder.record('test-name', 'GET', 'http://example.com');
        expect(recorder.hasCalls()).toBe(true);
    });

    test('resolvePolicy loads default policy from package if env var missing', async () => {
        const res = await createK6Recorder('test title', outDir);
        expect(res).not.toBeNull();
        expect(res?.recorder).toBeDefined();
    });

    test('wrap handles non-function properties', () => {
        const recorder = new K6Recorder(outDir, 'test-slug');
        const mockCtx: any = {
            someProp: 'value',
            post: jest.fn()
        };
        const wrapped = recorder.wrap(mockCtx);
        expect((wrapped as any).someProp).toBe('value');
    });

    test('wrap handles non-verb function properties', () => {
        const recorder = new K6Recorder(outDir, 'test-slug');
        const mockFn = jest.fn();
        const mockCtx: any = {
            otherFn: mockFn
        };
        const wrapped = recorder.wrap(mockCtx);
        (wrapped as any).otherFn('arg');
        expect(mockFn).toHaveBeenCalledWith('arg');
    });
});
