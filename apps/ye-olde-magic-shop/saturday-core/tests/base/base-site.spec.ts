
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { BaseSite } from '../../src/base/base-site';
import { BasePage } from '../../src/base/base-page';
import { BaseFlow } from '../../src/base/base-flow';
import { Page } from 'playwright';

// Mocks for Facades
class MockTrainersFacade {
  trainCurrentPage = vi.fn();
}
class MockDetectorsFacade {
  validateCurrentPage = vi.fn();
}
class MockModelsFacade {}

// Mock Page/Flow
class MockPage extends BasePage {
  containerSelector = '#app';
  path = '/mock';
  visit = vi.fn().mockResolvedValue(undefined);
  protected initializeElements() {}
  protected initializeFilters() {}
  // Overriding createContainer to avoid needing Locator
  protected createContainer() { (this as any).container = { waitFor: vi.fn(), isVisible: vi.fn() }; }
}

class MockFlow extends BaseFlow {
  async execute() {}
}

class TestSite extends BaseSite {
  protected initializePages(): void {
    this.registerPage('mockPage', MockPage);
  }
  protected initializeFlows(): void {
    this.registerFlow('mockFlow', MockFlow);
  }

  // Override factories to inject mocks
  protected createTrainersFacade(): any { return new MockTrainersFacade(); }
  protected createDetectorsFacade(): any { return new MockDetectorsFacade(); }
  protected createModelsFacade(): any { return new MockModelsFacade(); }
}

describe('BaseSite', () => {
  let mockPage: any;
  let testSite: TestSite;

  beforeEach(() => {
    mockPage = {
      goto: vi.fn(),
      url: vi.fn().mockReturnValue('http://localhost'),
      title: vi.fn().mockResolvedValue('Mock Title'),
      waitForLoadState: vi.fn(),
      locator: vi.fn(),
    };
    testSite = new TestSite(mockPage as Page, 'http://localhost');
  });

  it('should initialize correctly', () => {
    expect(testSite.getBaseUrl()).toBe('http://localhost');
    expect(testSite.pageObject).toBe(mockPage);
  });

  it('should register and get pages', () => {
    const page = testSite.getPage('mockPage');
    expect(page).toBeInstanceOf(MockPage);
    expect(() => testSite.getPage('nonExistent')).toThrow();
  });

  it('should register and get flows', () => {
    const flow = testSite.getFlow('mockFlow');
    expect(flow).toBeInstanceOf(MockFlow);
    expect(() => testSite.getFlow('nonExistent')).toThrow();
  });

  it('should lazy load facades', () => {
    const trainers1 = testSite.trainers;
    const trainers2 = testSite.trainers;
    expect(trainers1).toBeInstanceOf(MockTrainersFacade);
    expect(trainers1).toBe(trainers2); // Singleton per instance

    expect(testSite.detectors).toBeInstanceOf(MockDetectorsFacade);
    expect(testSite.models).toBeInstanceOf(MockModelsFacade);
  });

  it('should establish page baseline', async () => {
    const page = testSite.getPage('mockPage');
    await testSite.establishPageBaseline('mockPage', { fullPage: true });
    expect(page.visit).toHaveBeenCalled();
    expect(testSite.trainers.trainCurrentPage).toHaveBeenCalledWith('mockPage_baseline', { fullPage: true });
  });

  it('should validate page against baseline', async () => {
      const page = testSite.getPage('mockPage');
      await testSite.validatePageAgainstBaseline('mockPage');
      expect(page.visit).toHaveBeenCalled();
      expect(testSite.detectors.validateCurrentPage).toHaveBeenCalledWith('mockPage_baseline', undefined);
  });

  it('should delegate interactions to page', async () => {
      await testSite.visit();
      expect(mockPage.goto).toHaveBeenCalledWith('http://localhost');

      await testSite.getCurrentUrl();
      expect(mockPage.url).toHaveBeenCalled();

      await testSite.getCurrentTitle();
      expect(mockPage.title).toHaveBeenCalled();

      await testSite.waitForNavigation();
      expect(mockPage.waitForLoadState).toHaveBeenCalledWith('networkidle');
  });
});
