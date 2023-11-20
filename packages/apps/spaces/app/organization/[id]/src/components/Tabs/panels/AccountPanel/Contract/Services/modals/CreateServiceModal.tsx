'use client';
import { useRef, useEffect } from 'react';
import { useParams } from 'next/navigation';
import { useForm } from 'react-inverted-form';

import { produce } from 'immer';
import { useQueryClient } from '@tanstack/react-query';

import { Button } from '@ui/form/Button';
import { DataSource } from '@graphql/types';
import { FeaturedIcon } from '@ui/media/Icon';
import { Heading } from '@ui/typography/Heading';
import { toastError } from '@ui/presentation/Toast';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { Tab, Tabs, TabList, TabPanel, TabPanels } from '@ui/disclosure/Tabs';
import { useCreateServiceMutation } from '@organization/src/graphql/createService.generated';
import { billedTypeOptions } from '@organization/src/components/Tabs/panels/AccountPanel/utils';
import {
  GetContractsQuery,
  useGetContractsQuery,
} from '@organization/src/graphql/getContracts.generated';
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
import {
  ServiceDTO,
  ServiceForm,
} from '@organization/src/components/Tabs/panels/AccountPanel/Contract/Services/modals/Service.dto';
import { SubscriptionServiceFrom } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/Services/modals/SubscriptionServiceForm';

interface SubscriptionServiceModalProps {
  isOpen: boolean;
  contractId: string;
  onClose: () => void;
}

export const CreateServiceModal = ({
  isOpen,
  onClose,
  contractId,
}: SubscriptionServiceModalProps) => {
  const initialRef = useRef(null);
  const formId = `create-service-item`;
  const defaultValues = ServiceDTO.toForm();

  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const id = useParams()?.id as string;
  const { modal } = useAddServiceModalContext();
  const queryKey = useGetContractsQuery.getKey({ id });

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
      toastError('Failed to update contract', 'update-contract-error');
    },
    onSuccess: () => {
      modal.onClose();
    },
    onSettled: () => {
      queryClient.invalidateQueries(queryKey);
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

  const handleCreateService = () => {
    createService.mutate({
      input: { ...ServiceDTO.toPayload(state.values, contractId) },
    });
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
              borderRadius='md'
            >
              <Tab
                borderTopLeftRadius='md'
                borderBottomLeftRadius='md'
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
                  setDefaultValues({
                    ...defaultValues,
                  });
                }}
              >
                Subscription
              </Tab>
              <Tab
                borderTopRightRadius='md'
                borderBottomRightRadius='md'
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
                  setDefaultValues({
                    ...defaultValues,
                    billed: billedTypeOptions[0],
                  });
                }}
              >
                One-time
              </Tab>
            </TabList>

            <TabPanels>
              <TabPanel px={0} pb={2}>
                <SubscriptionServiceFrom formId={formId} />
              </TabPanel>
              <TabPanel px={0} pb={2}>
                <OneTimeServiceForm formId={formId} />
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
            onClick={handleCreateService}
          >
            Create
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
