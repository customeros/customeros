import { useRef, useState } from 'react';

import { Store } from '@store/store.ts';

import { Input } from '@ui/form/Input';
import { BankAccount } from '@graphql/types';

export const BankNameInput = ({ account }: { account: Store<BankAccount> }) => {
  const nameRef = useRef<HTMLInputElement | null>(null);
  const [name, setName] = useState(() => account?.value?.bankName ?? '');

  return (
    <Input
      ref={nameRef}
      autoComplete='off'
      value={name ?? ''}
      variant={'unstyled'}
      aria-label='Bank Name'
      placeholder='Bank name'
      className='text-md font-semibold'
      onFocus={(e) => e?.target?.select()}
      onChange={(e) => setName(e.target.value)}
      onBlur={(e) =>
        account?.update((acc) => {
          acc.bankName = e.target.value;

          return acc;
        })
      }
    />
  );
};
