import { ContractRenewalCycle } from '@graphql/types';
import { SelectOption } from '@shared/types/SelectOptions';
import { billingFrequencyOptions } from '@organization/src/components/Tabs/panels/AccountPanel/utils';

export interface TimeToRenewalData {
  name?: string;
  endedAt?: Date;
  signedAt?: Date;
  serviceStartedAt?: Date;
  renewalCycle?: ContractRenewalCycle;
}
export interface TimeToRenewalForm {
  name?: string;
  endedAt?: Date;
  signedAt?: Date;
  serviceStartedAt?: Date;
  renewalCycle?: SelectOption<ContractRenewalCycle> | null;
}

export class ContractDTO implements TimeToRenewalForm {
  signedAt?: Date;
  endedAt?: Date;
  serviceStartedAt?: Date;
  renewalCycle?: SelectOption<ContractRenewalCycle> | null;
  name?: string;

  constructor(data?: TimeToRenewalData | null) {
    this.renewalCycle =
      billingFrequencyOptions.find((o) => o.value === data?.renewalCycle) ??
      null;
    this.signedAt = data?.signedAt && new Date(data.signedAt);
    this.endedAt = data?.endedAt && new Date(data.endedAt);
    this.serviceStartedAt =
      data?.serviceStartedAt && new Date(data.serviceStartedAt);
    this.name = data?.name ?? '';
  }

  static toForm(data?: TimeToRenewalData | null): TimeToRenewalForm {
    return new ContractDTO(data);
  }

  static toPayload(data: TimeToRenewalForm): TimeToRenewalForm {
    return {
      serviceStartedAt: data?.serviceStartedAt,
      signedAt: data?.signedAt,
      endedAt: data?.endedAt,
      renewalCycle: data?.renewalCycle,
      name: data?.name,
    };
  }
}
