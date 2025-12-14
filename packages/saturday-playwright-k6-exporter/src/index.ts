import fs from 'node:fs';
import path from 'node:path';
import { request, type APIRequestContext } from '@playwright/test';

export type RecordedCall = {
  name: string;
  method: string;
  url: string;
  headers?: Record<string, string>;
  body?: unknown;
  status?: number;
  responseHeaders?: Record<string, string>;
};


export type ExporterOptions = {
  outDir?: string;
  /** Redaction policy plugin */
  policy?: import('@orieken/saturday-k6-redaction-basic').RedactionPolicy;
  /** Write .env.apis with discovered secrets (default: true) */
  writeEnvFile?: boolean;
  /** Update .env.example with missing keys (default: true) */
  updateEnvExample?: boolean;
  /** Path to env files (root project dir) */
  envDir?: string;
};


async function resolvePolicy(options: ExporterOptions): Promise<import('@orieken/saturday-k6-redaction-basic').RedactionPolicy | undefined> {
  if (options.policy) return options.policy as any;
  const spec = process.env.K6_REDACTION_POLICY;
  if (spec) {
    try {
      const mod = await import(spec);
      return (mod.default || mod.createDefaultRedactionPolicy?.() || mod.policy) as any;
    } catch {}
  }
  try {
    const mod = await import('@orieken/saturday-k6-redaction-basic');
    // prefer a default instance if exported, else a factory.
      return ((mod as any).default || (mod as any).createDefaultRedactionPolicy?.() || (mod as any).policy) as any;
  } catch {}
  return undefined;
}

export class K6Recorder {

  private calls: RecordedCall[] = [];
  
private findings = new Map<string, string>(); // envName -> value
  private resolvedPolicy: any | undefined;
  constructor(
  private readonly outDir: string,
  private readonly testSlug: string,
  private readonly options: ExporterOptions = {}
) {}


  wrap(ctx: APIRequestContext): APIRequestContext {
    const handler = {
      get: (target: APIRequestContext, prop: keyof APIRequestContext) => {
        const orig = (target as any)[prop];
        if (typeof orig !== 'function') return orig;
        return async (...args: any[]) => {
          const verbLike = ['get', 'post', 'put', 'patch', 'delete', 'head', 'fetch'];
          if (!verbLike.includes(String(prop))) {
            return await orig.apply(target, args);
          }
          let method = String(prop).toUpperCase();
          let url = args[0];
          let opt = (args[1] ?? {}) as any;
          if (prop === 'fetch' && opt?.method) method = String(opt.method).toUpperCase();

          let body = opt.data;
          
          // Capture response
          const resp = await orig.apply(target, args);
          const status = resp.status();
          const respHeaders = resp.headers();

          const headers: Record<string, string> = {};
          const inputHeaders = opt?.headers ?? {};
          if (Array.isArray(inputHeaders)) {
            for (const [k, v] of inputHeaders) headers[String(k)] = String(v);
          } else if (typeof inputHeaders === 'object') {
            for (const k of Object.keys(inputHeaders)) headers[String(k)] = String(inputHeaders[k]);
          }

          let redactedHeaders: Record<string,string> = { ...headers };
          let redactedBody: unknown = body;

          const policy = this.resolvedPolicy ?? this.options.policy;
          if (policy?.redactHeader) {
            for (const hk of Object.keys(redactedHeaders)) {
              const res = policy.redactHeader(hk, String(redactedHeaders[hk]));
              if (res) {
                redactedHeaders[hk] = String(res.value);
                if (res.finding) this.findings.set(res.finding.envName, String(res.finding.value));
              }
            }
          }

          if (policy?.redactBody && redactedBody && typeof redactedBody === 'object') {
            const walk = (obj: any, pathPrefix: string) => {
              if (obj == null) return;
              if (Array.isArray(obj)) {
                obj.forEach((v, i) => {
                  const keyPath = `${pathPrefix}.${i}`;
                  if (typeof v === 'object' && v !== null) walk(v, keyPath);
                  else {
                    const res = policy.redactBody(keyPath, v);
                    if (res) {
                      obj[i] = res.value;
                      if (res.finding) this.findings.set(res.finding.envName, String(res.finding.value));
                    }
                  }
                });
              } else {
                for (const k of Object.keys(obj)) {
                  const v = obj[k];
                  const keyPath = pathPrefix ? `${pathPrefix}.${k}` : k;
                  if (typeof v === 'object' && v !== null) walk(v, keyPath);
                  else {
                    const res = policy.redactBody(keyPath, v);
                    if (res) {
                      obj[k] = res.value;
                      if (res.finding) this.findings.set(res.finding.envName, String(res.finding.value));
                    }
                  }
                }
              }
            };
            try {
              const clone = JSON.parse(JSON.stringify(redactedBody));
              walk(clone, '');
              redactedBody = clone;
            } catch {}
          }

          this.calls.push({ name: this.testSlug, method, url, headers: redactedHeaders, body: redactedBody, status, responseHeaders: respHeaders });
          return resp;

        };
      }
    };
    return new Proxy(ctx, handler);
  }

