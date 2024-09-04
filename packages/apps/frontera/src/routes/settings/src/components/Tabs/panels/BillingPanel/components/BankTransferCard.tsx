import { useState } from 'react';

import { Store } from '@store/store.ts';
import { observer } from 'mobx-react-lite';
import { BankNameInput } from '@settings/components/Tabs/panels/BillingPanel/components/BankNameInput';
import { useBankTransferSelectionContext } from '@settings/components/Tabs/panels/BillingPanel/context/BankTransferSelectionContext';

import { Input } from '@ui/form/Input';
import { Currency, BankAccount } from '@graphql/types';
import { Textarea } from '@ui/form/Textarea/Textarea.tsx';
import { MaskedInput } from '@ui/form/Input/MaskedInput.tsx';
import { Card, CardHeader, CardContent } from '@ui/presentation/Card/Card';

import { BankTransferMenu } from './BankTransferMenu';
import { BankTransferCurrencySelect } from './BankTransferCurrencySelect';

export const BankTransferCard = observer(
  ({
    account,
    existingCurrencies,
  }: {
    account: Store<BankAccount>;
    existingCurrencies: Array<string>;
  }) => {
    const [bankAccount, setBankAccount] = useState('');
    const [value, setValue] = useState(account.value.iban ?? '');
    const [sortCode, setSortCode] = useState('');
    const { setFocusAccount, setHoverAccount } =
      useBankTransferSelectionContext();

    return (
      <>
        <Card
          onMouseLeave={() => setHoverAccount(null)}
          className='py-2 px-4 rounded-lg border-[1px] border-gray-200'
          onMouseEnter={() => setHoverAccount(account.value?.metadata?.id)}
        >
          <CardHeader className='p-0 pb-1 flex justify-between'>
            <BankNameInput account={account} />

            <div className='flex'>
              <BankTransferCurrencySelect
                onChange={() => null}
                currency={account.value.currency}
                existingCurrencies={existingCurrencies}
              />

              <BankTransferMenu id={account?.value?.metadata?.id} />
            </div>
          </CardHeader>
          <CardContent className='p-0 gap-2'>
            {account?.value?.currency !== 'USD' &&
              account?.value?.currency !== 'GBP' && (
                <>
                  <MaskedInput
                    value={value}
                    variant='unstyled'
                    placeholder='Enter IBAN'
                    prepare={(str) => str.toUpperCase()}
                    mask='SS00 AAAA 0000 0000 0000 9999 9999 9999 99'
                    onAccept={(val) => {
                      setValue(val.toUpperCase());

                      account.update((acc) => {
                        acc.iban = val.toUpperCase();

                        return acc;
                      });
                    }}
                    definitions={{
                      S: {
                        mask: /[A-Za-z]/,
                      },
                      A: {
                        mask: /[A-Za-z0-9]/,
                      },
                      '0': {
                        mask: /[0-9]/,
                      },
                      '9': {
                        mask: /[0-9]/,
                      },
                    }}
                  />
                  <Input
                    name='bic'
                    autoComplete='off'
                    variant='unstyled'
                    aria-label='BIC/Swift'
                    placeholder='BIC/Swift'
                    value={account.value.bic ?? ''}
                    onBlur={() => setFocusAccount(null)}
                    onFocus={() => setFocusAccount(account.value?.metadata?.id)}
                    onChange={(e) => {
                      account.update((acc) => {
                        acc.bic = e.target.value;

                        return acc;
                      });
                    }}
                  />
                </>
              )}
            <div className='flex pb-1 gap-2'>
              {account.value.currency === 'GBP' && (
                <>
                  <MaskedInput
                    lazy={false}
                    variant='unstyled'
                    placeholderChar='_'
                    className='w-[110px]'
                    value={sortCode ?? ''}
                    placeholder='Enter sort code'
                    mask={[{ mask: '00-00-00' }, { mask: /^[0-9]{0,6}$/ }]}
                    definitions={{
                      '0': /[0-9]/,
                    }}
                    onAccept={(val) => {
                      setSortCode(val);
                      account.update((acc) => {
                        acc.sortCode = val;

                        return acc;
                      });
                    }}
                  />
                  <MaskedInput
                    autoComplete='off'
                    variant='unstyled'
                    value={bankAccount}
                    name='accountNumber'
                    aria-label='Account number'
                    placeholder='Bank account #'
                    onBlur={() => setFocusAccount(null)}
                    mask='[XX] 00 0000 0000 0000 0000 0000 0000'
                    onFocus={() => setFocusAccount(account.value?.metadata?.id)}
                    definitions={{
                      X: /[A-Za-z]/,
                      '0': /[0-9]/,
                    }}
                    onAccept={(val) => {
                      setBankAccount(val);
                      account.update((acc) => {
                        acc.accountNumber = val;

                        return acc;
                      });
                    }}
                  />
                </>
              )}
            </div>
            {account.value.currency === 'USD' && (
              <>
                <Input
                  className='mb-1'
                  autoComplete='off'
                  variant='unstyled'
                  name='routingNumber'
                  aria-label='Routing number'
                  placeholder='Routing number'
                  onBlur={() => setFocusAccount(null)}
                  onFocus={() => setFocusAccount(account.value?.metadata?.id)}
                />
                <MaskedInput
                  autoComplete='off'
                  variant='unstyled'
                  value={bankAccount}
                  name='accountNumber'
                  aria-label='Account number'
                  placeholder='Bank account #'
                  onBlur={() => setFocusAccount(null)}
                  mask='[XX] 00 0000 0000 0000 0000 0000 0000'
                  onFocus={() => setFocusAccount(account.value?.metadata?.id)}
                  definitions={{
                    X: /[A-Za-z]/,
                    '0': /[0-9]/,
                  }}
                  onAccept={(val) => {
                    setBankAccount(val);
                    account.update((acc) => {
                      acc.accountNumber = val;

                      return acc;
                    });
                  }}
                />
              </>
            )}
            {account?.value?.allowInternational &&
              (account?.value?.currency === 'USD' ||
                account?.value?.currency === 'GBP') && (
                <Input
                  name='bic'
                  autoComplete='off'
                  variant='unstyled'
                  aria-label='BIC/Swift'
                  placeholder='BIC/Swift'
                  value={account.value.bic ?? ''}
                  onBlur={() => setFocusAccount(null)}
                  onFocus={() => setFocusAccount(account.value?.metadata?.id)}
                  onChange={(e) => {
                    account.update((acc) => {
                      acc.bic = e.target.value;

                      return acc;
                    });
                  }}
                />
              )}

            {(account?.value?.allowInternational ||
              ![Currency.Gbp, Currency.Usd, Currency.Eur].includes(
                account?.value?.currency as Currency,
              )) && (
              <Textarea
                autoComplete='off'
                variant='unstyled'
                name='otherDetails'
                aria-label='Other details'
                placeholder='Other details'
                onBlur={() => setFocusAccount(null)}
                value={account.value.otherDetails ?? ''}
                onFocus={() => setFocusAccount(account.value?.metadata?.id)}
                onChange={(e) => {
                  account.update((acc) => {
                    acc.otherDetails = e.target.value;

                    return acc;
                  });
                }}
              />
            )}
          </CardContent>
        </Card>
      </>
    );
  },
);
