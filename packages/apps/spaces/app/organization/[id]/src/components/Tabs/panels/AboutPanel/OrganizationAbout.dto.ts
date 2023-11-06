import { SelectOption } from '@shared/types/SelectOptions';
import { Social, Organization, OrganizationUpdateInput } from '@graphql/types';
import { OrganizationQuery } from '@organization/src/graphql/organization.generated';

import {
  industryOptions,
  employeesOptions,
  businessTypeOptions,
  relationshipOptions,
  lastFundingRoundOptions,
} from './util';

const mergedIndustryOptions = industryOptions.reduce(
  (acc, curr) => [...acc, ...curr.options],
  [] as SelectOption[],
);

export interface OrganizationAboutForm
  extends Pick<
    Organization,
    | 'id'
    | 'name'
    | 'note'
    | 'description'
    | 'website'
    | 'targetAudience'
    | 'valueProposition'
    | 'lastFundingAmount'
  > {
  industry: SelectOption | null;
  businessType: SelectOption | null;
  lastFundingRound: SelectOption | null;
  socials: Pick<Social, 'id' | 'url'>[];
  employees: SelectOption<number> | null;
  isCustomer?: SelectOption<boolean> | null;
}

export class OrganizationAboutFormDto implements OrganizationAboutForm {
  id: string;
  name: string;
  description: string;
  website: string;
  industry: SelectOption | null;
  targetAudience: string;
  valueProposition: string;
  employees: SelectOption<number> | null;
  businessType: SelectOption | null;
  lastFundingRound: SelectOption | null;
  lastFundingAmount?: string;
  socials: Pick<Social, 'id' | 'url'>[];
  isCustomer?: SelectOption<boolean> | null;

  constructor(data?: Partial<OrganizationQuery['organization']> | null) {
    this.id = data?.id || '';
    this.name = data?.name || '';
    this.description = data?.description || '';
    this.website = data?.website || '';

    // Display industry even when data is not matching GICS
    this.industry = data?.industry
      ? mergedIndustryOptions.find((i) => data?.industry === i.value) || {
          label: data.industry,
          value: data.industry,
        }
      : null;
    this.targetAudience = data?.targetAudience || '';
    this.valueProposition = data?.valueProposition || '';
    this.employees =
      employeesOptions.find((i) => data?.employees <= i.value) || null;
    this.businessType =
      businessTypeOptions.find((i) => data?.market === i.value) || null;
    this.lastFundingRound =
      lastFundingRoundOptions.find((i) => data?.lastFundingRound === i.value) ||
      null;
    this.lastFundingAmount = data?.lastFundingAmount ?? '';
    this.socials = data?.socials || [];
    this.isCustomer = relationshipOptions.find(
      (i) => data?.isCustomer === i.value,
    );
  }

  static toForm(data: OrganizationQuery) {
    return new OrganizationAboutFormDto(data.organization);
  }

  static toPayload(data: OrganizationAboutForm) {
    return {
      id: data.id,
      name: data.name,
      description: data.description,
      market: data.businessType?.value,
      website: data.website,
      note: data.note,
      industry: data.industry?.value,
      employees: data.employees?.value,
      targetAudience: data.targetAudience,
      valueProposition: data.valueProposition,
      lastFundingRound: data.lastFundingRound?.value,
      lastFundingAmount: data.lastFundingAmount,
      isCustomer: data.isCustomer?.value,
    } as OrganizationUpdateInput;
  }
}
