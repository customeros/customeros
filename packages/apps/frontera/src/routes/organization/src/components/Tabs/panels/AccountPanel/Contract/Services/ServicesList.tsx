import React from 'react';

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

const ServiceItem = ({
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
        onClick={() => onOpen(contractLineItem)}
      >
        {contractLineItem?.description && (
          <p>{contractLineItem?.description}</p>
        )}
        <div className='flex justify-between'>
          <p>
            {![BilledType.Usage, BilledType.None].includes(
              contractLineItem?.billingCycle,
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
            {getBilledTypeLabel(contractLineItem?.billingCycle)}
          </p>
        </div>
      </div>
    </>
  );
};

interface ServicesListProps {
  onModalOpen: () => void;
  currency?: string | null;
  data?: Array<ServiceLineItem>;
}

export const ServicesList = ({
  data,
  currency,
  onModalOpen,
}: ServicesListProps) => {
  const filteredData = data?.filter(({ serviceEnded }) => !serviceEnded) ?? [];
  const { subscription, once } = filteredData.reduce<{
    once: Array<ServiceLineItem>;
    subscription: Array<ServiceLineItem>;
  }>(
    (acc, service) => {
      const key: 'subscription' | 'once' = [
        BilledType.Monthly,
        BilledType.Quarterly,
        BilledType.Annually,
      ].includes(service.billingCycle)
        ? 'subscription'
        : 'once';

      acc[key].push(service);

      return acc;
    },
    { subscription: [], once: [] },
  );

  return (
    <div className='w-full flex flex-col gap-1 mt-2'>
      {subscription?.length > 0 && (
        <article className='mb-1'>
          <h1 className='font-semibold text-sm mb-1'>Subscriptions</h1>
          {subscription?.map((service) => (
            <React.Fragment key={`service-item-${service?.metadata?.id}`}>
              <ServiceItem
                id={service?.metadata?.id}
                onOpen={onModalOpen}
                currency={currency}
              />
            </React.Fragment>
          ))}
        </article>
      )}

      {once?.length > 0 && (
        <article>
          <h1 className='font-semibold text-sm mb-1'>One-time</h1>
          {once?.map((service) => (
            <React.Fragment key={`service-item-${service?.metadata?.id}`}>
              <ServiceItem
                id={service?.metadata?.id}
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
