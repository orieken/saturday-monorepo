
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { BaseSite } from '../../src/base/base-site';
import { BasePage } from '../../src/base/base-page';
import { BaseFlow } from '../../src/base/base-flow';

class TestPage extends BasePage {
    containerSelector = 'b';
    path = '/p';
    protected initializeElements() {}
    protected initializeFilters() {}
}

class TestFlow extends BaseFlow {
    async execute() {}
    public testGetPage(n: string) { return this.getPage(n); }
    public testGetFlow(n: string) { return this.getFlow(n); }
}

class TestSite extends BaseSite {
    constructor(page: any) { super(page, 'http://base'); }
    initializePages() {}
    initializeFlows() {}
    
    public testRegisterPage(n: string) { this.registerPage(n, TestPage); }
    public testRegisterFlow(n: string) { this.registerFlow(n, TestFlow); }
}

describe('BaseSite and BaseFlow Branch Coverage', () => {
    let mockPage: any;
    let site: TestSite;

    beforeEach(() => {
        mockPage = {
            locator: vi.fn().mockReturnValue({ waitFor: vi.fn(), isVisible: vi.fn() }),
            screenshot: vi.fn(),
            goto: vi.fn(),
            url: vi.fn().mockReturnValue('u'),
            title: vi.fn().mockResolvedValue('t'),
            waitForLoadState: vi.fn()
        };
        site = new TestSite(mockPage);
    });

    it('should handle page/flow registration and retrieval', () => {
        expect(() => site.getPage('missing')).toThrow();
        expect(() => site.getFlow('missing')).toThrow();

        site.testRegisterPage('p');
        site.testRegisterFlow('f');

        expect(site.getPage('p')).toBeInstanceOf(TestPage);
        expect(site.getFlow('f')).toBeInstanceOf(TestFlow);
    });

    it('should cover pageObject getter', () => {
        expect(site.pageObject).toBe(mockPage);
    });

    it('should cover facade lazy creation', () => {
        // First access
        const t1 = site.trainers;
        expect(t1).toBeDefined();
        // Second access (cached)
        const t2 = site.trainers;
        expect(t2).toBe(t1);

        const d1 = site.detectors;
        expect(d1).toBeDefined();
        expect(site.detectors).toBe(d1);

        const m1 = site.models;
        expect(m1).toBeDefined();
        expect(site.models).toBe(m1);
    });

    it('should cover baseline methods', async () => {
        site.testRegisterPage('p');
        const trainSpy = vi.spyOn(site.trainers, 'trainCurrentPage').mockResolvedValue({} as any);
        const detectSpy = vi.spyOn(site.detectors, 'validateCurrentPage').mockResolvedValue({} as any);
        
        await site.getPage('p').visit(); // mock visit triggers goto

        await site.establishPageBaseline('p');
        expect(trainSpy).toHaveBeenCalledWith('p_baseline', undefined);

        await site.validatePageAgainstBaseline('p');
        expect(detectSpy).toHaveBeenCalledWith('p_baseline', undefined);
    });

    it('should cover basic site methods', async () => {
        await site.visit();
        expect(mockPage.goto).toHaveBeenCalledWith('http://base');
        
        expect(site.getBaseUrl()).toBe('http://base');
        expect(await site.getCurrentUrl()).toBe('u');
        expect(await site.getCurrentTitle()).toBe('t');
        
        await site.waitForNavigation();
        expect(mockPage.waitForLoadState).toHaveBeenCalledWith('networkidle');
    });

    it('should cover BaseFlow proxy methods', () => {
        site.testRegisterPage('p');
        site.testRegisterFlow('f');
        const flow = site.getFlow<TestFlow>('f');
        
        expect(flow.testGetPage('p')).toBeDefined();
        expect(flow.testGetFlow('f')).toBeDefined();
    });
});
