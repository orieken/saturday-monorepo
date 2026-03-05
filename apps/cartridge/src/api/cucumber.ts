export type CucumberStepRef = {
  line: number;
  text: string;
};

export type CucumberScenarioRef = {
  id: string;
  line: number;
  name: string;
  file: string;
  tags?: string[];
  steps?: CucumberStepRef[];
};

export type CucumberFeatureRef = {
  id: string;
  name: string;
  file: string;
  description?: string;
  scenarios?: CucumberScenarioRef[];
};

export type CucumberIndex = {
  framework: string;
  suiteId: string;
  features: CucumberFeatureRef[];
};

const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:9001';

export default async function fetchCucumberIndex(): Promise<CucumberIndex[] | null> {
  const res = await fetch(`${API_BASE}/api/cucumber/index`);
  if (res.status === 404) return null;
  if (!res.ok) throw new Error('Failed to fetch cucumber index');
  return (await res.json()) as CucumberIndex[];
}
