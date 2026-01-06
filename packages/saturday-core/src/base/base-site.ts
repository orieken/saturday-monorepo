import { Page } from 'playwright';
import { BasePage } from './base-page';
import { BaseFlow } from './base-flow';
import { TrainersFacade } from '../ml/facades/trainers.facade';
import { DetectorsFacade } from '../ml/facades/detectors.facade';
import { ModelsFacade } from '../ml/facades/models.facade';
import { TrainingOptions, ValidationOptions, ValidationResult } from '../ml/interfaces/ml-types';

export abstract class BaseSite {
  protected page: Page;
  protected baseUrl: string;
  protected pages: Map<string, BasePage> = new Map();
  protected flows: Map<string, BaseFlow> = new Map();

  get pageObject(): Page {
    return this.page;
  }

  // ML facades - lazy loaded
  private _trainers?: TrainersFacade;
  private _detectors?: DetectorsFacade;
  private _models?: ModelsFacade;

  constructor(page: Page, baseUrl: string) {
    this.page = page;
    this.baseUrl = baseUrl;
    this.initializePages();
    this.initializeFlows();
  }

  protected abstract initializePages(): void;
  protected abstract initializeFlows(): void;

  protected registerPage<T extends BasePage>(name: string, pageClass: new (page: Page, site: BaseSite) => T): void {
    const pageInstance = new pageClass(this.page, this);
    this.pages.set(name, pageInstance);
  }

  protected registerFlow<T extends BaseFlow>(name: string, flowClass: new (site: BaseSite) => T): void {
    const flowInstance = new flowClass(this);
    this.flows.set(name, flowInstance);
  }

  getPage<T extends BasePage>(name: string): T {
    const page = this.pages.get(name);
    if (!page) {
      throw new Error(`Page '${name}' is not registered in ${this.constructor.name}`);
    }
    return page as T;
  }

  getFlow<T extends BaseFlow>(name: string): T {
    const flow = this.flows.get(name);
    if (!flow) {
      throw new Error(`Flow '${name}' is not registered in ${this.constructor.name}`);
    }
    return flow as T;
  }

  // ML Facade Getters
  get trainers(): TrainersFacade {
    if (!this._trainers) {
      this._trainers = this.createTrainersFacade();
    }
    return this._trainers;
  }

  get detectors(): DetectorsFacade {
    if (!this._detectors) {
      this._detectors = this.createDetectorsFacade();
    }
    return this._detectors;
  }

  get models(): ModelsFacade {
    if (!this._models) {
      this._models = this.createModelsFacade();
    }
    return this._models;
  }

  // Protected factory methods for ML facades
  protected createTrainersFacade(): TrainersFacade {
    return new TrainersFacade(this);
  }

  protected createDetectorsFacade(): DetectorsFacade {
    return new DetectorsFacade(this);
  }

  protected createModelsFacade(): ModelsFacade {
    return new ModelsFacade(this);
  }

  // Convenience methods for common ML operations
  async establishPageBaseline(pageName: string, options?: TrainingOptions): Promise<void> {
    const page = this.getPage(pageName);
    await page.visit();
    await this.trainers.trainCurrentPage(`${pageName}_baseline`, options);
  }

  async validatePageAgainstBaseline(pageName: string, options?: ValidationOptions): Promise<ValidationResult> {
    const page = this.getPage(pageName);
    await page.visit();
    return await this.detectors.validateCurrentPage(`${pageName}_baseline`, options);
  }

  // Existing methods
  async visit(): Promise<void> {
    await this.page.goto(this.baseUrl);
  }

  getBaseUrl(): string {
    return this.baseUrl;
  }

  async getCurrentUrl(): Promise<string> {
    return this.page.url();
  }

  async getCurrentTitle(): Promise<string> {
    return await this.page.title();
  }

  async waitForNavigation(): Promise<void> {
    await this.page.waitForLoadState('networkidle');
  }
}
