import { LaunchOptions } from 'playwright';
import { getConfig } from '../config';

/* eslint-disable @typescript-eslint/naming-convention */
export const browserOptions: LaunchOptions = {
  get slowMo() { return getConfig().browser.slowMo; },
  get headless() { return getConfig().browser.headless; },
  get args() { return getConfig().browser.args; },
  firefoxUserPrefs: {
    'media.navigator.streams.fake': true,
    'media.navigator.permission.disabled': true,
  },
};
/* eslint-enable @typescript-eslint/naming-convention */
