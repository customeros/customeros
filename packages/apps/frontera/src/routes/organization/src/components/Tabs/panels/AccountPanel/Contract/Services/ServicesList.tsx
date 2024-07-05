import React from 'react';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { BilledType, ServiceLineItem } from '@graphql/types';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber';

function getBilledTypeLabel(billedType: BilledType): string {
  switch (billedType) {
    case BilledType.Annually:
      return '/year';
    case BilledType.Monthly:
      return '/month';
    case BilledType.None:
      return '';
    case BilledType.Once:
      return '';
    case BilledType.Usage:
      return '/use';
    case BilledType.Quarterly:
      return '/quarter';
    default:
      return '';
  }
}

const ServiceItem = observer(
  ({
    onOpen,
    currency,
    id,
  }: {
    id: string;
    currency?: string | null;
    onOpen: (props: ServiceLineItem) => void;
  }) => {
    const store = useStore();
    const contractLineItem = store.contractLineItems?.value.get(id)?.value;

    const allowedFractionDigits =
      contractLineItem?.billingCycle === BilledType.Usage ? 4 : 2;

    return (
      <>
        <div
          className='flex w-full justify-between cursor-pointer text-sm focus:outline-none'
          onClick={() => onOpen(contractLineItem as ServiceLineItem)}
        >
          {contractLineItem?.description && (
            <p>{contractLineItem?.description}</p>
          )}
          <div className='flex justify-between'>
            <p>
              {![BilledType.Usage, BilledType.None].includes(
                contractLineItem?.billingCycle as BilledType,
              ) && (
                <>
                  {contractLineItem?.quantity}
                  <span className='text-sm mx-1'>×</span>
                </>
              )}

              {formatCurrency(
                contractLineItem?.price ?? 0,
                allowedFractionDigits,
                currency || 'USD',
              )}
              {getBilledTypeLabel(contractLineItem?.billingCycle as BilledType)}
            </p>
          </div>
        </div>
      </>
    );
  },
);

interface ServicesListProps {
  onModalOpen: () => void;
  currency?: string | null;
  data?: Array<ServiceLineItem>;
}

export const ServicesList = ({
  data = [],
  currency,
  onModalOpen,
}: ServicesListProps) => {
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
      return sortedGroups.filter((group) =>
        group.some((service) => service?.serviceEnded === null),
      );
    };

    return {
      subscription: getGroupedServices(subscription),
      once: getGroupedServices(once),
    };
  };

  const groupedServicesByParentId = groupServicesByParentId(data);

  return (
    <div className='w-full flex flex-col gap-1 mt-2'>
      {groupedServicesByParentId?.subscription?.length > 0 && (
        <article className='mb-1'>
          <h1 className='font-semibold text-sm mb-1'>Subscriptions</h1>
          {groupedServicesByParentId?.subscription?.map((service) => (
            <React.Fragment key={`service-item-${service?.[0]?.metadata?.id}`}>
              <ServiceItem
                id={service?.[0]?.metadata?.id}
                onOpen={onModalOpen}
                currency={currency}
              />
            </React.Fragment>
          ))}
        </article>
      )}

      {groupedServicesByParentId?.once?.length > 0 && (
        <article>
          <h1 className='font-semibold text-sm mb-1'>One-time</h1>
          {groupedServicesByParentId?.once?.map((service) => (
            <React.Fragment key={`service-item-${service?.[0]?.metadata?.id}`}>
              <ServiceItem
                id={service?.[0]?.metadata?.id}
                onOpen={onModalOpen}
                currency={currency}
              />
            </React.Fragment>
          ))}
        </article>
      )}
    </div>
  );
};
