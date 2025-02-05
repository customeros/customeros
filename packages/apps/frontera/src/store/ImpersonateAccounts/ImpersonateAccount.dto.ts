import merge from 'lodash/merge';
import { observable } from 'mobx';
import { Entity } from '@store/record';
import { Transport } from '@infra/transport';
import { TenantImpersonateDatum } from '@infra/repositories/tenant/tenant-impersonate.datum';

import { ImpersonateAccountsStore } from './ImpersonateAccounts.store.ts';

export class ImpersonateAccount extends Entity<TenantImpersonateDatum> {
  @observable accessor value: TenantImpersonateDatum =
    ImpersonateAccount.default();

  constructor(
    store: ImpersonateAccountsStore,
    data: TenantImpersonateDatum,
    public transport?: Transport,
  ) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    super(store as any, data);
  }

  static default(
    payload?: Partial<TenantImpersonateDatum>,
  ): TenantImpersonateDatum {
    return merge(
      {
        tenant: '',
        createdBy: '',
        personal: false,
      },
      payload ?? {},
    );
  }
}
