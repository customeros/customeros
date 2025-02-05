import { Page, expect } from '@playwright/test';

import { SkuType } from '../../../src/routes/src/types/__generated__/graphql.types';
import {
  ensureLocatorIsVisible,
  clickLocatorThatIsVisible,
} from '../../helper';

export class SettingsProductsPage {
  constructor(page) {
    this.page = page;
  }

  private page: Page;
  settingsProducts = 'button[data-test="sideNav-settings-products"]';
  productsHeader = 'p[data-test="products-header"]';
  addProductButton = 'button[data-test="add-product-button"]';
  newProductHeader = 'h1[data-test="new-product-header"]';
  addProductSubscription = 'span[data-test="add-product-subscription"]';
  addProductOneTime = 'span[data-test="add-product-one-time"]';
  skuProductName = 'input[data-test="sku-product-name"]';
  skuPrice = 'input[data-test="sku-price"]';
  addSkuCancel = 'button[data-test="add-sku-cancel"]';
  addSkuConfirm = 'button[data-test="add-sku-confirm"]';
  addSkuX = 'button[data-test="add-sku-x"]';

  async goToSettingsProductsPage() {
    await clickLocatorThatIsVisible(this.page, this.settingsProducts);
    await ensureLocatorIsVisible(this.page, this.productsHeader);
    await expect(this.page.locator(this.productsHeader)).toHaveText('Products');
  }

  async openAddProduct() {
    await clickLocatorThatIsVisible(this.page, this.addProductButton);
    await expect(this.page.locator(this.newProductHeader)).toHaveText(
      'New product',
    );
  }

  async addProduct(
    productType: string,
    productName: string,
    productPrice: string,
  ) {
    await this.openAddProduct();

    if (productType === SkuType.Subscription) {
      await clickLocatorThatIsVisible(this.page, this.addProductSubscription);
    } else if (productType === SkuType.OneTime) {
      await clickLocatorThatIsVisible(this.page, this.addProductOneTime);
    }
    await this.page.locator(this.skuProductName).fill(productName);
    await this.page.locator(this.skuPrice).fill(productPrice);
    await clickLocatorThatIsVisible(this.page, this.addSkuConfirm);
  }
}
