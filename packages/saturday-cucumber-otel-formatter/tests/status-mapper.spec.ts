import { StatusMapper } from '../src/status-mapper';
import { TestStepResultStatus } from '@cucumber/messages';

describe('StatusMapper', () => {
  it('should map PASSED to ok', () => {
    expect(StatusMapper.map(TestStepResultStatus.PASSED)).toBe('ok');
  });

  it('should map FAILED to error', () => {
    expect(StatusMapper.map(TestStepResultStatus.FAILED)).toBe('error');
  });

  it('should map others to unset', () => {
    expect(StatusMapper.map(TestStepResultStatus.SKIPPED)).toBe('unset');
    expect(StatusMapper.map(TestStepResultStatus.PENDING)).toBe('unset');
    expect(StatusMapper.map(TestStepResultStatus.UNDEFINED)).toBe('unset');
  });

  it('should map undefined to unset', () => {
    expect(StatusMapper.map(undefined)).toBe('unset');
  });
});
