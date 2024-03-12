import { SelectOption } from '@shared/types/SelectOptions';
import { countryOptions } from '@shared/util/countryOptions';
import { Currency, ContractUpdateInput } from '@graphql/types';
import { getCurrencyOptions } from '@shared/util/currencyOptions';
import { GetContractQuery } from '@organization/src/graphql/getContract.generated';

export interface BillingDetailsForm {
  zip?: string | null;
  locality?: string | null;
  payOnline?: boolean | null;
  contractUrl?: string | null;
  invoiceEmail?: string | null;
  addressLine1?: string | null;
  addressLine2?: string | null;
  canPayWithCard?: boolean | null;
  payAutomatically?: boolean | null;
  organizationLegalName?: string | null;
  country?: SelectOption<string> | null;
  canPayWithDirectDebit?: boolean | null;
  canPayWithBankTransfer?: boolean | null;
  currency?: SelectOption<Currency> | null;
}

export class BillingDetailsDto implements BillingDetailsForm {
  zip?: string | null;
  locality?: string | null;
  invoiceEmail?: string | null;
  addressLine1?: string | null;
  addressLine2?: string | null;
  organizationLegalName?: string | null;
  country?: SelectOption<string> | null;
  currency?: SelectOption<Currency> | null;
  contractUrl?: string | null;
  canPayWithCard?: boolean | null;
  canPayWithDirectDebit?: boolean | null;
  canPayWithBankTransfer?: boolean | null;
  payAutomatically?: boolean | null;
  payOnline?: boolean | null;

  constructor(data?: Partial<GetContractQuery['contract']> | null) {
    this.zip = data?.billingDetails?.postalCode ?? '';
    this.locality = data?.billingDetails?.locality ?? '';
    this.invoiceEmail = data?.billingDetails?.billingEmail ?? '';
    this.addressLine1 = data?.billingDetails?.addressLine1 ?? '';
    this.canPayWithCard = data?.billingDetails?.canPayWithCard;
    this.canPayWithDirectDebit = data?.billingDetails?.canPayWithDirectDebit;
    this.canPayWithBankTransfer = data?.billingDetails?.canPayWithBankTransfer;

    this.country = countryOptions.find(
      (i) => data?.billingDetails?.country === i.value,
    );
    this.currency = getCurrencyOptions().find(
      (i) => data?.currency === i.value,
    );
    this.addressLine2 = data?.billingDetails?.addressLine2 ?? '';
    this.organizationLegalName = data?.organizationLegalName ?? '';
    this.contractUrl = data?.contractUrl ?? '';
    this.payOnline = data?.billingDetails?.payOnline;
    this.payAutomatically = data?.billingDetails?.payAutomatically;
  }

  static toForm(
    data?: Partial<GetContractQuery['contract']> | null,
  ): BillingDetailsForm {
    const formData = new BillingDetailsDto(data);

    return {
      ...formData,
    };
  }

  static toPayload(
    data: BillingDetailsForm,
  ): Omit<ContractUpdateInput, 'contractId'> {
    return {
      contractUrl: data?.contractUrl,
      zip: data?.zip,
      locality: data?.locality,
      invoiceEmail: data?.invoiceEmail,
      addressLine1: data?.addressLine1,
      country: data?.country?.value ?? '',
      currency: data?.currency?.value,
      addressLine2: data?.addressLine2,
      organizationLegalName: data?.organizationLegalName,
      canPayWithCard: data?.canPayWithCard,
      canPayWithDirectDebit: data?.canPayWithDirectDebit,
      canPayWithBankTransfer: data?.canPayWithBankTransfer,
      billingDetails: {
        payOnline: data?.payOnline,
        payAutomatically: data?.payAutomatically,
        canPayWithBankTransfer: data?.canPayWithBankTransfer,
      },
      patch: true,
    };
  }
}
