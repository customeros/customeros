import { Text } from '@ui/typography/Text';
import { OnboardingPlanMilestoneStatus } from '@graphql/types';

import { getMilestoneDueDate } from '../utils';

interface PlanDueDateProps {
  value: string;
  isDone: boolean;
  status?: OnboardingPlanMilestoneStatus;
}

export const PlanDueDate = ({
  value,
  isDone,
  status = OnboardingPlanMilestoneStatus.NotStarted,
}: PlanDueDateProps) => {
  if (isDone) return null;

  return (
    <Text as='label' fontSize='sm' color='gray.500' whiteSpace='nowrap'>
      {getMilestoneDueDate(value, status, isDone).toLowerCase()}
    </Text>
  );
};
