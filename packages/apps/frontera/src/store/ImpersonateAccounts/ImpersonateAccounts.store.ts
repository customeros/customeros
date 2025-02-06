import { runInAction } from 'mobx';
import { Store } from '@store/_store';
import { RootStore } from '@store/root';
import { Transport } from '@infra/transport';
import { TenantImpersonateRepository } from '@infra/repositories/tenant';
import { TenantImpersonateDatum } from '@infra/repositories/tenant/tenant-impersonate.datum';

import { unwrap } from '@utils/unwrap.ts';

import { ImpersonateAccount } from './ImpersonateAccount.dto.ts';

export class ImpersonateAccountsStore extends Store<
  TenantImpersonateDatum,
  ImpersonateAccount
> {
  private tenantRepository = new TenantImpersonateRepository();

  constructor(public root: RootStore, public transport: Transport) {
    super(root, transport, {
      name: 'ImpersonateAccounts',
      getId: (data) => data?.tenant,
      factory: ImpersonateAccount,
    });

    this.hydrate();
  }

  public async bootstrap() {
    const [data, err] = await unwrap(
      this.tenantRepository.getTenantImpersonateList(),
    );

    if (err) {
      console.error('Error bootstrapping agents:', err);

      return;
    }
    runInAction(() => {
      data?.tenant_impersonateList.forEach((datum) => {
        this.value.set(datum.tenant, new ImpersonateAccount(this, datum));
      });
      this.isBootstrapped = true;
      this.isBootstrapping = false;
    });
  }
}
