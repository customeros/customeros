import { useForm } from 'react-inverted-form';
import { useRef, useMemo, useState, useEffect, MutableRefObject } from 'react';

import isEqual from 'lodash/isEqual';
import { useMergeRefs } from 'rooks';
import { CSS } from '@dnd-kit/utilities';
import { useSortable } from '@dnd-kit/sortable';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { useOutsideClick } from '@ui/utils';
import { IconButton } from '@ui/form/IconButton';
import { pulseOpacity } from '@ui/utils/keyframes';
import { Collapse } from '@ui/transitions/Collapse';
import { Card, CardBody } from '@ui/presentation/Card';
import { HandleDrag } from '@ui/media/icons/HandleDrag';
import { ChevronExpand } from '@ui/media/icons/ChevronExpand';
import { ChevronCollapse } from '@ui/media/icons/ChevronCollapse';
import { CheckSquareBroken } from '@ui/media/icons/CheckSquareBroken';

import { Tasks } from './Tasks';
import { MilestoneDatum } from '../../types';
import { MilestoneMenu } from './MilestoneMenu';
import { MilestoneName } from './MilestoneName';
import { MilestoneDuration } from './MilestoneDuration';

type MilestoneForm = {
  id: string;
  name: string;
  order: number;
  items: string[];
  duration: number;
  optional: boolean;
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
  shouldFocusNameRef?: MutableRefObject<boolean>;
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
  shouldFocusNameRef,
}: MilestoneProps) => {
  const cardRef = useRef<HTMLDivElement>(null);
  const [isHovered, setIsHovered] = useState(false);
  const isMutating = milestone.id.startsWith('temp');

  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
    setActivatorNodeRef,
  } = useSortable({
    id: milestone.id,
  });

  const mergedCardRefs = useMergeRefs(cardRef, setNodeRef);
  const transformStyle = `${CSS.Translate.toString(transform) ?? ''} ${
    isActiveItem ? 'scale(1.02) rotate(4deg)' : ''
  }`.trim();

  const defaultValues = useMemo(
    () => mapMilestoneToForm(milestone),
    [
      milestone.id,
      milestone.name,
      milestone.order,
      milestone.optional,
      JSON.stringify(milestone.items),
    ],
  );
  const formId = `${milestone.id}-milestone-form`;

  const { setDefaultValues, state } = useForm<MilestoneForm>({
    formId,
    defaultValues,
    stateReducer: (_, action, next) => {
      if (action.type === 'FIELD_BLUR') {
        const nextMilestone = mapFormToMilestone(next.values);
        if (!isEqual(nextMilestone, milestone)) {
          onSync?.(nextMilestone);
        }

        if (shouldFocusNameRef && shouldFocusNameRef.current !== null)
          shouldFocusNameRef.current = false;
      }

      return next;
    },
  });

  const hoveredProps = useMemo(
    () => ({
      opacity: isDragging || isHovered || isOpen ? 1 : 0,
      transition: 'opacity 0.2s ease-out',
    }),
    [isDragging, isHovered, isOpen],
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
    JSON.stringify(milestone.items),
  ]);

  useOutsideClick({ ref: cardRef, handler: handleToggle, enabled: isOpen });

  return (
    <Card
      w='full'
      ref={mergedCardRefs}
      variant='outlinedElevated'
      transition={transition}
      transform={transformStyle}
      cursor={isActiveItem ? 'grabbing' : undefined}
      boxShadow={isDragging ? 'unset' : undefined}
      pointerEvents={isMutating ? 'none' : undefined}
      borderColor={isDragging ? 'gray.100' : undefined}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      animation={
        isMutating
          ? `${pulseOpacity} 0.7s infinite alternate ease-in-out`
          : 'unset'
      }
    >
      <CardBody pl='6'>
        <Flex flexDir='column' justify='flex-start'>
          <Flex align='center'>
            <Flex
              left='1'
              cursor='grab'
              position='absolute'
              ref={setActivatorNodeRef}
              {...listeners}
              {...attributes}
              {...hoveredProps}
            >
              <HandleDrag color='gray.400' />
            </Flex>
            <MilestoneName
              formId={formId}
              isLast={isLast}
              isMilestoneOpen={isOpen}
              isActiveItem={isActiveItem}
              defaultValue={milestone.name}
              onToggleMilestone={handleToggle}
              shouldFocus={!isMutating && isLast && shouldFocusNameRef?.current}
            />
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
            <MilestoneMenu
              onRetire={handleRetire}
              onDuplicate={handleDuplicate}
              isOptional={milestone.optional}
              onMakeOptional={handleMakeOptional}
              {...hoveredProps}
            />
          </Flex>

          <Flex align='center' justify='space-between' mb='2'>
            <MilestoneDuration
              formId={formId}
              isMilestoneOpen={isOpen}
              isActiveItem={isActiveItem}
              onToggleMilestone={handleToggle}
              defaultValue={state.values.duration ?? milestone.durationHours}
            />
            {!!milestone?.items?.length && (
              <Flex align='center' gap='1.5' mr='0.25' {...hoveredProps}>
                <CheckSquareBroken color='gray.400' />
                <Text fontSize='sm' color='gray.500'>
                  {`0/${milestone?.items?.length ?? 0}`}
                </Text>
              </Flex>
            )}
          </Flex>

          <Collapse in={isOpen} animateOpacity style={{ overflow: 'visible' }}>
            <Tasks
              formId={formId}
              isActiveItem={isActiveItem}
              defaultValue={milestone.items}
            />
          </Collapse>
        </Flex>
      </CardBody>
      {isDragging && (
        <Flex
          position='absolute'
          top='0'
          left='0'
          right='0'
          bottom='0'
          bg='gray.100'
          borderRadius='7px'
        />
      )}
    </Card>
  );
};

const mapMilestoneToForm = (
  milestone?: MilestoneDatum | null,
): MilestoneForm => {
  return {
    id: milestone?.id ?? '',
    name: milestone?.name ?? '',
    items: milestone?.items ?? [],
    order: milestone?.order ?? 0,
    optional: milestone?.optional ?? false,
    duration: milestone?.durationHours ? milestone?.durationHours / 24 : 1,
  };
};

const mapFormToMilestone = (formValues: MilestoneForm): MilestoneDatum => {
  return {
    id: formValues.id,
    name: formValues.name,
    items: formValues.items,
    optional: formValues.optional,
    durationHours: formValues.duration * 24,
    order: formValues.order,
    retired: false,
  };
};
