import * as fs from 'fs';
import * as path from 'path';

const htmlTemplate = (data: any[]) => `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Saturday Playwright Heatmap</title>
    <style>
        body { margin: 0; padding: 20px; font-family: -apple-system, system-ui, sans-serif; background: #1a1a1a; color: #fff; }
        .test-card { background: #2a2a2a; border-radius: 8px; margin-bottom: 20px; padding: 20px; }
        .heatmap-container { position: relative; display: inline-block; border: 1px solid #444; }
        .heatmap-img { display: block; max-width: 100%; }
        .overlay { position: absolute; top: 0; left: 0; pointer-events: none; }
        h1, h2 { color: #f0f0f0; }
        .stats { display: flex; gap: 20px; margin-bottom: 20px; }
        .stat-box { background: #333; padding: 10px 20px; border-radius: 4px; }
        .legend { display: flex; gap: 15px; margin-top: 10px; font-size: 0.9em; }
        .dot { width: 10px; height: 10px; display: inline-block; border-radius: 50%; margin-right: 5px; }
    </style>
</head>
<body>
    <h1>Test Execution Heatmap</h1>
    
    <div id="app"></div>

    <script>
        const testData = ${JSON.stringify(data)};

        function render() {
            const app = document.getElementById('app');
            
            testData.forEach(test => {
                const card = document.createElement('div');
                card.className = 'test-card';
                
                const stats = document.createElement('div');
                stats.className = 'stats';
                stats.innerHTML = \`
                    <div class="stat-box"><strong>Test:</strong> \${test.testTitle}</div>
                    <div class="stat-box"><strong>Interactions:</strong> \${test.events.length}</div>
                    <div class="stat-box"><strong>Interactables Found:</strong> \${test.interactables.length}</div>
                \`;
                
                card.appendChild(stats);

                if (test.snapshot) {
                    const container = document.createElement('div');
                    container.className = 'heatmap-container';

                    const img = document.createElement('img');
                    img.className = 'heatmap-img';
                    img.src = 'data:image/png;base64,' + test.snapshot;
                    
                    const canvas = document.createElement('canvas');
                    canvas.className = 'overlay';
                    canvas.style.position = 'absolute';
                    canvas.style.top = '0';
                    canvas.style.left = '0';
                    
                    img.onload = () => {
                        canvas.width = img.width;
                        canvas.height = img.height;
                        const ctx = canvas.getContext('2d');
                        
                        // Draw interactables (dimmed rectangles)
                        test.interactables.forEach(el => {
                            ctx.strokeStyle = 'rgba(0, 255, 255, 0.8)';
                            ctx.lineWidth = 3;
                            ctx.strokeRect(el.rect.x, el.rect.y, el.rect.width, el.rect.height);
                        });

                        // Draw interactions (heat map spots)
                        test.events.forEach(evt => {
                            ctx.beginPath();
                            ctx.arc(evt.x, evt.y + window.scrollY, 10, 0, 2 * Math.PI);
                            ctx.fillStyle = 'rgba(255, 0, 0, 0.6)';
                            ctx.fill();
                        });
                    };

                    container.appendChild(img);
                    container.appendChild(canvas);
                    card.appendChild(container);
                    
                    const legend = document.createElement('div');
                    legend.className = 'legend';
                    legend.innerHTML = \`
                        <span><span class="dot" style="background: rgba(0, 255, 255, 0.3); border: 1px solid rgba(0,255,255,1)"></span> Interactable Element</span>
                        <span><span class="dot" style="background: rgba(255, 0, 0, 0.6)"></span> Interaction Recorded</span>
                    \`;
                    card.appendChild(legend);
                }

                app.appendChild(card);
            });
        }

        render();
    </script>
</body>
</html>
`;

export function generateReport(inputDir: string, outputFile: string) {
    if (!fs.existsSync(inputDir)) {
        console.error(`Input directory ${inputDir} does not exist.`);
        return;
    }

    const files = fs.readdirSync(inputDir).filter(f => f.endsWith('.json'));
    const allData = files.map(f => {
        try {
            return JSON.parse(fs.readFileSync(path.join(inputDir, f), 'utf-8'));
        } catch (err) {
            console.error(`Failed to parse ${f}`, err);
            return null;
        }
    }).filter(d => d !== null);

    const html = htmlTemplate(allData);
    fs.writeFileSync(outputFile, html);
    console.log(`Report generated at ${outputFile}`);
}

if (require.main === module) {
    const args = process.argv.slice(2);
    const inputDir = args[0] || 'heatmap-data';
    const outputFile = args[1] || 'heatmap-report.html';
    generateReport(path.resolve(process.cwd(), inputDir), path.resolve(process.cwd(), outputFile));
}
