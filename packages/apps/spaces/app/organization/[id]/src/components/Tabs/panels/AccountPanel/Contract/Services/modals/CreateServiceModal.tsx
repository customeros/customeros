'use client';
import { useParams } from 'next/navigation';
import { useForm } from 'react-inverted-form';
import { useRef, useState, useEffect } from 'react';

import { produce } from 'immer';
import { useSession } from 'next-auth/react';
import { useQueryClient } from '@tanstack/react-query';

import { Button } from '@ui/form/Button';
import { FeaturedIcon } from '@ui/media/Icon';
import { Heading } from '@ui/typography/Heading';
import { toastError } from '@ui/presentation/Toast';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { ActionType, BilledType, DataSource } from '@graphql/types';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { Tab, Tabs, TabList, TabPanel, TabPanels } from '@ui/disclosure/Tabs';
import { useTimelineMeta } from '@organization/src/components/Timeline/shared/state';
import { useCreateServiceMutation } from '@organization/src/graphql/createService.generated';
import { useInfiniteGetTimelineQuery } from '@organization/src/graphql/getTimeline.generated';
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
import { useAddServiceModalContext } from '@organization/src/components/Tabs/panels/AccountPanel/context/AccountModalsContext';
import { OneTimeServiceForm } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/Services/modals/OneTimeServiceForm';
import { RecurringServiceFrom } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/Services/modals/RecurringService';
import {
  ServiceDTO,
  ServiceForm,
} from '@organization/src/components/Tabs/panels/AccountPanel/Contract/Services/modals/Service.dto';

interface SubscriptionServiceModalProps {
  isOpen: boolean;
  contractId: string;
  onClose: () => void;
  contractName: string;
}

