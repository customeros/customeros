import { OrganizationOnboardingPlansQuery } from '@organization/src/graphql/organizationOnboardingPlans.generated';

export type PlanDatum =
  OrganizationOnboardingPlansQuery['organizationPlansForOrganization'][number];

export type MilestoneDatum = PlanDatum['milestones'][number];
export type TaskDatum = MilestoneDatum['items'][number];

export type NewMilestoneInput = {
  name: string;
  order: number;
  dueDate: string;
  items: string[];
};
