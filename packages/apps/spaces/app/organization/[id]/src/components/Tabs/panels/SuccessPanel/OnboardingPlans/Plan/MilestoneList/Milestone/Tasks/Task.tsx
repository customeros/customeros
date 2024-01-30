import { useField } from 'react-inverted-form';
import { memo, useState, useCallback, ChangeEvent } from 'react';

import { produce } from 'immer';

import { Flex } from '@ui/layout/Flex';
import { Input } from '@ui/form/Input';
import { Text } from '@ui/typography/Text';
import { Tooltip } from '@ui/overlay/Tooltip';
import { IconButton } from '@ui/form/IconButton';
import { SkipForward } from '@ui/media/icons/SkipForward';

import { TaskDatum } from '../../../../types';
import { StatusCheckbox } from '../StatusCheckbox';

interface TaskProps {
  index: number;
  formId: string;
  defaultValue: TaskDatum;
}

export const Task = memo(({ index, formId, defaultValue }: TaskProps) => {
  const [showSkip, setShowSkip] = useState(false);
  const { getInputProps } = useField('items', formId);
  const { value, onChange } = getInputProps();
  const itemValue = value?.[index] as TaskDatum;

  const milestoneDueDate = useField('dueDate', formId).getInputProps()
    ?.value as string;

  const taskStatus = itemValue?.status ?? 'NOT_DONE';
  const taskUpdatedAt = new Date(itemValue?.updatedAt).valueOf();
  const taskUpdatedAtDate = new Date(itemValue?.updatedAt).toLocaleDateString();
  const milestoneDueAt = new Date(milestoneDueDate).valueOf();

  const colorScheme = (() => {
    const isLate = taskUpdatedAt > milestoneDueAt;

    if (taskStatus === 'SKIPPED') return 'gray';

    if (['DONE'].includes(taskStatus)) {
      return isLate ? 'warning' : 'success';
    }

    return isLate ? 'warning' : 'gray';
  })();

  const handleChange = useCallback(
    (e: ChangeEvent<HTMLInputElement>) => {
      const nextItems = produce<TaskDatum[]>(value, (draft) => {
        const item = draft[index];
        if (!item) return;

        item.status = e.target.checked ? 'DONE' : 'NOT_DONE';
        item.updatedAt = new Date().toISOString();
      });

      onChange?.(nextItems);
    },
    [onChange, index],
  );

  const handleSkip = useCallback(() => {
    const nextItems = produce<TaskDatum[]>(value, (draft) => {
      const item = draft[index];
      if (!item) return;

      item.status = item.status !== 'SKIPPED' ? 'SKIPPED' : 'NOT_DONE';
      item.updatedAt = new Date().toISOString();
    });

    onChange?.(nextItems);
  }, [onChange, index]);

  // const handleBlur = useCallback(
  //   (e: ChangeEvent<HTMLInputElement>) => {
  //     setIsFocused(false);
  //     setShowRemove(false);

  //     const nextItems = produce<string[]>(value, (draft) => {
  //       const val = e.target.value;
  //       if (!val.length) {
  //         draft.splice(index, 1);

  //         return;
  //       }

  //       draft[index] = e.target.value;
  //     });

  //     onBlur?.(nextItems);
  //   },
  //   [onBlur, index],
  // );

  return (
    <Flex
      w='full'
      onMouseEnter={() => setShowSkip(true)}
      onMouseLeave={() => setShowSkip(false)}
    >
      <StatusCheckbox
        mr='2'
        size='md'
        onChange={handleChange}
        colorScheme={colorScheme}
        isChecked={value[index].status === 'DONE'}
      />
      <Input
        w='full'
        fontSize='sm'
        variant='unstyled'
        borderRadius='unset'
        value={value[index].text}
        placeholder='Task name'
        fontStyle={taskStatus === 'SKIPPED' ? 'italic' : 'normal'}
      />
      {taskStatus === 'DONE' && <Text>{taskUpdatedAtDate}</Text>}
      {taskStatus !== 'DONE' && (
        <Tooltip label={taskStatus === 'SKIPPED' ? 'Skipped' : 'Skip this'}>
          <IconButton
            size='xs'
            variant='ghost'
            onClick={handleSkip}
            opacity={showSkip ? 1 : 0}
            aria-label='Skip Milestone Task'
            icon={<SkipForward color='gray.400' />}
          />
        </Tooltip>
      )}
    </Flex>
  );
});
