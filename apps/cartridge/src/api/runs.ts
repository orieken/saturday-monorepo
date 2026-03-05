export type RunRequest = {
  framework: string;
  suiteId: string;
  scenarioId: string;
  executor?: 'docker' | 'k8s';
};

export type Run = {
  id: string;
  status: string; // running|passed|failed
  startedAt?: string;
  framework?: string;
  suiteId?: string;
  scenarioId?: string;
  reportUrl?: string;
};

const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:9001';

export async function startRun(req: RunRequest): Promise<Run> {
  const res = await fetch(`${API_BASE}/api/runs`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(req),
  });
  if (!res.ok) throw new Error('Failed to start run');
  return (await res.json()) as Run;
}

export async function getRun(runId: string): Promise<Run> {
  const res = await fetch(`${API_BASE}/api/runs/${encodeURIComponent(runId)}`);
  if (!res.ok) throw new Error('Failed to fetch run status');
  return (await res.json()) as Run;
}

/**
 * Fetch the Cucumber JSON report for a given run
 * @param runId - The run ID
 * @returns The Cucumber JSON report
 * @throws Error if the report cannot be fetched
 */
export async function fetchRunReport(runId: string): Promise<import('../types/cucumber').CucumberReport> {
  const res = await fetch(`${API_BASE}/api/runs/${encodeURIComponent(runId)}/report`);
  if (!res.ok) {
    throw new Error(`Failed to fetch report: ${res.statusText}`);
  }
  return (await res.json()) as import('../types/cucumber').CucumberReport;
}
