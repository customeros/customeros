import { memo, useRef } from 'react';
import { useField } from 'react-inverted-form';

import { produce } from 'immer';
import setHours from 'date-fns/setHours';

import { Button } from '@ui/form/Button';
import { VStack } from '@ui/layout/Stack';
import { Plus } from '@ui/media/icons/Plus';
import { OnboardingPlanMilestoneItemStatus } from '@graphql/types';

import { Task } from './Task';
import { TaskDatum } from '../../../../types';

interface TasksProps {
  formId: string;
  milestoneDueDate: string;
}

export const Tasks = memo(
  ({ formId, milestoneDueDate }: TasksProps) => {
    const { getInputProps } = useField('items', formId);
    const { value, onChange } = getInputProps();
    const shouldFocusRef = useRef<number | null>(null);

    const handleAddTask = () => {
      const nextItems = produce<TaskDatum[]>(value, (draft) => {
        const dueDate = setHours(new Date(milestoneDueDate), 0);
        const nowDate = setHours(new Date(), 0);
        const isPastDue = dueDate < nowDate;

        draft.push({
          text: '',
          uuid: crypto.randomUUID(),
          updatedAt: new Date().toISOString(),
          status: isPastDue
            ? OnboardingPlanMilestoneItemStatus.NotDoneLate
            : OnboardingPlanMilestoneItemStatus.NotDone,
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
            index={idx}
            key={m.uuid}
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
