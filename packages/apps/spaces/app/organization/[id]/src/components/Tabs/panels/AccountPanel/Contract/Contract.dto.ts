import { SelectOption } from '@shared/types/SelectOptions';
import {
  ContractUpdateInput,
  ContractRenewalCycle,
  ContractBillingCycle,
} from '@graphql/types';
import {
  billingFrequencyOptions,
  contractBillingCycleOptions,
} from '@organization/src/components/Tabs/panels/AccountPanel/utils';

export interface TimeToRenewalData {
  name?: string;
  endedAt?: Date;
  serviceStartedAt?: Date;
  invoicingStartDate?: Date;
  organizationName?: string;
  contractUrl?: string | null;

  renewalPeriods?: number | null;
  renewalCycle?: ContractRenewalCycle;
  billingCycle?: ContractBillingCycle | null;
}
export interface TimeToRenewalForm {
  name?: string;
  endedAt?: Date;
  serviceStartedAt?: Date;
  invoicingStartDate?: Date;
  contractUrl?: string | null;
  renewalPeriods?: string | null;
  country?: SelectOption<string> | null;
  organizationLegalName?: string | null;
  billingCycle?: SelectOption<ContractBillingCycle> | null;
  renewalCycle?: SelectOption<ContractRenewalCycle | 'MULTI_YEAR'> | null;
}

export class ContractDTO implements TimeToRenewalForm {
  endedAt?: Date;
  invoicingStartDate?: Date;
  serviceStartedAt?: Date;
  renewalCycle?: SelectOption<ContractRenewalCycle | 'MULTI_YEAR'> | null;
  name?: string;
  contractUrl?: string | null;
  renewalPeriods?: string | null;
  billingCycle?: SelectOption<ContractBillingCycle> | null;

  constructor(data?: TimeToRenewalData | null) {
    this.renewalCycle =
      [...billingFrequencyOptions].find(({ value }) =>
        (data?.renewalPeriods ?? 0) > 1
          ? value === 'MULTI_YEAR'
          : value === data?.renewalCycle,
      ) ?? undefined;
    this.billingCycle =
      [...contractBillingCycleOptions].find(
        ({ value }) => value === data?.billingCycle,
      ) ?? undefined;
    this.endedAt = data?.endedAt && new Date(data.endedAt);
    this.invoicingStartDate =
      data?.invoicingStartDate && new Date(data.invoicingStartDate);
    this.serviceStartedAt =
      data?.serviceStartedAt && new Date(data.serviceStartedAt);
    this.name = data?.name?.length
      ? data?.name
      : `${
          data?.organizationName?.length
            ? `${data?.organizationName}'s`
            : "Unnamed's"
        } contract`;
    this.contractUrl = data?.contractUrl ?? '';
    this.renewalPeriods = String(data?.renewalPeriods ?? 2);
  }

  static toForm(data?: TimeToRenewalData | null): TimeToRenewalForm {
    const formData = new ContractDTO(data);

    return {
      ...formData,
    };
  }

  static toPayload(
    data: TimeToRenewalForm,
  ): Omit<ContractUpdateInput, 'contractId'> {
    return {
      endedAt: data?.endedAt,
      invoicingStartDate: data?.invoicingStartDate,
      serviceStartedAt: data?.serviceStartedAt,
      billingCycle: data?.billingCycle?.value,
      renewalCycle:
        data?.renewalCycle?.value === 'MULTI_YEAR'
          ? ContractRenewalCycle.AnnualRenewal
          : data?.renewalCycle?.value,
      name: data?.name,
      contractUrl: data?.contractUrl,
      renewalPeriods:
        data?.renewalCycle?.value === 'MULTI_YEAR'
          ? parseInt(data?.renewalPeriods || '2')
          : data?.renewalPeriods
          ? parseInt(data?.renewalPeriods)
          : undefined,
    };
  }
}
