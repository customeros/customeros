import { FC } from 'react';

import { observer } from 'mobx-react-lite';

import { ContractStatus } from '@graphql/types';
import { Divider } from '@ui/presentation/Divider/Divider.tsx';
import { ServiceList } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/components/Services/components/ServiceList.tsx';
import { AddNewServiceMenu } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/components/Services/components/AddNewServiceMenu.tsx';

interface SubscriptionServiceModalProps {
  id: string;
  currency?: string;
  contractStatus?: ContractStatus | null;
}

export const Services: FC<SubscriptionServiceModalProps> = observer(
  ({ id, currency, contractStatus }) => {
    return (
      <>
        <div className='flex relative items-center h-8 '>
          <p className='text-sm text-gray-500 after:border-t-2 w-fit whitespace-nowrap mr-2'>
            Services
          </p>
          <Divider />
          <AddNewServiceMenu contractId={id} isInline={false} />
        </div>

        <ServiceList
          id={id}
          currency={currency}
          contractStatus={contractStatus}
        />
      </>
    );
  },
);
