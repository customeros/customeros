import { BillingDetails, Scalars, Maybe } from '@graphql/types';
import { frequencyOptions } from '@organization/components/Tabs/panels/AccountPanel/BillingDetailsCard/utils';

export interface OrganizationAccountBillingDetailsForm {
  amount?: number | null;
  renewalCycleStart: Date | null;
  renewalCycle: Maybe<{ label: string; value: string }> | undefined;
  frequency: Maybe<{ label: string; value: string }> | undefined;
}

export class OrganizationAccountBillingDetails
  implements OrganizationAccountBillingDetailsForm
{
  amount: number | undefined;
  renewalCycleStart: Date | null;
  renewalCycle: Maybe<{ label: string; value: string }> | undefined;
  frequency: Maybe<{ label: string; value: string }> | undefined;

  constructor(data: BillingDetails & { amount?: string | null }) {
    this.amount = data?.amount ? parseFloat(data.amount) : undefined;
    this.renewalCycleStart = data?.renewalCycleStart
      ? new Date(data?.renewalCycleStart)
      : null;
    this.renewalCycle = frequencyOptions.find(
      ({ value }) => value === data?.renewalCycle,
    );
    this.frequency = frequencyOptions.find(
      ({ value }) => value === data?.frequency,
    );
  }

  static toForm(data: any) {
    return new OrganizationAccountBillingDetails(data.organization);
  }

  static toPayload(data: OrganizationAccountBillingDetailsForm) {
    return {
      amount: data.amount,
      renewalCycleStart: data.renewalCycleStart,
      renewalCycle:
        typeof data.renewalCycle === 'string'
          ? data?.renewalCycle
          : data?.renewalCycle?.value,
      frequency:
        typeof data.frequency === 'string'
          ? data?.frequency
          : data?.frequency?.value,
    } as any;
  }
}
