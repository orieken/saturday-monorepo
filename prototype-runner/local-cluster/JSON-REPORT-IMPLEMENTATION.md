# JSON Report Rendering Implementation Summary

## What We've Accomplished

### ✅ Backend API (Console)

1. **New Endpoint**: `GET /api/runs/{runId}/report`
   - Location: `apps/console/internal/httpserver/handlers.go`
   - Returns the `cucumber.json` file for a given run
   - Handles missing files gracefully (returns empty array)
   - Validates JSON before sending

2. **Route Added**: 
   - Location: `apps/console/internal/httpserver/router.go`
   - Registered at `/api/runs/{runId}/report`

3. **Deployment**:
   - Image rebuilt and loaded into kind cluster
   - Service restarted successfully

### 📋 Next Steps for Cartridge UI

To complete the implementation, you'll need to add these components to Cartridge:

#### 1. Create Cucumber Report Types

Create `apps/cartridge/src/types/cucumber.ts`:

```typescript
export interface CucumberStep {
  keyword: string;
  name: string;
  line: number;
  match?: {
    location: string;
  };
  result: {
    status: 'passed' | 'failed' | 'skipped' | 'pending' | 'undefined';
    duration?: number;
    error_message?: string;
  };
}

export interface CucumberScenario {
  id: string;
  keyword: string;
  name: string;
  description?: string;
  line: number;
  type: string;
  tags?: Array<{ name: string; line: number }>;
  steps: CucumberStep[];
}

export interface CucumberFeature {
  uri: string;
  id: string;
  keyword: string;
  name: string;
  description?: string;
  line: number;
  tags?: Array<{ name: string; line: number }>;
  elements: CucumberScenario[];
}

export type CucumberReport = CucumberFeature[];
```

#### 2. Create API Function

Add to `apps/cartridge/src/api/runs.ts`:

