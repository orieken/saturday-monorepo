// @saturday/k6-redaction-basic
export type SecretFinding = { envName: string; value: string; source: 'header'|'body'; path: string };

export type RedactionPolicy = {
  redactHeader?: (name: string, value: string) => { value: string; finding?: SecretFinding } | null;
  redactBody?: (keyPath: string, value: unknown) => { value: unknown; finding?: SecretFinding } | null;
};

export type RedactionOptions = {
  envPrefix?: string;
  detectCommonSecrets?: boolean;
};

const DEFAULT_SECRET_HEADERS = [
  'authorization', 'proxy-authorization', 'x-api-key', 'x-apikey', 'x-api-token',
  'x-auth-token', 'x-access-token', 'x-amz-security-token', 'set-cookie'
];

function makeEnvName(prefix: string, name: string) {
  return (prefix + name.replace(/[^a-z0-9]+/gi, '_')).toUpperCase();
}

export function createDefaultRedactionPolicy(opts: RedactionOptions = {}): RedactionPolicy {
  const { envPrefix = 'K6_', detectCommonSecrets = true } = opts;
  return {
    redactHeader: (name, value) => {
      const lower = name.toLowerCase();
      if (!detectCommonSecrets) return null;
      if (DEFAULT_SECRET_HEADERS.includes(lower)) {
        const envName = makeEnvName(envPrefix, lower);
        return { value: `\${${envName}}`, finding: { envName, value, source: 'header', path: `headers.${lower}` } };
      }
      if (lower === 'cookie' && /session|token|auth/i.test(value)) {
        const envName = makeEnvName(envPrefix, 'cookie');
        return { value: `\${${envName}}`, finding: { envName, value, source: 'header', path: 'headers.cookie' } };
      }
      return null;
    },
    redactBody: (keyPath, value) => {
      if (value == null) return null;
      if (typeof value === 'string' && /token|secret|password|apikey|api_key/i.test(keyPath)) {
        const key = keyPath.split('.').pop() || 'secret';
        const envName = makeEnvName(envPrefix, key);
        return { value: `\${${envName}}`, finding: { envName, value, source: 'body', path: keyPath } };
      }
      return null;
    }
  };
}