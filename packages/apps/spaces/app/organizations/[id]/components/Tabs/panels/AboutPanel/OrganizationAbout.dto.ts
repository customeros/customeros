import { Organization } from '@graphql/types';
import { SelectOption } from '@shared/types/SelectOptions';

import { OrganizationQuery } from '../../graphql/organization.generated';
import { UpdateOrganizationMutationVariables } from '../../graphql/updateOrganization.generated';
import {
  stageOptions,
  industryOptions,
  employeesOptions,
  relationshipOptions,
  businessTypeOptions,
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
    | 'description'
    | 'website'
    | 'targetAudience'
    | 'valueProposition'
    | 'domains'
  > {
  businessType: SelectOption | null;
  lastFunding: string;
  relationship: SelectOption | null;
  stage: SelectOption | null;
  lastStage: string;
  industry: SelectOption | null;
  employees: SelectOption<number> | null;
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
  lastFunding: string;
  relationship: SelectOption | null;
  stage: SelectOption | null;
  lastStage: string;
  domains: string[];

  constructor(data?: Partial<OrganizationQuery['organization']> | null) {
    this.id = data?.id || '';
    this.name = data?.name || '';
    this.description = data?.description || '';
    this.website = data?.website || '';
    this.industry =
      mergedIndustryOptions.find((i) => data?.industry === i.value) || null;
    this.targetAudience = data?.targetAudience || '';
    this.valueProposition = data?.valueProposition || '';
    this.employees =
      employeesOptions.find((i) => data?.employees === i.value) || null;
    this.businessType =
      businessTypeOptions.find((i) => data?.market === i.value) || null;
    this.lastFunding = '';
    this.relationship =
      relationshipOptions.find(
        (i) => data?.relationshipStages?.[0]?.relationship === i.value,
      ) || null;
    this.stage =
      stageOptions.find(
        (i) => data?.relationshipStages?.[0]?.stage === i.value,
      ) || null;
    this.lastStage = '';
    this.domains = data?.domains || [];
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
      industry: data.industry?.value,
      employees: data.employees?.value,
      domains: data.domains,
    } as UpdateOrganizationMutationVariables['input'];
  }
}
