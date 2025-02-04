import { RootStore } from '@store/root.ts';
import { SkuService } from '@domain/services';
import { action, computed, observable } from 'mobx';

import { SkuType } from '@graphql/types';

export class EditSkuUsecase {
  private root = RootStore.getInstance();
  private service = new SkuService();
  @observable accessor productName: string = '';
  @observable accessor maskedPrice: string = '';
  @observable accessor price: number | undefined = undefined;
  @observable accessor errors = {
    productName: '',
    price: '',
  };
  private skuId: string;

  constructor(skuId: string) {
    this.skuId = skuId;
    this.editProductName = this.editProductName.bind(this);
    this.editPrice = this.editPrice.bind(this);
    this.setMaskedValue = this.setMaskedValue.bind(this);

    this.init();
  }

  @computed
  get sku() {
    if (!this.skuId) return;

    return this.root.skus.getById(this.skuId);
  }

  @action
  init() {
    this.maskedPrice = this.sku?.value?.price?.toString() ?? '';
  }

  @action
  resetErrors() {
    this.errors = {
      productName: '',
      price: '',
    };
  }

  @action
  resetPriceError() {
    this.errors = {
      ...this.errors,
      price: '',
    };
  }

  @action
  resetNameError() {
    this.errors = {
      ...this.errors,
      productName: '',
    };
  }

  @action
  reset() {
    this.resetErrors();

    this.productName = '';
    this.price = 0;
  }

  @action
  editProductName(name: string) {
    this.productName = name;
  }

  @action
  setMaskedValue(masked: string) {
    this.maskedPrice = masked;
  }

  @action
  setInitial() {
    this.price = this.sku?.value?.price ?? undefined;
    this.productName = this.sku?.value?.name ?? '';
  }

  @action
  editPrice(price: string) {
    const parsedValue = parseFloat(price);

    this.price = parsedValue;
  }

  @action
  validate() {
    if (!this.price) {
      this.errors.price = 'Please enter a product price';
    }

    if (typeof this.price !== 'number') {
      this.errors.price = 'Please enter a product price';

      return;
    }

    if (!this.productName.trim().length) {
      this.errors.productName = 'Please enter a product name';
    }
  }

  @action
  async editSku() {
    if (typeof this.price !== 'number') {
      return;
    }
    const [res, _err] = await this.service.editSku({
      id: this.skuId,
      type: this.sku?.value.type ?? SkuType.Subscription,
      name: this.productName,
      price: this.price,
    });

    if (res) {
      this.reset();
      this.root.ui.commandMenu.setOpen(false);
      this.root.ui.commandMenu.clearContext();
    }
  }
}
