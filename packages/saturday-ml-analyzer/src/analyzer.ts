import { kmeans } from 'ml-kmeans';
import * as fs from 'fs';
import * as path from 'path';

export interface InteractionEvent {
    type: string;
    x: number;
    y: number;
    selector: string;
    timestamp: number;
}

export interface InteractableElement {
    selector: string;
    tagName: string;
    text: string;
    rect: {
        x: number;
        y: number;
        width: number;
        height: number;
    };
}

export interface TestData {
    testId: string;
    testTitle: string;
    status: string;
    events: InteractionEvent[];
    interactables: InteractableElement[];
}

export interface AnalysisResult {
    testId: string;
    hotspots: { x: number, y: number }[];
    coldSpots: InteractableElement[]; // Elements far from any interaction
    coverageScore: number; // Percentage of interactables engaged
}

export class HeatmapAnalyzer {
    
    // Distance threshold in pixels to consider an element "interacted with"
    // if a cluster centroid or event is within this range.
    private readonly PROXIMITY_THRESHOLD = 50; 

    public analyze(data: TestData): AnalysisResult {
        const interactions = data.events.map(e => [e.x, e.y]);
        
        let hotspots: { x: number, y: number }[] = [];
        
        // If we have enough points, use K-Means to find centroids (Hotspots)
        if (interactions.length >= 3) {
            // Determine K dynamically? For now, let's try a heuristic: 1 cluster per 5 interactions
            const k = Math.max(1, Math.floor(interactions.length / 5));
            try {
                const result = kmeans(interactions, k, {});
                hotspots = result.centroids.map(c => ({ x: c[0], y: c[1] }));
            } catch (e) {
                // Fallback if clustering fails (e.g. all points same location)
                hotspots = interactions.map(i => ({ x: i[0], y: i[1] }));
            }
        } else {
             // Too few points for meaningful clustering, just use raw points
             hotspots = interactions.map(i => ({ x: i[0], y: i[1] }));
        }

        // Identify Cold Spots (Untested Elements)
        const coldSpots: InteractableElement[] = [];
        let interactedCount = 0;

        data.interactables.forEach(el => {
            const elCenter = {
                x: el.rect.x + (el.rect.width / 2),
                y: el.rect.y + (el.rect.height / 2)
            };

            // Check if any interaction event or hotspot is close enough
            // Using raw events for precision, but hotspots give us the "center of gravity"
            
            // Simple approach: Check distance to nearest interaction event
            const nearestDistance = data.events.reduce((min, evt) => {
                const dist = Math.sqrt(Math.pow(evt.x - elCenter.x, 2) + Math.pow(evt.y - elCenter.y, 2));
                return Math.min(min, dist);
            }, Infinity);

            if (nearestDistance > this.PROXIMITY_THRESHOLD) {
                coldSpots.push(el);
            } else {
                interactedCount++;
            }
        });

        const coverageScore = data.interactables.length > 0 
            ? (interactedCount / data.interactables.length) * 100 
            : 100;

        return {
            testId: data.testId,
            hotspots,
            coldSpots,
            coverageScore
        };
    }
}
