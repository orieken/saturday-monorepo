<template>
  <div class="cucumber-report">
    <!-- Loading State -->
    <div v-if="loading" class="text-center py-12">
      <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-sky-500 mb-4"></div>
      <div class="text-slate-400">Loading report...</div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="text-center py-12">
      <div class="text-rose-400 mb-2">
        <svg class="inline-block w-12 h-12 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>
      <div class="text-slate-300 font-medium">{{ error }}</div>
    </div>

    <!-- Empty State -->
    <div v-else-if="!report || report.length === 0" class="text-center py-12">
      <div class="text-slate-500 mb-2">
        <svg class="inline-block w-12 h-12 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
      </div>
      <div class="text-slate-400">No report data available</div>
      <div class="text-slate-500 text-sm mt-2">The test may not have generated a report yet</div>
    </div>

    <!-- Report Content -->
    <div v-else class="space-y-6">
      <!-- Report Summary -->
      <div v-if="stats" class="bg-slate-800 rounded-lg p-6 border border-slate-700">
        <h3 class="text-lg font-semibold text-slate-100 mb-4">Test Summary</h3>
        <div class="grid grid-cols-2 md:grid-cols-5 gap-4">
          <div class="text-center">
            <div class="text-2xl font-bold text-slate-300">{{ stats.totalScenarios }}</div>
            <div class="text-xs text-slate-500 uppercase">Total</div>
          </div>
          <div class="text-center">
            <div class="text-2xl font-bold text-emerald-400">{{ stats.passedScenarios }}</div>
            <div class="text-xs text-slate-500 uppercase">Passed</div>
          </div>
          <div class="text-center">
            <div class="text-2xl font-bold text-rose-400">{{ stats.failedScenarios }}</div>
            <div class="text-xs text-slate-500 uppercase">Failed</div>
          </div>
          <div class="text-center">
            <div class="text-2xl font-bold text-slate-400">{{ stats.skippedScenarios }}</div>
            <div class="text-xs text-slate-500 uppercase">Skipped</div>
          </div>
          <div class="text-center">
            <div class="text-2xl font-bold" :class="stats.successRate >= 80 ? 'text-emerald-400' : stats.successRate >= 50 ? 'text-amber-400' : 'text-rose-400'">
              {{ stats.successRate.toFixed(0) }}%
            </div>
            <div class="text-xs text-slate-500 uppercase">Success Rate</div>
          </div>
        </div>
      </div>

      <!-- Features -->
      <div v-for="feature in report" :key="feature.id" class="bg-slate-800 rounded-lg p-6 border border-slate-700">
        <!-- Feature Header -->
        <div class="mb-6">
          <div class="flex items-start gap-3 mb-2">
            <span class="text-xs font-mono text-sky-400 mt-1">{{ feature.keyword }}</span>
            <h3 class="text-xl font-semibold text-slate-100 flex-1">{{ feature.name }}</h3>
          </div>
          <p v-if="feature.description" class="text-sm text-slate-400 ml-16 whitespace-pre-wrap">{{ feature.description }}</p>
          <div v-if="feature.tags && feature.tags.length > 0" class="flex flex-wrap gap-2 mt-3 ml-16">
            <span v-for="tag in feature.tags" :key="tag.name" 
                  class="text-xs px-2 py-1 rounded bg-slate-700 text-slate-300 font-mono">
              {{ tag.name }}
            </span>
          </div>
        </div>

        <!-- Scenarios -->
        <div class="space-y-4">
          <div v-for="scenario in feature.elements" :key="scenario.id" 
               class="border-l-4 pl-6 py-3 rounded-r bg-slate-900/50"
               :class="getScenarioBorderClass(scenario)">
            <!-- Scenario Header -->
            <div class="flex items-center justify-between mb-3">
              <div class="flex-1">
                <div class="flex items-center gap-2">
                  <span class="text-xs font-mono text-slate-500">{{ scenario.keyword }}</span>
                  <span class="text-slate-200 font-medium">{{ scenario.name }}</span>
                </div>
                <div v-if="scenario.tags && scenario.tags.length > 0" class="flex flex-wrap gap-1 mt-2">
                  <span v-for="tag in scenario.tags" :key="tag.name" 
                        class="text-xs px-1.5 py-0.5 rounded bg-slate-800 text-slate-400 font-mono">
                    {{ tag.name }}
                  </span>
                </div>
              </div>
              <span class="text-xs px-3 py-1 rounded font-medium ml-4"
                    :class="getStatusBadgeClass(getScenarioStatus(scenario))">
                {{ getScenarioStatus(scenario).toUpperCase() }}
              </span>
            </div>

            <!-- Steps -->
            <div class="space-y-2">
              <div v-for="(step, idx) in scenario.steps" :key="idx"
                   class="flex items-start gap-3 text-sm py-2 px-3 rounded"
                   :class="getStepBackgroundClass(step.result.status)">
                <!-- Status Icon -->
                <span class="flex-shrink-0 w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold mt-0.5"
                      :class="getStepIconClass(step.result.status)">
                  <span v-if="step.result.status === 'passed'">✓</span>
                  <span v-else-if="step.result.status === 'failed'">✗</span>
                  <span v-else-if="step.result.status === 'skipped'">−</span>
                  <span v-else-if="step.result.status === 'pending'">⋯</span>
                  <span v-else>?</span>
                </span>

                <!-- Step Content -->
                <div class="flex-1 min-w-0">
                  <div :class="getStepTextClass(step.result.status)">
                    <span class="font-mono text-xs mr-2 text-slate-500">{{ step.keyword }}</span>
                    <span class="font-medium">{{ step.name }}</span>
                    <span v-if="step.result.duration" class="text-xs text-slate-500 ml-3">
                      {{ formatDuration(step.result.duration) }}
                    </span>
                  </div>

                  <!-- Error Message -->
                  <div v-if="step.result.error_message" 
                       class="mt-3 p-4 bg-rose-950/50 border border-rose-900 rounded text-xs font-mono text-rose-300 whitespace-pre-wrap overflow-x-auto">
                    {{ step.result.error_message }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import type { CucumberReport, CucumberScenario, ScenarioStatus, ReportStats } from '../types/cucumber';
import { fetchRunReport } from '../api/runs';

const props = defineProps<{
  runId: string;
}>();

const report = ref<CucumberReport | null>(null);
const loading = ref(true);
const error = ref<string | null>(null);

// Computed statistics
const stats = computed<ReportStats | null>(() => {
  if (!report.value || report.value.length === 0) return null;

  let totalScenarios = 0;
  let passedScenarios = 0;
  let failedScenarios = 0;
  let skippedScenarios = 0;
  let pendingScenarios = 0;
  let totalDuration = 0;

  for (const feature of report.value) {
    for (const scenario of feature.elements) {
      totalScenarios++;
      const status = getScenarioStatus(scenario);
      
      if (status === 'passed') passedScenarios++;
      else if (status === 'failed') failedScenarios++;
      else if (status === 'skipped') skippedScenarios++;
      else if (status === 'pending') pendingScenarios++;

      // Sum up step durations
      for (const step of scenario.steps) {
        if (step.result.duration) {
          totalDuration += step.result.duration / 1000000; // Convert nanoseconds to milliseconds
        }
      }
    }
  }

  const successRate = totalScenarios > 0 ? (passedScenarios / totalScenarios) * 100 : 0;

  return {
    totalScenarios,
    passedScenarios,
    failedScenarios,
    skippedScenarios,
    pendingScenarios,
    totalDuration,
    successRate,
  };
});

onMounted(async () => {
  try {
    report.value = await fetchRunReport(props.runId);
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Failed to load report';
  } finally {
    loading.value = false;
  }
});

function getScenarioStatus(scenario: CucumberScenario): ScenarioStatus {
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
  if (status === 'skipped') return 'border-slate-600';
  if (status === 'pending') return 'border-amber-500';
  return 'border-slate-700';
}

function getStatusBadgeClass(status: ScenarioStatus): string {
  if (status === 'passed') return 'bg-emerald-900/50 text-emerald-300 border border-emerald-700';
  if (status === 'failed') return 'bg-rose-900/50 text-rose-300 border border-rose-700';
  if (status === 'skipped') return 'bg-slate-700 text-slate-300 border border-slate-600';
  if (status === 'pending') return 'bg-amber-900/50 text-amber-300 border border-amber-700';
  return 'bg-slate-700 text-slate-400 border border-slate-600';
}

function getStepIconClass(status: string): string {
  if (status === 'passed') return 'bg-emerald-900/70 text-emerald-300 border border-emerald-700';
  if (status === 'failed') return 'bg-rose-900/70 text-rose-300 border border-rose-700';
  if (status === 'skipped') return 'bg-slate-700 text-slate-400 border border-slate-600';
  if (status === 'pending') return 'bg-amber-900/70 text-amber-300 border border-amber-700';
  return 'bg-slate-700 text-slate-400 border border-slate-600';
}

function getStepBackgroundClass(status: string): string {
  if (status === 'failed') return 'bg-rose-950/20';
  return 'bg-transparent';
}

function getStepTextClass(status: string): string {
  if (status === 'passed') return 'text-slate-300';
  if (status === 'failed') return 'text-rose-200';
  if (status === 'skipped') return 'text-slate-500';
  if (status === 'pending') return 'text-amber-300';
  return 'text-slate-400';
}

function formatDuration(ns: number): string {
  const ms = ns / 1000000;
  if (ms < 1000) return `${ms.toFixed(0)}ms`;
  return `${(ms / 1000).toFixed(2)}s`;
}
</script>

<style scoped>
.cucumber-report {
  @apply w-full;
}
</style>
