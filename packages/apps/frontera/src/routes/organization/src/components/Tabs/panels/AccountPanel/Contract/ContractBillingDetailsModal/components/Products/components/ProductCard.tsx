import { useState } from 'react';

import { observer } from 'mobx-react-lite';
import { ContractLineItemStore } from '@store/ContractLineItems/ContractLineItem.store';

import { cn } from '@ui/utils/cn';
import { DateTimeUtils } from '@utils/date';
import { ContractStatus } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { FlipBackward } from '@ui/media/icons/FlipBackward';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { ChevronExpand } from '@ui/media/icons/ChevronExpand';
import { ChevronCollapse } from '@ui/media/icons/ChevronCollapse';
import { Card, CardHeader, CardContent } from '@ui/presentation/Card/Card';

import { ProductItem } from './ProductItem';
import { ProductItemMenu } from './ProductItemMenu';

interface ProductCardProps {
  ids?: string[];
  currency: string;
  contractId: string;
  type: 'subscription' | 'one-time';
  contractStatus?: ContractStatus | null;
}

export const ProductCard = observer(
  ({ ids, type, contractId, currency, contractStatus }: ProductCardProps) => {
    const [showEnded, setShowEnded] = useState(false);
    const [allowIndividualRestore, setAllowIndividualRestore] = useState(true);
    const store = useStore();
    const contractLineItemsStore = store.contractLineItems;
    const contractLineItems = contractLineItemsStore.value;

    const thisGroupLineItems = ids?.map(
      (id) => contractLineItems.get(id) as ContractLineItemStore,
    );
    const endedServices = thisGroupLineItems?.filter((service) => {
      return (
        !!service?.tempValue?.serviceEnded &&
        DateTimeUtils.isPast(service?.tempValue?.serviceEnded)
      );
    });

    const liveServices = thisGroupLineItems?.filter(
      (service) =>
        !service?.tempValue?.serviceEnded ||
        !DateTimeUtils.isPast(service?.tempValue?.serviceEnded),
    );

    const closedServices = thisGroupLineItems?.filter(
      (service) => service?.tempValue?.closed,
    );

    const sku = liveServices?.[0]?.tempValue?.skuId
      ? store.skus.getById(liveServices[0].tempValue.skuId)
      : null;

    const isClosed = liveServices?.every(
      (service) => service?.tempValue?.closed,
    );

    const handleCloseChange = (closed: boolean) => {
      liveServices?.forEach((service) => {
        (service as ContractLineItemStore)?.updateTemp((prev) => ({
          ...prev,
          closed,
        }));
      });
      closedServices?.forEach((service) => {
        (service as ContractLineItemStore)?.updateTemp((prev) => ({
          ...prev,
          closed,
        }));
      });
      setAllowIndividualRestore(!closed);
    };

    const handlePauseChange = (paused: boolean) => {
      liveServices?.forEach((service) => {
        (service as ContractLineItemStore)?.updateTemp((prev) => ({
          ...prev,
          paused,
        }));
      });
      closedServices?.forEach((service) => {
        (service as ContractLineItemStore)?.updateTemp((prev) => ({
          ...prev,
          paused,
        }));
      });
      setAllowIndividualRestore(!closed);
    };

    return (
      <Card className='px-3 py-2 mb-2 rounded-lg'>
        <CardHeader className={cn('flex justify-between pb-0.5')}>
          <p
            title={
              sku?.value?.name || thisGroupLineItems?.[0]?.value?.description
            }
            className={cn(
              'text-gray-700 min-w-2.5 w-full min-h-0 border-none hover:border-none focus:border-none truncate overflow-hidden',
              {
                'text-gray-400 line-through': isClosed,
              },
            )}
          >
            {sku?.value?.name || thisGroupLineItems?.[0]?.value?.description}
          </p>

          <div className='flex items-baseline'>
            {endedServices && endedServices.length > 0 && (
              <IconButton
                size='xs'
                variant='ghost'
                className='p-0 px-1 text-gray-400'
                onClick={() => setShowEnded(!showEnded)}
                aria-label={
                  showEnded ? 'Hide ended services' : 'Show ended services'
                }
                icon={
                  showEnded ? (
                    <ChevronCollapse className='text-inherit' />
                  ) : (
                    <ChevronExpand className='text-inherit' />
                  )
                }
              />
            )}

            {isClosed ? (
              <>
                <IconButton
                  size='xs'
                  variant='ghost'
                  aria-label='Undo'
                  onClick={() => handleCloseChange(false)}
                  icon={<FlipBackward className='text-gray-400' />}
                  className='p-1  max-h-5 hover:bg-gray-100 rounded translate-x-1'
                />
              </>
            ) : (
              <ProductItemMenu
                contractId={contractId}
                handleCloseService={handleCloseChange}
                handlePauseService={handlePauseChange}
                closed={thisGroupLineItems?.[0]?.tempValue?.closed}
                paused={thisGroupLineItems?.[0]?.tempValue?.paused}
                allowPausing={
                  type !== 'one-time' && contractStatus === ContractStatus.Live
                }
                allowAddModification={
                  type !== 'one-time' &&
                  !!thisGroupLineItems?.[0]?.tempValue?.parentId
                }
                id={
                  thisGroupLineItems?.[0]?.tempValue?.parentId ||
                  thisGroupLineItems?.[0]?.tempValue?.metadata?.id ||
                  ''
                }
              />
            )}
          </div>
        </CardHeader>
        <CardContent className='text-sm p-0 gap-y-0.25 flex flex-col'>
          {showEnded &&
            endedServices?.map(
              (service, serviceIndex) =>
                service && (
                  <ProductItem
                    isEnded
                    type={type}
                    service={service}
                    currency={currency}
                    isModification={false}
                    contractStatus={contractStatus}
                    key={`ended-service-item-${serviceIndex}`}
                    allowIndividualRestore={allowIndividualRestore}
                  />
                ),
            )}
          {liveServices?.map(
            (service, serviceIndex) =>
              service && (
                <ProductItem
                  type={type}
                  service={service}
                  currency={currency}
                  contractStatus={contractStatus}
                  key={`service-item-${service.id}`}
                  allowIndividualRestore={allowIndividualRestore}
                  allServices={thisGroupLineItems as ContractLineItemStore[]}
                  isModification={
                    thisGroupLineItems &&
                    thisGroupLineItems?.length > 1 &&
                    serviceIndex !== 0
                  }
                />
              ),
          )}
        </CardContent>
      </Card>
    );
  },
);
