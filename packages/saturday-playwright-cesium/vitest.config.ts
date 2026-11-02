import { defineConfig } from 'vitest/config';

// Node-only: the geometry module takes no browser, so this suite runs
// without Playwright or a DOM, per architecture-notes.md's decision to keep
// distance/lineLength/pointInCircle package-local and dependency-free.
export default defineConfig({
  test: {
    globals: true,
    environment: 'node',
    include: ['tests/geometry.spec.ts'],
  },
});
