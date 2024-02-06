import {
  OnboardingPlanMilestoneStatus,
  OnboardingPlanMilestoneItemStatus,
} from '@graphql/types';

export type MilestoneForm = {
  id: string;
  name: string;
  dueDate: string;
  items: {
    text: string;
    updatedAt: string;
    status: OnboardingPlanMilestoneItemStatus;
  }[];
  statusDetails: {
    text: string;
    updatedAt: string;
    status: OnboardingPlanMilestoneStatus;
  };
};
