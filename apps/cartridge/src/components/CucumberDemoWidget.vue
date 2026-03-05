<template>
  <div class="space-y-4">
    <h2 class="text-xl font-semibold flex items-center gap-2">
      <span>Cucumber Suite:</span>
      <span class="font-mono text-emerald-400">{{ currentSuiteId }}</span>
    </h2>

    <p class="text-sm text-slate-400">
      Expand a feature to see its scenarios. Click a scenario row to toggle its steps, then hit
      <span class="font-semibold">Run</span> to start a run.
    </p>

    <div v-for="idx in store.indexes || []" :key="idx.suiteId">
      <div
        v-for="feature in idx.features"
        :key="feature.id"
        class="border border-slate-800 rounded-lg overflow-hidden bg-slate-900/60 mb-4"
      >
        <div
          role="button"
          tabindex="0"
          class="w-full flex items-center justify-between px-5 py-4 bg-slate-900 hover:bg-slate-800 text-left"
          @click="toggleFeature(feature.id)"
        >
          <div class="min-w-0 flex-1 mr-4">
            <div class="font-medium text-slate-50 truncate" :title="feature.name">{{ feature.name }}</div>
            <div v-if="feature.description" class="text-xs text-slate-400 truncate" :title="feature.description">{{ feature.description }}</div>
          </div>
          <div class="flex flex-col items-end text-xs text-slate-400">
            <span class="font-mono">{{ feature.file }}</span>
            <div class="flex items-center gap-3 mt-2" @click.stop>
                <select v-model="selectedExecutor[feature.id]" class="text-[10px] bg-slate-800 text-slate-200 px-1 py-0.5 rounded">
                    <option value="docker">Docker</option>
                    <option value="k8s">Kubernetes</option>
                </select>
                <button class="px-2 py-0.5 text-[10px] rounded bg-emerald-600 hover:bg-emerald-500 text-slate-100 font-semibold" @click.stop="onRunFeature(idx.suiteId, feature)">
                    Run Feature
                </button>
            </div>
            <span class="mt-1">
              <span class="inline-flex items-center justify-center w-5 h-5 rounded-full border border-slate-500 text-[10px]">
                {{ isFeatureOpen(feature.id) ? '-' : '+' }}
              </span>
            </span>
          </div>
        </div>

        <transition name="fade">
          <div v-if="isFeatureOpen(feature.id)" class="p-3 space-y-2 bg-slate-950/70">
            <div v-for="scenario in feature.scenarios" :key="scenarioKey(feature, scenario)" class="rounded-md bg-slate-900">
              <div
                class="flex flex-col md:flex-row md:items-center justify-between gap-4 px-4 py-3 cursor-pointer"
                @click="toggleScenario(feature, scenario)"
              >
                <div>
                  <div class="font-medium text-slate-50 flex items-center gap-2">
                    <span>{{ scenario.name }}</span>
                    <span class="text-[10px] text-slate-500">(click to {{ isScenarioOpenKey(feature, scenario) ? 'hide' : 'show' }} steps)</span>
                  </div>
                  <div class="text-xs text-slate-500 flex flex-wrap items-center gap-2 mt-1">
                    <span class="font-mono">{{ scenario.file }}:{{ scenario.line }}</span>
                    <span v-if="scenario.tags && scenario.tags.length" class="inline-flex flex-wrap gap-1">
                      <span v-for="tag in scenario.tags" :key="tag" class="px-1.5 py-0.5 rounded bg-slate-800 text-[10px] font-mono text-slate-300">{{ tag }}</span>
                    </span>
                  </div>
                </div>

                <div class="flex items-center gap-3">
                  <StatusBadge :status="runStatusForScenario(scenario)" />

                  <select v-model="selectedExecutor[scenario.id]" class="text-xs bg-slate-800 text-slate-200 px-2 py-1 rounded">
                    <option value="docker">Docker</option>
                    <option value="k8s">Kubernetes</option>
                  </select>

                  <button
                    class="px-3 py-1 text-sm rounded bg-emerald-500 hover:bg-emerald-400 text-slate-900 font-semibold"
                    @click.stop="onRun(idx.suiteId, scenario)"
                  >
                    Run
                  </button>

                  <a v-if="reportUrlFor(scenario.id) !== '#'" :href="reportUrlFor(scenario.id)" target="_blank" class="text-xs text-cyan-300 hover:underline" @click.stop>View report</a>
                </div>
              </div>

              <transition name="fade">
                <div v-if="isScenarioOpenKey(feature, scenario)" class="px-5 pb-3 space-y-1 border-t border-slate-800 bg-slate-950/80">
                  <div class="text-xs text-slate-400 mb-1">Steps:</div>
                  <div v-for="st in scenario.steps || []" :key="st.text + st.line" class="text-sm text-slate-200">
                    {{ st.text }}
                  </div>
                </div>
              </transition>
            </div>
          </div>
        </transition>
      </div>
    </div>

    <RunModal v-if="showRunModal" :run="selectedRun" :scenario="selectedScenario" @close="closeRunModal" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue';
import useCucumberStore from '../stores/cucumber';
import type { CucumberScenarioRef } from '../api/cucumber';
import { useRunsStore } from '../stores/runs';
import useToastStore from '../stores/toast';
import StatusBadge from './StatusBadge.vue';
import RunModal from './RunModal.vue';

const store = useCucumberStore();
const runsStore = useRunsStore();
const toast = useToastStore();

