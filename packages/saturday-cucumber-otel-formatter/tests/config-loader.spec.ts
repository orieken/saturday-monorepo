
import { loadConfig } from '../src/config-loader';
import * as fs from 'fs';
import * as path from 'path';

jest.mock('fs');
jest.mock('path', () => ({
    ...jest.requireActual('path'),
    join: (...args: string[]) => args.join('/'),
    resolve: (...args: string[]) => args.join('/')
}));


describe('Config Loader', () => {
    beforeEach(() => {
        jest.clearAllMocks();
    });

    it('should return default config if path is undefined', async () => {
        (fs.existsSync as jest.Mock).mockReturnValue(false);
        const config = await loadConfig();
        expect(config).toBeDefined();
    });

    it('should return default config if loading fails', async () => {
        // Suppress console.error for this test
        const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});
        (fs.existsSync as jest.Mock).mockReturnValue(true);
        // import() will fail for non-existent file anyway, triggering catch block
        
        const config = await loadConfig('./non-existent-config.mjs');
        expect(config).toBeDefined();
        expect(consoleSpy).toHaveBeenCalled();
        
        consoleSpy.mockRestore();
    });
});
