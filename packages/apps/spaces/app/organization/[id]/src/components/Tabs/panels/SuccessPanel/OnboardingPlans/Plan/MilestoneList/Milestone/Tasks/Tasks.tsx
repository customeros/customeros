import { useState, useCallback } from 'react';
import { useField } from 'react-inverted-form';

import { produce } from 'immer';

import { Button } from '@ui/form/Button';
import { VStack } from '@ui/layout/Stack';
import { Plus } from '@ui/media/icons/Plus';
import { OnboardingPlanMilestoneItemStatus } from '@graphql/types';

import { Task } from './Task';
import { TaskDatum } from '../../../../types';

interface TasksProps {
  formId: string;
  defaultValue: TaskDatum[];
}

export const Tasks = ({ formId, defaultValue }: TasksProps) => {
  const { getInputProps } = useField('items', formId);
  const { value, onChange } = getInputProps();
  const [focusedIndex, setFocusedIndex] = useState<number | null>(null);

  const handleTaskBlur = useCallback(
    () => setFocusedIndex(null),
    [setFocusedIndex],
  );
  const handleTaskFocus = useCallback(
    (index: number) => () => setFocusedIndex(index),
    [setFocusedIndex],
  );

  const handleAddTask = () => {
    const nextItems = produce<TaskDatum[]>(value, (draft) => {
      draft.push({
        text: '',
        updatedAt: new Date().toISOString(),
        status: OnboardingPlanMilestoneItemStatus.NotDone,
      });
    });

    onChange?.(nextItems);
  };

  return (
    <VStack align='flex-start' spacing='1' pl='6'>
      {(value as TaskDatum[])?.map((_, idx) => (
        <Task
          key={idx}
          index={idx}
          formId={formId}
          onInputBlur={handleTaskBlur}
          isFocused={focusedIndex === idx}
          onInputFocus={handleTaskFocus(idx)}
          shouldFocus={focusedIndex ? focusedIndex + 1 === idx : false}
        />
      ))}
      <Button
        ml='-12px'
        size='sm'
        variant='ghost'
        color='gray.500'
        fontWeight='normal'
        onClick={handleAddTask}
        leftIcon={<Plus color='gray.400' />}
      >
        Add task
      </Button>
    </VStack>
  );
};
