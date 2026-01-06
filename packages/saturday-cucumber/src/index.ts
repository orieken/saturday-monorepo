import { BeforeAll, AfterAll, Before, After } from '@cucumber/cucumber';
import { createBrowser, closeBrowser, createContext, closeContext } from './hooks/setup';
import { installSaturdayWorld } from './world/custom-world';

export * from './world/custom-world';
export * from './world/custom-world-options';
export * from './hooks/setup';
export * from './utils/browser-options';

export function installSaturdayHooks() {
  installSaturdayWorld();
  
  BeforeAll(createBrowser());
  AfterAll(closeBrowser());
  Before(createContext());
  After(closeContext());
}
