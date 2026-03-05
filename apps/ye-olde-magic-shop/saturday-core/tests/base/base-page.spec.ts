
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { BasePage } from '../../src/base/base-page';
import { BaseSite } from '../../src/base/base-site';
import { Page, Locator } from 'playwright';
import { BaseElement } from '../../src/base/base-element';
import { BaseFilter } from '../../src/base/base-filter';

class MockElement extends BaseElement {
  constructor(page: Page, selector: string) {
    super(page, selector);
  }
}

class MockFilter extends BaseFilter {
  constructor(page: Page) {
    super(page);
  }
  async matches(condition: any): Promise<boolean> { return true; }
  async getValue(): Promise<any> { return 'mock'; }
  async apply(value: any): Promise<void> { }
  async clear(): Promise<void> { }
}

class TestPage extends BasePage {
  get containerSelector() { return '#app'; }
  path?: string = '/test';

  protected initializeElements(): void {
    this.registerElement('testElement', MockElement, '.test');
  }

  protected initializeFilters(): void {
    this.registerFilter('testFilter', MockFilter);
  }
}


describe('BasePage', () => {
  let mockPage: any;
  let mockSite: any;
  let testPage: TestPage;
  let mockLocator: any;

  beforeEach(() => {
    mockLocator = {
      waitFor: vi.fn(),
      isVisible: vi.fn().mockResolvedValue(true),
    };

    mockPage = {
      locator: vi.fn().mockReturnValue(mockLocator),
      goto: vi.fn(),
      evaluate: vi.fn((fn) => fn()),
      screenshot: vi.fn().mockResolvedValue(Buffer.from('fake-screenshot')),
      url: vi.fn().mockReturnValue('http://localhost/test'),
    };

    mockSite = {
      getBaseUrl: vi.fn().mockReturnValue('http://localhost'),
      detectors: {
        anomalyDetector: {
          validateImage: vi.fn().mockResolvedValue({
            isValid: true,
            score: 0.99,
            anomalies: [],
            metadata: {}
          })
        }
      }
    };

    global.window = {
      scrollTo: vi.fn(),
      innerWidth: 1024,
      innerHeight: 768,
    } as any;
    
    vi.stubGlobal('navigator', {
        userAgent: 'test-agent',
        platform: 'test-platform',
        cookieEnabled: true
    });
    
    global.document = {
        body: {
            scrollHeight: 2000
        }
    } as any;

    testPage = new TestPage(mockPage as Page, mockSite as BaseSite);
  });

  it('should initialize correctly', () => {
    expect(testPage.getContainer()).toBe(mockLocator);
    expect(mockPage.locator).toHaveBeenCalledWith('#app');
  });

  it('should visit the page', async () => {
    await testPage.visit();
    expect(mockSite.getBaseUrl).toHaveBeenCalled();
    expect(mockPage.goto).toHaveBeenCalledWith('http://localhost/test');
    expect(mockLocator.waitFor).toHaveBeenCalledWith({ state: 'visible' });
  });

  it('should register and retrieve elements', () => {
    const element = testPage.getElement('testElement');
    expect(element).toBeInstanceOf(MockElement);
    expect(() => testPage.getElement('nonExistent')).toThrow();
  });

  it('should register and retrieve filters', () => {
    const filter = testPage.getFilter('testFilter');
    expect(filter).toBeInstanceOf(MockFilter);
    expect(() => testPage.getFilter('nonExistent')).toThrow();
  });

  it('should check if loaded', async () => {
    const loaded = await testPage.isLoaded();
    expect(loaded).toBe(true);
    expect(mockLocator.isVisible).toHaveBeenCalled();
  });

  it('should scroll to top and bottom', async () => {
    await testPage.scrollToTop();
    expect(mockPage.evaluate).toHaveBeenCalled();
    await testPage.scrollToBottom();
    expect(mockPage.evaluate).toHaveBeenCalled();
  });

  it('should throw if path is not defined', async () => {
    class NoPathPage extends TestPage {
      path = undefined;
    }
    const noPathPage = new NoPathPage(mockPage as Page, mockSite as BaseSite);
    await expect(noPathPage.visit()).rejects.toThrow('Path is not defined');
  });

  it('should capture for training', async () => {
    const result = await testPage.captureForTraining('test-label');
    expect(mockPage.screenshot).toHaveBeenCalled();
    expect(result.label).toBe('test-label');
    expect(result.metadata.pageUrl).toBe('http://localhost/test');
    expect(result.metadata.browserInfo.userAgent).toBe('test-agent');
  });
  
  it('should capture for training with options', async () => {
     const result = await testPage.captureForTraining('test-label', { fullPage: false });
     expect(mockPage.screenshot).toHaveBeenCalledWith(expect.objectContaining({ fullPage: false }));
  });

  it('should validate with ML', async () => {
    const result = await testPage.validateWithML('test-model');
    expect(mockPage.screenshot).toHaveBeenCalled();
    expect(mockSite.detectors.anomalyDetector.validateImage).toHaveBeenCalled();
    expect(result.isValid).toBe(true);
  });
});
