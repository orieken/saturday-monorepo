import * as fs from 'fs';
import * as path from 'path';
import { parseArgs, parseFeatureFile, walkFeatures } from '../src/index';

jest.mock('fs');

describe('Cucumber Indexer', () => {
  describe('parseArgs', () => {
    it('should parse named arguments', () => {
      const argv = ['node', 'script', '--features', './tests', '--out', './out.json'];
      const args = parseArgs(argv);
      expect(args.features).toBe('./tests');
      expect(args.out).toBe('./out.json');
    });

    it('should parse boolean flags', () => {
      const argv = ['node', 'script', '--verbose', '--dry-run'];
      const args = parseArgs(argv);
      expect(args.verbose).toBe(true);
      expect(args['dry-run']).toBe(true);
    });

    it('should handle mixed args', () => {
      const argv = ['node', 'script', '--features', './tests', '--verbose'];
      const args = parseArgs(argv);
      expect(args.features).toBe('./tests');
      expect(args.verbose).toBe(true);
    });
  });

  describe('walkFeatures', () => {
    it('should recursively find .feature files', () => {
      // Mock readdirSync to return structure
      (fs.readdirSync as jest.Mock).mockReturnValueOnce([
        { name: 'subdir', isDirectory: () => true, isFile: () => false },
        { name: 'test.feature', isDirectory: () => false, isFile: () => true },
        { name: 'ignore.txt', isDirectory: () => false, isFile: () => true },
      ]).mockReturnValueOnce([ // subdir call
        { name: 'nested.feature', isDirectory: () => false, isFile: () => true }
      ]);

      const files = walkFeatures('/root');
      expect(files).toEqual([
        path.join('/root', 'subdir', 'nested.feature'),
        path.join('/root', 'test.feature')
      ]);
    });
  });

  describe('parseFeatureFile', () => {
    it('should parse a simple feature file', () => {
      const featureContent = `
Feature: User Login
  As a user I want to login

  Scenario: Successful Login
    Given I visit the login page
    When I enter valid credentials
    Then I should be logged in
      `;

      (fs.readFileSync as jest.Mock).mockReturnValue(featureContent);
      
      const result = parseFeatureFile('/tests/login.feature', '/tests');
      
      expect(result).not.toBeNull();
      expect(result!.name).toBe('User Login');
      expect(result!.description).toContain('As a user I want to login');
      expect(result!.scenarios).toHaveLength(1);
      
      const scenario = result!.scenarios[0];
      expect(scenario.name).toBe('Successful Login');
      expect(scenario.steps).toHaveLength(3);
      expect(scenario.steps[0].text).toContain('Given I visit the login page');
    });

    it('should parse scenario tags', () => {
      const featureContent = `
Feature: Tagged Feature
  
  @smoke @critical
  Scenario: Critical Path
    Given something
      `;

      (fs.readFileSync as jest.Mock).mockReturnValue(featureContent);
      
      const result = parseFeatureFile('/tests/tagged.feature', '/tests');
      
      const scenario = result!.scenarios[0];
      expect(scenario.tags).toEqual(['@smoke', '@critical']);
    });

    it('should handle Scenario Outline', () => {
        const featureContent = `
Feature: Outline Feature
  
  Scenario Outline: Data Driven
    Given I have <quantity> items
    Then I should have <result>
      `;
  
        (fs.readFileSync as jest.Mock).mockReturnValue(featureContent);
        
        const result = parseFeatureFile('/tests/outline.feature', '/tests');
        
        const scenario = result!.scenarios[0];
        expect(scenario.name).toBe('Data Driven');
        expect(scenario.steps).toHaveLength(2);
    });

    it('should return null if file has no feature', () => {
        (fs.readFileSync as jest.Mock).mockReturnValue('Just some random text');
        const result = parseFeatureFile('/tests/empty.feature', '/tests');
        expect(result).toBeNull();
    });
  });
});
