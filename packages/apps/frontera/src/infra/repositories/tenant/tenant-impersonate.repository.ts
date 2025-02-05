import { Transport } from './../../transport';
import TenantImpersonateListDocument from './queries/impersonateList.graphql';
import { TenantImpersonateListQuery } from './queries/impersonateList.generated';

export class TenantImpersonateRepository {
  static instance: TenantImpersonateRepository | null = null;
  private transport = Transport.getInstance();

  public static getInstance() {
    if (!TenantImpersonateRepository.instance) {
      TenantImpersonateRepository.instance = new TenantImpersonateRepository();
    }

    return TenantImpersonateRepository.instance;
  }

  async getTenantImpersonateList() {
    return this.transport.graphql.request<TenantImpersonateListQuery>(
      TenantImpersonateListDocument,
    );
  }
}
