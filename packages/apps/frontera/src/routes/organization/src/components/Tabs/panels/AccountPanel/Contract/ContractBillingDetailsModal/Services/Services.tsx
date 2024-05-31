import { FC } from 'react';

import { observer } from 'mobx-react-lite';

import { ContractStatus } from '@graphql/types';
import { Divider } from '@ui/presentation/Divider/Divider';
import { ServiceList } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/Services/components/ServiceList';
import { AddNewServiceMenu } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/Services/components/AddNewServiceMenu';

interface SubscriptionServiceModalProps {
  currency?: string;
  billingEnabled?: boolean;
  contractStatus?: ContractStatus | null;
}

export const Services: FC<SubscriptionServiceModalProps> = observer(
  ({ currency, contractStatus, billingEnabled = false }) => {
    return (
      <>
        <div className='flex relative items-center h-8 '>
          <p className='text-sm text-gray-500 after:border-t-2 w-fit whitespace-nowrap mr-2'>
            Services
          </p>
          <Divider />
          <AddNewServiceMenu isInline={false} />
        </div>

        <ServiceList
          currency={currency}
          contractStatus={contractStatus}
          billingEnabled={billingEnabled as boolean}
        />
      </>
    );
  },
);
