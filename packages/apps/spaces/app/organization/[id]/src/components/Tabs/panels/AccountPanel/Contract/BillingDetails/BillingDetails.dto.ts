import { SelectOption } from '@shared/types/SelectOptions';
import { Currency, ContractUpdateInput } from '@graphql/types';
import { GetContractQuery } from '@organization/src/graphql/getContract.generated';
import {
  countryOptions,
  getCurrencyOptions,
} from '@organization/src/components/Tabs/panels/AccountPanel/Contract/BillingDetails/utils';

export interface BillingDetailsForm {
  zip?: string | null;
  locality?: string | null;
  contractUrl?: string | null;
  invoiceEmail?: string | null;
  addressLine1?: string | null;
  addressLine2?: string | null;
  organizationLegalName?: string | null;
  country?: SelectOption<string> | null;
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

  constructor(data?: Partial<GetContractQuery['contract']> | null) {
    this.zip = data?.zip;
    this.locality = data?.locality;
    this.invoiceEmail = data?.invoiceEmail;
    this.addressLine1 = data?.addressLine1;
    this.country = countryOptions.find((i) => data?.country === i.value);
    this.currency = getCurrencyOptions().find(
      (i) => data?.currency === i.value,
    );
    this.addressLine2 = data?.addressLine2;
    this.organizationLegalName = data?.organizationLegalName;
    this.contractUrl = data?.contractUrl;
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
      patch: true,
    };
  }
}
