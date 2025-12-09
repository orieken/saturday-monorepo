#!/usr/bin/env node
  import fs from 'node:fs';
  import path from 'node:path';

  function ensureDir(p: string) { fs.mkdirSync(p, { recursive: true }); }
  function writeIfMissing(p: string, content: string) { if (!fs.existsSync(p)) fs.writeFileSync(p, content, 'utf-8'); }

  const root = process.cwd();
  const e2eApiDir = path.join(root, 'e2e', 'api');
  ensureDir(e2eApiDir);

  const example = `import { expect } from '@playwright/test';
import { test } from '@orieken/saturday-playwright-k6-exporter/fixture';
import { createK6Recorder } from '@orieken/saturday-playwright-k6-exporter';

test('GET/POST to httpbin @k6', async ({}, testInfo) => {
  const setup = await createK6Recorder(testInfo.title);
  if (!setup) test.skip(true, 'Set K6_EXPORT=1 to export k6 script');
  const { ctx, recorder } = setup!;

  const res = await ctx.get('https://httpbin.org/get?hello=world', {
    headers: { Accept: 'application/json' },
    k6Name: 'Get hello world'
  });
  expect(res.status()).toBe(200);

  const post = await ctx.post('https://httpbin.org/post', {
    headers: { 'Content-Type': 'application/json' },
    data: { foo: 'bar' },
    k6Name: 'Post foo bar'
  });
  expect(post.status()).toBe(200);

  await recorder.flushToK6();
  await ctx.dispose();
});
`;

  const dest = path.join(e2eApiDir, 'k6-export-example.spec.ts');
  writeIfMissing(dest, example);

  // Suggest package.json scripts if present
  const pkgPath = path.join(root, 'package.json');
  if (fs.existsSync(pkgPath)) {
    const pkg = JSON.parse(fs.readFileSync(pkgPath, 'utf-8'));
    pkg.scripts = pkg.scripts || {};
    pkg.scripts["k6:export"] = pkg.scripts["k6:export"] || "K6_EXPORT=1 playwright test -g \"@k6\"";
    pkg.scripts["k6:export:all"] = pkg.scripts["k6:export:all"] || "K6_EXPORT=1 playwright test";
    fs.writeFileSync(pkgPath, JSON.stringify(pkg, null, 2));
    console.log("Added scripts k6:export and k6:export:all to package.json");
  }

  console.log("Scaffolded example at e2e/api/k6-export-example.spec.ts");
  console.log("Run: npm run k6:export   # to generate perf/k6/*.k6.js");