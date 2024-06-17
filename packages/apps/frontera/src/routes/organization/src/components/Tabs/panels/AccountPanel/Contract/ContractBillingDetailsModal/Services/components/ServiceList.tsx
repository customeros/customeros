import React from 'react';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { BilledType, ContractStatus, ServiceLineItem } from '@graphql/types';

import { ServiceCard } from './ServiceCard';

export const ServiceList: React.FC<{
  id: string;
  currency?: string;
  billingEnabled: boolean;
  contractStatus?: ContractStatus | null;
}> = observer(({ id, currency, contractStatus, billingEnabled }) => {
  const store = useStore();
  const ids = store.contracts.value
    .get(id)
    ?.value?.contractLineItems?.map((item) => item?.metadata?.id);

  const serviceLineItems =
    ids?.map((id) => store.contractLineItems?.value.get(id)?.value) || [];

  const groupServicesByParentId = (services: ServiceLineItem[]) => {
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
        const parentId = service?.parentId || service?.metadata?.id;
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

  const groupedServicesByParentId = groupServicesByParentId(
    serviceLineItems as ServiceLineItem[],
  );

  return (
    <div className='flex flex-col'>
      {groupedServicesByParentId.subscription.length !== 0 && (
        <p className='text-sm font-medium mb-2'>Subscriptions</p>
      )}
      {groupedServicesByParentId.subscription.map((data, i) => (
        <React.Fragment
          key={`subscription-card-item-${data[0]?.parentId}-${data[0].description}-${i}`}
        >
          <ServiceCard
            contractId={id}
            ids={data.map((e) => e?.metadata?.id)}
            currency={currency ?? 'USD'}
            type='subscription'
            contractStatus={contractStatus}
            billingEnabled={billingEnabled}
          />
        </React.Fragment>
      ))}
      {groupedServicesByParentId.once.length !== 0 && (
        <p className='text-sm font-medium mb-2'>One-time</p>
      )}
      {groupedServicesByParentId.once.map((data, i) => (
        <React.Fragment
          key={`one-time-card-item${data[0]?.parentId}-${data[0].description}-${i}`}
        >
          <ServiceCard
            contractId={id}
            type='one-time'
            ids={data.map((e) => e?.metadata?.id)}
            currency={currency ?? 'USD'}
            contractStatus={contractStatus}
            billingEnabled={billingEnabled}
          />
        </React.Fragment>
      ))}
    </div>
  );
});
