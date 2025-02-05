import { TenantImpersonateListQuery } from './queries/impersonateList.generated.ts';

export type TenantImpersonateDatum = NonNullable<
  TenantImpersonateListQuery['tenant_impersonateLis'][0]
>;
