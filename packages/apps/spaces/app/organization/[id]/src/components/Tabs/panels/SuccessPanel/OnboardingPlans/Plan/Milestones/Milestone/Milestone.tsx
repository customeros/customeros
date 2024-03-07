import { useForm } from 'react-inverted-form';
import { useMemo, useState, useEffect } from 'react';

import { useDebounce } from 'rooks';
import isEqual from 'lodash/isEqual';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';
import { pulseOpacity } from '@ui/utils/keyframes';
import { Collapse } from '@ui/transitions/Collapse';
import { Card, CardBody } from '@ui/presentation/Card';
import { ChevronExpand } from '@ui/media/icons/ChevronExpand';
import { ChevronCollapse } from '@ui/media/icons/ChevronCollapse';
import { CheckSquareBroken } from '@ui/media/icons/CheckSquareBroken';
import {
  OnboardingPlanMilestoneStatus,
  OnboardingPlanMilestoneItemStatus,
} from '@graphql/types';

import { Tasks } from './Tasks';
import { MilestoneForm } from './types';
import { TaskDatum, MilestoneDatum } from '../../../types';
import {
  checkMilestoneDone,
  checkMilestoneLate,
  computeMilestoneStatus,
} from './utils';
import {
  MilestoneMenu,
  MilestoneName,
  MilestoneDueDate,
  MilestoneCheckbox,
} from './components';

interface MilestoneProps {
  isLast?: boolean;
  isOpen?: boolean;
  isActiveItem?: boolean;
  onToggle?: (id: string) => void;
  onRemove?: (id: string) => void;
  onSync?: (milestone: MilestoneDatum) => void;
  milestone: MilestoneDatum & { items: TaskDatum[] };
}

