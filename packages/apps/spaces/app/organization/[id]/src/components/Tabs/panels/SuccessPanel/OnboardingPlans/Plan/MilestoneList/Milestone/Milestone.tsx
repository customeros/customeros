import { useForm, useField } from 'react-inverted-form';
import { useRef, useMemo, useState, useEffect, ChangeEvent } from 'react';

import { produce } from 'immer';
import isEqual from 'lodash/isEqual';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { useOutsideClick } from '@ui/utils';
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
import { MilestoneMenu } from './MilestoneMenu';
import { MilestoneName } from './MilestoneName';
import { MilestoneDatum } from '../../../types';
import { StatusCheckbox } from './StatusCheckbox';
import { MilestoneDueDate } from './MilestoneDueDate';

type MilestoneForm = {
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

interface MilestoneProps {
  isLast?: boolean;
  isOpen?: boolean;
  isActiveItem?: boolean;
  milestone: MilestoneDatum;
  onToggle?: (id: string) => void;
  onRemove?: (id: string) => void;
  onDuplicate?: (id: string) => void;
  onMakeOptional?: (id: string) => void;
  onSync?: (milestone: MilestoneDatum) => void;
}

export const Milestone = ({
  isLast,
  isOpen,
  onSync,
  onToggle,
  onRemove,
  milestone,
  onDuplicate,
  isActiveItem,
  onMakeOptional,
}: MilestoneProps) => {
  const cardRef = useRef<HTMLDivElement>(null);
  const [isHovered, setIsHovered] = useState(false);
  const isMutating = milestone.id.startsWith('temp');
  const hasTasks = milestone.items?.length > 0;

  const defaultValues = useMemo(
    () => mapMilestoneToForm(milestone),
    [
      milestone.id,
      milestone.name,
      milestone.order,
      milestone.optional,
      JSON.stringify(milestone.items || []),
      JSON.stringify(milestone.statusDetails || {}),
    ],
  );
  const formId = `${milestone.id}-plan-milestone-form`;

  const { setDefaultValues, state } = useForm<MilestoneForm>({
    formId,
    defaultValues,
    stateReducer: (_, action, next) => {
      const nextMilestone = mapFormToMilestone(milestone, next.values);

      if (action.type === 'FIELD_CHANGE') {
        if (action.payload.name === 'items') {
          const nextValues = produce(next, (draft) => {
            draft.values.statusDetails.status = computeMilestoneStatus(
              draft.values,
            );
          });

          const nextMilestone = mapFormToMilestone(
            milestone,
            nextValues.values,
          );
          if (!isEqual(nextMilestone, milestone)) {
            onSync?.(nextMilestone);
          }

          return nextValues;
        }

        if (!isEqual(nextMilestone, milestone)) {
          onSync?.(nextMilestone);
        }
      }

      return next;
    },
  });

  const isPastDueDate = new Date(state?.values?.dueDate) < new Date();
  const isLate = useMemo(() => {
    if (
      [
        OnboardingPlanMilestoneStatus.DoneLate,
        OnboardingPlanMilestoneStatus.StartedLate,
        OnboardingPlanMilestoneStatus.NotStartedLate,
      ].includes(milestone?.statusDetails?.status)
    )
      return true;

    if (!state?.values?.items?.length) return isPastDueDate;

    return state?.values?.items?.some((item) =>
      [
        OnboardingPlanMilestoneItemStatus.DoneLate,
        OnboardingPlanMilestoneItemStatus.NotDoneLate,
        OnboardingPlanMilestoneItemStatus.SkippedLate,
      ].includes(item.status),
    );
  }, [state.values.items, milestone?.statusDetails?.status, isPastDueDate]);

  const isChecked = useMemo(() => {
    if (
      [
        OnboardingPlanMilestoneStatus.Done,
        OnboardingPlanMilestoneStatus.DoneLate,
      ].includes(milestone?.statusDetails?.status)
    )
      return true;

    if (!state?.values?.items?.length) return false;

    return state?.values?.items?.every((item) =>
      [
        OnboardingPlanMilestoneItemStatus.Done,
        OnboardingPlanMilestoneItemStatus.DoneLate,
        OnboardingPlanMilestoneItemStatus.Skipped,
        OnboardingPlanMilestoneItemStatus.SkippedLate,
      ].includes(item.status),
    );
  }, [state.values.items, milestone?.statusDetails?.status]);

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
  const handleDuplicate = () => onDuplicate?.(milestone.id);
  const handleMakeOptional = () => onMakeOptional?.(milestone.id);

  useEffect(() => {
    setDefaultValues(defaultValues);
  }, [
    milestone.id,
    milestone.name,
    milestone.order,
    milestone.optional,
    JSON.stringify(milestone.items || []),
    JSON.stringify(milestone.statusDetails || {}),
  ]);

  useOutsideClick({ ref: cardRef, handler: handleToggle, enabled: isOpen });

  return (
    <Card
      w='full'
      ref={cardRef}
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
      <CardBody pl='3' pb='4'>
        <Flex flexDir='column' justify='flex-start'>
          <Flex align='center'>
            <Flex>
              <PlanStatusCheckbox
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
              onDuplicate={handleDuplicate}
              isOptional={milestone.optional}
              onMakeOptional={handleMakeOptional}
              {...hoveredProps}
            />
          </Flex>

          <Flex align='center' justify='space-between' pl='6' mb='2'>
            <MilestoneDueDate
              isDone={isChecked}
              value={state.values.dueDate ?? milestone.dueDate}
            />
            {!!milestone?.items?.length && (
              <Flex align='center' gap='6px' mr='0.5' {...hoveredProps}>
                <CheckSquareBroken color='gray.400' />
                <Text fontSize='sm' color='gray.500'>
                  {`${doneTakskCount}/${milestone?.items?.length ?? 0}`}
                </Text>
              </Flex>
            )}
          </Flex>

          <Collapse in={isOpen} animateOpacity style={{ overflow: 'visible' }}>
            <Tasks formId={formId} defaultValue={milestone.items} />
          </Collapse>
        </Flex>
      </CardBody>
    </Card>
  );
};

interface PlanStatusCheckboxProps {
  formId: string;
  readOnly?: boolean;
  colorScheme: string;
  showCustomIcon?: boolean;
  onToggleMilestone?: () => void;
}

const PlanStatusCheckbox = ({
  formId,
  readOnly,
  colorScheme,
  showCustomIcon,
  onToggleMilestone,
}: PlanStatusCheckboxProps) => {
  const { getInputProps } = useField('statusDetails', formId);
  const { value, onChange, onBlur, ...inputProps } = getInputProps();

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (readOnly) {
      onToggleMilestone?.();

      return;
    }

    onChange?.({
      ...value,
      status: e.target.checked ? 'DONE' : 'NOT_STARTED',
      updatedAt: new Date().toISOString(),
    });
  };

  const handleBlur = (e: ChangeEvent<HTMLInputElement>) => {
    onChange?.({
      ...value,
      status: e.target.checked ? 'DONE' : 'NOT_STARTED',
      updatedAt: new Date().toISOString(),
    });
  };

  return (
    <StatusCheckbox
      mr='2'
      size='md'
      onBlur={handleBlur}
      onChange={handleChange}
      colorScheme={colorScheme}
      showCustomIcon={showCustomIcon}
      isChecked={[
        OnboardingPlanMilestoneStatus.Done,
        OnboardingPlanMilestoneStatus.DoneLate,
      ].includes(value.status as unknown as OnboardingPlanMilestoneStatus)}
      {...inputProps}
    />
  );
};

