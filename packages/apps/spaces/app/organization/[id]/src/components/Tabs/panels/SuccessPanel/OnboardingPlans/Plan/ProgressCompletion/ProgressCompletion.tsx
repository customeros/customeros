import { forwardRef } from 'react';

import fill from 'lodash/fill';
import chunk from 'lodash/chunk';

import { Flex } from '@ui/layout/Flex';
import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { wave } from '@ui/utils/keyframes';
import { Check } from '@ui/media/icons/Check';
import { Box, BoxProps } from '@ui/layout/Box';
import { Portal, useDisclosure } from '@ui/utils';
import { DateTimeUtils } from '@spaces/utils/date';
import { Card, CardBody } from '@ui/presentation/Card';
import { CheckSquareBroken } from '@ui/media/icons/CheckSquareBroken';
import {
  OnboardingPlanStatus,
  OnboardingPlanMilestoneStatus,
} from '@graphql/types';
import {
  Popover,
  PopoverBody,
  PopoverArrow,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover';

import { getMilestoneDueDate } from '../utils';
import { PlanDatum, MilestoneDatum } from '../../types';

const getCompletion = (milestone: MilestoneDatum) => {
  switch (milestone.statusDetails.status) {
    case OnboardingPlanMilestoneStatus.NotStarted:
    case OnboardingPlanMilestoneStatus.NotStartedLate:
      return 0;
    case OnboardingPlanMilestoneStatus.Started:
    case OnboardingPlanMilestoneStatus.StartedLate:
      return 50;
    case OnboardingPlanMilestoneStatus.Done:
    case OnboardingPlanMilestoneStatus.DoneLate:
      return 100;
    default:
      return 0;
  }
};

const getColorScheme = (milestone: MilestoneDatum) => {
  switch (milestone.statusDetails.status) {
    case OnboardingPlanMilestoneStatus.NotStarted:
      return 'gray';
    case OnboardingPlanMilestoneStatus.NotStartedLate:
    case OnboardingPlanMilestoneStatus.StartedLate:
    case OnboardingPlanMilestoneStatus.DoneLate:
      return 'warning';
    case OnboardingPlanMilestoneStatus.Done:
    case OnboardingPlanMilestoneStatus.Started:
      return 'success';
    default:
      return 'gray';
  }
};

interface ProgressCompletionProps {
  plan: PlanDatum;
}

export const ProgressCompletion = ({ plan }: ProgressCompletionProps) => {
  const milestones = plan.milestones.filter((m) => !m.retired);

  const rows = chunk(milestones, 7).map((row, index, arr) => {
    if (arr.length === 1) return row;
    if (index !== arr.length - 1) {
      return row;
    }

    if (row.length === 7) {
      return row;
    }

    const remaining = 7 - row.length;

    return [...row, ...fill(new Array(remaining), 'EMPTY')];
  });

  const isDone = [
    OnboardingPlanStatus.Done,
    OnboardingPlanStatus.DoneLate,
  ].includes(plan.statusDetails.status);

  const isDoneLate =
    plan.statusDetails.status === OnboardingPlanStatus.DoneLate;

  const doneOnLabel = (() => {
    const days = Math.abs(
      DateTimeUtils.differenceInDays(
        plan.createdAt,
        plan.statusDetails.updatedAt,
      ),
    );
    const formattedDate = DateTimeUtils.format(
      plan.statusDetails.updatedAt,
      DateTimeUtils.dateWithAbreviatedMonth,
    );

    const prefix = isDoneLate ? 'Done late on ' : 'Done on';
    const sufix = days === 1 ? 'day' : 'days';

    return `${prefix} ${formattedDate} • Took ${days} ${sufix}`;
  })();

  if (isDone)
    return (
      <Flex mx='1' mb='2'>
        <Text>{doneOnLabel}</Text>
      </Flex>
    );

  return (
    <Card variant='outlinedElevated' mb='2' mx='1' mt='3'>
      <CardBody>
        <VStack spacing='2' w='full'>
          {rows.map((row, index) => (
            <Flex
              w='full'
              gap='0.5'
              key={index}
              flexWrap='wrap'
              align='center'
              position='relative'
              justify={'space-between'}
            >
              {row.map((milestone, index, arr) => {
                const isHidden = typeof milestone === 'string';
                const firstHiddenIndex = arr.indexOf('EMPTY');
                const lastMilestoneIndex =
                  firstHiddenIndex === -1 ? null : firstHiddenIndex - 1;

                const isLast = lastMilestoneIndex
                  ? index === lastMilestoneIndex
                  : index === row.length - 1;

                const itemsCount = !isHidden ? milestone?.items?.length : 0;
                const completedItemsCount = !isHidden
                  ? milestone?.items?.filter((item) =>
                      [
                        OnboardingPlanMilestoneStatus.Done,
                        OnboardingPlanMilestoneStatus.DoneLate,
                      ].includes(milestone.statusDetails.status),
                    ).length
                  : 0;

                return (
                  <ProgressStep
                    isLast={isLast}
                    isHidden={isHidden}
                    itemsCount={itemsCount}
                    key={isHidden ? index : milestone.id}
                    completedItemsCount={completedItemsCount}
                    dueDate={isHidden ? '' : milestone?.dueDate}
                    completion={isHidden ? 0 : getCompletion(milestone)}
                    label={isHidden ? 'hidden milestone' : milestone?.name}
                    colorScheme={isHidden ? 'gray' : getColorScheme(milestone)}
                  />
                );
              })}
            </Flex>
          ))}
        </VStack>
      </CardBody>
    </Card>
  );
};

interface ProgressStep {
  label?: string;
  isLast: boolean;
  dueDate?: string;
  isHidden?: boolean;
  colorScheme: string;
  itemsCount?: number;
  completion?: 0 | 50 | 100;
  completedItemsCount?: number;
}

const ProgressStep = ({
  label,
  isLast,
  dueDate,
  isHidden,
  completion,
  colorScheme,
  itemsCount = 0,
  completedItemsCount,
}: ProgressStep) => {
  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <>
      <Popover isOpen={isOpen} onOpen={onOpen} onClose={onClose}>
        <PopoverTrigger>
          <ProgressCompletionCircle
            onMouseEnter={onOpen}
            onMouseLeave={onClose}
            completion={completion}
            colorScheme={colorScheme}
            visibility={isHidden ? 'hidden' : 'visible'}
          />
        </PopoverTrigger>
        <Portal>
          <PopoverContent
            bg='gray.700'
            color='gray.25'
            width='200px'
            borderRadius='8px'
          >
            <PopoverArrow bg='gray.700' />
            <PopoverBody>
              <Text noOfLines={1} fontWeight='medium'>
                {label}
              </Text>
              <Flex justify='space-between'>
                {dueDate && (
                  <Text color='gray.300'>
                    {getMilestoneDueDate(dueDate, completion === 100)}
                  </Text>
                )}
                <Flex align='center' gap='2'>
                  <CheckSquareBroken color='gray.300' />
                  {itemsCount > 0 && (
                    <Text color='gray.300'>{`${completedItemsCount}/${itemsCount}`}</Text>
                  )}
                </Flex>
              </Flex>
            </PopoverBody>
          </PopoverContent>
        </Portal>
      </Popover>
      {!isLast && (
        <Flex
          h='2px'
          flex='1'
          bg='gray.200'
          visibility={isHidden ? 'hidden' : 'visible'}
        />
      )}
    </>
  );
};