export const Milestone = ({
  isLast,
  isOpen,
  onSync,
  onToggle,
  onRemove,
  milestone,
  isActiveItem,
}: MilestoneProps) => {
  const [isHovered, setIsHovered] = useState(false);
  const [isDueDateOpen, setIsDueDateOpen] = useState(false);
  const isMutating = milestone.id.startsWith('temp');
  const hasTasks = milestone.items?.length > 0;

  const defaultValues = useMemo(
    () => mapMilestoneToForm(milestone),
    [
      milestone.id,
      milestone.name,
      milestone.order,
      milestone.dueDate,
      milestone.optional,
      JSON.stringify(milestone.items || []),
      JSON.stringify(milestone.statusDetails || {}),
    ],
  );
  const formId = `${milestone.id}-plan-milestone-form`;

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const debouncedSync = useDebounce(onSync as any, 500);

  const { setDefaultValues, state } = useForm<MilestoneForm>({
    formId,
    defaultValues,
    stateReducer: (_, action, next) => {
      const nextMilestone = mapFormToMilestone(milestone, next.values);

      if (action.type === 'FIELD_CHANGE') {
        if (action.payload.name === 'items') {
          const nextValues = {
            ...next.values,
            statusDetails: {
              ...next.values.statusDetails,
              status: computeMilestoneStatus(next.values),
            },
          };

          const nextMilestone = mapFormToMilestone(milestone, nextValues);
          if (!isEqual(nextMilestone, milestone)) {
            debouncedSync?.(nextMilestone);
          }

          return {
            ...next,
            values: nextValues,
          };
        }

        if (!isEqual(nextMilestone, milestone)) {
          onSync?.(nextMilestone);
        }
      }

      if (action.type === 'FIELD_BLUR') {
        if (action.payload.name === 'items') {
          const nextValues = {
            ...next.values,
            statusDetails: {
              ...next.values.statusDetails,
              status: computeMilestoneStatus(next.values),
            },
          };

          const nextMilestone = mapFormToMilestone(milestone, nextValues);

          if (!isEqual(nextMilestone, milestone)) {
            onSync?.(nextMilestone);
          }
        }
      }

      return next;
    },
  });

  const isLate = useMemo(
    () => checkMilestoneLate(state.values),
    [state.values],
  );

  const isChecked = useMemo(
    () => checkMilestoneDone(state.values),
    [state.values],
  );

  const mostRecentDoneTask = useMemo(() => {
    const doneTasks = milestone?.items?.filter(
      (i) =>
        i.status === OnboardingPlanMilestoneItemStatus.Done ||
        i.status === OnboardingPlanMilestoneItemStatus.DoneLate,
    );

    return doneTasks?.reduce((acc, curr) => {
      if (!acc) return curr;
      const accDate = new Date(acc?.updatedAt as string);
      const currDate = new Date(curr?.updatedAt as string);

      return accDate > currDate ? acc : curr;
    }, null as TaskDatum | null);
  }, [milestone.items]);

  const checkboxColorScheme = useMemo(() => {
    if (isChecked) return isLate ? 'warning' : 'success';
    if (isLate) return 'warning';

    return 'gray';
  }, [isChecked, isLate]);

  const doneTakskCount = useMemo(
    () =>
      milestone?.items?.filter((i) =>
        [
          OnboardingPlanMilestoneItemStatus.Done,
          OnboardingPlanMilestoneItemStatus.DoneLate,
        ].includes(i.status),
      ).length,
    [milestone.items],
  );

  const hoveredProps = useMemo(
    () => ({
      opacity: isHovered || isOpen ? 1 : 0,
      transition: 'opacity 0.2s ease-out',
    }),
    [isHovered, isOpen],
  );

  const handleToggle = () => onToggle?.(milestone.id);
  const handleRetire = () => onRemove?.(milestone.id);
  const handleOpenDueDate = () => setIsDueDateOpen(true);
  const handleCloseDueDate = () => setIsDueDateOpen(false);

  useEffect(() => {
    setDefaultValues(defaultValues);
  }, [
    milestone.id,
    milestone.name,
    milestone.order,
    milestone.dueDate,
    milestone.optional,
  ]);

  return (
    <Card
      w='full'
      variant='outlinedElevated'
      cursor={isActiveItem ? 'grabbing' : undefined}
      pointerEvents={isMutating ? 'none' : undefined}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      animation={
        isMutating
          ? `${pulseOpacity} 0.7s infinite alternate ease-in-out`
          : 'unset'
      }
    >
      <CardBody pb='4' py='3'>
        <Flex flexDir='column' justify='flex-start'>
          <Flex align='center'>
            <Flex>
              <MilestoneCheckbox
                formId={formId}
                onToggleMilestone={handleToggle}
                colorScheme={checkboxColorScheme}
                readOnly={milestone.items?.length > 0}
                showCustomIcon={milestone.items?.length > 0 && !isChecked}
              />
            </Flex>
            <MilestoneName
              isActiveItem
              formId={formId}
              isLast={isLast}
              isMilestoneOpen={isOpen}
              defaultValue={milestone.name}
              onToggleMilestone={handleToggle}
            />
            {hasTasks && (
              <IconButton
                size='xs'
                variant='ghost'
                aria-label={
                  isOpen ? 'Collapse Master Plan' : 'Expand Master Plan'
                }
                icon={
                  isOpen ? (
                    <ChevronCollapse color='gray.400' />
                  ) : (
                    <ChevronExpand color='gray.400' />
                  )
                }
                onClick={handleToggle}
                {...hoveredProps}
              />
            )}
            <MilestoneMenu
              onRetire={handleRetire}
              isMilestoneDone={isChecked}
              onSetDueDate={handleOpenDueDate}
              {...hoveredProps}
            />
          </Flex>

          <Flex
            pl='6'
            align='center'
            justify='space-between'
            mb={isOpen ? '2' : '0'}
            transition='margin 0.2s ease-in-out'
          >
            <MilestoneDueDate
              formId={formId}
              isDone={isChecked}
              isOpen={isDueDateOpen}
              onOpen={handleOpenDueDate}
              onClose={handleCloseDueDate}
              minDate={mostRecentDoneTask?.updatedAt}
              status={state?.values?.statusDetails?.status}
            />
            {!!milestone?.items?.length && (
              <Flex align='center' gap='6px' mr='0.5'>
                <CheckSquareBroken color='gray.400' />
                <Text fontSize='sm' color='gray.500'>
                  {`${doneTakskCount}/${milestone?.items?.length ?? 0}`}
                </Text>
              </Flex>
            )}
          </Flex>

          <Collapse in={isOpen} animateOpacity style={{ overflow: 'visible' }}>
            <Tasks
              formId={formId}
              milestoneDueDate={state?.values?.dueDate ?? ''}
            />
          </Collapse>
        </Flex>
      </CardBody>
    </Card>
  );
};

const mapMilestoneToForm = (
  milestone?: (MilestoneDatum & { items: TaskDatum[] }) | null,
): MilestoneForm => {
  return {
    id: milestone?.id ?? '',
    name: milestone?.name ?? '',
    items: milestone?.items ?? [],
    dueDate: milestone?.dueDate ?? null,
    statusDetails: milestone?.statusDetails ?? {
      text: '',
      status: OnboardingPlanMilestoneStatus.NotStarted,
      updatedAt: '',
    },
  };
};

const mapFormToMilestone = (
  milestoneDatum: MilestoneDatum,
  formValues: MilestoneForm,
): MilestoneDatum => {
  return {
    ...milestoneDatum,
    name: formValues?.name,
    dueDate: formValues?.dueDate,
    items: formValues?.items ?? [],
    statusDetails: formValues?.statusDetails,
  };
};
