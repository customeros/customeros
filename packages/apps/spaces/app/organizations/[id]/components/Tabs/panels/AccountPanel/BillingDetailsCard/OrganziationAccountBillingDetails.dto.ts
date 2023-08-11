import { FrequencyOptions } from '@organization/components/Tabs/panels/AccountPanel/BillingDetailsCard/utils';

export interface OrganizationAccountBillingDetailsForm {
  billingDetailsAmount: number | null;
  billingDetailsRenewalCycleStart: Date | null;
  billingDetailsRenewalCycle: FrequencyOptions | null;
  billingDetailsFrequency: FrequencyOptions | null;
}

export class OrganizationAccountBillingDetails
  implements OrganizationAccountBillingDetailsForm
{
  billingDetailsAmount: number | null;
  billingDetailsRenewalCycleStart: Date | null;
  billingDetailsRenewalCycle: FrequencyOptions;
  billingDetailsFrequency: FrequencyOptions;

  constructor(data?: any) {
    this.billingDetailsAmount = data?.billingDetailsAmount;
    this.billingDetailsRenewalCycleStart =
      data?.billingDetailsRenewalCycleStart;
    this.billingDetailsRenewalCycle = data?.billingDetailsRenewalCycle;
    this.billingDetailsFrequency = data?.billingDetailsFrequency;
  }

  static toForm(data: any) {
    return new OrganizationAccountBillingDetails(data.organization);
  }

  static toPayload(data: OrganizationAccountBillingDetailsForm) {
    return {
      billingDetailsAmount: data.billingDetailsAmount,
      billingDetailsRenewalCycleStart: data.billingDetailsRenewalCycleStart,
      billingDetailsRenewalCycle: data.billingDetailsRenewalCycle,
      billingDetailsFrequency: data.billingDetailsFrequency,
    } as any;
  }
}
