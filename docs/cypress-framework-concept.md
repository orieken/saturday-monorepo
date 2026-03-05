# Architecture Design: Cypress Automation Framework Concept

**Role:** `@architect`
**Date:** 2026-03-03
**Subject:** High-level architectural design for translating Saturday Framework concepts cleanly into a Cypress-native paradigm.

## Executive Summary
Based on the `@analyst` report (`cypress-framework-concept-analysis.md`), we must respect Cypress's unique constraint: **queued, single-tab execution**. Attempting to force Playwright/Node-native async/await patterns (like dynamic decorators and multi-tab managers) into Cypress leads to brittle, unidiomatic tests. 

This document outlines an architectural blueprint for a Cypress framework that achieves the *goals* of the Saturday Framework—high structure, ML-readiness, and observability—using strictly *Cypress-native* patterns.

---

## 1. The Core Architecture: The Cypress Site Façade

We retain the Site-Centric model, but adapt it slightly so it interacts correctly with the Cypress command queue. Tests will never instantiate individual pages; they will instantiate a `Site` that serves as the root boundary.

### Core Structure
```typescript
// Base architecture classes
abstract class CypressBasePage {
    protected path: string;
    public visit() { cy.visit(this.path); }
}

abstract class CypressBaseFlow {
    // Flows interact with multiple pages via the Site
    constructor(protected site: any) {}
}

// Concrete Implementation
class StoreSite {
    // Lazy-loaded getters
    get loginPage() { return new LoginPage(); }
    get checkoutFlow() { return new CheckoutFlow(this); }
    get visualAI() { return new VisualAIFacade(); }
}
```

### Usage in Spec
```typescript
describe('E-Commerce Site', () => {
    const site = new StoreSite();

    it('processes a checkout', () => {
        site.loginPage.visit();
        site.checkoutFlow.executeStandardPurchase('user1');
        site.visualAI.assertGoldenBaseline('checkout_success');
    });
});
```

---

## 2. Handling State & Guards (Replacing Decorators)

Since standard TypeScript decorators (`@RequiresFilter`) evaluate too early in the Cypress queue, we migrate to **Guarded Cypress Commands** and **Assertion-First Page Methods**.

Instead of trying to catch state asynchronously outside of Cypress, we bake state checks into the beginning of the Page Object methods using Cypress's built-in retryability.

```typescript
class CartPage extends CypressBasePage {
    
    // Guarded execution: We rely on Cypress to retry until the guard passes 
    // or the default timeout is reached.
    public checkout() {
        // The Guard
        cy.get('@userSession').should('exist', 'User must be authenticated to checkout.');
        
        // The Action
        cy.get('[data-test="checkout-btn"]').click();
    }
}
```

---

## 3. Element Modeling (Dropping `BaseElement`)

Due to Cypress's chaining mechanics, wrapping elements in custom `BaseElement` classes creates friction. Instead, our framework dictates that **Locators are encapsulated strictly within Page Objects as private getters returning Cypress Chains.**

```typescript
class HeaderComponent {
    // Return the cy.Chainable directly
    private get userProfileMenu() { return cy.get('[data-cy="profile-menu"]'); }
    private get logoutBtn() { return cy.get('[data-cy="logout-btn"]'); }

    public openProfile() {
        this.userProfileMenu.click();
    }
}
```

---

## 4. Machine Learning & Backend Integration (Tasks)

Cypress operates inside the browser, so it cannot natively interact with our ML Python microservices or filesystem models. To bridge this, the framework utilizes a structured library of **Cypress Tasks**.

The `Site` facade will route ML calls through `cy.task()`.

```typescript
class VisualAIFacade {
    public assertGoldenBaseline(imageName: string) {
        // 1. Capture screen within Cypress
        cy.screenshot(`temp/${imageName}`);
        
        // 2. Pass off to Node.js backend for ML Analysis
        cy.task('analyzeVisualAnomaly', { targetImage: imageName }).then((result) => {
            expect(result.isMatch).to.be.true;
        });
    }
}
```

---

## 5. Scope Boundaries (Handling Cypress Limitations)

To effectively replace the Saturday Framework's robust `TabManager` and `SiteManager`, we must establish strict operational boundaries to prevent automation engineers from attempting impossible Cypress operations.

1. **The "Single Tab" Rule:** Tests mapping to multi-tab flows *must* run linearly in one tab. Clicks that pop a new tab must be augmented with `invoke('removeAttr', 'target')` in the Page Object layer.
2. **The "API Spoofing" Rule:** If a test requires simulating a secondary user (e.g., admin approving a document while user waits), the admin action *must* be performed via `cy.request()` to the backend API, rather than attempting to launch an `AdminSite` in a separate browser context.

---

## Summary for Development Team

The `cypress-framework` will provide strong typing, reusable flows, and direct ML integration while remaining idiomatically correct for Cypress. **Do not attempt to port the Playwright Manager ecosystem.** Rely on Cypress plugins for telemetry, `cy.task` for ML, and the Façade pattern for test structure.
