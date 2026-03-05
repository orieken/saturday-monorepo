
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ButtonElement } from '../../src/elements/button-element';
import { InputElement } from '../../src/elements/input-element';
import { LinkElement } from '../../src/elements/link-element';
import { Page } from 'playwright';

describe('Elements', () => {
  let mockPage: any;
  let mockLocator: any;

  beforeEach(() => {
    mockLocator = {
      isDisabled: vi.fn(),
      isEnabled: vi.fn(),
      type: vi.fn(),
      fill: vi.fn(),
      clear: vi.fn(),
      inputValue: vi.fn(),
      getAttribute: vi.fn(),
    };
    mockPage = {
      locator: vi.fn().mockReturnValue(mockLocator),
    };
  });

  describe('ButtonElement', () => {
      it('should check status', async () => {
          const btn = new ButtonElement(mockPage as Page, 'button');
          mockLocator.isDisabled.mockResolvedValue(true);
          expect(await btn.isDisabled()).toBe(true);

          mockLocator.isEnabled.mockResolvedValue(true);
          expect(await btn.isEnabled()).toBe(true);
      });
  });

  describe('InputElement', () => {
    it('should interact with input', async () => {
       const input = new InputElement(mockPage as Page, 'input');
       await input.type('hello', { delay: 10 });
       expect(mockLocator.type).toHaveBeenCalledWith('hello', { delay: 10 });

       await input.fill('hello');
       expect(mockLocator.fill).toHaveBeenCalledWith('hello');

       await input.clear();
       expect(mockLocator.clear).toHaveBeenCalled();

       mockLocator.inputValue.mockResolvedValue('value');
       expect(await input.getValue()).toBe('value');
    });
  });

  describe('LinkElement', () => {
      it('should get attributes', async () => {
          const link = new LinkElement(mockPage as Page, 'a');
          
          mockLocator.getAttribute.mockImplementation((attr: string) => {
              if(attr === 'href') return 'http://google.com';
              if(attr === 'target') return '_blank';
              return null;
          });

          expect(await link.getHref()).toBe('http://google.com');
          expect(await link.getTarget()).toBe('_blank');
          expect(await link.opensInNewTab()).toBe(true);
      });

      it('should handle non-blank target', async () => {
        const link = new LinkElement(mockPage as Page, 'a');
        mockLocator.getAttribute.mockResolvedValue('_self');
        expect(await link.opensInNewTab()).toBe(false);
      });
  });
});
