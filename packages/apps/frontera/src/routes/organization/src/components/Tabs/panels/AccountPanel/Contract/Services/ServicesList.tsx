import React from 'react';

import { observer } from 'mobx-react-lite';
import { ContractStore } from '@store/Contracts/Contract.store.ts';

import { useStore } from '@shared/hooks/useStore';
import { BilledType, ServiceLineItem } from '@graphql/types';
import { PauseCircle } from '@ui/media/icons/PauseCircle.tsx';
import { groupServicesByParentId } from '@organization/components/Tabs/panels/AccountPanel/Contract/Services/utils.ts';

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

function formatCurrency(amount: number, currency: string = 'USD'): string {
  return Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(amount);
}

const ServiceItem = observer(
  ({
    onOpen,
    currency,
    isPaused,
    id,
  }: {
    id: string;
    isPaused?: boolean;
    currency?: string | null;
    onOpen: (props: ServiceLineItem) => void;
  }) => {
    const store = useStore();
    const contractLineItem = store.contractLineItems?.value.get(id)?.value;

    return (
      <>
        <div
          onClick={() => onOpen(contractLineItem as ServiceLineItem)}
          className='flex w-full justify-between cursor-pointer text-sm focus:outline-none'
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

              {formatCurrency(contractLineItem?.price ?? 0, currency || 'USD')}
              {getBilledTypeLabel(contractLineItem?.billingCycle as BilledType)}
              {isPaused && (
                <PauseCircle className='ml-2 text-gray-500 size-4' />
              )}
            </p>
          </div>
        </div>
      </>
    );
  },
);

interface ServicesListProps {
  id: string;
  onModalOpen: () => void;
  currency?: string | null;
}

export const ServicesList = observer(
  ({ currency, onModalOpen, id }: ServicesListProps) => {
    const store = useStore();
    const contractStore = store.contracts.value.get(id) as ContractStore;
    const data = contractStore?.contractLineItems
      ?.map((e) => e.value)
      ?.filter((e) => !e?.metadata?.id?.includes('new'));
    const groupedServicesByParentId = groupServicesByParentId(data);

    return (
      <div className='w-full flex flex-col gap-1 mt-2'>
        {groupedServicesByParentId?.subscription?.length > 0 && (
          <article className='mb-1'>
            <h1
              className='font-semibold text-sm mb-1'
              data-test='account-panel-contract-subscription'
            >
              Subscriptions
            </h1>
            {groupedServicesByParentId?.subscription
              ?.sort((a, b) => {
                const aDate = new Date(a.currentLineItem?.serviceStarted || 0);
                const bDate = new Date(b.currentLineItem?.serviceStarted || 0);

                if (aDate.getTime() !== bDate.getTime()) {
                  return bDate.getTime() - aDate.getTime();
                }

                return (
                  (b.currentLineItem?.price || 0) -
                  (a.currentLineItem?.price || 0)
                );
              })
              .map((service) => (
                <React.Fragment
                  key={`service-item-${service?.currentLineItem?.metadata?.id}`}
                >
                  <ServiceItem
                    currency={currency}
                    onOpen={onModalOpen}
                    id={service?.currentLineItem?.metadata?.id}
                    isPaused={service?.currentLineItem?.paused}
                  />
                </React.Fragment>
              ))}
          </article>
        )}

        {groupedServicesByParentId?.once?.length > 0 && (
          <article>
            <h1
              className='font-semibold text-sm mb-1'
              data-test='account-panel-contract-one-time'
            >
              One-time
            </h1>
            {groupedServicesByParentId?.once
              ?.sort((a, b) => {
                const aDate = new Date(a.currentLineItem?.serviceStarted || 0);
                const bDate = new Date(b.currentLineItem?.serviceStarted || 0);

                if (aDate.getTime() !== bDate.getTime()) {
                  return bDate.getTime() - aDate.getTime();
                }

                return (
                  (b.currentLineItem?.price || 0) -
                  (a.currentLineItem?.price || 0)
                );
              })
              .map((service) => (
                <React.Fragment
                  key={`service-item-${service?.currentLineItem?.metadata?.id}`}
                >
                  <ServiceItem
                    currency={currency}
                    onOpen={onModalOpen}
                    id={service?.currentLineItem?.metadata?.id}
                  />
                </React.Fragment>
              ))}
          </article>
        )}
      </div>
    );
  },
);
