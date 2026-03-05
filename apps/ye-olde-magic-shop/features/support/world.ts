import { setWorldConstructor, IWorldOptions } from '@cucumber/cucumber';
import { SaturdayWorld } from '@orieken/saturday-cucumber';
import { MagicShop } from '../../lib/site/MagicShop';

export interface TestWorld extends SaturdayWorld {
  site?: MagicShop;
}

export class WorldImpl extends SaturdayWorld implements TestWorld {
  public site?: MagicShop;

  constructor(options: IWorldOptions) {
    super(options);
  }

  // We can hook into the manager initialization if we want to auto-create our specific site
  override initializeManagers(): void {
    super.initializeManagers();
    if (this.page) {
        this.site = new MagicShop(this.page);
        // Also register it with the manager for good measure
        this.siteManager.register('magic-shop', this.site);
    }
  }
}

setWorldConstructor(WorldImpl);
