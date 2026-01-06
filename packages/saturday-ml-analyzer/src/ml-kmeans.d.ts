declare module 'ml-kmeans' {
    export interface KMeansResult {
        clusters: number[];
        centroids: number[][];
        iterations: number;
    }
    
    export interface KMeansOptions {
        maxIterations?: number;
        tolerance?: number;
        withIterations?: boolean;
        initialization?: 'kmeans++' | 'random' | 'mostDistant';
        seed?: number;
    }

    export function kmeans(data: number[][], k: number, options?: KMeansOptions): KMeansResult;
}
