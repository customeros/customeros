import React from 'react';

import { observer } from 'mobx-react-lite';

import { ServiceCard } from './ServiceCard';
import { useEditContractModalStores } from '../../stores/EditContractModalStores';

export const ServiceList: React.FC<{
  currency?: string;
}> = observer(({ currency }) => {
  const { serviceFormStore } = useEditContractModalStores();

  return (
    <div className='flex flex-col'>
      {serviceFormStore.subscriptionServices.length !== 0 && (
        <p className='text-sm font-medium mb-2'>Subscriptions</p>
      )}
      {serviceFormStore.subscriptionServices.map((data, i) => (
        <React.Fragment
          key={`subscription-card-item-${data[0]?.serviceLineItem?.parentId}-${data[0].serviceLineItem?.description}-${i}`}
        >
          <ServiceCard data={data} currency={currency} type='subscription' />
        </React.Fragment>
      ))}
      {serviceFormStore.oneTimeServices.length !== 0 && (
        <p className='text-sm font-medium mb-2'>One-time</p>
      )}
      {serviceFormStore.oneTimeServices.map((data, i) => (
        <React.Fragment
          key={`one-time-card-item${data[0]?.serviceLineItem?.parentId}-${data[0].serviceLineItem?.description}-${i}`}
        >
          <ServiceCard data={data} type='one-time' currency={currency} />
        </React.Fragment>
      ))}
    </div>
  );
});
