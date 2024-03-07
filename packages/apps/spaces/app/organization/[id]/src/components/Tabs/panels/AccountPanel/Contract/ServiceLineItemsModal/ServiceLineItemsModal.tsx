'use client';
import { useParams } from 'next/navigation';
import React, { useRef, useState, Fragment, useEffect } from 'react';

import { produce } from 'immer';
import { useSession } from 'next-auth/react';
import { useQueryClient } from '@tanstack/react-query';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { Plus } from '@ui/media/icons/Plus';
import { FeaturedIcon } from '@ui/media/Icon';
// import { Grid, GridItem } from '@ui/layout/Grid';
import { Heading } from '@ui/typography/Heading';
import { toastError } from '@ui/presentation/Toast';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { AutoresizeTextarea } from '@ui/form/Textarea';
// import { Invoice } from '@shared/components/Invoice/Invoice';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useTimelineMeta } from '@organization/src/components/Timeline/shared/state';
import { useInfiniteGetTimelineQuery } from '@organization/src/graphql/getTimeline.generated';
import { useUpdateServicesMutation } from '@organization/src/graphql/updateServiceLineItems.generated';
import {
  GetContractsQuery,
  useGetContractsQuery,
} from '@organization/src/graphql/getContracts.generated';
import { useUpdateCacheWithNewEvent } from '@organization/src/components/Timeline/hooks/updateCacheWithNewEvent';
import {
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal';
import {
  BilledType,
  DataSource,
  InputMaybe,
  // InvoiceLine,
  ServiceLineItem,
  ServiceLineItemBulkUpdateItem,
} from '@graphql/types';
import { updateTimelineCacheAfterServiceLineItemChange } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/ServiceLineItemsModal/utils';
import {
  ServiceLineItemsDTO,
  BulkUpdateServiceLineItem,
} from '@organization/src/components/Tabs/panels/AccountPanel/Contract/ServiceLineItemsModal/ServiceLineItemsModal.dto';

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
  const initialRef = useRef(null);
  const client = getGraphQLClient();
  const id = useParams()?.id as string;
  const [timelineMeta] = useTimelineMeta();
  const queryKey = useGetContractsQuery.getKey({ id });
  const timelineQueryKey = useInfiniteGetTimelineQuery.getKey(
    timelineMeta.getTimelineVariables,
  );
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const queryClient = useQueryClient();
  const session = useSession();
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
                (contract) => contract.id === input.contractId,
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
    onError: (err, __, context) => {
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
        user: session?.data?.user?.name ?? '',
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
            // vatRate: e.vatRate,
            serviceStarted: e.serviceStarted,
          })),
      },
    });
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      initialFocusRef={initialRef}
      size='2xl'
      closeOnOverlayClick
    >
      <ModalOverlay />
      <ModalContent borderRadius='2xl'>
        {/*<Grid h='100%' templateColumns='1fr' gap={4} overflow='scroll'>*/}
        {/*545px 1fr*/}
        {/*<GridItem*/}
        {/*  rowSpan={1}*/}
        {/*  colSpan={1}*/}
        {/*  h='100%'*/}
        {/*  display='flex'*/}
        {/*  flexDir='column'*/}
        {/*  justifyContent='space-between'*/}
        {/*  bg='gray.25'*/}
        {/*  borderRight='1px solid'*/}
        {/*  borderColor='gray.200'*/}
        {/*  borderRadius='2xl'*/}
        {/*  // borderTopLeftRadius='2xl'*/}
        {/*  // borderBottomLeftRadius='2xl'*/}
        {/*  backgroundImage='/backgrounds/organization/circular-bg-pattern.png'*/}
        {/*  backgroundRepeat='no-repeat'*/}
        {/*  sx={{*/}
        {/*    backgroundPositionX: '1px',*/}
        {/*    backgroundPositionY: '-7px',*/}
        {/*  }}*/}
        {/*>*/}
        <ModalHeader>
          <FeaturedIcon size='lg' colorScheme='primary'>
            <DotSingle color='primary.700' />
          </FeaturedIcon>
          <Heading fontSize='lg' mt='4'>
            Modify contract service line items
          </Heading>
        </ModalHeader>
        <ModalBody pb='0' display='flex' flexDir='column' flex={1}>
          <Flex
            justifyContent='space-between'
            alignItems='center'
            pr='20px'
            borderBottom='1px solid'
            borderColor='gray.300'
            pb={1}
          >
            <Text fontSize='sm' fontWeight='medium' w='20%'>
              Name
            </Text>
            <Text fontSize='sm' fontWeight='medium' w='15%'>
              Type
            </Text>
            <Text fontSize='sm' fontWeight='medium' w='10%'>
              Qty
            </Text>
            <Text fontSize='sm' fontWeight='medium' w='15%'>
              Unit Price
            </Text>
            <Text fontSize='sm' fontWeight='medium' w='15%'>
              Recurring
            </Text>{' '}
            <Text fontSize='sm' fontWeight='medium' w='15%'>
              Service Start
            </Text>
          </Flex>

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
          <Box>
            <Button
              leftIcon={<Plus />}
              variant='ghost'
              size='sm'
              px={2}
              my={1}
              color='gray.500'
              fontWeight='regular'
              onClick={handleAddService}
            >
              New item
            </Button>
          </Box>

          <AutoresizeTextarea
            label='Note'
            isLabelVisible
            labelProps={{
              fontSize: 'sm',
              mb: 0,
              fontWeight: 'semibold',
            }}
            name='invoiceNote'
            textOverflow='ellipsis'
            placeholder='Customer invoice note'
            value={invoiceNote}
            onChange={(event) => setInvoiceNote(event.target.value)}
          />
        </ModalBody>
        <ModalFooter p='6'>
          <Button
            variant='outline'
            w='full'
            onClick={onClose}
            isDisabled={updateServices.isPending}
          >
            Cancel
          </Button>
          <Button
            ml='3'
            w='full'
            variant='outline'
            colorScheme='primary'
            loadingText='Applying changes...'
            isLoading={updateServices.isPending}
            onClick={handleApplyChanges}
          >
            Apply changes
          </Button>
        </ModalFooter>
        {/*</GridItem>*/}
        {/*<GridItem pr={3}>*/}
        {/*  <Invoice*/}
        {/*    isDraft*/}
        {/*    tax={10}*/}
        {/*    note={note}*/}
        {/*    from={{*/}
        {/*      address: '',*/}
        {/*      address2: '',*/}
        {/*      city: '',*/}
        {/*      country: '',*/}
        {/*      name: '',*/}
        {/*      zip: '',*/}
        {/*      email: '',*/}
        {/*    }}*/}
        {/*    total={1400}*/}
        {/*    dueDate={new Date()}*/}
        {/*    subtotal={2002}*/}
        {/*    lines={services as unknown as InvoiceLine[]}*/}
        {/*    issueDate='10.01.2024'*/}
        {/*    billedTo={{*/}
        {/*      address: '',*/}
        {/*      address2: '',*/}
        {/*      city: '',*/}
        {/*      country: '',*/}
        {/*      name: '',*/}
        {/*      zip: '',*/}
        {/*      email: '',*/}
        {/*    }}*/}
        {/*    invoiceNumber='INV-001'*/}
        {/*  />*/}
        {/*</GridItem>*/}
        {/*</Grid>*/}
      </ModalContent>
    </Modal>
  );
};
