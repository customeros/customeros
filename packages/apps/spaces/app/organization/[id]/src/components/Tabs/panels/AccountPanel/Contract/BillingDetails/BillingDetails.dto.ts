import { SelectOption } from '@shared/types/SelectOptions';
import { Currency, ContractUpdateInput } from '@graphql/types';
import {
  countryOptions,
  currencyOptions,
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

  constructor(data?: BillingDetailsForm | null) {
    this.zip = data?.zip;
    this.locality = data?.locality;
    this.invoiceEmail = data?.invoiceEmail;
    this.addressLine1 = data?.addressLine1;
    this.country = countryOptions.find((i) => data?.country?.value === i.value);
    this.currency = currencyOptions.find(
      (i) => data?.currency?.value === i.value,
    );
    this.addressLine2 = data?.addressLine2;
    this.organizationLegalName = data?.organizationLegalName;
    this.contractUrl = data?.contractUrl;
  }

  static toForm(data?: BillingDetailsForm | null): BillingDetailsForm {
    const formData = new BillingDetailsDto(data);

    return {
      ...formData,
    };
  }

  static toPayload(
    data: BillingDetailsForm,
  ): Omit<ContractUpdateInput, 'contractId'> {
    return {
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
