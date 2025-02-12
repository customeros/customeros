import merge from 'lodash/merge';
import { Entity } from '@store/record';
import { action, computed, observable } from 'mobx';
import { SkuDatum } from '@infra/repositories/sku/sku.datum';

import { SkuType, SkuInput } from '@graphql/types';

import { SkusStore } from './Skus.store.ts';

export class Sku extends Entity<SkuDatum> {
  @observable accessor value: SkuDatum = Sku.default();

  constructor(store: SkusStore, data: SkuDatum) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    super(store as any, data);
  }

  @computed
  get id() {
    return this.value.id;
  }

  @computed
  get typeLabel() {
    return this.value.type === SkuType.Subscription
      ? 'Subscription'
      : 'One-time';
  }

  @computed
  get formattedPrice() {
    const currency =
      this.store.root.settings.tenant.value?.baseCurrency ?? 'USD';

    return Intl.NumberFormat('en-US', {
      style: 'currency',
      currency,
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(this.value.price);
  }

  @action
  public archiveSku() {
    this.draft();
    this.value.archived = true;
    this.commit({ syncOnly: true });
  }

  static default(payload?: SkuInput): SkuDatum {
    return merge(
      {
        id: crypto.randomUUID(),
        name: '',
        price: 0,
        type: SkuType.Subscription,
        archived: false,
      },
      payload ?? {},
    );
  }
}
