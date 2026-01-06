import * as fs from 'fs';
import * as path from 'path';
import { HeatmapAnalyzer, TestData } from './analyzer';

export function runAnalysis(inputDir: string) {
    if (!fs.existsSync(inputDir)) {
        console.error(`Input directory ${inputDir} does not exist.`);
        return;
    }

    const analyzer = new HeatmapAnalyzer();
    const files = fs.readdirSync(inputDir).filter(f => f.endsWith('.json'));
    
    console.log(`Analyzing ${files.length} test files in ${inputDir}...\n`);

    files.forEach(file => {
        try {
            const rawData = fs.readFileSync(path.join(inputDir, file), 'utf-8');
            const data: TestData = JSON.parse(rawData);
            
            const result = analyzer.analyze(data);
            
            console.log(`Test: ${data.testTitle} (${data.status})`);
            console.log(`  Coverage Score: ${result.coverageScore.toFixed(1)}%`);
            console.log(`  Interactables: ${data.interactables.length}`);
            console.log(`  Hotspots (Clusters): ${result.hotspots.length}`);
            console.log(`  Cold Spots (Untested): ${result.coldSpots.length}`);
            
            if (result.coldSpots.length > 0) {
                console.log(`  Top 3 Cold Spots:`);
                result.coldSpots.slice(0, 3).forEach(spot => {
                    console.log(`    - ${spot.tagName}${spot.text ? ` "${spot.text}"` : ''} (${spot.selector})`);
                });
            }
            console.log('-'.repeat(40));

        } catch (err) {
            console.error(`Error analyzing ${file}:`, err);
        }
    });
}

if (require.main === module) {
    const args = process.argv.slice(2);
    const inputDir = args[0] || 'heatmap-data';
    runAnalysis(path.resolve(process.cwd(), inputDir));
}
