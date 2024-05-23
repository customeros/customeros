import { SelectOption } from '@shared/types/SelectOptions';
import { OrganizationQuery } from '@organization/graphql/organization.generated';
import {
  Social,
  Organization,
  OrganizationStage,
  OrganizationUpdateInput,
  OrganizationRelationship,
} from '@graphql/types';

import {
  stageOptions,
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
  stage?: SelectOption<OrganizationStage> | null;
  relationship?: SelectOption<OrganizationRelationship> | null;
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
  stage?: SelectOption<OrganizationStage> | null;
  relationship?: SelectOption<OrganizationRelationship> | null;

  constructor(data?: Partial<OrganizationQuery['organization']> | null) {
    this.id = data?.id || '';
    this.name = data?.name ?? 'Unnamed';
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
    this.stage = stageOptions.find((i) => data?.stage === i.value);
    this.relationship = relationshipOptions.find(
      (i) => data?.relationship === i.value,
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
      stage: data.stage?.value,
      relationship: data.relationship?.value,
    } as OrganizationUpdateInput;
  }
}
