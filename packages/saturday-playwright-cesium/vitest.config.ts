import { defineConfig } from 'vitest/config';

// Node-only: the geometry module takes no browser, so this suite runs
// without Playwright or a DOM. Ratified in
// docs/features/saturday-playwright-cesium/architecture-notes.md: kept
// distance/lineLength/pointInCircle package-local and dependency-free.
export default defineConfig({
  test: {
    globals: true,
    environment: 'node',
    include: ['tests/geometry.spec.ts'],
  },
});
