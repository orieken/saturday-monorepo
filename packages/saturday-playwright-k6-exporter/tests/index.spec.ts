import { createK6Recorder, K6Recorder, makeSlugFromTitle } from '../src/index';
import { request } from '@playwright/test';
import fs from 'node:fs';
import path from 'node:path';

// Mock fs
jest.mock('node:fs');
jest.mock('node:path', () => {
    const original = jest.requireActual('node:path');
    return {
        ...original,
        join: (...args: string[]) => args.join('/'),
    };
});


describe('Playwright K6 Exporter', () => {
  const mockOutDir = '/tmp/k6-out';

  beforeEach(() => {
    jest.clearAllMocks();
    (fs.existsSync as jest.Mock).mockReturnValue(false);
  });

  describe('makeSlugFromTitle', () => {
    it('should sanitize titles', () => {
      expect(makeSlugFromTitle('My Test Case')).toBe('my-test-case');
      expect(makeSlugFromTitle('Test @tag 123')).toBe('test-tag-123');
    });
  });

  describe('createK6Recorder', () => {
    it('should return null if K6_EXPORT is not set', async () => {
      delete process.env.K6_EXPORT;
      const result = await createK6Recorder('test');
      expect(result).toBeNull();
    });

    it('should return recorder and context if K6_EXPORT is set', async () => {
      process.env.K6_EXPORT = '1';
      (request.newContext as jest.Mock).mockResolvedValue({});
      const result = await createK6Recorder('test', mockOutDir);
      expect(result).not.toBeNull();
      expect(result?.recorder).toBeInstanceOf(K6Recorder);
    });
  });

  describe('K6Recorder', () => {
    let recorder: K6Recorder;
    let mockContext: any;

    beforeEach(() => {
        recorder = new K6Recorder(mockOutDir, 'test-slug');
        mockContext = {
            get: jest.fn(),
            post: jest.fn(),
            fetch: jest.fn(),
        }
    });

    it('should wrap context and capture requests', async () => {
        const wrapped = recorder.wrap(mockContext);
        
        // Mock response
        const mockResponse = { 
            status: () => 200, 
            headers: () => ({ 'content-type': 'application/json' }) 
        };
        mockContext.get.mockResolvedValue(mockResponse);

        await wrapped.get('https://example.com/api', { headers: { 'X-Foo': 'bar' } });

        expect(mockContext.get).toHaveBeenCalledWith('https://example.com/api', expect.anything());
        expect(recorder.hasCalls()).toBe(true);
    });

    it('should generate k6 script', async () => {
         const wrapped = recorder.wrap(mockContext);
         mockContext.post.mockResolvedValue({ 
            status: () => 201, 
            headers: () => ({}) 
        });

        await wrapped.post('https://api.com/users', { 
            data: { name: 'John' },
            headers: { 'Authorization': 'Bearer secret' }
        });

        const filePath = await recorder.flushToK6();
        
        expect(fs.mkdirSync).toHaveBeenCalledWith(mockOutDir, expect.objectContaining({ recursive: true }));
        expect(fs.writeFileSync).toHaveBeenCalledWith(
            '/tmp/k6-out/test-slug.k6.js', 
            expect.stringContaining('http.request(\'POST\', "https://api.com/users"'),
            'utf-8'
        );
    });

    it('should redact secrets if policy provided', async () => {
        const policy = {
            redactHeader: jest.fn().mockReturnValue({ value: '${SECRET_HEADER}', finding: { envName: 'SECRET', value: 'secret' } }),
            redactBody: jest.fn()
        };
        const recorderWithPolicy = new K6Recorder(mockOutDir, 'test', { policy: policy as any });
        const wrapped = recorderWithPolicy.wrap(mockContext);

        mockContext.get.mockResolvedValue({ status: () => 200, headers: () => ({}) });

        await wrapped.get('http://foo.com', { headers: { 'Auth': 'bad-secret' } });
        
        await recorderWithPolicy. flushToK6();
        expect(fs.writeFileSync).toHaveBeenCalledWith(
            expect.stringContaining('.env.apis'),
            expect.stringContaining('SECRET=secret'),
            expect.anything()
        );
    });

    it('should redact body content recursively', async () => {
         const policy = {
            redactBody: jest.fn().mockImplementation((path, val) => {
                if (val === 'secret-val') return { value: '${SECRET_BODY}', finding: { envName: 'BODY_KEY', value: 'secret-val' } };
                return null;
            })
        };
        const recorderWithPolicy = new K6Recorder(mockOutDir, 'test-body', { policy: policy as any });
        const wrapped = recorderWithPolicy.wrap(mockContext);

        mockContext.post.mockResolvedValue({ status: () => 200, headers: () => ({}) });

        await wrapped.post('http://api.com', { 
            data: { 
                user: { 
                    name: 'John', 
                    details: { secret: 'secret-val' } 
                },
                items: [{ id: 1 }, { token: 'secret-val' }]
            } 
        });
        
        await recorderWithPolicy.flushToK6();

        expect(policy.redactBody).toHaveBeenCalled();
        expect(fs.writeFileSync).toHaveBeenCalledWith(
             expect.stringContaining('test-body.k6.js'),
             expect.stringContaining('${SECRET_BODY}'),
             'utf-8'
        );
        // Check env file for findings
        expect(fs.writeFileSync).toHaveBeenCalledWith(
             expect.stringContaining('.env.apis'),
             expect.stringContaining('BODY_KEY=secret-val'),
             expect.anything()
        );
    });
  });
});
