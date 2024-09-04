import { useRef } from 'react';

import { Store } from '@store/store.ts';

import { Input } from '@ui/form/Input';
import { BankAccount } from '@graphql/types';

export const BankNameInput = ({ account }: { account: Store<BankAccount> }) => {
  const nameRef = useRef<HTMLInputElement | null>(null);

  return (
    <Input
      ref={nameRef}
      autoComplete='off'
      variant={'unstyled'}
      aria-label='Bank Name'
      placeholder='Bank name'
      className='text-md font-semibold'
      onFocus={(e) => e?.target?.select()}
      value={account?.value?.bankName ?? ''}
      onChange={(e) =>
        account?.update((acc) => {
          acc.bankName = e.target.value;

          return acc;
        })
      }
    />
  );
};
