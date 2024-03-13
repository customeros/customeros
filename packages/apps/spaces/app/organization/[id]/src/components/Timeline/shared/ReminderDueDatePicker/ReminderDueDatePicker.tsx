import { useField } from 'react-inverted-form';

import set from 'date-fns/set';
import getHours from 'date-fns/getHours';
import getMinutes from 'date-fns/getMinutes';

import { Portal } from '@ui/utils/';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { DateTimeUtils } from '@spaces/utils/date';
import { Input, InputProps } from '@ui/form/Input';
import { InlineDatePicker } from '@ui/form/DatePicker';
import {
  Popover,
  PopoverBody,
  PopoverFooter,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover';

interface DueDatePickerProps {
  name: string;
  formId: string;
}

export const ReminderDueDatePicker = ({ name, formId }: DueDatePickerProps) => {
  const { getInputProps } = useField(name, formId);
  const { onChange, ...inputProps } = getInputProps();

  const time = (() => {
    const dateStr = inputProps.value;
    const date = dateStr ? new Date(dateStr) : new Date();

    const hours = (() => {
      const h = String(getHours(date));

      return h.length === 1 ? `0${h}` : h;
    })();
    const minutes = (() => {
      const h = String(getMinutes(date));

      return h.length === 1 ? `0${h}` : h;
    })();

    return `${hours}:${minutes}`;
  })();

  const handleChange = (date: Date | null) => {
    if (!date) return;
    const [hours, minutes] = time.split(':').map(Number);
    const _date = set(date, { hours, minutes });

    onChange(_date.toISOString());
  };

  return (
    <Flex justify='flex-start' align='center'>
      <Popover placement='top-start'>
        <PopoverTrigger>
          <Text whiteSpace='pre' pb='1px'>{`${DateTimeUtils.format(
            inputProps.value,
            DateTimeUtils.date,
          )} • `}</Text>
        </PopoverTrigger>
        <Portal>
          <PopoverContent>
            <PopoverBody>
              <InlineDatePicker {...inputProps} onChange={handleChange} />
            </PopoverBody>
            <PopoverFooter>PLM</PopoverFooter>
          </PopoverContent>
        </Portal>
      </Popover>
      <TimeInput
        value={time}
        onChange={(v) => {
          const [hours, minutes] = v.split(':').map(Number);
          const date = set(new Date(inputProps.value), { hours, minutes });

          onChange(date.toISOString());
        }}
      />
    </Flex>
  );
};

interface TimeInputProps extends Omit<InputProps, 'value' | 'onChange'> {
  value?: string;
  onChange?: (value: string) => void;
}

const TimeInput = ({ onChange, value, ...rest }: TimeInputProps) => {
  return (
    <Input
      p='0'
      type='time'
      list='hidden'
      value={value}
      lineHeight='1'
      h='min-content'
      w='fit-content'
      onChange={(e) => {
        const val = e.target.value;
        onChange?.(val);
      }}
      sx={{
        '&::-webkit-calendar-picker-indicator': {
          display: 'none',
        },
      }}
      {...rest}
    />
  );
};
