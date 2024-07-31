import { useDeepCompareEffect } from 'rooks';
import { useBankAccountsQuery } from '@settings/graphql/getBankAccounts.generated';
import { useBankTransferSelectionContext } from '@settings/components/Tabs/panels/BillingPanel/context/BankTransferSelectionContext';

import { BankAccount } from '@graphql/types';
import { FormSwitch } from '@ui/form/Switch/FormSwitch';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

import { AddAccountButton } from './AddAccountButton';
import { BankTransferCard } from './BankTransferCard';

export const BankTransferAccountList = ({
  formId,
  legalName,
}: {
  formId: string;
  legalName?: string | null;
}) => {
  const client = getGraphQLClient();
  const { data } = useBankAccountsQuery(client);
  const { setAccounts } = useBankTransferSelectionContext();
  const existingAccountCurrencies = data?.bankAccounts.map(
    (account) => account.currency as string,
  );

  useDeepCompareEffect(() => {
    if (data?.bankAccounts) {
      setAccounts((data?.bankAccounts as BankAccount[]) ?? []);
    }
  }, [data?.bankAccounts]);

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

      {data?.bankAccounts?.map((account) => (
        <div key={account.metadata.id}>
          <BankTransferCard
            account={account as BankAccount}
            existingCurrencies={existingAccountCurrencies ?? []}
          />
        </div>
      ))}
    </>
  );
};
