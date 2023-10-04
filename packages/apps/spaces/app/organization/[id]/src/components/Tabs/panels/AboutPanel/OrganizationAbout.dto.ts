import { Organization, Social } from '@graphql/types';
import { SelectOption } from '@shared/types/SelectOptions';

import { OrganizationQuery } from '../../../../graphql/organization.generated';
import { UpdateOrganizationMutationVariables } from '../../../../graphql/updateOrganization.generated';
import {
  industryOptions,
  employeesOptions,
  otherStageOptions,
  relationshipOptions,
  businessTypeOptions,
  customerStageOptions,
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
  businessType: SelectOption | null;
  relationship: SelectOption | null;
  lastFundingRound: SelectOption | null;
  stage: SelectOption | null;
  industry: SelectOption | null;
  employees: SelectOption<number> | null;
  socials: Pick<Social, 'id' | 'url'>[];
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
  relationship: SelectOption | null;
  stage: SelectOption | null;
  lastFundingRound: SelectOption | null;
  lastFundingAmount?: string;
  socials: Pick<Social, 'id' | 'url'>[];

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
    this.relationship =
      relationshipOptions.find(
        (i) => data?.relationshipStages?.[0]?.relationship === i.value,
      ) || null;
    this.stage =
      (data?.relationshipStages?.[0]?.relationship === 'CUSTOMER'
        ? customerStageOptions
        : otherStageOptions
      ).find((i) => data?.relationshipStages?.[0]?.stage === i.value) || null;
    this.socials = data?.socials || [];
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
      patch: true,
    } as UpdateOrganizationMutationVariables['input'];
  }
}