onMounted(() => {
  store.load();
  // fetch server config to get default executor
  (async () => {
    try {
      const res = await fetch('/api/config');
      if (!res.ok) return;
      const js = await res.json();
      if (js && js.defaultExecutor) {
        defaultExecutor.value = js.defaultExecutor === 'k8s' ? 'k8s' : 'docker';
      }
    } catch (e) {
      // ignore
    }
  })();
});

const expandedFeatures = ref<Set<string>>(new Set());
const expandedScenarios = ref<Set<string>>(new Set());

// per-scenario executor selection map
const selectedExecutor = ref<Record<string, 'docker' | 'k8s'>>({});
// server default executor (fetched from /api/config)
const defaultExecutor = ref<'docker' | 'k8s'>('k8s');

const showRunModal = ref(false);
const selectedRunId = ref<string | null>(null);
const selectedScenario = ref<CucumberScenarioRef | null>(null);
const selectedRun = computed(() => {
  if (!selectedRunId.value) return null;
  return runsStore.getRunById(selectedRunId.value) || null;
});

// initialize selectedExecutor entries whenever indexes load OR defaultExecutor changes
watch(
  [() => store.indexes, defaultExecutor],
  ([idx, defExec]) => {
    if (!idx) return;
    for (const suite of idx) {
      for (const feature of suite.features || []) {
        if (!(feature.id in selectedExecutor.value)) {
          selectedExecutor.value[feature.id] = defExec;
        }
        for (const scenario of feature.scenarios || []) {
          if (!(scenario.id in selectedExecutor.value)) {
            selectedExecutor.value[scenario.id] = defExec;
          }
        }
      }
    }
  },
  { immediate: true }
);

function isFeatureOpen(id: string) {
  return expandedFeatures.value.has(id);
}
function toggleFeature(id: string) {
  const set = new Set(expandedFeatures.value);
  if (set.has(id)) set.delete(id);
  else set.add(id);
  expandedFeatures.value = set;
}

function scenarioKey(feature: any, scenario: CucumberScenarioRef) {
  return `${feature.id}:${scenario.file}:${scenario.line}`;
}

function isScenarioOpenKey(feature: any, scenario: CucumberScenarioRef) {
  return expandedScenarios.value.has(scenarioKey(feature, scenario));
}

function toggleScenario(feature: any, scenario: CucumberScenarioRef) {
  const key = scenarioKey(feature, scenario);
  const set = new Set(expandedScenarios.value);
  if (set.has(key)) set.delete(key);
  else set.add(key);
  expandedScenarios.value = set;
}

const latestRunsByScenario = computed(() => {
  const map: Record<string, any> = {};
  const all = Object.values(runsStore.runs);
  // sort ascending by time so that iterating assigns the latest one last
  all.sort((a, b) => {
    const tA = a.startedAt ? new Date(a.startedAt).getTime() : 0;
    const tB = b.startedAt ? new Date(b.startedAt).getTime() : 0;
    return tA - tB;
  });
  for (const r of all) {
    if (r.scenarioId) {
      map[r.scenarioId] = r;
    }
  }
  return map;
});

function runStatusForScenario(scenario: CucumberScenarioRef) {
  const found = latestRunsByScenario.value[scenario.id];
  if (!found) return 'idle' as const;
  return (found.status as 'running' | 'passed' | 'failed' | 'idle');
}

function reportUrlFor(scenarioId: string) {
  const found = latestRunsByScenario.value[scenarioId];
  if (!found) return '#';
  let url = found.reportUrl || '#';
  if (url.startsWith('/')) return `http://localhost:9001${url}`;
  return url;
}

async function onRun(suiteId: string, scenario: CucumberScenarioRef) {
  // prefer per-scenario selection, fall back to server default
  const execChoice = (selectedExecutor.value && selectedExecutor.value[scenario.id]) || defaultExecutor.value || 'docker';
  const req = { framework: 'cucumber', suiteId, scenarioId: scenario.id, executor: execChoice };
  try {
    const run = await runsStore.createRun(req);
    // show modal for this run; use run id and compute the latest run from the store
    selectedRunId.value = run.id;
    selectedScenario.value = scenario;
    showRunModal.value = true;
    // start polling is handled by the runs store
  } catch (e: any) {
    console.error('failed to start run', e);
    toast.add(`Failed to start run: ${e?.message || e}`, 'error', 6000);
  }
}

async function onRunFeature(suiteId: string, feature: any) {
  const execChoice = (selectedExecutor.value && selectedExecutor.value[feature.id]) || defaultExecutor.value || 'docker';
  const req = { framework: 'cucumber', suiteId, scenarioId: feature.id, executor: execChoice };
  try {
    const run = await runsStore.createRun(req);
    const dummyScenario = {
        id: feature.id,
        name: feature.name,
        file: feature.file,
        line: 0,
        tags: [],
        steps: []
    };
    selectedRunId.value = run.id;
    selectedScenario.value = dummyScenario;
    showRunModal.value = true;
  } catch (e: any) {
    console.error('failed to start run', e);
    toast.add(`Failed to start run: ${e?.message || e}`, 'error', 6000);
  }
}

function closeRunModal() {
  showRunModal.value = false;
  selectedRunId.value = null;
  selectedScenario.value = null;
}

const currentSuiteId = computed(() => (store.indexes && store.indexes.length ? store.indexes[0].suiteId : ''));
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
