import {
  OnboardingPlanMilestoneStatus,
  OnboardingPlanMilestoneItemStatus,
} from '@graphql/types';

export type MilestoneForm = {
  id: string;
  name: string;
  dueDate: string;
  statusDetails: {
    text: string;
    updatedAt: string;
    status: OnboardingPlanMilestoneStatus;
  };
  items: {
    uuid: string;
    text: string;
    updatedAt: string;
    status: OnboardingPlanMilestoneItemStatus;
  }[];
};
