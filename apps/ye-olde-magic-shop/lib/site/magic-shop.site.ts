
import { BaseSite } from '@orieken/saturday-core';
import { Page } from '@playwright/test';
import { HomePage } from './home.page';
import { CartPage } from './cart.page';
import { ProductDetailsPage } from './product-details.page';
import { LoginPage } from './login.page';
import { AccountPage } from './account.page';
import { ReagentMapPage } from './reagent-map.page';
import { DrawLineFlow } from './flows/draw-line.flow';

export class MagicShop extends BaseSite {
  constructor(page: Page) {
    super(page, process.env.BASE_URL || 'http://localhost:8000');
  }

  protected initializePages(): void {
    this.registerPage('homePage', HomePage);
    this.registerPage('cartPage', CartPage);
    this.registerPage('productDetailsPage', ProductDetailsPage);
    this.registerPage('loginPage', LoginPage);
    this.registerPage('accountPage', AccountPage);
    this.registerPage('reagentMapPage', ReagentMapPage);
  }

  protected initializeFlows(): void {
    this.registerFlow('drawLine', DrawLineFlow);
  }

  get homePage(): HomePage {
    return this.getPage<HomePage>('homePage');
  }

  get cartPage(): CartPage {
    return this.getPage<CartPage>('cartPage');
  }

  get productDetailsPage(): ProductDetailsPage {
    return this.getPage<ProductDetailsPage>('productDetailsPage');
  }

  get loginPage(): LoginPage {
    return this.getPage<LoginPage>('loginPage');
  }

  get accountPage(): AccountPage {
    return this.getPage<AccountPage>('accountPage');
  }

  get reagentMapPage(): ReagentMapPage {
    return this.getPage<ReagentMapPage>('reagentMapPage');
  }

  get drawLine(): DrawLineFlow {
    return this.getFlow<DrawLineFlow>('drawLine');
  }
}
