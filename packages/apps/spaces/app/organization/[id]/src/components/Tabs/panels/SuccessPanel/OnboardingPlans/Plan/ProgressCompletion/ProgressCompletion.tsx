import { Box } from '@ui/layout/Box';
import { wave } from '@ui/utils/keyframes';
import { Card, CardBody } from '@ui/presentation/Card';
import { OnboardingPlanMilestoneStatus } from '@shared/types/__generated__/graphql.types';
import {
  Step,
  Stepper,
  useSteps,
  StepIcon,
  StepStatus,
  StepSeparator,
  StepIndicator,
} from '@ui/navigation/Stepper';

import { PlanDatum } from '../../types';

interface ProgressCompletionProps {
  plan: PlanDatum;
}

export const ProgressCompletion = ({ plan: _ }: ProgressCompletionProps) => {
  const { activeStep } = useSteps({
    index: 1,
    count: 3,
  });

  return (
    <Card variant='outlinedElevated' mb='2' mx='1'>
      <CardBody>
        <Stepper colorScheme='success' index={activeStep}>
          <ProgressStep
            index={0}
            completion={50}
            milestoneStatus={OnboardingPlanMilestoneStatus.Done}
          />
          <ProgressStep
            index={1}
            completion={50}
            milestoneStatus={OnboardingPlanMilestoneStatus.Started}
          />
          <ProgressStep
            index={2}
            completion={50}
            milestoneStatus={OnboardingPlanMilestoneStatus.NotStarted}
          />
        </Stepper>
      </CardBody>
    </Card>
  );
};

interface ProgressStep {
  index: number;
  completion?: number;
  milestoneStatus: OnboardingPlanMilestoneStatus;
}

const ProgressStep = ({ index, completion, milestoneStatus }: ProgressStep) => {
  return (
    <Step key={index}>
      <StepIndicator border='unset'>
        <StepStatus
          complete={<StepIcon />}
          incomplete={<ProgressCompletionCircle completion={completion} />}
          active={
            <ProgressCompletionCircle completion={15} colorScheme='success' />
          }
        />
      </StepIndicator>
      <StepSeparator />
    </Step>
  );
};

interface ProgressCompletionCircleProps {
  completion?: number;
  colorScheme?: string;
}

const ProgressCompletionCircle = ({
  completion = 100,
  colorScheme = 'gray',
}: ProgressCompletionCircleProps) => {
  return (
    <Box
      w='24px'
      h='24px'
      bg='white'
      position='relative'
      border='1px solid'
      borderColor={`${colorScheme}.500`}
      borderRadius='50%'
      overflow='hidden'
    >
      <Box
        bg={`${colorScheme}.300`}
        position='absolute'
        top={`${completion}%`}
        h='200%'
        w='200%'
        borderRadius='38%'
        left='-50%'
        transform='rotate(360deg)'
        transition='all 5s ease'
        animation={`${wave} 30s linear infinite`}
      />
      <Box
        bg={`${colorScheme}.500`}
        position='absolute'
        top={`${completion}%`}
        h='200%'
        w='200%'
        borderRadius='38%'
        left='-50%'
        transform='rotate(360deg)'
        transition='all 5s ease'
        animation={`${wave} 45s linear infinite`}
      />
    </Box>
  );
};
