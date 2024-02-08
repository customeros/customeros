import { memo, useRef } from 'react';
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
  defaultValue?: string[];
}

export const Tasks = memo(
  ({ formId }: TasksProps) => {
    const { getInputProps } = useField('items', formId);
    const { value, onChange } = getInputProps();
    const shouldFocusRef = useRef<number | null>(null);

    const handleAddTask = () => {
      const nextItems = produce<TaskDatum[]>(value, (draft) => {
        draft.push({
          text: '',
          updatedAt: new Date().toISOString(),
          status: OnboardingPlanMilestoneItemStatus.NotDone,
        });
      });

      if (shouldFocusRef) {
        shouldFocusRef.current = value.length;
      }
      onChange?.(nextItems);
    };

    return (
      <VStack align='flex-start' spacing='1' pl='6'>
        {(value as TaskDatum[])?.map((m, idx) => (
          <Task
            key={m.updatedAt}
            index={idx}
            formId={formId}
            defaultValue={m.text}
            shouldFocusRef={shouldFocusRef}
            isLast={idx === (value as TaskDatum[]).length - 1}
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
  },
  (prev, next) => JSON.stringify(prev) === JSON.stringify(next),
);
