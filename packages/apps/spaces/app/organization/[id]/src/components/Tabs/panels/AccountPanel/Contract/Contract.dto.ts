import { utcToZonedTime } from 'date-fns-tz';

import { SelectOption } from '@shared/types/SelectOptions';
import { UpdateContractMutationVariables } from '@organization/src/graphql/updateContract.generated';
import {
  Contract,
  ContractUpdateInput,
  ContractRenewalCycle,
  ContractBillingCycle,
} from '@graphql/types';
import {
  paymentDueOptions,
  autorenewalOptions,
  billingFrequencyOptions,
  contractBillingCycleOptions,
} from '@organization/src/components/Tabs/panels/AccountPanel/utils';

export interface TimeToRenewalForm {
  name?: string;
  endedAt?: Date;
  serviceStarted?: Date;
  invoicingStartDate?: Date;
  contractUrl?: string | null;
  committedPeriods?: string | null;
  dueDays?: SelectOption<number> | null;
  country?: SelectOption<string> | null;
  organizationLegalName?: string | null;
  autoRenew?: SelectOption<boolean> | null;
  billingEnabled?: SelectOption<boolean> | null;
  billingCycle?: SelectOption<ContractBillingCycle> | null;
  contractRenewalCycle?: SelectOption<
    ContractRenewalCycle | 'MULTI_YEAR'
  > | null;
}

export class ContractDTO implements TimeToRenewalForm {
  endedAt?: Date;
  invoicingStartDate?: Date;
  serviceStarted?: Date;
  contractRenewalCycle?: SelectOption<
    ContractRenewalCycle | 'MULTI_YEAR'
  > | null;
  name?: string;
  contractUrl?: string | null;
  renewalPeriods?: string | null;
  billingCycle?: SelectOption<ContractBillingCycle> | null;
  billingEnabled?: SelectOption<boolean> | null;
  autoRenew?: SelectOption<boolean> | null;
  dueDays?: SelectOption<number> | null;

  constructor(data?: Contract | null) {
    this.contractRenewalCycle =
      [...billingFrequencyOptions].find(({ value }) =>
        (data?.committedPeriods ?? 0) > 1
          ? value === 'MULTI_YEAR'
          : value === data?.contractRenewalCycle,
      ) ?? undefined;
    this.billingEnabled = data?.billingEnabled
      ? { label: 'Enabled', value: true }
      : { label: 'Disabled', value: false };
    this.billingCycle =
      [...contractBillingCycleOptions].find(
        ({ value }) => value === data?.billingDetails?.billingCycle,
      ) ?? undefined;
    this.endedAt =
      data?.contractEnded && utcToZonedTime(data?.contractEnded, 'UTC');
    this.invoicingStartDate =
      data?.billingDetails?.invoicingStarted &&
      utcToZonedTime(data?.billingDetails?.invoicingStarted, 'UTC');

    this.serviceStarted =
      data?.serviceStarted && utcToZonedTime(data?.serviceStarted, 'UTC');

    this.name = data?.contractName?.length
      ? data?.contractName
      : `${
          data?.billingDetails?.organizationLegalName?.length
            ? `${data?.billingDetails?.organizationLegalName}'s`
            : "Unnamed's"
        } contract`;
    this.contractUrl = data?.contractUrl ?? '';
    this.dueDays = paymentDueOptions.find(
      (e) => e.value === data?.billingDetails?.dueDays,
    );
    this.renewalPeriods = String(data?.committedPeriods ?? 2);
    this.autoRenew = autorenewalOptions.find(
      (e) => e.value === data?.autoRenew,
    );
  }

  static toForm(data?: Contract | null): TimeToRenewalForm {
    const formData = new ContractDTO(data);

    return {
      ...formData,
    };
  }

  static toPayload(
    data: Partial<ContractUpdateInput> & { contractId: string },
  ): UpdateContractMutationVariables {
    return {
      input: {
        patch: true,
        ...data,
      },
    };
  }
}
