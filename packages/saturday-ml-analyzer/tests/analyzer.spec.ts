
import { describe, it, expect, vi } from 'vitest';
import { HeatmapAnalyzer, TestData } from '../src/analyzer';

// Mock ml-kmeans since it's a computation heavy library and we want to test logic
vi.mock('ml-kmeans', () => {
    return {
        kmeans: vi.fn((data, k) => {
             // Simple mock that returns the first k points as centroids
             return {
                 centroids: data.slice(0, k),
                 clusters: new Array(data.length).fill(0)
             };
        })
    };
});

describe('HeatmapAnalyzer', () => {
    const analyzer = new HeatmapAnalyzer();
    const mockInteractable = {
        selector: '.btn',
        tagName: 'BUTTON',
        text: 'Click Me',
        rect: { x: 100, y: 100, width: 20, height: 10 }
    };

    it('should analyze interactions correctly', () => {
        const data: TestData = {
            testId: 't1',
            testTitle: 'test',
            status: 'passed',
            events: [
                { type: 'click', x: 110, y: 105, selector: '.btn', timestamp: 100 },
                { type: 'click', x: 112, y: 106, selector: '.btn', timestamp: 200 }
            ],
            interactables: [mockInteractable]
        };

        const result = analyzer.analyze(data);
        expect(result.testId).toBe('t1');
        expect(result.coverageScore).toBe(100);
        expect(result.coldSpots).toEqual([]);
        expect(result.hotspots.length).toBe(2);
    });

    it('should identify cold spots', () => {
        const data: TestData = {
            testId: 't2',
            testTitle: 'cold spot test',
            status: 'failed',
            events: [
                { type: 'click', x: 0, y: 0, selector: 'body', timestamp: 1 } // Far away
            ],
            interactables: [mockInteractable] // at 100,100
        };

        const result = analyzer.analyze(data);
        expect(result.coverageScore).toBe(0);
        expect(result.coldSpots).toHaveLength(1);
        expect(result.coldSpots[0].selector).toBe('.btn');
    });

    it('should handle clustering with many points', () => {
         const events = [];
         for(let i=0; i<20; i++) {
             events.push({ type: 'click', x: 100+i, y: 100+i, selector: '.btn', timestamp: i });
         }
         
         const data: TestData = {
             testId: 'cluster',
             testTitle: 'cluster',
             status: 'passed',
             events,
             interactables: [mockInteractable]
         };

         const result = analyzer.analyze(data);
         // 20 points / 5 = 4 clusters expected
         expect(result.hotspots).toHaveLength(4);
         expect(result.coverageScore).toBe(100);
    });

    it('should handle zero interactables', () => {
        const data: TestData = {
            testId: 'empty',
            testTitle: 'empty',
            status: 'passed',
            events: [],
            interactables: []
        };
        const result = analyzer.analyze(data);
        expect(result.coverageScore).toBe(100);
    });
});
