
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { BaseElement } from '../../src/base/base-element';

class TestElement extends BaseElement {
    constructor(page: any) { super(page, '.test'); }
}

describe('BaseElement Coverage', () => {
    let mockPage: any;
    let element: TestElement;

    beforeEach(() => {
        mockPage = {
            locator: vi.fn().mockReturnValue({
                waitFor: vi.fn(),
                isVisible: vi.fn().mockResolvedValue(true),
                count: vi.fn().mockResolvedValue(1),
                textContent: vi.fn().mockResolvedValue('text'),
                innerText: vi.fn().mockResolvedValue('inner'),
                getAttribute: vi.fn().mockResolvedValue('attr'),
                click: vi.fn(),
                dblclick: vi.fn(),
                hover: vi.fn(),
                scrollIntoViewIfNeeded: vi.fn(),
                boundingBox: vi.fn().mockResolvedValue({ x:0,y:0,width:10,height:10 }),
                evaluate: vi.fn().mockResolvedValue('div')
            }),
            waitForTimeout: vi.fn(),
            screenshot: vi.fn().mockResolvedValue(Buffer.from('img')),
            url: vi.fn().mockReturnValue('http://localhost'),
            evaluate: vi.fn().mockResolvedValue({})
        };
        element = new TestElement(mockPage);
    });

    it('should cover all getters and simple checking methods', async () => {
        expect(element.getLocator()).toBeDefined();
        expect(element.getSelector()).toBe('.test');
        await element.waitForVisible();
        await element.waitForHidden();
        expect(await element.isVisible()).toBe(true);
        expect(await element.exists()).toBe(true);
        expect(await element.getText()).toBe('text');
        expect(await element.getInnerText()).toBe('inner');
        expect(await element.getAttribute('class')).toBe('attr');
        await element.doubleClick();
        await element.rightClick();
        await element.hover();
        await element.scrollIntoView();
        expect(await element.getBoundingBox()).toBeDefined();
    });

    it('should cover clickUntil logic', async () => {
        // Condition met immediately
        const condTrue = vi.fn().mockResolvedValue(true);
        await element.clickUntil(condTrue);
        expect(mockPage.locator().click).not.toHaveBeenCalled();

        // Condition met after 1 click
        mockPage.locator().click.mockClear();
        const condFalseTrue = vi.fn()
            .mockResolvedValueOnce(false)
            .mockResolvedValueOnce(true);
        await element.clickUntil(condFalseTrue, { maxRetries: 2, interval: 10 });
        expect(mockPage.locator().click).toHaveBeenCalledTimes(1);
        expect(mockPage.waitForTimeout).toHaveBeenCalledWith(10);

        // Condition never met
        mockPage.locator().click.mockClear();
        const condFalse = vi.fn().mockResolvedValue(false);
        await expect(element.clickUntil(condFalse, { maxRetries: 2, interval: 1 })).rejects.toThrow(/not met/);
        expect(mockPage.locator().click).toHaveBeenCalledTimes(2);
    });

    it('should cover ML getters', async () => {
        const capture = await element.captureForTraining('lbl');
        expect(capture.label).toBe('lbl');
        
        mockPage.locator().boundingBox.mockResolvedValue(null);
        await expect(element.captureForTraining('lbl')).rejects.toThrow('visible');
    });

    it('should cover detectAnomalies', async () => {
         const res = await element.detectAnomalies('lbl');
         expect(res.isValid).toBe(true);
         
         mockPage.locator().boundingBox.mockResolvedValue(null);
         await expect(element.detectAnomalies('lbl')).rejects.toThrow('visible');
    });
});
