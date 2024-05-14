import { FC } from 'react';

import { observer } from 'mobx-react-lite';

import { Divider } from '@ui/presentation/Divider/Divider';
import { ServiceList } from '@organization/components/Tabs/panels/AccountPanel/ContractNew/ContractBillingDetailsModal/Services/components/ServiceList';
import { AddNewServiceMenu } from '@organization/components/Tabs/panels/AccountPanel/ContractNew/ContractBillingDetailsModal/Services/components/AddNewServiceMenu';

interface SubscriptionServiceModalProps {
  currency?: string;
}

export const Services: FC<SubscriptionServiceModalProps> = observer(
  ({ currency }) => {
    return (
      <>
        <div className='flex relative items-center h-8 '>
          <p className='text-sm text-gray-500 after:border-t-2 w-fit whitespace-nowrap mr-2'>
            Services
          </p>
          <Divider />
          <AddNewServiceMenu isInline={false} />
        </div>

        <ServiceList currency={currency} />
      </>
    );
  },
);
