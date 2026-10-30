#!/usr/bin/env node
// Baked-in WebGL smoke test for the Saturday runner image.
// Launches Chrome and Edge headed, verifies a webgl2 context in each, and
// writes the renderer strings, browser versions, and installed packages to
// /work/evidence/webgl.txt. Exits non-zero if either context is null.
'use strict';

const fs = require('fs');
const path = require('path');
const { execFileSync } = require('child_process');
const { chromium } = require('playwright-core');

const EVIDENCE_PATH = process.env.RUNNER_CHECK_EVIDENCE || '/work/evidence/webgl.txt';

// Headed Chromium under Xvfb still needs an explicit software GL backend in
// a container with no real GPU — without these, webgl2 context creation
// returns null even headed (contrary to new-headless's documented
// --enable-unsafe-swiftshader need, the same gap shows up here).
const SOFTWARE_GL_ARGS = ['--use-gl=angle', '--use-angle=swiftshader'];

const BROWSERS = [
  { label: 'Chrome', channel: 'chrome' },
  { label: 'Edge', channel: 'msedge' },
];

async function checkWebgl2(channel) {
  const browser = await chromium.launch({ channel, headless: false, args: SOFTWARE_GL_ARGS });
  try {
    const page = await browser.newPage();
    const renderer = await page.evaluate(() => {
      const canvas = document.createElement('canvas');
      const gl = canvas.getContext('webgl2');
      if (!gl) return null;
      const info = gl.getExtension('WEBGL_debug_renderer_info');
      return info ? gl.getParameter(info.UNMASKED_RENDERER_WEBGL) : gl.getParameter(gl.RENDERER);
    });
    return { version: browser.version(), renderer };
  } finally {
    await browser.close();
  }
}

function listInstalledPackages() {
  const output = execFileSync('dpkg-query', ['-W', '-f=${Package}=${Version}\n'], {
    encoding: 'utf8',
  });
  return output.trim().split('\n').sort().join('\n');
}

function formatEvidence(results, packages) {
  const lines = [`# runner-check evidence — ${new Date().toISOString()}`, ''];
  for (const result of results) {
    lines.push(`${result.label}:`);
    lines.push(`  version: ${result.version}`);
    lines.push(`  webgl2 renderer: ${result.renderer === null ? 'NONE (context is null)' : result.renderer}`);
    lines.push('');
  }
  lines.push('installed packages:');
  lines.push(packages);
  return lines.join('\n') + '\n';
}

async function runChecks() {
  const results = [];
  for (const browser of BROWSERS) {
    const { version, renderer } = await checkWebgl2(browser.channel);
    results.push({ label: browser.label, version, renderer });
  }
  return results;
}

async function main() {
  const results = await runChecks();
  const packages = listInstalledPackages();
  const evidence = formatEvidence(results, packages);

  fs.mkdirSync(path.dirname(EVIDENCE_PATH), { recursive: true });
  fs.writeFileSync(EVIDENCE_PATH, evidence);

  const failed = results.filter((result) => result.renderer === null);
  if (failed.length > 0) {
    for (const result of failed) {
      console.error(`runner-check: ${result.label} has no webgl2 context`);
    }
    process.exitCode = 1;
    return;
  }
  for (const result of results) {
    console.log(`runner-check: ${result.label} ${result.version} — ${result.renderer}`);
  }
}

main().catch((err) => {
  console.error('runner-check: failed with an error', err);
  process.exitCode = 1;
});
