import { useField } from 'react-inverted-form';

import { Portal } from '@ui/utils';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { InlineDatePicker } from '@ui/form/DatePicker';
import { OnboardingPlanMilestoneStatus } from '@graphql/types';
import {
  Popover,
  PopoverBody,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover';

import { getMilestoneDueDate } from '../../../utils';

interface MilestoneDueDateProps {
  formId: string;
  isDone?: boolean;
  isOpen?: boolean;
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
}: MilestoneDueDateProps) => {
  const { getInputProps } = useField('dueDate', formId);
  const { id, name, onChange, value } = getInputProps();

  return (
    <Flex w='full'>
      <Popover
        isOpen={isOpen}
        onClose={onClose}
        onOpen={onOpen}
        closeOnEsc
        closeOnBlur
        matchWidth
        isLazy
      >
        <PopoverTrigger>
          <Text
            as='label'
            fontSize='sm'
            color='gray.500'
            whiteSpace='nowrap'
            cursor='pointer'
            pointerEvents={isOpen ? 'none' : 'auto'}
          >
            {getMilestoneDueDate(value, status, isDone)}
          </Text>
        </PopoverTrigger>

        <Portal>
          <PopoverContent width='fit-content'>
            <PopoverBody w='fit-content'>
              <InlineDatePicker
                id={id}
                name={name}
                selected={new Date(value)}
                onChange={(date) => {
                  onChange(date?.toISOString());
                }}
              />
            </PopoverBody>
          </PopoverContent>
        </Portal>
      </Popover>
    </Flex>
  );
};
