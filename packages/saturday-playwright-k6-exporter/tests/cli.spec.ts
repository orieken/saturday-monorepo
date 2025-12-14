import fs from 'node:fs';
import path from 'node:path';

// Mock fs and path
jest.mock('node:fs');
jest.mock('node:path', () => {
    return {
        join: (...args: string[]) => args.join('/'),
        dirname: (dir: string) => dir,
    };
});

describe('k6-exporter-init CLI', () => {
    let consoleSpy: jest.SpyInstance;

    beforeEach(() => {
        jest.clearAllMocks();
        consoleSpy = jest.spyOn(console, 'log').mockImplementation(() => {});
        // Mock process.cwd to return a fixed path
         jest.spyOn(process, 'cwd').mockReturnValue('/root');
    });

    afterEach(() => {
        consoleSpy.mockRestore();
        jest.restoreAllMocks();
    });

    it('should scaffold example test file', async () => {
        // Trigger CLI execution by importing it. 
        // Note: CLI executes on import, so we need to isolate modules or use require + jest.resetModules
        
        jest.isolateModules(() => {
            require('../src/cli');
        });

        const expectedDir = '/root/e2e/api';
        const expectedFile = '/root/e2e/api/k6-export-example.spec.ts';

        expect(fs.mkdirSync).toHaveBeenCalledWith(expectedDir, { recursive: true });
        // It tries to write the file, first check if exists
        expect(fs.existsSync).toHaveBeenCalledWith(expectedFile);
        // It should write if not exists. We mocked existsSync to return false (default mock behavior for bool?)
        // jest.mock returns undefined by default. undefined is falsy.
        expect(fs.writeFileSync).toHaveBeenCalledWith(
            expectedFile, 
            expect.stringContaining('@orieken/saturday-playwright-k6-exporter'), 
            'utf-8'
        );
    });

    it('should update package.json provided it exists', () => {
        (fs.existsSync as jest.Mock).mockImplementation((path) => {
            if (path === '/root/package.json') return true;
            return false;
        });
        (fs.readFileSync as jest.Mock).mockReturnValue('{"scripts":{}}');

        jest.isolateModules(() => {
            require('../src/cli');
        });

        expect(fs.writeFileSync).toHaveBeenCalledWith(
            '/root/package.json',
            expect.stringContaining('"k6:export": "K6_EXPORT=1 playwright test -g \\"@k6\\""')
        );
        // Check CLI.ts line 49: fs.writeFileSync(pkgPath, JSON.stringify(pkg, null, 2)); -> NO 3rd arg in source code!
    });
    
    it('should not update package.json if it does not exist', () => {
        (fs.existsSync as jest.Mock).mockReturnValue(false);

        jest.isolateModules(() => {
             require('../src/cli');
        });

        expect(fs.writeFileSync).not.toHaveBeenCalledWith(
            '/root/package.json',
            expect.any(String)
        );
    });
});
