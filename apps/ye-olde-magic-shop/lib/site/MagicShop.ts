
import { BaseSite } from '@orieken/saturday-core';
import { Page } from '@playwright/test';
import { HomePage } from './HomePage';
import { CartPage } from './CartPage';
import { ProductDetailsPage } from './ProductDetailsPage';
import { LoginPage } from './LoginPage';
import { AccountPage } from './AccountPage';

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
  }

  protected initializeFlows(): void {
    // No flows
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
}
