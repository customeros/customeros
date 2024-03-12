import { useState } from 'react';
import { useField } from 'react-inverted-form';

import getHours from 'date-fns/getHours';
import getMinutes from 'date-fns/getMinutes';

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

  const defaultTime = (() => {
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

  const [time, setTime] = useState(defaultTime);

  return (
    <Flex justify='flex-start' align='center'>
      <Popover>
        <PopoverTrigger>
          <Text whiteSpace='pre' pb='1px'>{`${DateTimeUtils.format(
            inputProps.value,
            DateTimeUtils.date,
          )} • `}</Text>
        </PopoverTrigger>
        <PopoverContent>
          <PopoverBody>
            <InlineDatePicker {...inputProps} onChange={onChange} />
          </PopoverBody>
          <PopoverFooter>PLM</PopoverFooter>
        </PopoverContent>
      </Popover>
      <TimeInput value={time} onChange={(v) => setTime(v)} />
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
