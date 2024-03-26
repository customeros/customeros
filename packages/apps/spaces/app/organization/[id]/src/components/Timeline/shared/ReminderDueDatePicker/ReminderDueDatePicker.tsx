import { useRef } from 'react';
import { useField } from 'react-inverted-form';

import set from 'date-fns/set';
import addDays from 'date-fns/addDays';
import getHours from 'date-fns/getHours';
import getMinutes from 'date-fns/getMinutes';

import { Portal } from '@ui/utils/';
import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
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
  const containerRef = useRef<HTMLDivElement>(null);

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
    const _date = set(date, { hours, minutes, seconds: 0, milliseconds: 0 });

    onChange(_date.toISOString());
  };

  const handleClickTomorrow = () => {
    const date = set(addDays(new Date(), 1), {
      hours: 9,
      minutes: 0,
      seconds: 0,
      milliseconds: 0,
    });
    onChange(date.toISOString());
  };

  return (
    <Flex ref={containerRef} justify='flex-start' align='center'>
      <Popover placement='top-start' matchWidth>
        {({ isOpen, onClose }) => (
          <>
            <PopoverTrigger>
              <Text
                cursor='pointer'
                whiteSpace='pre'
                pb='1px'
                fontSize='sm'
                borderTop='1px solid transparent'
                color={isOpen ? 'gray.700' : 'gray.500'}
                _hover={{ color: 'gray.700' }}
              >{`${DateTimeUtils.format(
                inputProps.value,
                DateTimeUtils.date,
              )} • `}</Text>
            </PopoverTrigger>
            <Portal>
              <PopoverContent w='fit-content'>
                <PopoverBody w='fit-content'>
                  <InlineDatePicker
                    {...inputProps}
                    onChange={(date) => {
                      handleChange(date);
                      onClose();
                    }}
                    minDate={new Date()}
                  />
                </PopoverBody>
                <PopoverFooter
                  display='flex'
                  alignItems='center'
                  justifyContent='right'
                  px='6'
                >
                  <Button
                    variant='outline'
                    borderRadius='full'
                    onClick={() => {
                      handleClickTomorrow();
                      onClose();
                    }}
                  >
                    Tomorrow
                  </Button>
                </PopoverFooter>
              </PopoverContent>
            </Portal>
          </>
        )}
      </Popover>
      <TimeInput
        color='gray.500'
        value={time}
        fontSize='sm'
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
      _focus={{
        color: 'gray.700',
        boxShadow: 'unset',
        borderColor: 'primary.500',
      }}
      cursor='text'
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
