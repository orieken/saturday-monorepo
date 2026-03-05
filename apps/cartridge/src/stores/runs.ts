import { defineStore } from 'pinia';
import { ref, watch } from 'vue';
import type { Run, RunRequest } from '../api/runs';
import { startRun, getRun } from '../api/runs';
import useToastStore from './toast';

export const useRunsStore = defineStore('runs', () => {
  const STORAGE_KEY = 'test-runner-ui:runs';

  // load persisted runs from localStorage (best-effort)
  let initial: Record<string, Run> = {};
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (raw) initial = JSON.parse(raw) as Record<string, Run>;
  } catch (e) {
    // ignore parse errors
  }

  const runs = ref<Record<string, Run>>(initial || {});
  const polling = ref<Record<string, number | undefined>>({});
  const toast = useToastStore();

  // persist runs to localStorage whenever they change
  watch(
    runs,
    (val) => {
      try {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(val));
      } catch (e) {
        // ignore quota errors
      }
    },
    { deep: true }
  );

  async function createRun(req: RunRequest) {
    const run = await startRun(req);
    runs.value[run.id] = run;
    toast.add(`Started run ${run.id}`, 'info', 4000);
    startPolling(run.id);
    return run;
  }

  function startPolling(runId: string) {
    if (polling.value[runId]) return;
    const intervalId = window.setInterval(async () => {
      try {
        const updated = await getRun(runId);
        runs.value[runId] = updated;
        if (updated.status !== 'running') {
          stopPolling(runId);
          // notify on completion
          if (updated.status === 'passed') {
            toast.add(`Run ${runId} passed`, 'success', 6000);
          } else {
            toast.add(`Run ${runId} ${updated.status}`, 'error', 8000);
          }
        }
      } catch (e) {
        // ignore
      }
    }, 2000);
    polling.value[runId] = intervalId as unknown as number;
  }

  function stopPolling(runId: string) {
    const id = polling.value[runId];
    if (id) {
      clearInterval(id);
      polling.value[runId] = undefined;
    }
  }

  function getRunById(runId: string) {
    return runs.value[runId];
  }

  function listRuns() {
    return Object.values(runs.value).sort((a,b) => (a.startedAt || '') < (b.startedAt || '') ? 1 : -1);
  }

  // If there are persisted runs that are still running, resume polling for them
  for (const r of Object.values(runs.value)) {
    if (r && r.status === 'running' && r.id) {
      startPolling(r.id);
    }
  }

  return { runs, polling, createRun, startPolling, stopPolling, getRunById, listRuns };
});
