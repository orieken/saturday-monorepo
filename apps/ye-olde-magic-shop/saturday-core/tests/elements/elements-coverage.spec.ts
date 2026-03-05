
import { describe, it, expect, vi } from 'vitest';
import { ButtonElement } from '../../src/elements/button-element';
import { InputElement } from '../../src/elements/input-element';
import { LinkElement } from '../../src/elements/link-element';

describe('Elements Coverage', () => {
    let mockPage: any;

    beforeEach(() => {
        mockPage = {
            locator: vi.fn().mockReturnValue({ 
                click: vi.fn(), 
                fill: vi.fn(), 
                getAttribute: vi.fn(),
                type: vi.fn(),
                clear: vi.fn(),
                inputValue: vi.fn()
            })
        };
    });

    it('should cover ButtonElement', async () => {
        const btn = new ButtonElement(mockPage, 'btn');
        expect(btn).toBeDefined();
    });

    it('should cover InputElement', async () => {
        const input = new InputElement(mockPage, 'inp');
        await input.fill('val');
        expect(mockPage.locator().fill).toHaveBeenCalledWith('val');
        
        await input.type('t');
        await input.clear();

        mockPage.locator().inputValue = vi.fn().mockResolvedValue('val');
        expect(await input.getValue()).toBe('val');
    });

    it('should cover LinkElement', async () => {
        const link = new LinkElement(mockPage, 'a');
        mockPage.locator().getAttribute.mockResolvedValue('h');
        expect(await link.getHref()).toBe('h');
    });
});
