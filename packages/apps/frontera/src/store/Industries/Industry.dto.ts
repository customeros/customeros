import merge from 'lodash/merge';
import { observable } from 'mobx';
import { Entity } from '@store/record';

import type { GetTenantIndustriesQuery } from './__service__/getTenantIndustries.generated';

import { IndustriesStore } from './Industries.store.ts';

export type IndustryDatum = NonNullable<
  GetTenantIndustriesQuery['industries_InUse'][number]
>;

export class Industry extends Entity<IndustryDatum> {
  @observable accessor value: IndustryDatum = Industry.default();

  constructor(store: IndustriesStore, data: IndustryDatum) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    super(store as any, data);
  }

  static default(): IndustryDatum {
    return merge({
      name: '',
      code: '',
    });
  }
}
