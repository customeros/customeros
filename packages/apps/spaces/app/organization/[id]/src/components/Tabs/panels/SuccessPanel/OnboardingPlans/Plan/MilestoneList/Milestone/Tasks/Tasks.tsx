import { useField } from 'react-inverted-form';

// import { produce } from 'immer';

// import { Button } from '@ui/form/Button';
import { VStack } from '@ui/layout/Stack';
// import { Plus } from '@ui/media/icons/Plus';

import { Task } from './Task';
import { TaskDatum } from '../../../../types';

interface TasksProps {
  formId: string;
  defaultValue: TaskDatum[];
}

export const Tasks = ({ formId, defaultValue }: TasksProps) => {
  const { getInputProps } = useField('items', formId);
  const { value } = getInputProps();

  // const handleAddTask = () => {
  //   const nextItems = produce<string[]>(value, (draft) => {
  //     draft.push('');
  //   });

  //   onChange?.(nextItems);
  //   onBlur?.(nextItems);
  // };

  return (
    <VStack align='flex-start' spacing='1' pl='6'>
      {(value as TaskDatum[])?.map((item, idx) => (
        <Task key={idx} index={idx} formId={formId} />
      ))}
      {/* <Button
        ml='-12px'
        size='sm'
        variant='ghost'
        color='gray.500'
        fontWeight='normal'
        onClick={handleAddTask}
        leftIcon={<Plus color='gray.400' />}
      >
        Add task
      </Button> */}
    </VStack>
  );
};
