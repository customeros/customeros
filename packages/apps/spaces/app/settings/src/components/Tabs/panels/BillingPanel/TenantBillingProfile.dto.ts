import { SelectOption } from '@shared/types/SelectOptions';
import { countryOptions } from '@shared/util/countryOptions';
import {
  TenantBillingProfile,
  TenantBillingProfileInput,
} from '@graphql/types';

export interface TenantBillingDetails {
  zip?: string | null;
  email?: string | null;
  phone?: string | null;
  locality?: string | null;
  legalName?: string | null;
  addressLine1?: string | null;
  addressLine2?: string | null;
  addressLine3?: string | null;
  country?: SelectOption<string> | null;
  domesticPaymentsBankInfo?: string | null;
  internationalPaymentsBankInfo?: string | null;
}

export class TenantBillingDetailsDto implements TenantBillingDetails {
  email?: string | null;
  phone?: string | null;
  addressLine1?: string | null;
  addressLine2?: string | null;
  addressLine3?: string | null;
  locality?: string | null;
  country?: SelectOption<string> | null;
  zip?: string | null;
  legalName?: string | null;
  domesticPaymentsBankInfo?: string | null;
  internationalPaymentsBankInfo?: string | null;

  constructor(data?: TenantBillingProfile | null) {
    this.email = data?.email;
    this.phone = data?.phone;
    this.addressLine1 = data?.addressLine1;
    this.addressLine2 = data?.addressLine2;
    this.addressLine3 = data?.addressLine3;
    this.locality = data?.locality;
    this.country = countryOptions.find((i) => data?.country === i.value);
    this.zip = data?.zip;
    this.legalName = data?.legalName;
    this.domesticPaymentsBankInfo = data?.domesticPaymentsBankInfo;
    this.internationalPaymentsBankInfo = data?.internationalPaymentsBankInfo;
  }

  static toForm(data?: TenantBillingProfile): TenantBillingDetails {
    const formData = new TenantBillingDetailsDto(data);

    return {
      ...formData,
    };
  }

  static toPayload(data: TenantBillingDetails): TenantBillingProfileInput {
    return {
      email: data?.email,
      phone: data?.phone,
      addressLine1: data?.addressLine1,
      addressLine2: data?.addressLine2,
      addressLine3: data?.addressLine3,
      locality: data?.locality,
      country: data?.country?.value ?? '',
      zip: data?.zip,
      legalName: data?.legalName,
      domesticPaymentsBankInfo: data?.domesticPaymentsBankInfo,
      internationalPaymentsBankInfo: data?.internationalPaymentsBankInfo,
      canPayWithDirectDebitACH: false, // todo change when implemented
      canPayWithDirectDebitSEPA: false,
      canPayWithDirectDebitBacs: false,
      canPayWithCard: false,
      canPayWithPigeon: false,
      sendInvoicesFrom: '',
      vatNumber: '',
    };
  }
}
