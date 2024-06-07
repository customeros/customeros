import { useRef, useState } from 'react';

import { observer } from 'mobx-react-lite';

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { ResizableInput } from '@ui/form/Input/ResizableInput';

export interface CommittedPeriodInputProps {
  contractId: string;
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

export const CommittedPeriodInput = observer(
  ({ contractId }: CommittedPeriodInputProps) => {
    const [isFocused, setIsFocused] = useState(false);
    const inputRef = useRef<HTMLInputElement>(null);
    const store = useStore();
    const contractStore = store.contracts.value.get(contractId);

    const handleFocus = () => {
      setIsFocused(true);

      setTimeout(() => {
        inputRef.current?.focus();
      }, 500);
    };

    const handleBlur = () => {
      setIsFocused(false);
    };

    const committedPeriodLabel = getCommittedPeriodLabel(
      contractStore?.value?.committedPeriodInMonths,
    );

    return (
      <>
        {isFocused && (
          <div className='flex mr-1 items-baseline'>
            <ResizableInput
              ref={inputRef}
              defaultValue={contractStore?.value?.committedPeriodInMonths ?? 1}
              value={contractStore?.value?.committedPeriodInMonths ?? 1}
              onChange={(e) =>
                contractStore?.update(
                  (contract) => ({
                    ...contract,
                    committedPeriodInMonths: e.target.value,
                  }),
                  { mutate: false },
                )
              }
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
