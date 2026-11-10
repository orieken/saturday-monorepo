/**
 * @issue local-fonts
 * AC: Cinzel (400, 700) and Inter (400, 600) load from the app's own origin
 * via @fontsource/cinzel and @fontsource/inter, and render as the named
 * families rather than a system fallback.
 *
 * Uses document.fonts.load, not document.fonts.check: check returned true here
 * even with the font imports removed. load returns the faces it matched, and
 * can match a neighbouring weight, so the weight is compared as well.
 */
import { test, expect } from '@playwright/test';

const FONTS = [
  { family: 'Cinzel', weight: '400' },
  { family: 'Cinzel', weight: '700' },
  { family: 'Inter', weight: '400' },
  { family: 'Inter', weight: '600' },
];

test('headings render in self-hosted Cinzel and body text in self-hosted Inter', async ({ page }) => {
  await page.goto('/');

  const missing = await page.evaluate(async (fonts) => {
    const results = await Promise.all(
      fonts.map(async ({ family, weight }) => {
        const faces = await document.fonts.load(`${weight} 1em ${family}`);
        return faces.some((face) => face.weight === weight) ? null : `${family} ${weight}`;
      })
    );
    return results.filter(Boolean);
  }, FONTS);

  expect(missing).toEqual([]);
});
