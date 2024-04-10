import { utcToZonedTime } from 'date-fns-tz';

import { SelectOption } from '@shared/types/SelectOptions';
import { currencyOptions } from '@shared/util/currencyOptions';
import { GetContractQuery } from '@organization/src/graphql/getContract.generated';
import {
  Currency,
  ContractUpdateInput,
  ContractBillingCycle,
} from '@graphql/types';

import { paymentDueOptions, contractBillingCycleOptions } from '../../utils';

export interface ContractDetailsForm {
  payOnline?: boolean | null;
  autoRenew?: boolean | null;
  contractEnded?: Date | null;
  serviceStarted?: Date | null;
  invoicingStarted?: Date | null;
  canPayWithCard?: boolean | null;
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

  constructor(data?: Partial<GetContractQuery['contract']> | null) {
    this.committedPeriodInMonths = data?.committedPeriodInMonths;
    this.currency = currencyOptions.find((i) => data?.currency === i.value);
    this.contractEnded = data?.contractEnded;
    this.billingCycle =
      [...contractBillingCycleOptions].find(
        ({ value }) => value === data?.billingDetails?.billingCycle,
      ) ?? undefined;
    this.invoicingStarted = data?.billingDetails?.invoicingStarted
      ? utcToZonedTime(data.billingDetails.invoicingStarted, 'UTC')
      : null;
    this.serviceStarted =
      data?.serviceStarted && utcToZonedTime(data?.serviceStarted, 'UTC');
    this.canPayWithBankTransfer = data?.billingDetails?.canPayWithBankTransfer;
    this.canPayWithDirectDebit = data?.billingDetails?.canPayWithDirectDebit;
    this.canPayWithCard = data?.billingDetails?.canPayWithCard;
    this.payOnline = data?.billingDetails?.payOnline;
    this.payAutomatically = data?.billingDetails?.payAutomatically;
    this.dueDays = paymentDueOptions.find(
      (e) => e.value === data?.billingDetails?.dueDays,
    );
    this.autoRenew = data?.autoRenew;
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
      canPayWithDirectDebit: data?.canPayWithDirectDebit,
      canPayWithBankTransfer: data?.canPayWithBankTransfer,
      autoRenew: data?.autoRenew,
      serviceStarted: data?.serviceStarted,
      billingDetails: {
        payOnline: data?.payOnline,
        payAutomatically: data?.payAutomatically,
        canPayWithBankTransfer: data?.canPayWithBankTransfer,
        canPayWithDirectDebit: data?.canPayWithDirectDebit,
        dueDays: data?.dueDays?.value,
        invoicingStarted: data?.invoicingStarted,
        billingCycle: data?.billingCycle?.value,
      },
      patch: true,
    };
  }
}
