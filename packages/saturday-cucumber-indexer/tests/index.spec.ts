
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import * as fs from 'fs';
import * as path from 'path';
import { parseArgs, walkFeatures, parseFeatureFile } from '../src/index';

// Mock fs to avoid actual file system usage
vi.mock('fs', async () => {
    const actual = await vi.importActual('fs');
    return {
        ...actual as any,
        readdirSync: vi.fn(),
        statSync: vi.fn(),
        readFileSync: vi.fn(),
        existsSync: vi.fn(),
        mkdirSync: vi.fn(),
        writeFileSync: vi.fn()
    };
});

describe('cucumber-indexer', () => {

    describe('parseArgs', () => {
        it('should parse named arguments', () => {
            const argv = ['node', 'script', '--foo', 'bar', '--baz', 'true'];
            const args = parseArgs(argv);
            expect(args.foo).toBe('bar');
            expect(args.baz).toBe('true');
        });

        it('should parse flags as true', () => {
            const argv = ['node', 'script', '--flag', '--other', 'val'];
            const args = parseArgs(argv);
            expect(args.flag).toBe(true);
            expect(args.other).toBe('val');
        });
        
        it('should handle flag at end', () => {
             const argv = ['node', 'script', '--flag'];
             const args = parseArgs(argv);
             expect(args.flag).toBe(true);
        });

        it('should handle empty args', () => {
            const argv = ['node', 'script'];
            const args = parseArgs(argv);
            expect(Object.keys(args)).toHaveLength(0);
        });

        it('should skip undefined tokens', () => {
            const argv = ['node', 'script', undefined as any, '--foo', 'bar'];
            const args = parseArgs(argv);
            expect(args.foo).toBe('bar');
        });
    });

    describe('walkFeatures', () => {
        it('should walk recursively', () => {
             vi.mocked(fs.readdirSync).mockImplementation((p: any) => {
                 if (p === 'root') {
                     return [
                         { name: 'file1.feature', isFile: () => true, isDirectory: () => false },
                         { name: 'subdir', isFile: () => false, isDirectory: () => true }
                     ] as any;
                 }
                 if (p === 'root/subdir') {
                      return [
                         { name: 'file2.feature', isFile: () => true, isDirectory: () => false }
                      ] as any;
                 }
                 return [];
             });

             const files = walkFeatures('root');
             expect(files).toHaveLength(2);
             expect(files).toContain('root/file1.feature');
             expect(files).toContain('root/subdir/file2.feature');
        });

        it('should skip non-feature files', () => {
            vi.mocked(fs.readdirSync).mockReturnValue([
                { name: 'readme.md', isFile: () => true, isDirectory: () => false },
                { name: 'test.txt', isFile: () => true, isDirectory: () => false }
            ] as any);

            const files = walkFeatures('root');
            expect(files).toHaveLength(0);
        });
    });

    describe('parseFeatureFile', () => {
        it('should parse feature content', () => {
             const content = 
`@tag1
Feature: Login
  Description line 1
  Description line 2

  @sanity
  Scenario: Valid Login
    Given I open login page
    When I submit valid creds
    Then I am logged in
`;
             vi.mocked(fs.readFileSync).mockReturnValue(content);
             
             const result = parseFeatureFile('/path/to/login.feature', '/path/to');
             expect(result).not.toBeNull();
             expect(result?.name).toBe('Login');
             expect(result?.scenarios).toHaveLength(1);
             expect(result?.scenarios[0].name).toBe('Valid Login');
             expect(result?.scenarios[0].tags).toContain('@sanity');
             expect(result?.scenarios[0].steps).toHaveLength(3);
             expect(result?.description).toContain('Description line 1');
        });

        it('should handle Feature with no name', () => {
             vi.mocked(fs.readFileSync).mockReturnValue('');
             const result = parseFeatureFile('file', 'root');
             expect(result).toBeNull();
        });

        it('should handle Scenario Outline', () => {
            const content = 
`Feature: Test
  Scenario Outline: Test scenario
    Given step
`;
            vi.mocked(fs.readFileSync).mockReturnValue(content);
            const result = parseFeatureFile('test.feature', 'root');
            expect(result?.scenarios[0].name).toBe('Test scenario');
        });

        it('should handle multiple scenarios', () => {
            const content = 
`Feature: Multi
  Scenario: First
    Given step1
  
  @tag2
  Scenario: Second
    When step2
    Then step3
`;
            vi.mocked(fs.readFileSync).mockReturnValue(content);
            const result = parseFeatureFile('multi.feature', 'root');
            expect(result?.scenarios).toHaveLength(2);
            expect(result?.scenarios[0].name).toBe('First');
            expect(result?.scenarios[1].name).toBe('Second');
            expect(result?.scenarios[1].tags).toContain('@tag2');
        });

        it('should handle And and But steps', () => {
            const content = 
`Feature: Steps
  Scenario: All steps
    Given initial step
    And another given
    When action happens
    But not this
    Then result
`;
            vi.mocked(fs.readFileSync).mockReturnValue(content);
            const result = parseFeatureFile('steps.feature', 'root');
            expect(result?.scenarios[0].steps).toHaveLength(5);
            expect(result?.scenarios[0].steps[1].text).toBe('And another given');
            expect(result?.scenarios[0].steps[3].text).toBe('But not this');
        });

        it('should handle multiple tags on same line', () => {
            const content = 
`Feature: Tags
  @smoke @regression @critical
  Scenario: Tagged
    Given step
`;
            vi.mocked(fs.readFileSync).mockReturnValue(content);
            const result = parseFeatureFile('tags.feature', 'root');
            expect(result?.scenarios[0].tags).toHaveLength(3);
            expect(result?.scenarios[0].tags).toContain('@smoke');
            expect(result?.scenarios[0].tags).toContain('@regression');
            expect(result?.scenarios[0].tags).toContain('@critical');
        });

        it('should ignore empty lines and whitespace', () => {
            const content = 
`Feature: Whitespace

  
  Scenario: Test
  
    Given step
    
`;
            vi.mocked(fs.readFileSync).mockReturnValue(content);
            const result = parseFeatureFile('ws.feature', 'root');
            expect(result?.scenarios).toHaveLength(1);
        });

        it('should handle Windows line endings', () => {
            const content = "Feature: Windows\r\n  Scenario: Test\r\n    Given step\r\n";
            vi.mocked(fs.readFileSync).mockReturnValue(content);
            const result = parseFeatureFile('win.feature', 'root');
            expect(result?.scenarios).toHaveLength(1);
        });
    });
});
