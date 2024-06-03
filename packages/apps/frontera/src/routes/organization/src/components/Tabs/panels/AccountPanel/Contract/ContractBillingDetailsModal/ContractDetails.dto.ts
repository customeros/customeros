import { toZonedTime } from 'date-fns-tz';

import { DateTimeUtils } from '@utils/date';
import { SelectOption } from '@shared/types/SelectOptions';
import { currencyOptions } from '@shared/util/currencyOptions';
import { GetContractQuery } from '@organization/graphql/getContract.generated';
import {
  Currency,
  ContractUpdateInput,
  ContractBillingCycle,
} from '@graphql/types';

import { paymentDueOptions, contractBillingCycleOptions } from '../../utils';

export interface ContractDetailsForm {
  check?: boolean | null;
  payOnline?: boolean | null;
  autoRenew?: boolean | null;
  contractEnded?: Date | null;
  serviceStarted?: Date | null;
  contractName?: string | null;
  invoicingStarted?: Date | null;
  canPayWithCard?: boolean | null;
  billingEnabled?: boolean | null;
  payAutomatically?: boolean | null;
  dueDays?: SelectOption<number> | null;
  canPayWithDirectDebit?: boolean | null;
  committedPeriodInMonths?: number | null;
  canPayWithBankTransfer?: boolean | null;
  currency?: SelectOption<Currency> | null;
  billingCycle?: SelectOption<ContractBillingCycle> | null;
}

export class ContractDetailsDto implements ContractDetailsForm {
  committedPeriodInMonths?: number | null;
  currency?: SelectOption<Currency> | null;
  contractEnded?: Date | null;
  billingCycle?: SelectOption<ContractBillingCycle> | null;
  invoicingStarted?: Date | null;
  serviceStarted?: Date | null;
  canPayWithBankTransfer?: boolean | null;
  canPayWithDirectDebit?: boolean | null;
  canPayWithCard?: boolean | null;
  payOnline?: boolean | null;
  payAutomatically?: boolean | null;
  dueDays?: SelectOption<number> | null;
  autoRenew?: boolean | null;
  check?: boolean | null;
  billingEnabled?: boolean | null;
  contractName?: string | null;

  constructor(data?: Partial<GetContractQuery['contract']> | null) {
    this.committedPeriodInMonths = data?.committedPeriodInMonths;
    this.currency = currencyOptions.find((i) => data?.currency === i.value);
    this.contractEnded = data?.contractEnded;
    this.billingCycle =
      [...contractBillingCycleOptions].find(
        ({ value }) => value === data?.billingDetails?.billingCycle,
      ) ?? contractBillingCycleOptions[0];
    this.invoicingStarted = data?.billingDetails?.invoicingStarted
      ? toZonedTime(data.billingDetails.invoicingStarted, 'UTC')
      : data?.serviceStarted
      ? DateTimeUtils.addMonth(data?.serviceStarted, 1)
      : null;
    this.serviceStarted =
      data?.serviceStarted && toZonedTime(data?.serviceStarted, 'UTC');
    this.canPayWithBankTransfer = data?.billingDetails?.canPayWithBankTransfer;
    this.canPayWithDirectDebit = data?.billingDetails?.canPayWithDirectDebit;
    this.canPayWithCard = data?.billingDetails?.canPayWithCard;
    this.payOnline = data?.billingDetails?.payOnline;
    this.payAutomatically = data?.billingDetails?.payAutomatically;
    this.dueDays = paymentDueOptions.find(
      (e) => e.value === data?.billingDetails?.dueDays,
    );
    this.autoRenew = data?.autoRenew;
    this.billingEnabled = data?.billingEnabled;
    this.contractName = data?.contractName;
    this.check = data?.billingDetails?.check ?? false;
  }

  static toForm(
    data?: Partial<GetContractQuery['contract']> | null,
  ): ContractDetailsForm {
    const formData = new ContractDetailsDto(data);

    return {
      ...formData,
    };
  }

  static toPayload(
    data: ContractDetailsForm,
  ): Omit<ContractUpdateInput, 'contractId'> {
    return {
      currency: data?.currency?.value,
      contractName: data?.contractName,
      canPayWithDirectDebit: data?.canPayWithDirectDebit,
      canPayWithBankTransfer: data?.canPayWithBankTransfer,
      autoRenew: data?.autoRenew,
      serviceStarted: data?.serviceStarted
        ? DateTimeUtils.getUTCDateAtMidnight(data.serviceStarted)
        : undefined,
      committedPeriodInMonths: data?.committedPeriodInMonths,
      billingEnabled: data?.billingEnabled,
      billingDetails: {
        payOnline: data?.payOnline,
        payAutomatically: data?.payAutomatically,
        canPayWithBankTransfer: data?.canPayWithBankTransfer,
        canPayWithDirectDebit: data?.canPayWithDirectDebit,
        dueDays: data?.dueDays?.value,
        invoicingStarted: data?.invoicingStarted
          ? DateTimeUtils.getUTCDateAtMidnight(data.invoicingStarted)
          : undefined,
        billingCycle: data?.billingCycle?.value,
        check: data?.check,
      },
      patch: true,
    };
  }
}
