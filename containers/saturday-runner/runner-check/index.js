#!/usr/bin/env node
// Baked-in WebGL smoke test for the Saturday runner image.
// Launches Chrome and Edge headed, verifies a webgl2 context in each, and
// writes the renderer strings, browser versions, and installed packages to
// /work/evidence/webgl.txt. Exits non-zero if either context is null, if
// the display's X pointer is not parked in its corner (see x-pointer.py), or
// if Playwright's ffmpeg, which records videos, is missing.
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
// entrypoint.sh parks the pointer at the bottom-right of Xvfb's 1920x1080
// screen. Left at the centre, it sends stray mouse moves into tests' drags.
const PARKED_POINTER = '1919,1079';

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

// dpkg on the debian target, rpm on the ubi9 target.
function listInstalledPackages() {
  const [command, args] = fs.existsSync('/usr/bin/dpkg-query')
    ? ['dpkg-query', ['-W', '-f=${Package}=${Version}\n']]
    : ['rpm', ['-qa', '--qf', '%{NAME}=%{EPOCH}:%{VERSION}-%{RELEASE}\n']];
  const output = execFileSync(command, args, { encoding: 'utf8' });
  return output.trim().split('\n').sort().join('\n');
}

// The Containerfile installs ffmpeg under PLAYWRIGHT_BROWSERS_PATH. Without it,
// every page fails to open once RUNNER_VIDEO=1 turns recording on.
function findFfmpeg() {
  const root = process.env.PLAYWRIGHT_BROWSERS_PATH || '';
  const dirs = fs.existsSync(root) ? fs.readdirSync(root).filter((name) => name.startsWith('ffmpeg-')) : [];
  const found = dirs.map((dir) => path.join(root, dir, 'ffmpeg-linux')).find((file) => fs.existsSync(file));
  return found || null;
}

function readPointer() {
  return execFileSync('/usr/local/bin/x-pointer.py', ['where'], { encoding: 'utf8' }).trim();
}

function formatEvidence(results, packages, pointer, ffmpeg) {
  const lines = [`# runner-check evidence — ${new Date().toISOString()}`, ''];
  lines.push(`X pointer: ${pointer} (expected ${PARKED_POINTER})`);
  lines.push(`ffmpeg: ${ffmpeg || 'MISSING'}`, '');
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
  const pointer = readPointer();
  const ffmpeg = findFfmpeg();
  const evidence = formatEvidence(results, packages, pointer, ffmpeg);

  fs.mkdirSync(path.dirname(EVIDENCE_PATH), { recursive: true });
  fs.writeFileSync(EVIDENCE_PATH, evidence);

  const failed = results.filter((result) => result.renderer === null);
  for (const result of failed) {
    console.error(`runner-check: ${result.label} has no webgl2 context`);
  }
  if (pointer !== PARKED_POINTER) {
    console.error(`runner-check: X pointer is at ${pointer}, not parked at ${PARKED_POINTER}`);
  }
  if (!ffmpeg) {
    console.error(`runner-check: no ffmpeg under PLAYWRIGHT_BROWSERS_PATH (${process.env.PLAYWRIGHT_BROWSERS_PATH || 'unset'}); videos cannot be recorded`);
  }
  if (failed.length > 0 || pointer !== PARKED_POINTER || !ffmpeg) {
    process.exitCode = 1;
    return;
  }
  for (const result of results) {
    console.log(`runner-check: ${result.label} ${result.version} — ${result.renderer}`);
  }
  console.log(`runner-check: X pointer parked at ${pointer}`);
  console.log(`runner-check: ffmpeg at ${ffmpeg}`);
}

main().catch((err) => {
  console.error('runner-check: failed with an error', err);
  process.exitCode = 1;
});
