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
import { Tooltip } from '@ui/overlay/Tooltip';
import { IconButton } from '@ui/form/IconButton';
import { DateTimeUtils } from '@spaces/utils/date';
import { SkipForward } from '@ui/media/icons/SkipForward';
import { OnboardingPlanMilestoneItemStatus } from '@graphql/types';

import { TaskDatum } from '../../../../types';
import { StatusCheckbox } from '../StatusCheckbox';

interface TaskProps {
  index: number;
  formId: string;
  isFocused?: boolean;
  shouldFocus?: boolean;
  onInputBlur?: () => void;
  onInputFocus?: () => void;
}

export const Task = memo(
  ({
    index,
    formId,
    isFocused,
    shouldFocus,
    onInputFocus,
    onInputBlur,
  }: TaskProps) => {
    const ref = useRef<HTMLInputElement>(null);

    const [showSkip, setShowSkip] = useState(false);

    const { getInputProps } = useField('items', formId);
    const { value, onChange, onBlur } = getInputProps();
    const itemValue = value?.[index] as TaskDatum;

    const milestoneDueDate = useField('dueDate', formId).getInputProps()
      ?.value as string;

    const taskStatus =
      itemValue?.status ?? OnboardingPlanMilestoneItemStatus.NotDone;
    const taskUpdatedAt = new Date(itemValue?.updatedAt).valueOf();
    const taskUpdatedAtDate = DateTimeUtils.format(
      itemValue?.updatedAt,
      DateTimeUtils.dateWithShortYear,
    );
    const milestoneDueAt = new Date(milestoneDueDate).valueOf();

    const isDone = [
      OnboardingPlanMilestoneItemStatus.Done,
      OnboardingPlanMilestoneItemStatus.DoneLate,
    ].includes(taskStatus);
    const isSkipped = [
      OnboardingPlanMilestoneItemStatus.Skipped,
      OnboardingPlanMilestoneItemStatus.SkippedLate,
    ].includes(taskStatus);

    const colorScheme = (() => {
      switch (taskStatus) {
        case OnboardingPlanMilestoneItemStatus.NotDone:
        case OnboardingPlanMilestoneItemStatus.Skipped:
        case OnboardingPlanMilestoneItemStatus.SkippedLate:
          return 'gray';
        case OnboardingPlanMilestoneItemStatus.Done:
          return 'success';
        case OnboardingPlanMilestoneItemStatus.DoneLate:
        case OnboardingPlanMilestoneItemStatus.NotDoneLate:
          return 'warning';
        default:
          return 'gray';
      }
    })();

    const handleChange = useCallback(
      (e: ChangeEvent<HTMLInputElement>) => {
        if (!itemValue) return;
        const isChecked = e.target.checked;

        const updatedItem = produce<TaskDatum>(itemValue, (draft) => {
          const isLate = new Date().valueOf() > milestoneDueAt;

          if (isChecked) {
            draft.status = isLate
              ? OnboardingPlanMilestoneItemStatus.DoneLate
              : OnboardingPlanMilestoneItemStatus.Done;
          } else {
            draft.status = isLate
              ? OnboardingPlanMilestoneItemStatus.NotDoneLate
              : OnboardingPlanMilestoneItemStatus.NotDone;
          }

          draft.updatedAt = new Date().toISOString();
        });

        const next = (value as TaskDatum[]).map((v, i) =>
          i === index ? updatedItem : v,
        );

        onChange?.(next);
      },
      [onChange, index, value, taskUpdatedAt, milestoneDueAt],
    );

    const handleInputChange = useCallback(
      (e: ChangeEvent<HTMLInputElement>) => {
        const nextItems = produce<TaskDatum[]>(value, (draft) => {
          const item = draft?.[index];
          if (!item) return;

          item.text = e.target.value;
        });

        onChange?.(nextItems);
      },
      [onChange, index, value],
    );

    const handleInputBlur = useCallback(
      (e: ChangeEvent<HTMLInputElement>) => {
        onInputBlur?.();

        const nextItems = produce<TaskDatum[]>(value, (draft) => {
          const item = draft?.[index];
          if (!item) return;

          item.text = e.target.value;
        });

        onBlur?.(nextItems);
      },
      [onBlur, onChange, index, value],
    );

    const handleRemove = () => {
      const nextItems = produce<TaskDatum[]>(value, (draft) => {
        draft.splice(index, 1);
      });

      onChange?.(nextItems);
    };

    const handleSkip = useCallback(() => {
      const nextItems = produce<TaskDatum[]>(value, (draft) => {
        const item = draft?.[index];
        if (!item) return;

        const isPastDueDate = taskUpdatedAt > milestoneDueAt;

        item.status = (() => {
          switch (item.status) {
            case OnboardingPlanMilestoneItemStatus.NotDone:
              return isPastDueDate
                ? OnboardingPlanMilestoneItemStatus.SkippedLate
                : OnboardingPlanMilestoneItemStatus.Skipped;
            case OnboardingPlanMilestoneItemStatus.Skipped:
            case OnboardingPlanMilestoneItemStatus.SkippedLate:
              return isPastDueDate
                ? OnboardingPlanMilestoneItemStatus.NotDoneLate
                : OnboardingPlanMilestoneItemStatus.NotDone;
            default:
              return item.status;
          }
        })();

        item.updatedAt = new Date().toISOString();
      });

      onChange?.(nextItems);
    }, [onChange, index, value, taskUpdatedAt, milestoneDueAt]);

    const handleAdd = () => {
      const nextItems = produce<TaskDatum[]>(value, (draft) => {
        const isPastDueDate = taskUpdatedAt > milestoneDueAt;

        draft.splice(index + 1, 0, {
          text: '',
          updatedAt: new Date().toISOString(),
          status: isPastDueDate
            ? OnboardingPlanMilestoneItemStatus.NotDoneLate
            : OnboardingPlanMilestoneItemStatus.NotDone,
        });
      });

      onChange?.(nextItems);
    };

    useKey('Enter', () => {
      if (isFocused) {
        handleAdd();
      }
    });
    useKey('Backspace', () => {
      if (isFocused && !ref.current?.value) {
        handleRemove();
      }
    });

    useEffect(() => {
      if (shouldFocus && !isFocused) {
        ref?.current?.focus();
      }
    }, [shouldFocus, isFocused]);

    return (
      <Flex
        w='full'
        onMouseEnter={() => (!isSkipped ? setShowSkip(true) : undefined)}
        onMouseLeave={() => (!isSkipped ? setShowSkip(false) : undefined)}
      >
        <StatusCheckbox
          mr='2'
          size='md'
          isChecked={isDone}
          onChange={handleChange}
          colorScheme={colorScheme}
        />
        <Input
          w='full'
          ref={ref}
          fontSize='sm'
          variant='unstyled'
          borderRadius='unset'
          placeholder='Task name'
          onBlur={handleInputBlur}
          value={value?.[index]?.text}
          onChange={handleInputChange}
          onFocus={() => {
            onInputFocus?.();
          }}
          fontStyle={isSkipped ? 'italic' : 'normal'}
        />
        {isDone && (
          <Text fontSize='sm' color='gray.500' whiteSpace='nowrap'>
            {taskUpdatedAtDate}
          </Text>
        )}
        {!isDone && (
          <Tooltip label={isSkipped ? 'Skipped' : 'Skip this'}>
            <IconButton
              size='xs'
              variant='ghost'
              onClick={handleSkip}
              opacity={showSkip || isSkipped ? 1 : 0}
              aria-label='Skip Milestone Task'
              icon={<SkipForward color='gray.400' />}
            />
          </Tooltip>
        )}
      </Flex>
    );
  },
);