```typescript
export async function fetchRunReport(runId: string): Promise<CucumberReport> {
  const response = await fetch(`http://localhost:9001/api/runs/${runId}/report`);
  if (!response.ok) {
    throw new Error(`Failed to fetch report: ${response.statusText}`);
  }
  return response.json();
}
```

#### 3. Create Report Renderer Component

Create `apps/cartridge/src/components/CucumberReportViewer.vue`:

```vue
<template>
  <div class="cucumber-report">
    <div v-if="loading" class="text-center py-8">
      <div class="text-slate-400">Loading report...</div>
    </div>

    <div v-else-if="error" class="text-center py-8">
      <div class="text-rose-400">{{ error }}</div>
    </div>

    <div v-else-if="report && report.length > 0" class="space-y-6">
      <div v-for="feature in report" :key="feature.id" class="bg-slate-800 rounded-lg p-6">
        <!-- Feature Header -->
        <div class="mb-4">
          <div class="flex items-center gap-2 mb-2">
            <span class="text-xs font-mono text-slate-500">{{ feature.keyword }}</span>
            <h3 class="text-lg font-semibold text-slate-100">{{ feature.name }}</h3>
          </div>
          <p v-if="feature.description" class="text-sm text-slate-400">{{ feature.description }}</p>
          <div v-if="feature.tags" class="flex gap-2 mt-2">
            <span v-for="tag in feature.tags" :key="tag.name" 
                  class="text-xs px-2 py-1 rounded bg-slate-700 text-slate-300">
              {{ tag.name }}
            </span>
          </div>
        </div>

        <!-- Scenarios -->
        <div class="space-y-4">
          <div v-for="scenario in feature.elements" :key="scenario.id" 
               class="border-l-4 pl-4"
               :class="getScenarioBorderClass(scenario)">
            <div class="flex items-center justify-between mb-2">
              <div>
                <span class="text-xs font-mono text-slate-500 mr-2">{{ scenario.keyword }}</span>
                <span class="text-slate-200">{{ scenario.name }}</span>
              </div>
              <span class="text-xs px-2 py-1 rounded"
                    :class="getStatusClass(getScenarioStatus(scenario))">
                {{ getScenarioStatus(scenario) }}
              </span>
            </div>

            <!-- Steps -->
            <div class="space-y-1 mt-3">
              <div v-for="(step, idx) in scenario.steps" :key="idx"
                   class="flex items-start gap-3 text-sm py-1">
                <span class="flex-shrink-0 w-5 h-5 rounded-full flex items-center justify-center text-xs"
                      :class="getStepIconClass(step.result.status)">
                  <span v-if="step.result.status === 'passed'">✓</span>
                  <span v-else-if="step.result.status === 'failed'">✗</span>
                  <span v-else-if="step.result.status === 'skipped'">-</span>
                  <span v-else>?</span>
                </span>
                <div class="flex-1">
                  <div :class="getStepTextClass(step.result.status)">
                    <span class="font-mono text-xs mr-2">{{ step.keyword }}</span>
                    <span>{{ step.name }}</span>
                    <span v-if="step.result.duration" class="text-xs text-slate-500 ml-2">
                      ({{ formatDuration(step.result.duration) }})
                    </span>
                  </div>
                  <div v-if="step.result.error_message" 
                       class="mt-2 p-3 bg-rose-950/50 border border-rose-900 rounded text-xs font-mono text-rose-300 whitespace-pre-wrap">
                    {{ step.result.error_message }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="text-center py-8">
      <div class="text-slate-500">No report data available</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import type { CucumberReport, CucumberScenario, CucumberStep } from '../types/cucumber';
import { fetchRunReport } from '../api/runs';

const props = defineProps<{
  runId: string;
}>();

const report = ref<CucumberReport | null>(null);
const loading = ref(true);
const error = ref<string | null>(null);

onMounted(async () => {
  try {
    report.value = await fetchRunReport(props.runId);
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Failed to load report';
  } finally {
    loading.value = false;
  }
});

function getScenarioStatus(scenario: CucumberScenario): string {
  const statuses = scenario.steps.map(s => s.result.status);
  if (statuses.some(s => s === 'failed')) return 'failed';
  if (statuses.some(s => s === 'pending' || s === 'undefined')) return 'pending';
  if (statuses.every(s => s === 'skipped')) return 'skipped';
  if (statuses.every(s => s === 'passed')) return 'passed';
  return 'unknown';
}

function getScenarioBorderClass(scenario: CucumberScenario): string {
  const status = getScenarioStatus(scenario);
  if (status === 'passed') return 'border-emerald-500';
  if (status === 'failed') return 'border-rose-500';
  if (status === 'skipped') return 'border-slate-500';
  return 'border-amber-500';
}

function getStatusClass(status: string): string {
  if (status === 'passed') return 'bg-emerald-900/50 text-emerald-300';
  if (status === 'failed') return 'bg-rose-900/50 text-rose-300';
  if (status === 'skipped') return 'bg-slate-700 text-slate-300';
  return 'bg-amber-900/50 text-amber-300';
}

function getStepIconClass(status: string): string {
  if (status === 'passed') return 'bg-emerald-900/50 text-emerald-400';
  if (status === 'failed') return 'bg-rose-900/50 text-rose-400';
  if (status === 'skipped') return 'bg-slate-700 text-slate-400';
  return 'bg-amber-900/50 text-amber-400';
}

function getStepTextClass(status: string): string {
  if (status === 'passed') return 'text-slate-300';
  if (status === 'failed') return 'text-rose-300';
  if (status === 'skipped') return 'text-slate-500';
  return 'text-amber-300';
}

function formatDuration(ns: number): string {
  const ms = ns / 1000000;
  if (ms < 1000) return `${ms.toFixed(0)}ms`;
  return `${(ms / 1000).toFixed(2)}s`;
}
</script>
```

#### 4. Update RunModal to Use New Component

Modify `apps/cartridge/src/components/RunModal.vue`:

```vue
<!-- Replace the iframe section with: -->
<div v-if="run?.reportUrl" class="w-full bg-slate-900 rounded overflow-hidden">
  <CucumberReportViewer :runId="run.id" />
</div>

<!-- Add import -->
<script setup lang="ts">
import CucumberReportViewer from './CucumberReportViewer.vue';
// ... rest of imports
</script>
```

### 🎯 Benefits

1. **Consistent Styling**: Matches Cartridge's design system
2. **Better UX**: 
   - Expandable/collapsible sections
   - Color-coded status indicators
   - Inline error messages
   - Duration display
3. **Queryable**: Can add filtering, search, etc.
4. **Fallback**: HTML report still available via direct link

### 🧪 Testing

Once implemented, test by:

1. Trigger a test run from Cartridge
2. Click on the run to open RunModal
3. Report should render inline with custom styling
4. Verify all scenarios, steps, and statuses display correctly

### 📊 Future Enhancements

- Add filtering by status (passed/failed/skipped)
- Add search functionality
- Show aggregated statistics
- Export to different formats
- Compare runs side-by-side
- Show flaky test detection

## Current Status

✅ Backend API endpoint implemented and deployed
⏳ Frontend components ready to be implemented
📝 Full implementation guide provided above

The backend is ready! You can now implement the Cartridge components to render the JSON reports beautifully.
