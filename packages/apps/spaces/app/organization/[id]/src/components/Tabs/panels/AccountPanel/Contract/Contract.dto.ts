import { SelectOption } from '@shared/types/SelectOptions';
import { ContractUpdateInput, ContractRenewalCycle } from '@graphql/types';
import { billingFrequencyOptions } from '@organization/src/components/Tabs/panels/AccountPanel/utils';

export interface TimeToRenewalData {
  name?: string;
  endedAt?: Date;
  signedAt?: Date;
  serviceStartedAt?: Date;
  contractUrl?: string | null;
  renewalCycle?: ContractRenewalCycle;
}
export interface TimeToRenewalForm {
  name?: string;
  endedAt?: Date;
  signedAt?: Date;
  serviceStartedAt?: Date;
  contractUrl?: string | null;
  renewalCycle?: SelectOption<ContractRenewalCycle> | null;
}

export class ContractDTO implements TimeToRenewalForm {
  signedAt?: Date;
  endedAt?: Date;
  serviceStartedAt?: Date;
  renewalCycle?: SelectOption<ContractRenewalCycle> | null;
  name?: string;
  contractUrl?: string | null;

  constructor(data?: TimeToRenewalData | null) {
    this.renewalCycle =
      billingFrequencyOptions.find((o) => o.value === data?.renewalCycle) ??
      null;
    this.signedAt = data?.signedAt && new Date(data.signedAt);
    this.endedAt = data?.endedAt && new Date(data.endedAt);
    this.serviceStartedAt =
      data?.serviceStartedAt && new Date(data.serviceStartedAt);
    this.name = data?.name ?? '';
    this.contractUrl = data?.contractUrl ?? '';
  }

  static toForm(
    organizationName: string,
    data?: TimeToRenewalData | null,
  ): TimeToRenewalForm {
    const formData = new ContractDTO(data);

    return {
      ...formData,
      name: formData.name?.length
        ? formData.name
        : `${organizationName}'s contract`,
    };
  }

  static toPayload(
    data: TimeToRenewalForm,
  ): Omit<ContractUpdateInput, 'contractId'> {
    return {
      serviceStartedAt: data?.serviceStartedAt,
      signedAt: data?.signedAt,
      endedAt: data?.endedAt,
      renewalCycle: data?.renewalCycle?.value,
      name: data?.name,
      contractUrl: data?.contractUrl,
    };
  }
}
