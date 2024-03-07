import { useField } from 'react-inverted-form';

import { Portal } from '@ui/utils';
import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { DateTimeUtils } from '@spaces/utils/date';
import { InlineDatePicker } from '@ui/form/DatePicker';
import { OnboardingPlanMilestoneStatus } from '@graphql/types';
import {
  Popover,
  PopoverBody,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover';

import { getMilestoneDueDate, getMilestoneDoneDate } from '../../../utils';

interface MilestoneDueDateProps {
  formId: string;
  isDone?: boolean;
  isOpen?: boolean;
  minDate?: string;
  onOpen?: () => void;
  onClose?: () => void;
  status: OnboardingPlanMilestoneStatus;
}

export const MilestoneDueDate = ({
  status,
  isDone,
  formId,
  isOpen,
  onOpen,
  onClose,
  minDate,
}: MilestoneDueDateProps) => {
  const dueDate = useField('dueDate', formId).getInputProps();
  const statusDetails = useField('statusDetails', formId).getInputProps();
  const date = isDone ? statusDetails?.value?.updatedAt : dueDate?.value;

  const handleChange = (date: Date | null) => {
    if (!date) return;

    if (isDone) {
      statusDetails?.onChange({
        ...statusDetails.value,
        updatedAt: DateTimeUtils.toISOMidnight(date),
      });
    } else {
      dueDate?.onChange(DateTimeUtils.toISOMidnight(date));
    }
  };

  const displayDate = isDone
    ? getMilestoneDoneDate(date, status)
    : getMilestoneDueDate(date, status, isDone);

  const hoverDate = DateTimeUtils.format(date, DateTimeUtils.dateWithShortYear);

  return (
    <Flex w='full'>
      <Popover
        isLazy
        matchWidth
        closeOnEsc
        closeOnBlur
        isOpen={isOpen}
        onOpen={onOpen}
        onClose={onClose}
      >
        <PopoverTrigger>
          <Text
            as='label'
            fontSize='sm'
            cursor='pointer'
            whiteSpace='nowrap'
            pointerEvents={isOpen || isDone ? 'none' : 'auto'}
            color={isOpen ? 'primary.500' : 'gray.500'}
            _hover={
              !isDone
                ? {
                    '& #display-date': {
                      display: 'none',
                    },
                    '& #hover-date': {
                      display: 'block',
                    },
                  }
                : {}
            }
          >
            <Box
              as='span'
              id='display-date'
              display={isOpen && !isDone ? 'none' : 'block'}
            >
              {displayDate}
            </Box>
            <Box
              as='span'
              id='hover-date'
              display={isOpen && !isDone ? 'block' : 'none'}
            >
              {`Due on ${hoverDate}`}
            </Box>
          </Text>
        </PopoverTrigger>

        <Portal>
          <PopoverContent width='fit-content'>
            <PopoverBody w='fit-content'>
              <InlineDatePicker
                onChange={handleChange}
                selected={new Date(date)}
                minDate={minDate ? new Date(minDate) : undefined}
              />
            </PopoverBody>
          </PopoverContent>
        </Portal>
      </Popover>
    </Flex>
  );
};
