import { IWorldOptions } from '@cucumber/cucumber';

export interface SaturdayWorldOptions extends IWorldOptions {
  parameters: { [key: string]: string };
}
