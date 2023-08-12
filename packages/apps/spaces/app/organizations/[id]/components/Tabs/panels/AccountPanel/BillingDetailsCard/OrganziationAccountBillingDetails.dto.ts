import { BillingDetails, Scalars, Maybe, RenewalCycle } from '@graphql/types';

export interface OrganizationAccountBillingDetailsForm {
  amount?: number | null;
  renewalCycleStart: Date | null;
  renewalCycle?: Maybe<RenewalCycle>;
  frequency?: Maybe<RenewalCycle>;
}

export class OrganizationAccountBillingDetails
  implements OrganizationAccountBillingDetailsForm
{
  amount: Maybe<Scalars['Float']> | undefined;
  renewalCycleStart: Date | null;
  renewalCycle: Maybe<RenewalCycle> | undefined;
  frequency: Maybe<RenewalCycle> | undefined;

  constructor(data: BillingDetails) {
    this.amount = data?.amount;
    this.renewalCycleStart = data?.renewalCycleStart;
    this.renewalCycle = data?.renewalCycle;
    this.frequency = data?.frequency;
  }

  static toForm(data: any) {
    return new OrganizationAccountBillingDetails(data.organization);
  }

  static toPayload(data: OrganizationAccountBillingDetailsForm) {
    return {
      amount: data.amount,
      renewalCycleStart: data.renewalCycleStart,
      renewalCycle: data.renewalCycle,
      frequency: data.frequency,
    } as any;
  }
}
