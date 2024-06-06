import React from 'react';

import { observer } from 'mobx-react-lite';

import { DateTimeUtils } from '@utils/date.ts';
import { useStore } from '@shared/hooks/useStore';
import { BilledType, ContractStatus, ServiceLineItem } from '@graphql/types';
import ServiceLineItemStore from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/stores/Service.store.ts';

import { ServiceCard } from './ServiceCard';
import { useEditContractModalStores } from '../../stores/EditContractModalStores';

export const ServiceList: React.FC<{
  id: string;
  currency?: string;
  billingEnabled: boolean;
  contractStatus?: ContractStatus | null;
}> = observer(({ id, currency, contractStatus, billingEnabled }) => {
  const { serviceFormStore } = useEditContractModalStores();
  const store = useStore();
  const ids = store.contracts.value
    .get(id)
    ?.value?.contractLineItems?.map((item) => item?.metadata?.id);

  const serviceLineItems = ids?.map(
    (id) => store.contractLineItems?.value.get(id)?.value,
  );

  const groupServicesByParentId = (
    services: ServiceLineItem[],
  ): Array<ServiceLineItem[]> => {
    const { subscription, once } = services.reduce<{
      once: ServiceLineItem[];
      subscription: ServiceLineItem[];
    }>(
      (acc, item) => {
        const key: 'subscription' | 'once' = [
          BilledType.Monthly,
          BilledType.Quarterly,
          BilledType.Annually,
        ].includes(item.billingCycle)
          ? 'subscription'
          : 'once';

        acc[key].push(item);

        return acc;
      },
      { subscription: [], once: [] },
    );

    const getGroupedServices = (services: ServiceLineItem[]) => {
      const grouped: Record<string, ServiceLineItem[]> = {};

      services.forEach((service) => {
        const parentId = service?.parentId;
        if (parentId) {
          if (!grouped[parentId]) {
            grouped[parentId] = [];
          }
          grouped[parentId].push(service);
        }
      });
      const sortedGroups = Object.values(grouped).map((group) =>
        group.sort(
          (a, b) =>
            new Date(a?.serviceStarted).getTime() -
            new Date(b?.serviceStarted).getTime(),
        ),
      );

      // Filtering groups to exclude those where all items have 'serviceEnded' as null
      const filtered = sortedGroups.filter((group) =>
        group.some((service) => service?.serviceEnded === null),
      );

      return filtered;
    };

    return {
      subscription: getGroupedServices(subscription),
      once: getGroupedServices(once),
    };
  };

  const groupedServicesByParentId = groupServicesByParentId(serviceLineItems);
  console.log('🏷️ ----- : ', groupedServicesByParentId);

  return (
    <div className='flex flex-col'>
      {serviceFormStore.subscriptionServices.length !== 0 && (
        <p className='text-sm font-medium mb-2'>Subscriptions</p>
      )}
      {serviceFormStore.subscriptionServices.map((data, i) => (
        <React.Fragment
          key={`subscription-card-item-${data[0]?.serviceLineItem?.parentId}-${data[0].serviceLineItem?.description}-${i}`}
        >
          <ServiceCard
            data={data}
            currency={currency}
            type='subscription'
            contractStatus={contractStatus}
            billingEnabled={billingEnabled}
          />
        </React.Fragment>
      ))}
      {serviceFormStore.oneTimeServices.length !== 0 && (
        <p className='text-sm font-medium mb-2'>One-time</p>
      )}
      {serviceFormStore.oneTimeServices.map((data, i) => (
        <React.Fragment
          key={`one-time-card-item${data[0]?.serviceLineItem?.parentId}-${data[0].serviceLineItem?.description}-${i}`}
        >
          <ServiceCard
            data={data}
            type='one-time'
            currency={currency}
            contractStatus={contractStatus}
            billingEnabled={billingEnabled}
          />
        </React.Fragment>
      ))}
    </div>
  );
});