interface ProgressCompletionCircleProps extends BoxProps {
  colorScheme?: string;
  completion?: 0 | 50 | 100;
  visibility?: 'hidden' | 'visible';
}

const ProgressCompletionCircle = forwardRef<
  HTMLDivElement,
  ProgressCompletionCircleProps
>(
  (
    { completion = 0, colorScheme = 'gray', visibility = 'visible', ...props },
    ref,
  ) => {
    const _completion = completion === 0 ? 110 : completion === 100 ? 0 : 50;

    const borderColor = (() => {
      switch (colorScheme) {
        case 'gray':
          return 'gray.400';
        case 'warning':
          return 'warning.500';
        case 'success':
          return 'success.500';
        default:
          return 'gray.400';
      }
    })();

    return (
      <Box
        w='24px'
        h='24px'
        ref={ref}
        bg='white'
        maxW='24px'
        overflow='hidden'
        borderRadius='50%'
        border='1px solid'
        position='relative'
        visibility={visibility}
        borderColor={borderColor}
        {...props}
      >
        {completion === 100 ? (
          <Flex
            w='full'
            h='full'
            align='center'
            justify='center'
            bg={`${colorScheme}.500`}
          >
            <Check color='white' mr='1px' />
          </Flex>
        ) : (
          <>
            <Box
              bg={`${colorScheme}.300`}
              position='absolute'
              top={`${_completion}%`}
              h='200%'
              w='200%'
              borderRadius='35%'
              left='-50%'
              transform='rotate(340deg)'
              transition='all 5s ease'
              animation={`${wave} 45s linear infinite`}
            />
            <Box
              bg={`${colorScheme}.500`}
              position='absolute'
              top={`${_completion}%`}
              h='200%'
              w='200%'
              borderRadius='42%'
              left='-50%'
              transform='rotate(295deg)'
              transition='all 5s ease'
              animation={`${wave} 30s linear infinite`}
            />
          </>
        )}
      </Box>
    );
  },
);
