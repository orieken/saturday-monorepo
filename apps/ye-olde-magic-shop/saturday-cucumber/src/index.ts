import { BeforeAll, AfterAll, Before, After } from '@cucumber/cucumber';
import { createBrowser, closeBrowser, createContext, closeContext } from './hooks/setup';
import { installSaturdayWorld } from './world/saturday-world';

export * from './world/saturday-world';
export * from './world/saturday-world-options';
export * from './hooks/setup';
export * from './utils/browser-options';
export * from './support/managers';
export * from './config';

export function installSaturdayHooks() {
  installSaturdayWorld();
  
  BeforeAll(createBrowser());
  AfterAll(closeBrowser());
  Before(createContext());
  After(closeContext());
}
