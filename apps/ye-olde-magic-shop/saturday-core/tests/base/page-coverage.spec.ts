
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { BasePage } from '../../src/base/base-page';
import { BaseElement } from '../../src/base/base-element';
import { BaseFilter } from '../../src/base/base-filter';

class TestElement extends BaseElement {
    constructor(page: any, selector: string) { super(page, selector); }
}

class TestFilter extends BaseFilter {
    constructor(page: any) { super(page); }
    async apply() {}
    async clear() {}
    async isActive() { return false; }
}

class TestPage extends BasePage {
    containerSelector = 'body';
    path?: string = '/test';
    
    constructor(page: any, site: any) { super(page, site); }
    protected initializeElements() {}
    protected initializeFilters() {}
    
    public testRegisterElement(name: string, sel: string) { this.registerElement(name, TestElement, sel); }
    public testRegisterFilter(name: string) { this.registerFilter(name, TestFilter); }
    public setPath(p?: string) { this.path = p; }
}

describe('BasePage Branch Coverage', () => {
    let mockPage: any;
    let mockSite: any;
    let page: TestPage;

    beforeEach(() => {
        mockPage = {
            locator: vi.fn().mockReturnValue({ waitFor: vi.fn(), isVisible: vi.fn() }),
            screenshot: vi.fn().mockResolvedValue(Buffer.from('img')),
            evaluate: vi.fn().mockResolvedValue({}),
            url: vi.fn().mockReturnValue('u'),
            goto: vi.fn()
        };
        mockSite = {
            detectors: { anomalyDetector: { validateImage: vi.fn().mockResolvedValue({}) } },
            getBaseUrl: vi.fn().mockReturnValue('http://localhost')
        };
        page = new TestPage(mockPage, mockSite);
    });

    it('should handle element registration and retrieval', () => {
        expect(() => page.getElement('missing')).toThrow();
        page.testRegisterElement('e', '.e');
        expect(page.getElement('e')).toBeInstanceOf(TestElement);
    });

    it('should handle filter registration and retrieval', () => {
         expect(() => page.getFilter('missing')).toThrow();
         page.testRegisterFilter('f');
         expect(page.getFilter('f')).toBeInstanceOf(TestFilter);
    });

    it('should handle visit constraints', async () => {
        page.setPath(undefined);
        await expect(page.visit()).rejects.toThrow('Path is not defined');
        
        page.setPath('/p');
        await page.visit();
        expect(mockSite.getBaseUrl).toHaveBeenCalled();
        expect(mockPage.goto).toHaveBeenCalledWith('http://localhost/p');
    });

    it('should cover ML methods', async () => {
        const capture = await page.captureForTraining('l');
        expect(capture).toBeDefined();
        
        const val = await page.validateWithML('l');
        expect(val).toBeDefined();
    });

    it('should cover scroll utilities', async () => {
        await page.scrollToTop();
        await page.scrollToBottom();
        expect(mockPage.evaluate).toHaveBeenCalledTimes(2); // browser info info calls + 2 scrolls. 
        // captureForTraining calls evaluate for browserInfo viewport.
    });
});
