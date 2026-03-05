#!/usr/bin/env node
import { spawn } from 'node:child_process';
import { fileURLToPath } from 'node:url';
import path from 'node:path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

function run(cmd: string, args: string[], env: NodeJS.ProcessEnv = {}) {
  return new Promise<void | number>((resolve, reject) => {
    const child = spawn(cmd, args, { stdio: 'inherit', env: { ...process.env, ...env } });
    child.on('close', (code) => code === 0 ? resolve(0) : reject(new Error(`${cmd} exited ${code}`)));
  });
}

const [,, sub, ...rest] = process.argv;

async function main() {
  switch (sub) {
    case 'init':
      // Delegate to the exporter scaffold command if present
      try {
        await run('npx', ['k6-exporter-init']);
      } catch {
        console.error('Could not run k6-exporter-init. Ensure @saturday/playwright-k6-exporter is installed.');
        process.exit(1);
      }
      break;

    case 'export':
      // Pass-through to playwright with K6_EXPORT=1
      // Example: saturday-k6 export --grep @k6
      await run('npx', ['playwright', 'test', ...rest], { K6_EXPORT: '1' });
      break;

    case 'policy':
      // policy ls: show which policy resolves
      if (rest[0] === 'ls') {
        // Try resolution via env override (K6_REDACTION_POLICY) or fallback
        const spec = process.env.K6_REDACTION_POLICY || '@saturday/k6-redaction-basic';
        console.log(`Policy resolution order: options.policy -> K6_REDACTION_POLICY -> @saturday/k6-redaction-basic`);
        console.log(`Candidate: ${spec}`);
        try {
          // eslint-disable-next-line @typescript-eslint/no-var-requires
          const mod = await import(spec);
          const hasFactory = !!(mod.createDefaultRedactionPolicy);
          console.log(`Loaded: ${spec} (factory: ${hasFactory ? 'yes' : 'no'})`);
        } catch (e) {
          console.log(`Could not load: ${spec}`);
        }
      } else {
        console.log('Usage: saturday-k6 policy ls');
      }
      break;

    default:
      console.log(`saturday-k6 <command>
  init                 Scaffold example tests and scripts
  export [--grep ...]  Run Playwright tests with K6_EXPORT=1
  policy ls            Show detected/loaded redaction policy
`);
      process.exit(1);
  }
}

main().catch((e) => { console.error(e); process.exit(1); });
