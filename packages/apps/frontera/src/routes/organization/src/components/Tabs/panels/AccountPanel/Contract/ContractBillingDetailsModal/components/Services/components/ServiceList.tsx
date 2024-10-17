import React from 'react';

import { observer } from 'mobx-react-lite';
import { ContractStore } from '@store/Contracts/Contract.store';
import { ContractLineItemStore } from '@store/ContractLineItems/ContractLineItem.store';

import { useStore } from '@shared/hooks/useStore';
import { BilledType, ContractStatus, ServiceLineItem } from '@graphql/types';

import { ServiceCard } from './ServiceCard';

interface ServiceListProps {
  id: string;
  currency?: string;
  contractStatus?: ContractStatus | null;
}

export const ServiceList = observer(
  ({ id, currency, contractStatus }: ServiceListProps) => {
    const store = useStore();
    const ids = (
      store.contracts.value.get(id) as ContractStore
    )?.tempValue?.contractLineItems?.map((item) => item?.metadata?.id);

    const serviceLineItems =
      ids
        ?.map(
          (id) =>
            (store.contractLineItems?.value.get(id) as ContractLineItemStore)
              ?.tempValue,
        )
        .filter((e) => Boolean(e)) || [];

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
        {groupedServicesByParentId.subscription
          .sort((a, b) => {
            const aDate = new Date(a[0]?.serviceStarted || 0);
            const bDate = new Date(b[0]?.serviceStarted || 0);

            const dateComparison = bDate.getTime() - aDate.getTime();

            if (dateComparison === 0) {
              const aTimestamp = new Date(a[0]?.serviceStarted || 0).getTime();
              const bTimestamp = new Date(b[0]?.serviceStarted || 0).getTime();

              return bTimestamp - aTimestamp;
            }

            return dateComparison;
          })
          .map((data) => (
            <React.Fragment
              key={`subscription-card-item-${data[0]?.parentId}-${data[0].description}-${data[0].metadata.id}`}
            >
              <ServiceCard
                contractId={id}
                type='subscription'
                currency={currency ?? 'USD'}
                contractStatus={contractStatus}
                ids={data.map((e) => e?.metadata?.id)}
              />
            </React.Fragment>
          ))}

        {groupedServicesByParentId.once.length !== 0 && (
          <p className='text-sm font-medium mb-2'>One-time</p>
        )}
        {groupedServicesByParentId.once
          .sort((a, b) => {
            const aDate = new Date(a[0]?.serviceStarted || 0);
            const bDate = new Date(b[0]?.serviceStarted || 0);

            const dateComparison = bDate.getTime() - aDate.getTime();

            if (dateComparison === 0) {
              const aTimestamp = new Date(a[0]?.serviceStarted || 0).getTime();
              const bTimestamp = new Date(b[0]?.serviceStarted || 0).getTime();

              return bTimestamp - aTimestamp;
            }

            return dateComparison;
          })
          .map((data, i) => (
            <React.Fragment
              key={`one-time-card-item-${data[0]?.parentId}-${data[0].description}-${i}`}
            >
              <ServiceCard
                contractId={id}
                type='one-time'
                currency={currency ?? 'USD'}
                contractStatus={contractStatus}
                ids={data.map((e) => e?.metadata?.id)}
              />
            </React.Fragment>
          ))}
      </div>
    );
  },
);
