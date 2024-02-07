import { useForm, useField } from 'react-inverted-form';
import { useRef, useMemo, useState, useEffect, ChangeEvent } from 'react';

import isEqual from 'lodash/isEqual';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { useOutsideClick } from '@ui/utils';
import { IconButton } from '@ui/form/IconButton';
import { pulseOpacity } from '@ui/utils/keyframes';
import { Collapse } from '@ui/transitions/Collapse';
import { Card, CardBody } from '@ui/presentation/Card';
import { useThrottle } from '@shared/hooks/useThrottle';
import { ChevronExpand } from '@ui/media/icons/ChevronExpand';
import { ChevronCollapse } from '@ui/media/icons/ChevronCollapse';
import { CheckSquareBroken } from '@ui/media/icons/CheckSquareBroken';
import {
  OnboardingPlanMilestoneStatus,
  OnboardingPlanMilestoneItemStatus,
} from '@graphql/types';

import { Tasks } from './Tasks';
import { MilestoneForm } from './types';
import { MilestoneMenu } from './MilestoneMenu';
import { MilestoneName } from './MilestoneName';
import { MilestoneDatum } from '../../../types';
import { StatusCheckbox } from './StatusCheckbox';
import { MilestoneDueDate } from './MilestoneDueDate';
import {
  checkMilestoneDone,
  checkMilestoneLate,
  computeMilestoneStatus,
} from './utils';

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
      milestone.dueDate,
      milestone.optional,
      JSON.stringify(milestone.items || []),
      JSON.stringify(milestone.statusDetails || {}),
    ],
  );
  const formId = `${milestone.id}-plan-milestone-form`;

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const torttledSync = useThrottle(onSync as any, 500, []);

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
            torttledSync?.(nextMilestone);
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
    milestone.dueDate,
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
      <CardBody pb='4' py='3'>
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

          <Flex
            pl='6'
            align='center'
            justify='space-between'
            mb={isOpen ? '2' : '0'}
          >
            <MilestoneDueDate isDone={isChecked} value={state.values.dueDate} />
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

  return (
    <StatusCheckbox
      mr='2'
      size='md'
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
    name: formValues?.name,
    items: formValues?.items,
    dueDate: formValues?.dueDate,
    statusDetails: formValues?.statusDetails,
  };
};