  record(name: string, method: string, url: string, options?: { headers?: Record<string, string>; body?: unknown }) {
    this.calls.push({ name, method, url, headers: options?.headers, body: options?.body });
  }

  hasCalls(): boolean { return this.calls.length > 0; }

  
  private writeEnvArtifacts() {
    const { writeEnvFile = true, updateEnvExample = true, envDir = process.cwd() } = this.options as any;
    if (this.findings.size === 0) return;

    const lines: string[] = [];
    for (const [envName, value] of this.findings.entries()) {
      lines.push(`${envName}=${value}`);
    }
  try {
    const envApis = path.join(envDir, '.env.apis');
    const content = lines.join('\n') + '\n';
    fs.writeFileSync(envApis, content, { encoding: 'utf-8', flag: 'w' });
  } catch {}

  try {
    const examplePath = path.join(envDir, '.env.example');
    const existing = fs.existsSync(examplePath) ? fs.readFileSync(examplePath, 'utf-8') : '';
    const need = Array.from(this.findings.keys()).filter(k => !new RegExp(`^\\s*${k}=`,'m').test(existing));
    if (need.length) {
      const append = need.map(k => `${k}=`).join('\n') + '\n';
      fs.writeFileSync(examplePath, existing + (existing && !existing.endsWith('\n') ? '\n' : '') + append, 'utf-8');
    }
  } catch {}
}

async flushToK6(): Promise<string | null> {

    if (!this.calls.length) return null;
    const safeSlug = this.testSlug.replace(/[^a-z0-9\-_.]/gi, '_').toLowerCase();
    const file = path.join(this.outDir, `${safeSlug}.k6.js`);
    const code = this.toK6(this.calls);
      this.writeEnvArtifacts();
    fs.mkdirSync(this.outDir, { recursive: true });
    fs.writeFileSync(file, code, 'utf-8');
    return file;
  }

  private toK6(calls: RecordedCall[]): string {
    const uniq = (obj?: Record<string,string>) => {
      if (!obj) return '{}';
      const out: Record<string,string> = {};
      for (const k of Object.keys(obj)) {
        const key = k.toLowerCase();
        if (!(key in out)) out[key] = String(obj[k]);
      }
      return JSON.stringify(out, null, 2);
    };
    const lines: string[] = [];
    lines.push(`import http from 'k6/http';`);
    lines.push(`import { check, sleep } from 'k6';`);
    lines.push(`export const options = { vus: 1, iterations: 1 };`);
    lines.push(`export default function () {`);
    for (const [i, c] of calls.entries()) {
      const varName = `res${i+1}`;
      const hdrs = uniq(c.headers);
      const body = (c.body === undefined || c.body === null) ? 'null' : JSON.stringify(c.body, null, 2);
      const method = c.method.toUpperCase();
      const url = JSON.stringify(c.url);
      const params = `{ headers: ${hdrs} }`;
      if (['GET','DELETE','HEAD'].includes(method)) {
        lines.push(`  // ${c.name}`);
        lines.push(`  const ${varName} = http.request('${method}', ${url}, null, ${params});`);
      } else {
        lines.push(`  // ${c.name}`);
        lines.push(`  const ${varName} = http.request('${method}', ${url}, ${body}, ${params});`);
      }
      lines.push(`  check(${varName}, { 'status is 2xx/3xx': (r) => r.status >= 200 && r.status < 400 });`);
      lines.push(`  sleep(0.1);`);
    }
    lines.push(`}`);
    return lines.join('\n');
  }
}

export function makeSlugFromTitle(title: string): string {
  return title.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)/g, '');
}

export async function createK6Recorder(testTitle: string, outDir = path.join(process.cwd(), 'perf', 'k6'), options?: ExporterOptions): Promise<{ ctx: APIRequestContext; recorder: K6Recorder } | null> {
  if (!process.env.K6_EXPORT) return null;
  const testSlug = makeSlugFromTitle(testTitle);
  /* istanbul ignore next */
  if (!options) options = {};
  if (!options.policy) options.policy = await resolvePolicy(options);
  const recorder = new K6Recorder(outDir, testSlug, options);
  const ctx = await request.newContext();
  const wrapped = recorder.wrap(ctx);
  return { ctx: wrapped, recorder };
}