export const CreateServiceModal = ({
  isOpen,
  onClose,
  contractId,
  contractName,
}: SubscriptionServiceModalProps) => {
  const initialRef = useRef(null);
  const formId = `create-service-item-${contractId}`;
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const defaultValues = ServiceDTO.toForm();
  const [activeTab, setActiveTab] = useState('RECURRING');
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const id = useParams()?.id as string;
  const { modal } = useAddServiceModalContext();
  const queryKey = useGetContractsQuery.getKey({ id });
  const updateTimelineCache = useUpdateCacheWithNewEvent();
  const [timelineMeta] = useTimelineMeta();
  const timelineQueryKey = useInfiniteGetTimelineQuery.getKey(
    timelineMeta.getTimelineVariables,
  );
  const session = useSession();
  const createService = useCreateServiceMutation(client, {
    onMutate: ({ input }) => {
      queryClient.cancelQueries({ queryKey });
      queryClient.setQueryData<GetContractsQuery>(queryKey, (currentCache) => {
        return produce(currentCache, (draft) => {
          const previousContracts = draft?.['organization']?.['contracts'];
          const updatedContractIndex = previousContracts?.findIndex(
            (contract) => contract.id === input.contractId,
          );
          if (!draft) return;

          const newItem = {
            id: Math.random().toString(),
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
            name: input.name,
            billed: input.billed,
            price: input.price,
            quantity: input.quantity,
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
                amount: input.price,
                maxAmount: input.price,
                name: '',
                renewalLikelihood: 'HIGH',
                renewalUpdatedByUserId: '',
                renewalUpdatedByUserAt: new Date().toISOString(),
                renewedAt: new Date().toISOString(),
              },
            ],
          };

          if (draft?.['organization']?.['contracts']) {
            draft['organization']['contracts']?.map((contractData, index) => {
              if (index !== updatedContractIndex) {
                return contractData;
              }

              return {
                ...contractData,
                serviceLineItems: [
                  ...(contractData.serviceLineItems ?? []),
                  newItem,
                ],
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
      toastError('Failed to create service', 'update-service-error');
    },
    onSuccess: (_, variables) => {
      modal.onClose();
      const isRecurring = [
        BilledType.Annually,
        BilledType.Monthly,
        BilledType.Quarterly,
      ].includes(variables?.input?.billed as BilledType);
      const metadata = JSON.stringify({
        price: variables?.input?.price,
        billedType: variables?.input?.billed,
      });
      const user = session?.data?.user?.name ?? '';
      const actionType = isRecurring
        ? ActionType.ServiceLineItemBilledTypeRecurringCreated
        : variables?.input?.billed === BilledType.Usage
        ? ActionType.ServiceLineItemBilledTypeUsageCreated
        : ActionType.ServiceLineItemBilledTypeOnceCreated;
      updateTimelineCache(
        {
          __typename: 'Action',
          id: Math.random().toString(),
          createdAt: new Date(),
          updatedAt: '',
          actionType,
          appSource: 'customeros-optimistic-update',
          source: 'customeros-optimistic-update',
          metadata,
          actionCreatedBy: {
            firstName: user,
            lastName: '',
          },
          content: `${user} added a ${
            isRecurring
              ? 'recurring'
              : variables?.input.billed === BilledType.Usage
              ? 'use based'
              : 'one-time'
          } service to ${contractName}: ${
            variables.input.name
          } , at ${formatCurrency(variables.input.price ?? 0)}`,
        },
        timelineQueryKey,
      );
    },
    onSettled: () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
      timeoutRef.current = setTimeout(() => {
        queryClient.invalidateQueries(queryKey);
        queryClient.invalidateQueries(timelineQueryKey);
      }, 1000);
    },
  });
  const { setDefaultValues, state } = useForm<ServiceForm>({
    formId,
    defaultValues,
    stateReducer: (_, action, next) => {
      return next;
    },
  });
  useEffect(() => {
    setDefaultValues(defaultValues);
  }, [isOpen]);

  const determineBillingDetails = () => {
    if (activeTab !== 'RECURRING') {
      const billedTypeValue =
        activeTab === BilledType.Once ? BilledType.Once : BilledType.Usage;

      return {
        quantity: undefined,
        billed: {
          value: billedTypeValue,
          label: '',
        },
      };
    }

    return {};
  };

  const handleServiceCreation = () => {
    const billingDetails = determineBillingDetails();

    const serviceInputPayload = ServiceDTO.toPayload(
      { ...state.values, ...billingDetails },
      contractId,
    );

    createService.mutate({ input: serviceInputPayload });
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} initialFocusRef={initialRef}>
      <ModalOverlay />
      <ModalContent
        borderRadius='2xl'
        backgroundImage='/backgrounds/organization/circular-bg-pattern.png'
        backgroundRepeat='no-repeat'
        sx={{
          backgroundPositionX: '1px',
          backgroundPositionY: '-7px',
        }}
      >
        <ModalHeader>
          <FeaturedIcon size='lg' colorScheme='primary'>
            <DotSingle color='primary.600' />
          </FeaturedIcon>
          <Heading fontSize='lg' mt='4'>
            Add a new service
          </Heading>
        </ModalHeader>
        <ModalBody pb='0'>
          <Tabs isFitted>
            <TabList
              bg='white'
              border='1px solid'
              borderColor='gray.300'
              borderRadius='8px'
            >
              <Tab
                borderTopLeftRadius='8px'
                borderBottomLeftRadius='8px'
                borderBottom='none'
                flex={1}
                borderRight='1px solid'
                borderRightColor='gray.300 !important'
                color='gray.500'
                bg='gray.50'
                mb={0}
                _selected={{
                  color: 'gray.500',
                  bg: 'white',
                  fontWeight: 'semibold',
                }}
                onClick={() => {
                  setActiveTab('RECURRING');
                }}
              >
                Recurring
              </Tab>
              <Tab
                borderRight='1px solid'
                borderRightColor='gray.300 !important'
                borderBottom='none'
                flex={1}
                mb={0}
                color='gray.500'
                bg='gray.50'
                _selected={{
                  color: 'gray.500',
                  bg: 'white',
                  fontWeight: 'semibold',
                }}
                onClick={() => {
                  setActiveTab('USAGE');
                }}
              >
                Per use
              </Tab>
              <Tab
                borderTopRightRadius='8px'
                borderBottomRightRadius='8px'
                borderRadius='md'
                borderBottom='none'
                flex={1}
                mb={0}
                color='gray.500'
                bg='gray.50'
                _selected={{
                  color: 'gray.500',
                  bg: 'white',
                  fontWeight: 'semibold',
                }}
                onClick={() => {
                  setActiveTab('ONCE');
                }}
              >
                One-time
              </Tab>
            </TabList>

            <TabPanels>
              <TabPanel px={0} pb={2}>
                <RecurringServiceFrom formId={formId} />
              </TabPanel>
              <TabPanel px={0} pb={2}>
                <OneTimeServiceForm
                  formId={formId}
                  billedType={BilledType.Usage}
                />
              </TabPanel>
              <TabPanel px={0} pb={2}>
                <OneTimeServiceForm
                  formId={formId}
                  billedType={BilledType.Once}
                />
              </TabPanel>
            </TabPanels>
          </Tabs>
        </ModalBody>
        <ModalFooter p='6'>
          <Button variant='outline' w='full' onClick={onClose}>
            Cancel
          </Button>
          <Button
            ml='3'
            w='full'
            variant='outline'
            colorScheme='primary'
            isLoading={createService.status === 'loading'}
            loadingText='Creating...'
            onClick={handleServiceCreation}
          >
            Create
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
