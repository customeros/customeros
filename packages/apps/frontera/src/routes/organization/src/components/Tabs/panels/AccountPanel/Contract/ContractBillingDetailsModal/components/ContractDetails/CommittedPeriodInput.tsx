import { useRef, useState, ChangeEvent } from 'react';

import { observer } from 'mobx-react-lite';
import { ContractStore } from '@store/Contracts/Contract.store.ts';

import { useStore } from '@shared/hooks/useStore';
import { Button } from '@ui/form/Button/Button.tsx';
import { ResizableInput } from '@ui/form/Input/ResizableInput.tsx';

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
    const contractStore = store.contracts.value.get(
      contractId,
    ) as ContractStore;

    const handleFocus = () => {
      setIsFocused(true);

      setTimeout(() => {
        inputRef.current?.focus();
      }, 500);
    };

    const handleBlur = (e: ChangeEvent<HTMLInputElement>) => {
      setIsFocused(false);

      if (!e.target.value) {
        contractStore?.updateTemp((contract) => ({
          ...contract,
          committedPeriodInMonths: 1,
        }));
      }
    };

    const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
      contractStore?.updateTemp((contract) => ({
        ...contract,
        committedPeriodInMonths:
          e.target.value && Number(e.target.value) >= 9999
            ? contract.committedPeriodInMonths
            : e.target.value,
      }));
    };

    const committedPeriodLabel = getCommittedPeriodLabel(
      contractStore?.tempValue?.committedPeriodInMonths,
    );

    return (
      <>
        {isFocused && (
          <div className='flex mr-1 items-baseline'>
            <ResizableInput
              min={1}
              size='xs'
              max={999}
              type='number'
              ref={inputRef}
              onBlur={handleBlur}
              onFocus={handleFocus}
              onChange={handleChange}
              className='text-base min-w-2.5 min-h-0 max-h-4'
              value={contractStore?.tempValue?.committedPeriodInMonths ?? 1}
              defaultValue={
                contractStore?.tempValue?.committedPeriodInMonths ?? 1
              }
            />
            <span> -month</span>
          </div>
        )}

        {!isFocused && (
          <Button
            size='sm'
            variant='ghost'
            onClick={handleFocus}
            className='font-normal text-base p-0 mr-1 relative text-gray-500 hover:bg-transparent focus:bg-transparent underline color-gray-500'
          >
            {committedPeriodLabel}
          </Button>
        )}
      </>
    );
  },
);
