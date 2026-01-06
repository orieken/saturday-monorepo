#!/usr/bin/env node
import * as fs from 'fs';
import * as path from 'path';
import * as http from 'http';
import * as https from 'https';

export interface CucumberStepRef {
  text: string;
  line: number;
}

export interface CucumberScenarioRef {
  id: string;
  line: number;
  name: string;
  file: string;
  tags: string[];
  steps: CucumberStepRef[];
}

export interface CucumberFeatureRef {
  id: string;
  name: string;
  file: string;
  description: string;
  scenarios: CucumberScenarioRef[];
}

export interface CucumberIndex {
  framework: string;
  suiteId: string;
  features: CucumberFeatureRef[];
}

export function parseArgs(argv: string[]) {
  const args: Record<string, string | boolean> = {};
  for (let i = 2; i < argv.length; i++) {
    const token = argv[i];
    if (!token) continue;
    if (token.startsWith('--')) {
      const key = token.slice(2);
      const next = argv[i + 1];
      if (!next || next.startsWith('--')) {
        args[key] = true; // boolean flag
      } else {
        args[key] = next;
        i++;
      }
    }
  }
  return args as Record<string, string> & Record<string, any>;
}

export function walkFeatures(dir: string): string[] {
  const out: string[] = [];
  const entries = fs.readdirSync(dir, { withFileTypes: true });
  for (const e of entries) {
    const full = path.join(dir, e.name);
    if (e.isDirectory()) {
      out.push(...walkFeatures(full));
    } else if (e.isFile() && e.name.endsWith('.feature')) {
      out.push(full);
    }
  }
  return out;
}

export function parseFeatureFile(file: string, root: string): CucumberFeatureRef | null {
  const rel = path.relative(root, file);
  const data = fs.readFileSync(file, 'utf8').split(/\r?\n/);
  let featureName = '';
  const featureDesc: string[] = [];
  const scenarios: CucumberScenarioRef[] = [];
  let currentScenario: CucumberScenarioRef | null = null;
  let pendingTags: string[] = [];

  for (let i = 0; i < data.length; i++) {
    const lineNo = i + 1;
    const raw = data[i];
    const line = raw.trim();

    if (line.startsWith('@')) {
      const tags = line.split(/\s+/).filter(t => t.startsWith('@'));
      pendingTags.push(...tags);
      continue;
    }

    if (line.startsWith('Feature:')) {
      featureName = line.replace('Feature:', '').trim();
      pendingTags = []; // Ignore feature tags for now
      continue;
    }
    if (line.startsWith('Scenario Outline:') || line.startsWith('Scenario:') || line.startsWith('Scenario')) {
      if (currentScenario) {
        scenarios.push(currentScenario);
      }
      const name = line.replace('Scenario Outline:', '').replace('Scenario:', '').trim();
      currentScenario = {
        id: `${rel}:${lineNo}`,
        line: lineNo,
        name,
        file: path.basename(file),
        tags: pendingTags,
        steps: []
      };
      pendingTags = [];
      continue;
    }
    
    if (/^(Given|When|Then|And|But)\s+/.test(line) && currentScenario) {
      currentScenario.steps.push({
        line: lineNo,
        text: line
      });
      continue;
    }
    // Lines directly after Feature: as description until first scenario
    if (!currentScenario && featureName && raw.length > 0 && !/^\s*$/.test(raw)) {
      featureDesc.push(raw);
    }
  }
  if (currentScenario) {
    scenarios.push(currentScenario);
  }
  if (!featureName) return null;

  return {
    id: rel,
    name: featureName,
    file: path.basename(file),
    description: featureDesc.join('\n'),
    scenarios
  };
}

function postJSON(urlStr: string, body: any): Promise<{ status: number; text: string }> {
  return new Promise((resolve, reject) => {
    try {
      const u = new URL(urlStr);
      const isHttps = u.protocol === 'https:';
      const lib = isHttps ? https : http;
      const req = lib.request(
        {
          method: 'POST',
          hostname: u.hostname,
          port: u.port || (isHttps ? 443 : 80),
          path: u.pathname + u.search,
          headers: {
            'Content-Type': 'application/json'
          }
        },
        (res) => {
          let data = '';
          res.on('data', (chunk) => (data += chunk));
          res.on('end', () => resolve({ status: res.statusCode || 0, text: data }));
        }
      );
      req.on('error', reject);
      req.write(JSON.stringify(body));
      req.end();
    } catch (e) {
      reject(e);
    }
  });
}

async function main() {
  const args = parseArgs(process.argv);
  const featuresDir = (args['features'] as string) || path.join(process.cwd(), 'features');
  const out = (args['out'] as string) || path.join(process.cwd(), 'cucumber-index.json');
  const suiteId = (args['suiteId'] as string) || 'demo-ecommerce';
  const backend = (args['backend'] as string) || '';

  if (!fs.existsSync(featuresDir)) {
    console.error(`Features directory not found: ${featuresDir}`);
    process.exit(2);
  }

  const files = walkFeatures(featuresDir);
  const features: CucumberFeatureRef[] = [];
  for (const f of files) {
    const parsed = parseFeatureFile(f, featuresDir);
    if (parsed) features.push(parsed);
  }

  const index: CucumberIndex = {
    framework: 'cucumber',
    suiteId,
    features
  };

  fs.mkdirSync(path.dirname(out), { recursive: true });
  fs.writeFileSync(out, JSON.stringify(index, null, 2), 'utf8');
  console.log(`Wrote index for suite '${suiteId}' with ${features.length} features to ${out}`);

  if (backend) {
    const url = backend.endsWith('/') ? `${backend}api/cucumber/index` : `${backend}/api/cucumber/index`;
    try {
      const res = await postJSON(url, index);
      if (res.status >= 200 && res.status < 300) {
        console.log(`POSTed index to ${url} – status ${res.status}`);
      } else {
        console.error(`Failed to POST index to ${url} – status ${res.status} – body: ${res.text}`);
        process.exitCode = 1;
      }
    } catch (err: any) {
      console.error(`Error POSTing index to ${url}:`, err?.message || err);
      process.exitCode = 1;
    }
  }
}

if (require.main === module) {
  // eslint-disable-next-line @typescript-eslint/no-floating-promises
  main();
}
