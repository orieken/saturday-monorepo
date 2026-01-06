
import { describe, it, expect } from 'vitest';
import { generateSelector, trackerScript } from '../src/tracker';

describe('tracker.ts', () => {
    describe('generateSelector', () => {
        it('should use ID if present', () => {
             const el = { id: 'test-id', tagName: 'DIV' } as unknown as Element;
             expect(generateSelector(el)).toBe('#test-id');
        });

        it('should use basic tag for body/html', () => {
             const body = { id: '', tagName: 'BODY' } as unknown as Element;
             expect(generateSelector(body)).toBe('body');
        });

        it('should use classes if present', () => {
             const el = { id: '', tagName: 'DIV', className: 'foo bar' } as unknown as Element;
             expect(generateSelector(el)).toBe('div.foo.bar');
        });

        it('should use nth-of-type for siblings', () => {
             const parent = { children: [] } as unknown as Element;
             const el1 = { id:'', tagName:'LI', parentElement: parent } as unknown as Element;
             const el2 = { id:'', tagName:'LI', parentElement: parent } as unknown as Element;
             const el3 = { id:'', tagName:'LI', parentElement: parent } as unknown as Element;
             
             (parent as any).children = [el1, el2, el3];
             
             expect(generateSelector(el2)).toBe('li:nth-of-type(2)');
        });
        
        it('should fallback to tagname', () => {
             const el = { id:'', tagName:'SPAN', parentElement: { children: [{ tagName: 'SPAN' }] } } as unknown as Element;
             expect(generateSelector(el)).toBe('span');
        });
    });

    describe('trackerScript', () => {
        it('should be a string containing the IIFE', () => {
            expect(trackerScript).toContain('(function() {');
            expect(trackerScript).toContain('window.__SATURDAY_HEATMAP_EVENTS__');
        });
    });
});
