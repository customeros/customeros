import { TenantBillingProfile } from '@graphql/types';
import { SelectOption } from '@shared/types/SelectOptions';
import { countryOptions } from '@shared/util/countryOptions';
import { getCurrencyOptions } from '@shared/util/currencyOptions';

export interface TenantBillingDetails {
  zip?: string | null;
  phone?: string | null;
  locality?: string | null;
  legalName?: string | null;
  vatNumber?: string | null;
  addressLine1?: string | null;
  addressLine2?: string | null;
  addressLine3?: string | null;
  canPayWithCard?: boolean | null;
  sendInvoicesBcc?: string | null;
  sendInvoicesFrom?: string | null;
  canPayWithPigeon?: boolean | null;
  country?: SelectOption<string> | null;
  canPayWithBankTransfer?: boolean | null;
  baseCurrency?: SelectOption<string> | null;
}

export class TenantBillingDetailsDto implements TenantBillingDetails {
  phone?: string | null;
  addressLine1?: string | null;
  addressLine2?: string | null;
  addressLine3?: string | null;
  baseCurrency?: SelectOption<string> | null;
  locality?: string | null;
  country?: SelectOption<string> | null;
  zip?: string | null;
  legalName?: string | null;
  domesticPaymentsBankInfo?: string | null;
  internationalPaymentsBankInfo?: string | null;
  canPayWithBankTransfer;
  canPayWithPigeon;
  sendInvoicesFrom;
  sendInvoicesBcc;
  vatNumber;

  constructor(
    data?: (TenantBillingProfile & { baseCurrency?: string | null }) | null,
  ) {
    this.phone = data?.phone;
    this.addressLine1 = data?.addressLine1;
    this.addressLine2 = data?.addressLine2;
    this.addressLine3 = data?.addressLine3;
    this.locality = data?.locality;
    this.country = countryOptions.find((i) => data?.country === i.value);
    this.zip = data?.zip;
    this.legalName = data?.legalName;
    this.canPayWithBankTransfer = data?.canPayWithBankTransfer;
    this.canPayWithPigeon = data?.canPayWithPigeon;
    this.sendInvoicesFrom = data?.sendInvoicesFrom?.split('@')[0] ?? '';
    this.sendInvoicesBcc = data?.sendInvoicesBcc ?? '';
    this.vatNumber = data?.vatNumber;
    this.baseCurrency = getCurrencyOptions().find(
      (i) => data?.baseCurrency === i.value,
    );
  }
}
