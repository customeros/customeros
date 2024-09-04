import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { FormSwitch } from '@ui/form/Switch/FormSwitch';

import { AddAccountButton } from './AddAccountButton';
import { BankTransferCard } from './BankTransferCard';

export const BankTransferAccountList = observer(
  ({ formId, legalName }: { formId: string; legalName?: string | null }) => {
    const store = useStore();

    const bankAccounts = store.settings.bankAccounts?.toArray();

    const existingAccountCurrencies = bankAccounts?.map(
      (account) => account.value?.currency as string,
    );

    return (
      <>
        <div className='flex justify-between items-center'>
          <span className='text-sm whitespace-nowrap font-semibold'>
            Bank transfer
          </span>
          <div className='flex items-center'>
            <AddAccountButton
              legalName={legalName}
              existingCurrencies={existingAccountCurrencies ?? []}
            />
            <div>
              <FormSwitch
                size='sm'
                formId={formId}
                name='canPayWithBankTransfer'
              />
            </div>
          </div>
        </div>

        {bankAccounts?.map((account) => (
          <div key={account.value?.metadata?.id}>
            <BankTransferCard
              account={account}
              existingCurrencies={existingAccountCurrencies ?? []}
            />
          </div>
        ))}
      </>
    );
  },
);
