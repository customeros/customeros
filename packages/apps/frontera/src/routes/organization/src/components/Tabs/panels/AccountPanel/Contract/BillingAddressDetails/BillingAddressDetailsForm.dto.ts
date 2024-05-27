import { ContractUpdateInput } from '@graphql/types';
import { SelectOption } from '@shared/types/SelectOptions';
import { countryOptions } from '@shared/util/countryOptions';
import { GetContractQuery } from '@organization/graphql/getContract.generated';

export interface BillingAddressDetailsFormDto {
  region?: string | null;
  locality?: string | null;
  postalCode?: string | null;
  addressLine1?: string | null;
  addressLine2?: string | null;
  billingEmail?: string | null;
  organizationLegalName?: string | null;
  country?: SelectOption<string> | null;
  billingEmailCC?: Array<SelectOption<string>> | null;
  billingEmailBCC?: Array<SelectOption<string>> | null;
}

export class BillingDetailsDto implements BillingAddressDetailsFormDto {
  postalCode?: string | null;
  region?: string | null;
  locality?: string | null;
  addressLine1?: string | null;
  addressLine2?: string | null;
  organizationLegalName?: string | null;
  country?: SelectOption<string> | null;
  billingEmail?: string | null;
  billingEmailCC?: Array<SelectOption<string>> | null;
  billingEmailBCC?: Array<SelectOption<string>> | null;

  constructor(data?: Partial<GetContractQuery['contract']> | null) {
    this.postalCode = data?.billingDetails?.postalCode ?? '';
    this.locality = data?.billingDetails?.locality ?? '';
    this.addressLine1 = data?.billingDetails?.addressLine1 ?? '';
    this.country = countryOptions.find(
      (i) => data?.billingDetails?.country === i.value,
    );
    this.addressLine2 = data?.billingDetails?.addressLine2 ?? '';
    this.organizationLegalName = data?.organizationLegalName ?? '';
    this.region = data?.billingDetails?.region;
    this.billingEmail = data?.billingDetails?.billingEmail;
    this.billingEmailCC = data?.billingDetails?.billingEmailCC?.map((e) => ({
      label: e,
      value: e,
    }));
    this.billingEmailBCC = data?.billingDetails?.billingEmailBCC?.map((e) => ({
      label: e,
      value: e,
    }));
  }

  static toForm(
    data?: Partial<GetContractQuery['contract']> | null,
  ): BillingAddressDetailsFormDto {
    const formData = new BillingDetailsDto(data);

    return {
      ...formData,
    };
  }

  static toPayload(
    data: BillingAddressDetailsFormDto,
  ): Omit<ContractUpdateInput, 'contractId'> {
    return {
      billingDetails: {
        organizationLegalName: data?.organizationLegalName,
        addressLine1: data?.addressLine1,
        addressLine2: data?.addressLine2,
        region: data?.region,
        locality: data?.locality,
        country: data?.country?.value ?? '',
        postalCode: data?.postalCode,
        billingEmail: data?.billingEmail,
        billingEmailCC: data?.billingEmailCC?.map((e) => e?.value),
        billingEmailBCC: data?.billingEmailBCC?.map((e) => e?.value),
      },
      patch: true,
    };
  }
}
