import { useForm } from 'react-inverted-form';
import { useRef, useMemo, useState, useEffect } from 'react';

import { useMergeRefs } from 'rooks';
import isEqual from 'lodash/isEqual';
import { CSS } from '@dnd-kit/utilities';
import { useSortable } from '@dnd-kit/sortable';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { useOutsideClick } from '@ui/utils';
import { IconButton } from '@ui/form/IconButton';
import { Collapse } from '@ui/transitions/Collapse';
import { Card, CardBody } from '@ui/presentation/Card';
import { HandleDrag } from '@ui/media/icons/HandleDrag';
import { ChevronExpand } from '@ui/media/icons/ChevronExpand';
import { FormInput, FormResizableInput } from '@ui/form/Input';
import { ChevronCollapse } from '@ui/media/icons/ChevronCollapse';
import { CheckSquareBroken } from '@ui/media/icons/CheckSquareBroken';

import { Tasks } from './Tasks';
import { MilestoneDatum } from './types';
import { MilestoneMenu } from './MilestoneMenu';

type MilestoneForm = {
  id: string;
  name: string;
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
  const nameInputRef = useRef<HTMLInputElement>(null);
  const [isHovered, setIsHovered] = useState(false);

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
    [milestone.id],
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
      }

      return next;
    },
  });

  const hoveredProps = useMemo(
    () => ({
      opacity: isDragging || isHovered || isOpen ? 1 : 0,
      transition: 'opacity 0.2s ease-out',
    }),
    [isHovered, isOpen],
  );

  const handleToggle = () => onToggle?.(milestone.id);
  const handleRetire = () => onRemove?.(milestone.id);
  const handleDuplicate = () => onDuplicate?.(milestone.id);
  const handleMakeOptional = () => onMakeOptional?.(milestone.id);

  useEffect(() => {
    if (isLast && nameInputRef?.current) {
      nameInputRef?.current?.focus();
      setTimeout(() => nameInputRef?.current?.select(), 0);
    }
  }, [isLast]);

  useEffect(() => {
    setDefaultValues(defaultValues);
  }, [milestone.id]);

  useOutsideClick({ ref: cardRef, handler: handleToggle, enabled: isOpen });

  return (
    <Card
      w='full'
      ref={mergedCardRefs}
      variant='outlinedElevated'
      transition={transition}
      transform={transformStyle}
      opacity={isDragging ? 0.4 : undefined}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      <CardBody pl='6'>
        <Flex flexDir='column' justify='flex-start'>
          <Flex align='center'>
            <Flex
              left='1'
              position='absolute'
              ref={setActivatorNodeRef}
              {...listeners}
              {...attributes}
              {...hoveredProps}
            >
              <HandleDrag color='gray.400' />
            </Flex>
            <FormInput
              name='name'
              formId={formId}
              ref={nameInputRef}
              variant='unstyled'
              fontWeight='medium'
              borderRadius='unset'
              placeholder='Milestone name'
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
              onMakeOptional={handleMakeOptional}
              {...hoveredProps}
            />
          </Flex>

          <Flex align='center' justify='space-between' mb='2'>
            <Flex align='center' gap='1'>
              <Text
                as='label'
                fontSize='sm'
                color='gray.500'
                whiteSpace='nowrap'
                htmlFor='duration-input'
              >
                Max duration:
              </Text>
              <FormResizableInput
                min={1}
                size='sm'
                type='number'
                name='duration'
                formId={formId}
                variant='unstyled'
                id='duration-input'
                borderRadius='unset'
              />
              <Text fontSize='sm' color='gray.500' whiteSpace='nowrap'>
                {state.values?.duration === 1 ? 'day' : 'days'}
              </Text>
            </Flex>
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
            <Tasks formId={formId} />
          </Collapse>
        </Flex>
      </CardBody>
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
    optional: milestone?.optional ?? false,
    duration: milestone?.durationHours ?? 0,
  };
};

const mapFormToMilestone = (formValues: MilestoneForm): MilestoneDatum => {
  return {
    id: formValues.id,
    name: formValues.name,
    items: formValues.items,
    optional: formValues.optional,
    durationHours: formValues.duration,
    order: 0,
    retired: false,
  };
};
