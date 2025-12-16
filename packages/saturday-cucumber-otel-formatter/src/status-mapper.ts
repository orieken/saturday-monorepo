import { TestStepResultStatus } from '@cucumber/messages';

export class StatusMapper {
  static map(status?: TestStepResultStatus): string {
    switch (status) {
      case TestStepResultStatus.PASSED:
        return 'passed';
      case TestStepResultStatus.FAILED:
        return 'error';
      case TestStepResultStatus.SKIPPED:
      case TestStepResultStatus.PENDING:
      case TestStepResultStatus.UNDEFINED:
      case TestStepResultStatus.AMBIGUOUS:
      case TestStepResultStatus.UNKNOWN:
        return 'unset';
      default:
        return 'unset';
    }
  }
}
