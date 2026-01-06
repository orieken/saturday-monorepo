
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { BaseElement } from '../../src/base/base-element';
import { Page, Locator } from 'playwright';

class TestElement extends BaseElement {
  constructor(page: Page, selector: string) {
    super(page, selector);
  }
}

describe('BaseElement', () => {
  let mockPage: any;
  let mockLocator: any;
  let testElement: TestElement;

  beforeEach(() => {
    mockLocator = {
      waitFor: vi.fn(),
      isVisible: vi.fn().mockResolvedValue(true),
      count: vi.fn().mockResolvedValue(1),
      textContent: vi.fn().mockResolvedValue('Text Content'),
      innerText: vi.fn().mockResolvedValue('Inner Text'),
      getAttribute: vi.fn(),
      click: vi.fn(),
      dblclick: vi.fn(),
      hover: vi.fn(),
      scrollIntoViewIfNeeded: vi.fn(),
      boundingBox: vi.fn().mockResolvedValue({ x: 0, y: 0, width: 100, height: 100 }),
      evaluate: vi.fn((fn) => fn({ tagName: 'DIV' })),
    };

    mockPage = {
      locator: vi.fn().mockReturnValue(mockLocator),
      waitForTimeout: vi.fn(),
      evaluate: vi.fn((fn) => fn()),
      screenshot: vi.fn().mockResolvedValue(Buffer.from('fake-screenshot')),
      url: vi.fn().mockReturnValue('http://localhost/test'),
    };

    global.window = {
      innerWidth: 1024,
      innerHeight: 768,
    } as any;

    vi.stubGlobal('navigator', {
        userAgent: 'test-agent',
        platform: 'test-platform',
        cookieEnabled: true
    });

    testElement = new TestElement(mockPage as Page, '.test-element');
  });

  it('should initialize correctly', () => {
    expect(testElement.getLocator()).toBe(mockLocator);
    expect(testElement.getSelector()).toBe('.test-element');
    expect(mockPage.locator).toHaveBeenCalledWith('.test-element');
  });

  it('should wait for visible/hidden', async () => {
    await testElement.waitForVisible();
    expect(mockLocator.waitFor).toHaveBeenCalledWith({ state: 'visible' });
    await testElement.waitForHidden();
    expect(mockLocator.waitFor).toHaveBeenCalledWith({ state: 'hidden' });
  });

  it('should check existence and visibility', async () => {
    expect(await testElement.isVisible()).toBe(true);
    expect(await testElement.exists()).toBe(true);
  });

  it('should get text and attribute', async () => {
    expect(await testElement.getText()).toBe('Text Content');
    expect(await testElement.getInnerText()).toBe('Inner Text');
    await testElement.getAttribute('class');
    expect(mockLocator.getAttribute).toHaveBeenCalledWith('class');
  });

  it('should perform interactions', async () => {
    await testElement.click();
    expect(mockLocator.click).toHaveBeenCalled();
    await testElement.doubleClick();
    expect(mockLocator.dblclick).toHaveBeenCalled();
    await testElement.rightClick();
    expect(mockLocator.click).toHaveBeenCalledWith({ button: 'right' });
    await testElement.hover();
    expect(mockLocator.hover).toHaveBeenCalled();
    await testElement.scrollIntoView();
    expect(mockLocator.scrollIntoViewIfNeeded).toHaveBeenCalled();
  });

  it('should clickUntil', async () => {
    let attempts = 0;
    const condition = vi.fn().mockImplementation(async () => {
        attempts++;
        return attempts > 2;
    });

    await testElement.clickUntil(condition, { maxRetries: 5, interval: 10 });
    expect(attempts).toBe(3);
    expect(mockLocator.click).toHaveBeenCalledTimes(2);
  });

  it('should capture for training', async () => {
    const result = await testElement.captureForTraining('test-label');
    expect(mockPage.screenshot).toHaveBeenCalled();
    expect(result.label).toBe('test-label');
    expect(result.metadata.elementType).toBe('div');
  });

  it('should detect anomalies', async () => {
    const result = await testElement.detectAnomalies('model-label');
    expect(result.isValid).toBe(true);
    expect(result.elementSelector).toBe('.test-element');
  });

  it('should throw if element not visible for capture', async () => {
    mockLocator.boundingBox.mockResolvedValueOnce(null);
    await expect(testElement.captureForTraining('label')).rejects.toThrow('not visible for capture');
  });

    it('should throw if element not visible for detection', async () => {
    mockLocator.boundingBox.mockResolvedValueOnce(null);
    await expect(testElement.detectAnomalies('label')).rejects.toThrow('not visible for detection');
  });
});
