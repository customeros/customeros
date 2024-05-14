import { useField } from 'react-inverted-form';
import { memo, useRef, useState } from 'react';

import { Button } from '@ui/form/Button/Button';
import { ResizableInput } from '@ui/form/Input/ResizableInput';

export interface CommittedPeriodInputProps {
  formId: string;
}

export function getCommittedPeriodLabel(months: string | number) {
  if (`${months}` === '1') {
    return 'Monthly';
  }
  if (`${months}` === '3') {
    return 'Quarterly';
  }

  if (`${months}` === '12') {
    return 'Annual';
  }

  return `${months}-month`;
}

export const CommittedPeriodInput = memo(
  ({ formId }: CommittedPeriodInputProps) => {
    const [isFocused, setIsFocused] = useState(false);
    const inputRef = useRef<HTMLInputElement>(null);
    const { getInputProps } = useField('committedPeriodInMonths', formId);
    const { onChange, value } = getInputProps();
    const handleFocus = () => {
      setIsFocused(true);

      setTimeout(() => {
        inputRef.current?.focus();
      }, 500);
    };

    const handleBlur = () => {
      setIsFocused(false);
    };

    const committedPeriodLabel = getCommittedPeriodLabel(value);

    return (
      <>
        {isFocused && (
          <div className='flex mr-1 items-baseline'>
            <ResizableInput
              value={value}
              ref={inputRef}
              onChange={onChange}
              onFocus={handleFocus}
              onBlur={handleBlur}
              size='xs'
              className='text-base min-w-2.5 min-h-0 max-h-4'
            />
            <span> -month</span>
          </div>
        )}

        {!isFocused && (
          <Button
            variant='ghost'
            size='sm'
            className='font-normal text-base p-0 mr-1 relative text-gray-500 hover:bg-transparent focus:bg-transparent underline color-gray-500'
            onClick={handleFocus}
          >
            {committedPeriodLabel}
          </Button>
        )}
      </>
    );
  },
);
