import type { RootStore } from '@store/root';
import type { Transport } from '@infra/transport';

import { Store } from '@store/_store';
import { Sku } from '@store/Sku/Sku.dto';
import { action, computed, runInAction } from 'mobx';
import { SkuRepository } from '@infra/repositories/sku';
import { SkuDatum } from '@infra/repositories/sku/sku.datum';

import { SkuInput } from '@graphql/types';

export class SkusStore extends Store<SkuDatum, Sku> {
  private service: SkuRepository = SkuRepository.getInstance();

  constructor(public root: RootStore, public transport: Transport) {
    super(root, transport, {
      name: 'Skus',
      getId: (data) => data?.id,
      factory: Sku,
    });
  }

  public async bootstrap() {
    if (this.isBootstrapped) return;

    try {
      this.isLoading = true;

      const { skus } = await this.service.getSkus();

      runInAction(() => {
        skus.forEach((raw) => {
          this.value.set(raw.id, new Sku(this, raw));
        });
        this.isBootstrapped = true;
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  @action
  public add(data: Omit<SkuInput, 'id'>) {
    const record = new Sku(this, Sku.default(data));

    this.value.set(record.id, record);

    return record.id;
  }

  @action
  public edit(data: SkuInput & { id: string }) {
    const record = this.getById(data.id)?.value;

    if (record) {
      record.price = data.price;
      record.name = data.name;
    }
  }

  @action
  public delete(id: string) {
    this.value.delete(id);
  }

  @computed
  get isFullyLoaded() {
    return this.totalElements === this.value.size;
  }
}
