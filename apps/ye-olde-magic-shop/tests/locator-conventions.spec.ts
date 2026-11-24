/**
 * @issue golden-test-ids
 * AC: every locator string under lib/site/ is a [data-testid="..."] attribute
 * selector, a role-based locator (getByRole/getByTestId), or a bare tag used
 * outside any registerElement call — never has-text(, :text(, a CSS :nth-
 * pseudo-class, or a selector string built from element nesting (a parent
 * selector combined with a nested element via a descendant combinator).
 *
 * Static source scan, not a browser test — it fails the moment a locator
 * string regresses to one of the forbidden forms, without needing the app
 * running. The behavioral proof that the surviving locators actually resolve
 * is the existing checkout (checkout.feature), cart (cart.spec.ts) and
 * item-details (item-details.feature) suites.
 */
import { test, expect } from '@playwright/test';
import { readFileSync, readdirSync } from 'node:fs';
import { join } from 'node:path';

const SITE_DIR = join(__dirname, '..', 'lib', 'site');

function findTsFiles(dir: string): string[] {
  return readdirSync(dir, { withFileTypes: true }).flatMap((entry) => {
    const entryPath = join(dir, entry.name);
    if (entry.isDirectory()) return findTsFiles(entryPath);
    return entry.name.endsWith('.ts') ? [entryPath] : [];
  });
}

const STRING_LITERAL = /'(?:[^'\\]|\\.)*'|"(?:[^"\\]|\\.)*"|`(?:[^`\\]|\\.)*`/g;

type Violation = { file: string; literal: string; reason: string };

function findViolations(filePath: string, source: string): Violation[] {
  const violations: Violation[] = [];
  for (const match of source.matchAll(STRING_LITERAL)) {
    const literal = match[0];
    const content = literal.slice(1, -1);

    if (content.includes('has-text(')) {
      violations.push({ file: filePath, literal, reason: 'uses has-text(' });
    } else if (content.includes(':text(')) {
      violations.push({ file: filePath, literal, reason: 'uses :text(' });
    } else if (/:nth-[a-z]/.test(content)) {
      violations.push({ file: filePath, literal, reason: 'uses a CSS :nth- pseudo-selector' });
    } else if (/\]\s+\S/.test(content)) {
      violations.push({ file: filePath, literal, reason: 'nests a selector via a descendant combinator' });
    }
  }
  return violations;
}

test.describe('locator conventions in lib/site/ (golden-test-ids AC1, AC3)', () => {
  test('no locator string uses has-text, :text, :nth-, or element nesting', () => {
    const violations = findTsFiles(SITE_DIR).flatMap((filePath) =>
      findViolations(filePath, readFileSync(filePath, 'utf-8'))
    );

    expect(violations).toEqual([]);
  });

  test("CartPage's complete-transaction button is located by data-testid=\"complete-transaction-btn\" (AC2)", () => {
    const source = readFileSync(join(SITE_DIR, 'CartPage.ts'), 'utf-8');

    expect(source).toContain(
      "this.registerElement('completeTransactionBtn', ButtonElement, '[data-testid=\"complete-transaction-btn\"]')"
    );
  });

  test("HomePage's clickFirstProductDetails uses a flat data-testid locator, not nested item-card markup", () => {
    const source = readFileSync(join(SITE_DIR, 'HomePage.ts'), 'utf-8');
    const methodMatch = source.match(/async clickFirstProductDetails\(\) \{([\s\S]*?)\n {2}\}/);

    expect(methodMatch).not.toBeNull();
    expect(methodMatch![1]).toContain('[data-testid="item-name"]');
    expect(methodMatch![1]).not.toContain('item-card');
  });
});
