'use client';
import React, { FC } from 'react';

import { observer } from 'mobx-react-lite';

import { TenantBillingProfile } from '@graphql/types';
import { Divider } from '@ui/presentation/Divider/Divider';
import { ServiceList } from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/ContractBillingDetailsModal/Services/components/ServiceList';
import { AddNewServiceMenu } from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/ContractBillingDetailsModal/Services/components/AddNewServiceMenu';

interface SubscriptionServiceModalProps {
  formId: string;
  currency?: string;
  renewedAt?: string;
  contractId: string;
  billingEnabled?: boolean;
  payAutomatically?: boolean | null;
  tenantBillingProfile?: TenantBillingProfile | null;
}

export const Services: FC<SubscriptionServiceModalProps> = observer(() => {
  return (
    <>
      <div className='flex relative items-center h-8 '>
        <p className='text-sm text-gray-500 after:border-t-2 w-fit whitespace-nowrap mr-2'>
          Services
        </p>
        <Divider />
        <AddNewServiceMenu isInline={false} />
      </div>

      <ServiceList />
    </>
  );
});
