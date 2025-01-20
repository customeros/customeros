import type { RootStore } from '@store/root';
import type { Transport } from '@infra/transport';

import { Store } from '@store/_store';
import { action, observable, runInAction } from 'mobx';
import { Industry, IndustryDatum } from '@store/Industries/Industry.dto';
import { IndustriesService } from '@store/Industries/__service__/Industries.service.ts';

export class IndustriesStore extends Store<IndustryDatum, Industry> {
  chunkSize = 50;
  private service = IndustriesService.getInstance();
  @observable accessor searchResults: Map<string, string[]> = new Map();
  @observable accessor cursors: Map<string, number> = new Map();
  @observable accessor availableCounts: Map<string, number> = new Map();

  constructor(public root: RootStore, public transport: Transport) {
    super(root, transport, {
      name: 'Industries',
      getId: (data) => data?.code,
      factory: Industry,
    });

    this.hydrate();
  }

  @action
  public async bootstrap() {
    try {
      const { industries_InUse } = await this.service.getIndustries();

      if (!industries_InUse) return;

      runInAction(() => {
        industries_InUse.forEach((raw) => {
          const foundRecord = this.value.get(raw.code);

          if (foundRecord) {
            Object.assign(foundRecord.value, raw);
          } else {
            const record = new Industry(this, raw);

            this.value.set(record.id, record);
          }
        });

        this.size = this.value.size;
        this.version++;
      });
    } catch (e) {
      console.error('Failed invalidating unsed industries list ');
    }
  }

  @action
  public async invalidate() {
    try {
      const { industries_InUse } = await this.service.getIndustries();

      if (!industries_InUse) return;

      runInAction(() => {
        industries_InUse.forEach((raw) => {
          const foundRecord = this.value.get(raw.code);

          if (foundRecord) {
            Object.assign(foundRecord.value, raw);
          } else {
            const record = new Industry(this, raw);

            this.value.set(record.id, record);
          }
        });

        this.size = this.value.size;
        this.version++;
      });
    } catch (e) {
      console.error('Failed invalidating unsed industries list ');
    }
  }
}
