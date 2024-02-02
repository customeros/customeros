import { Box } from '@ui/layout/Box';
import { wave } from '@ui/utils/keyframes';
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

// import { Flex } from '@ui/layout/Flex';

interface ProgressCompletionProps {
  plan: PlanDatum;
}

export const ProgressCompletion = ({ plan: _ }: ProgressCompletionProps) => {
  const { activeStep } = useSteps({
    index: 1,
    count: 3,
  });

  return (
    <Stepper index={activeStep}>
      <Step key={0}>
        <StepIndicator>
          <StepStatus
            complete={<StepIcon />}
            incomplete={<ProgressCompletionCircle completion={50} />}
            active={
              <ProgressCompletionCircle completion={15} colorScheme='success' />
            }
          />
        </StepIndicator>
        <StepSeparator />
      </Step>

      <Step key={1}>
        <StepIndicator>
          <StepStatus
            complete={<StepIcon />}
            incomplete={<ProgressCompletionCircle completion={50} />}
            active={
              <ProgressCompletionCircle completion={15} colorScheme='success' />
            }
          />
        </StepIndicator>
        <StepSeparator />
      </Step>

      <Step key={2}>
        <StepIndicator>
          <StepStatus
            complete={<StepIcon />}
            incomplete={<ProgressCompletionCircle completion={50} />}
            active={
              <ProgressCompletionCircle completion={15} colorScheme='success' />
            }
          />
        </StepIndicator>
        <StepSeparator />
      </Step>
    </Stepper>
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
