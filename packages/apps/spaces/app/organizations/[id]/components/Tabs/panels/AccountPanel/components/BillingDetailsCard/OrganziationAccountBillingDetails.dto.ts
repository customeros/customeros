import { FrequencyOptions } from '@organization/components/Tabs/panels/AccountPanel/components/BillingDetailsCard/utils';

export interface OrganizationAccountBillingDetailsForm {
  billingDetailsAmount: number | null | undefined;
  billingDetailsRenewalCycleStart: Date | null | undefined;
  billingDetailsRenewalCycle: FrequencyOptions | null | undefined;
  billingDetailsFrequency: FrequencyOptions | null | undefined;
}

export class OrganizationAccountBillingDetails
  implements OrganizationAccountBillingDetailsForm
{
  billingDetailsAmount: number | null | undefined;
  billingDetailsRenewalCycleStart: Date | null | undefined;
  billingDetailsRenewalCycle: FrequencyOptions | undefined;
  billingDetailsFrequency: FrequencyOptions | undefined;

  constructor(data?: any) {
    this.billingDetailsAmount = data?.billingDetailsAmount;
    this.billingDetailsRenewalCycleStart = data?.billingDetailsRenewalCycleStart;
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
