import React from 'react';

import { observer } from 'mobx-react-lite';
import { ContractStore } from '@store/Contracts/Contract.store.ts';

import { useStore } from '@shared/hooks/useStore';
import { BilledType, ServiceLineItem } from '@graphql/types';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber';
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
            <h1 className='font-semibold text-sm mb-1'>Subscriptions</h1>
            {groupedServicesByParentId?.subscription?.map((service) => (
              <React.Fragment
                key={`service-item-${service?.currentLineItem?.metadata?.id}`}
              >
                <ServiceItem
                  id={service?.currentLineItem?.metadata?.id}
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
              <React.Fragment
                key={`service-item-${service?.currentLineItem?.metadata?.id}`}
              >
                <ServiceItem
                  id={service?.currentLineItem?.metadata?.id}
                  onOpen={onModalOpen}
                  currency={currency}
                />
              </React.Fragment>
            ))}
          </article>
        )}
      </div>
    );
  },
);
