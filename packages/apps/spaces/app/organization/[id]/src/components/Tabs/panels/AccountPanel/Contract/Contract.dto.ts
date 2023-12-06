import { SelectOption } from '@shared/types/SelectOptions';
import { ContractUpdateInput, ContractRenewalCycle } from '@graphql/types';
import { billingFrequencyOptions } from '@organization/src/components/Tabs/panels/AccountPanel/utils';

export interface TimeToRenewalData {
  name?: string;
  endedAt?: Date;
  signedAt?: Date;
  serviceStartedAt?: Date;
  organizationName?: string;
  contractUrl?: string | null;
  renewalPeriods?: number | null;
  renewalCycle?: ContractRenewalCycle;
}
export interface TimeToRenewalForm {
  name?: string;
  endedAt?: Date;
  signedAt?: Date;
  serviceStartedAt?: Date;
  contractUrl?: string | null;
  renewalPeriods?: string | null;
  renewalCycle?: SelectOption<ContractRenewalCycle | 'MULTI_YEAR'> | null;
}

export class ContractDTO implements TimeToRenewalForm {
  signedAt?: Date;
  endedAt?: Date;
  serviceStartedAt?: Date;
  renewalCycle?: SelectOption<ContractRenewalCycle | 'MULTI_YEAR'> | null;
  name?: string;
  contractUrl?: string | null;
  renewalPeriods?: string | null;

  constructor(data?: TimeToRenewalData | null) {
    this.renewalCycle =
      [...billingFrequencyOptions].find(({ value }) =>
        data?.renewalPeriods ?? 0 > 1
          ? value === 'MULTI_YEAR'
          : value === data?.renewalCycle,
      ) ?? undefined;
    this.signedAt = data?.signedAt && new Date(data.signedAt);
    this.endedAt = data?.endedAt && new Date(data.endedAt);
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
      serviceStartedAt: data?.serviceStartedAt,
      signedAt: data?.signedAt,
      endedAt: data?.endedAt,
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
