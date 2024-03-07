import { useField } from 'react-inverted-form';
import {
  memo,
  useRef,
  useState,
  useEffect,
  useCallback,
  ChangeEvent,
} from 'react';

import { useKey } from 'rooks';
import { produce } from 'immer';

import { Flex } from '@ui/layout/Flex';
import { Input } from '@ui/form/Input';
import { Text } from '@ui/typography/Text';
import { Checkbox } from '@ui/form/Checkbox';
import { Delete } from '@ui/media/icons/Delete';
import { IconButton } from '@ui/form/IconButton';

interface TaskProps {
  index: number;
  formId: string;
  isLast?: boolean;
  defaultValue: string;
  isActiveItem?: boolean;
}

export const Task = memo(
  ({ index, formId, isLast, defaultValue, isActiveItem }: TaskProps) => {
    const ref = useRef<HTMLInputElement>(null);
    const [showRemove, setShowRemove] = useState(false);
    const [isFocused, setIsFocused] = useState(false);

    const { getInputProps } = useField('items', formId);
    const { value, onChange, onBlur } = getInputProps();

    const handleAddTask = () => {
      const nextItems = produce<string[]>(value, (draft) => {
        draft.push('');
      });

      onChange?.(nextItems);
      onBlur?.(nextItems);
    };

    const handleChange = useCallback(
      (e: ChangeEvent<HTMLInputElement>) => {
        const nextItems = produce<string[]>(value, (draft) => {
          draft[index] = e.target.value;
        });

        onChange?.(nextItems);
      },
      [onChange, index],
    );

    const handleBlur = useCallback(
      (e: ChangeEvent<HTMLInputElement>) => {
        setIsFocused(false);
        setShowRemove(false);

        const nextItems = produce<string[]>(value, (draft) => {
          const val = e.target.value;
          if (!val.length) {
            draft.splice(index, 1);

            return;
          }

          draft[index] = e.target.value;
        });

        onBlur?.(nextItems);
      },
      [onBlur, index],
    );

    const handleRemove = useCallback(() => {
      const nextItems = produce<string[]>(value, (draft) => {
        draft.splice(index, 1);
      });
      onChange?.(nextItems);
    }, [onChange, index]);

    useEffect(() => {
      if (isLast && !value?.[index]?.length) {
        ref?.current?.focus();
        setIsFocused(true);
        setShowRemove(true);
      }
    }, [isLast]);

    useKey(['Enter'], handleAddTask, { when: isFocused && isLast });
    useKey(
      ['Escape'],
      () => {
        setIsFocused(false);
        setShowRemove(false);
        ref.current?.blur();
      },
      { when: isFocused && isLast },
    );

    return (
      <Flex
        w='full'
        onMouseEnter={() => setShowRemove(true)}
        onMouseLeave={() => (isFocused ? undefined : setShowRemove(false))}
      >
        <Checkbox size='md' disabled readOnly mr='2' />
        {isActiveItem ? (
          <Text fontSize='sm' w='full'>
            {defaultValue}
          </Text>
        ) : (
          <Input
            w='full'
            ref={ref}
            fontSize='sm'
            variant='unstyled'
            onBlur={handleBlur}
            borderRadius='unset'
            value={value[index]}
            onChange={handleChange}
            placeholder='Task name'
            onFocus={() => setIsFocused(true)}
          />
        )}
        <IconButton
          size='xs'
          variant='ghost'
          onClick={handleRemove}
          opacity={showRemove ? 1 : 0}
          aria-label='Remove Milestone Task'
          icon={<Delete color='gray.400' />}
        />
      </Flex>
    );
  },
);
