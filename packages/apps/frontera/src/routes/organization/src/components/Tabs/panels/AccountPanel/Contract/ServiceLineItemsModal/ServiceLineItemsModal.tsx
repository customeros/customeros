import { useParams } from 'react-router-dom';
import { useRef, useState, Fragment, useEffect } from 'react';

import { produce } from 'immer';
import { useQueryClient } from '@tanstack/react-query';

import { Plus } from '@ui/media/icons/Plus';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { toastError } from '@ui/presentation/Toast';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { AutoresizeTextarea } from '@ui/form/Textarea/AutoresizeTextarea';
import { useTimelineMeta } from '@organization/components/Timeline/state';
import { useInfiniteGetTimelineQuery } from '@organization/graphql/getTimeline.generated';
import { useUpdateServicesMutation } from '@organization/graphql/updateServiceLineItems.generated';
import {
  GetContractsQuery,
  useGetContractsQuery,
} from '@organization/graphql/getContracts.generated';
import { useUpdateCacheWithNewEvent } from '@organization/components/Timeline/PastZone/hooks/updateCacheWithNewEvent';
import {
  BilledType,
  DataSource,
  InputMaybe,
  ServiceLineItem,
  ServiceLineItemBulkUpdateItem,
} from '@graphql/types';
import {
  Modal,
  ModalBody,
  ModalHeader,
  ModalFooter,
  ModalPortal,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal/Modal';
import { updateTimelineCacheAfterServiceLineItemChange } from '@organization/components/Tabs/panels/AccountPanel/Contract/ServiceLineItemsModal/utils';
import {
  ServiceLineItemsDTO,
  BulkUpdateServiceLineItem,
} from '@organization/components/Tabs/panels/AccountPanel/Contract/ServiceLineItemsModal/ServiceLineItemsModal.dto';

import { ServiceLineItemRow } from './ServiceLineItemRow';

interface SubscriptionServiceModalProps {
  isOpen: boolean;
  contractId: string;
  onClose: () => void;
  contractName: string;
  notes?: string | null;
  currency?: string | null;
  organizationName: string;
  contractLineItems: Array<ServiceLineItem>;
}

interface SubscriptionServiceModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const defaultValue = {
  name: 'Unnamed',
  quantity: 1,
  vatRate: 0,
  price: 0,
  billed: BilledType.Monthly,
  type: 'RECURRING',
  isDeleted: false,
};

const getNewItem = (input: InputMaybe<ServiceLineItemBulkUpdateItem>) => ({
  id: Math.random().toString(),
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString(),
  name: input?.name,
  billed: input?.billed,
  price: input?.price,
  quantity: input?.quantity,
  createdBy: '',
  source: DataSource.Openline,
  sourceOfTruth: '',
  appSource: DataSource.Openline,
  externalLinks: [],
  opportunities: [
    {
      comments: '',
      owner: null,
      internalStage: 'OPEN',
      internalType: 'RENEWAL',
      amount: input?.price,
      maxAmount: input?.price,
      name: '',
      renewalLikelihood: 'HIGH',
      renewalUpdatedByUserId: '',
      renewalUpdatedByUserAt: new Date().toISOString(),
      renewedAt: new Date().toISOString(),
    },
  ],
});
export const ServiceLineItemsModal = ({
  isOpen,
  onClose,
  contractLineItems,
  contractId,
  contractName,
  notes = '',
  currency,
}: SubscriptionServiceModalProps) => {
  const client = getGraphQLClient();
  const id = useParams()?.id as string;
  const [timelineMeta] = useTimelineMeta();
  const queryKey = useGetContractsQuery.getKey({ id });
  const timelineQueryKey = useInfiniteGetTimelineQuery.getKey(
    timelineMeta.getTimelineVariables,
  );
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const queryClient = useQueryClient();
  const store = useStore();
  const updateTimelineCache = useUpdateCacheWithNewEvent();

  const updateServices = useUpdateServicesMutation(client, {
    onMutate: ({ input }) => {
      queryClient.cancelQueries({ queryKey });
      queryClient.setQueryData<GetContractsQuery>(queryKey, (currentCache) => {
        return produce(currentCache, (draft) => {
          const previousContracts = draft?.['organization']?.['contracts'];
          if (draft?.['organization']?.['contracts']) {
            draft['organization']['contracts']?.map((contractData, index) => {
              const updatedContractIndex = previousContracts?.findIndex(
                (contract) => contract.metadata.id === input.contractId,
              );
              if (!draft) return;
              if (index !== updatedContractIndex) {
                return contractData;
              }

              return {
                ...contractData,
                invoiceNote: input.invoiceNote,
                serviceLineItems: {
                  ...input.serviceLineItems.map((e) => getNewItem(e)),
                },
              };
            });
          }
        });
      });

      const previousEntries =
        queryClient.getQueryData<GetContractsQuery>(queryKey);

      return { previousEntries };
    },
    onError: (_, __, context) => {
      queryClient.setQueryData<GetContractsQuery>(
        queryKey,
        context?.previousEntries,
      );
      toastError('Failed to update services', 'update-service-error');
    },
    onSuccess: (_, variables) => {
      updateTimelineCacheAfterServiceLineItemChange({
        timelineQueryKey,
        contractName,
        user: store.session.value.profile.name ?? '',
        updateTimelineCache,
        prevServiceLineItems: contractLineItems,
        newServiceLineItems: variables.input.serviceLineItems,
      });

      onClose();
    },
    onSettled: () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
      timeoutRef.current = setTimeout(() => {
        queryClient.invalidateQueries({ queryKey });
        queryClient.invalidateQueries({ queryKey: timelineQueryKey });
      }, 1000);
    },
  });
  const [services, setServices] = useState<Array<BulkUpdateServiceLineItem>>(
    [],
  );
  const [invoiceNote, setInvoiceNote] = useState<string>(`${notes}`);
  useEffect(() => {
    if (isOpen) {
      if (contractLineItems.length) {
        setServices(
          contractLineItems
            .filter((e) => !e.serviceEnded)
            .map((e) => {
              return ServiceLineItemsDTO.toPayload(e);
            }),
        );
      }
      if (!contractLineItems.length) {
        setServices([defaultValue]);
      }
    }
  }, [isOpen]);
  const handleAddService = () => {
    setServices([...services, defaultValue]);
  };

  const handleUpdateService = (
    index: number,
    updatedService: BulkUpdateServiceLineItem,
  ) => {
    const updatedServices = [...services];

    updatedServices[index] = updatedService;
    setServices(updatedServices);
  };

  const handleApplyChanges = async () => {
    updateServices.mutate({
      input: {
        contractId,
        invoiceNote,
        serviceLineItems: services
          .filter((e) => !e.isDeleted)
          .map((e) => ({
            serviceLineItemId: e?.serviceLineItemId ?? '',
            billed: e.billed,
            name: e.name,
            price: e.price,
            quantity: e.quantity,
            vatRate: e.vatRate,
            serviceStarted: e.serviceStarted,
          })),
      },
    });
  };

  return (
    <Modal open={isOpen} onOpenChange={onClose}>
      <ModalPortal>
        <ModalOverlay className='z-50' />

        <ModalContent className='min-w-[768px] rounded-2xl'>
          <ModalHeader>
            <FeaturedIcon
              size='lg'
              colorScheme='primary'
              className='mt-4 ml-[12px] mb-[30px]'
            >
              <DotSingle color='primary.700' />
            </FeaturedIcon>
            <h2 className='text-lg mt-4'>Modify contract service line items</h2>
          </ModalHeader>
          <ModalBody className='pb-0 flex flex-col flex-1'>
            <div className='flex justify-between items-center pr-5 border-b border-gray-300 pb-1'>
              <p className='text-sm font-medium w-[15%]'>Name</p>
              <p className='text-sm font-medium w-[15%]'>Type</p>
              <p className='text-sm font-medium w-[10%]'>Qty</p>
              <p className='text-sm font-medium w-[15%]'>Unit Price</p>
              <p className='text-sm font-medium w-[10%]'>Recurring</p>{' '}
              <p className='text-sm font-medium w-[10%]'>VAT</p>{' '}
              <p className='text-sm font-medium w-[15%]'>Service Start</p>
            </div>

            {services.map((service, index) => (
              <Fragment key={`service-line-item-${index}`}>
                <ServiceLineItemRow
                  service={service}
                  index={index}
                  currency={currency}
                  onChange={(data) => handleUpdateService(index, data)}
                  prevServiceLineItemData={contractLineItems.find(
                    (e) => e.metadata.id === service?.serviceLineItemId,
                  )}
                />
              </Fragment>
            ))}
            <div>
              <Button
                className='px-2 my-1 text-gray-500 font-base'
                leftIcon={<Plus />}
                variant='ghost'
                size='sm'
                onClick={handleAddService}
              >
                New item
              </Button>
            </div>

            <AutoresizeTextarea
              label='Note'
              className='overflow-ellipsis '
              labelProps={{
                className: 'text-sm mb-0 font-semibold',
              }}
              size='sm'
              name='invoiceNote'
              placeholder='Customer invoice note'
              value={invoiceNote}
              onChange={(event) => setInvoiceNote(event.target.value)}
            />
          </ModalBody>
          <ModalFooter className='flex p-6'>
            <Button
              variant='outline'
              className='w-full'
              onClick={onClose}
              isDisabled={updateServices.isPending}
            >
              Cancel
            </Button>
            <Button
              className='w-full ml-3'
              variant='outline'
              colorScheme='primary'
              isLoading={updateServices.isPending}
              onClick={handleApplyChanges}
            >
              {updateServices.isPending
                ? 'Applying changes...'
                : 'Apply changes'}
            </Button>
          </ModalFooter>
        </ModalContent>
      </ModalPortal>
    </Modal>
  );
};
