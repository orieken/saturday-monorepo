import { createDefaultRedactionPolicy } from '../src/index';

describe('createDefaultRedactionPolicy', () => {
  const policy = createDefaultRedactionPolicy();

  describe('redactHeader', () => {
    it('should redact known secret headers (authorization)', () => {
      const result = policy.redactHeader!('authorization', 'Bearer token123');
      expect(result).not.toBeNull();
      expect(result?.value).toBe('${K6_AUTHORIZATION}');
      expect(result?.finding).toEqual({
        envName: 'K6_AUTHORIZATION',
        value: 'Bearer token123',
        source: 'header',
        path: 'headers.authorization'
      });
    });

    it('should redact known secret headers (x-api-key)', () => {
      const result = policy.redactHeader!('x-api-key', '12345');
      expect(result?.value).toBe('${K6_X_API_KEY}');
      expect(result?.finding?.envName).toBe('K6_X_API_KEY');
    });

    it('should not redact common headers (content-type)', () => {
      const result = policy.redactHeader!('content-type', 'application/json');
      expect(result).toBeNull();
    });

    it('should redact cookies with sensitive content', () => {
      const result = policy.redactHeader!('cookie', 'session_id=abc; auth_token=xyz');
      expect(result).not.toBeNull();
      expect(result?.value).toBe('${K6_COOKIE}');
    });

    it('should NOT redact cookies without sensitive content', () => {
      // cookie without session/token/auth
      const result = policy.redactHeader!('cookie', 'flavor=chocolate');
      // The current regex is /session|token|auth/i
      expect(result).toBeNull();
    });

    it('should support custom envPrefix', () => {
      const customPolicy = createDefaultRedactionPolicy({ envPrefix: 'TEST_' });
      const result = customPolicy.redactHeader!('authorization', 'foo');
      expect(result?.value).toBe('${TEST_AUTHORIZATION}');
      expect(result?.finding?.envName).toBe('TEST_AUTHORIZATION');
    });

    it('should respect detectCommonSecrets: false', () => {
      const laxPolicy = createDefaultRedactionPolicy({ detectCommonSecrets: false });
      const result = laxPolicy.redactHeader!('authorization', 'foo');
      expect(result).toBeNull();
    });
  });

  describe('redactBody', () => {
    it('should redact fields containing "token"', () => {
      const result = policy.redactBody!('auth.accessToken', 'secret-val');
      expect(result).not.toBeNull();
      expect(result?.value).toBe('${K6_ACCESSTOKEN}');
      expect(result?.finding).toEqual({
        envName: 'K6_ACCESSTOKEN',
        value: 'secret-val',
        source: 'body',
        path: 'auth.accessToken'
      });
    });

    it('should redact fields containing "password"', () => {
      const result = policy.redactBody!('user.password', 'pass123');
      expect(result?.value).toBe('${K6_PASSWORD}');
    });

     it('should redact fields containing "secret"', () => {
      const result = policy.redactBody!('clientSecret', 'xyz');
      expect(result?.value).toBe('${K6_CLIENTSECRET}');
    });

    it('should return null for non-sensitive fields', () => {
      const result = policy.redactBody!('user.name', 'John');
      expect(result).toBeNull();
    });

    it('should return null for null values', () => {
       const result = policy.redactBody!('user.password', null);
       expect(result).toBeNull();
    });
    it('should handle edge case with trailing dot in keyPath', () => {
       const result = policy.redactBody!('secret.', 'val');
       // key defaults to 'secret', envName becomes K6_SECRET
       expect(result?.value).toBe('${K6_SECRET}');
    });
  });
});
