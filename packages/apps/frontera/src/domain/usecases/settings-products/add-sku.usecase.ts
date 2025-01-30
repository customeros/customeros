import { action, observable } from 'mobx';
import { RootStore } from '@store/root.ts';
import { SkuService } from '@domain/services';

import { SkuType } from '@graphql/types';

export class AddSkuUsecase {
  private root = RootStore.getInstance();
  private service = new SkuService();
  @observable accessor type: SkuType = SkuType.Subscription;
  @observable accessor productName: string = '';
  @observable accessor price: number | undefined = undefined;
  @observable accessor errors = {
    productName: '',
    price: '',
  };

  constructor() {
    this.editProductName = this.editProductName.bind(this);
    this.editType = this.editType.bind(this);
    this.editPrice = this.editPrice.bind(this);
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
    this.type = SkuType.Subscription;
  }

  @action
  editProductName(name: string) {
    this.productName = name;
  }

  @action
  editType(type: SkuType) {
    this.type = type;
  }

  @action
  editPrice(price: string) {
    const parsedValue = price.length ? parseFloat(price) : undefined;

    this.price = parsedValue;
  }

  @action
  validate() {
    if (!this.price) {
      this.errors.price = 'Please enter a product price';
    }

    if (!this.productName.trim().length) {
      this.errors.productName = 'Please enter a product name';
    }

    if (typeof this.price !== 'number') {
      this.errors.price = 'Please enter a product price';

      return;
    }
  }

  @action
  async createSku() {
    if (typeof this.price !== 'number') {
      return;
    }
    const [res, _err] = await this.service.createSku({
      name: this.productName,
      type: this.type,
      price: this.price,
    });

    if (res) {
      this.reset();
      this.root.ui.commandMenu.setOpen(false);
      this.root.ui.commandMenu.clearContext();
    }
  }
}
