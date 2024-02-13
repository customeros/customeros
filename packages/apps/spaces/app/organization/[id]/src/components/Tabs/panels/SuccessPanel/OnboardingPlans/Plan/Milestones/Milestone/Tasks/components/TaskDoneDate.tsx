import { Portal } from '@ui/utils';
import { Text } from '@ui/typography/Text';
import { DateTimeUtils } from '@spaces/utils/date';
import { InlineDatePicker } from '@ui/form/DatePicker';
import {
  Popover,
  PopoverBody,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover';

interface TaskDoneDateProps {
  value?: string;
  isOpen?: boolean;
  onOpen?: () => void;
  onClose?: () => void;
  onChange: (value: string) => void;
}

export const TaskDoneDate = ({
  isOpen,
  onOpen,
  onClose,
  onChange,
  value = new Date().toISOString(),
}: TaskDoneDateProps) => {
  const taskUpdatedAtDate = DateTimeUtils.format(
    value,
    DateTimeUtils.dateWithShortYear,
  );

  return (
    <Popover
      isOpen={isOpen}
      onClose={onClose}
      onOpen={onOpen}
      placement='bottom-end'
      closeOnEsc
      closeOnBlur
      isLazy
    >
      <PopoverTrigger>
        <Text
          as='label'
          fontSize='sm'
          whiteSpace='nowrap'
          cursor='pointer'
          pointerEvents={isOpen ? 'none' : 'auto'}
          color={isOpen ? 'primary.500' : 'gray.500'}
        >
          {taskUpdatedAtDate}
        </Text>
      </PopoverTrigger>

      <Portal>
        <PopoverContent width='fit-content'>
          <PopoverBody w='fit-content'>
            <InlineDatePicker
              selected={new Date(value)}
              maxDate={new Date()}
              onChange={(date) => {
                date && onChange(date?.toISOString());
              }}
            />
          </PopoverBody>
        </PopoverContent>
      </Portal>
    </Popover>
  );
};
