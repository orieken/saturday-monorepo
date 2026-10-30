---
applyTo: "**/*.vue,**/*.tsx,**/*.jsx"
---
# Vue 3 + Tailwind Frontend Conventions

ALWAYS use Vue 3 Composition API with <script setup>.
NEVER use Options API in new components.
ALWAYS use Tailwind CSS utility classes unless custom CSS is required.
ALWAYS use TypeScript with strict mode and never use raw any types.
ALWAYS co-locate component tests alongside components.
CRITICAL: Components must remain below 100 lines.