const mapMilestoneToForm = (
  milestone?: MilestoneDatum | null,
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
    name: formValues.name,
    items: formValues.items,
    dueDate: formValues.dueDate,
    statusDetails: formValues.statusDetails,
  };
};

/// Status Logic

const checkTaskDone = (status: OnboardingPlanMilestoneItemStatus) => {
  return [
    OnboardingPlanMilestoneItemStatus.Done,
    OnboardingPlanMilestoneItemStatus.DoneLate,
    OnboardingPlanMilestoneItemStatus.Skipped,
    OnboardingPlanMilestoneItemStatus.SkippedLate,
  ].includes(status);
};

const checkMilestoneDone = (status: OnboardingPlanMilestoneStatus) => {
  return [
    OnboardingPlanMilestoneStatus.Done,
    OnboardingPlanMilestoneStatus.DoneLate,
  ].includes(status);
};

const computeMilestoneStatus = (milestone: MilestoneForm) => {
  const isPastDueDate = new Date(milestone?.dueDate) < new Date();
  const allTasksDone = milestone?.items?.every((i) => checkTaskDone(i.status));
  const isMilestonePreviouslyDone = checkMilestoneDone(
    milestone.statusDetails.status,
  );

  if (allTasksDone && !isMilestonePreviouslyDone) {
    return isPastDueDate
      ? OnboardingPlanMilestoneStatus.DoneLate
      : OnboardingPlanMilestoneStatus.Done;
  }
  if (!allTasksDone && isMilestonePreviouslyDone) {
    return isPastDueDate
      ? OnboardingPlanMilestoneStatus.NotStartedLate
      : OnboardingPlanMilestoneStatus.NotStarted;
  }

  return milestone.statusDetails.status;
};
