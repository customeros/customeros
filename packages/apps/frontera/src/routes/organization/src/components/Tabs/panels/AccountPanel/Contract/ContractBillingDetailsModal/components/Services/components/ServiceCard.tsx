import React, { useState, ChangeEvent } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn.ts';
import { ContractStatus } from '@graphql/types';
import { Input } from '@ui/form/Input/Input.tsx';
import { FlipBackward } from '@ui/media/icons/FlipBackward.tsx';
import { IconButton } from '@ui/form/IconButton/IconButton.tsx';
import { ChevronExpand } from '@ui/media/icons/ChevronExpand.tsx';
import { ChevronCollapse } from '@ui/media/icons/ChevronCollapse.tsx';
import { Card, CardHeader, CardContent } from '@ui/presentation/Card/Card.tsx';
// import { Highlighter } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/Services/components/highlighters';

import { ContractLineItemStore } from '@store/ContractLineItems/ContractLineItem.store.ts';

import { DateTimeUtils } from '@utils/date.ts';
import { useStore } from '@shared/hooks/useStore';

import { ServiceItem } from './ServiceItem';
import { ServiceItemMenu } from './ServiceItemMenu.tsx';

interface ServiceCardProps {
  ids?: string[];
  currency: string;
  contractId: string;
  type: 'subscription' | 'one-time';
  contractStatus?: ContractStatus | null;
}

export const ServiceCard: React.FC<ServiceCardProps> = observer(
  ({ ids, type, contractId, currency, contractStatus }) => {
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

    const [description, setDescription] = useState(
      liveServices?.[0]?.tempValue?.description || '',
    );

    const isClosed = liveServices?.every(
      (service) => service?.tempValue?.closed,
    );
    const handleDescriptionChange = (e: ChangeEvent<HTMLInputElement>) => {
      if (!e.target.value?.length) {
        setDescription('Unnamed');
      }
      const newName = !e.target.value?.length ? 'Unnamed' : e.target.value;

      liveServices?.forEach((service) => {
        (service as ContractLineItemStore)?.updateTemp((prev) => ({
          ...prev,
          description: newName,
        }));
      });
    };
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

    return (
      <Card className='px-3 py-2 mb-2 rounded-lg'>
        <CardHeader className={cn('flex justify-between pb-0.5')}>
          {/*<Highlighter*/}
          {/*  highlightVersion={descriptionLI?.uiMetadata?.shapeVariant}*/}
          {/*  backgroundColor={*/}
          {/*    liveServices.length === 1 &&*/}
          {/*    descriptionLI?.isNewlyAdded &&*/}
          {/*    !isClosed*/}
          {/*      ? descriptionLI.uiMetadata?.color*/}
          {/*      : undefined*/}
          {/*  }*/}
          {/*>*/}
          <Input
            value={description ?? ''}
            onChange={(e) => setDescription(e.target.value)}
            onBlur={handleDescriptionChange}
            onFocus={(e) => e.target.select()}
            size='xs'
            placeholder='Service name'
            className={cn(
              'text-base text-gray-700 min-w-2.5 w-full min-h-0 border-none hover:border-none focus:border-none flex-1',
              {
                'text-gray-400 line-through': isClosed,
              },
            )}
          />
          {/*</Highlighter>*/}

          <div className='flex items-baseline'>
            {endedServices && endedServices.length > 0 && (
              <IconButton
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
                variant='ghost'
                size='xs'
                className='p-0 px-1 text-gray-400'
                onClick={() => setShowEnded(!showEnded)}
              />
            )}

            {isClosed ? (
              <>
                <IconButton
                  aria-label='Undo'
                  icon={<FlipBackward className='text-gray-400' />}
                  size='xs'
                  className='p-1  max-h-5 hover:bg-gray-100 rounded translate-x-1'
                  variant='ghost'
                  onClick={() => handleCloseChange(false)}
                />
              </>
            ) : (
              <ServiceItemMenu
                id={
                  thisGroupLineItems?.[0]?.tempValue?.parentId ||
                  thisGroupLineItems?.[0]?.tempValue?.metadata?.id ||
                  ''
                }
                contractId={contractId}
                closed={thisGroupLineItems?.[0]?.tempValue?.closed}
                handleCloseService={handleCloseChange}
                allowAddModification={
                  type !== 'one-time' &&
                  !!thisGroupLineItems?.[0]?.tempValue?.parentId
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
                  <ServiceItem
                    key={`ended-service-item-${serviceIndex}`}
                    service={service}
                    currency={currency}
                    isEnded
                    contractStatus={contractStatus}
                    isModification={false}
                    type={type}
                    allowIndividualRestore={allowIndividualRestore}
                  />
                ),
            )}
          {liveServices?.map(
            (service, serviceIndex) =>
              service && (
                <ServiceItem
                  key={`service-item-${serviceIndex}`}
                  currency={currency}
                  service={service}
                  type={type}
                  isModification={
                    thisGroupLineItems &&
                    thisGroupLineItems?.length > 1 &&
                    serviceIndex !== 0
                  }
                  contractStatus={contractStatus}
                  allowIndividualRestore={allowIndividualRestore}
                  allServices={thisGroupLineItems as ContractLineItemStore[]}
                />
              ),
          )}
        </CardContent>
      </Card>
    );
  },
);
