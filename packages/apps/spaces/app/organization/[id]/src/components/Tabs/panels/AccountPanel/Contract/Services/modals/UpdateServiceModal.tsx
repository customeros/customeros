'use client';
import { useRef, useEffect } from 'react';
import { useParams } from 'next/navigation';
import { useForm } from 'react-inverted-form';

import { produce } from 'immer';
import { useQueryClient } from '@tanstack/react-query';

import { Button } from '@ui/form/Button';
import { FeaturedIcon } from '@ui/media/Icon';
import { Heading } from '@ui/typography/Heading';
import { toastError } from '@ui/presentation/Toast';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { BilledType, ServiceLineItem } from '@graphql/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useUpdateServiceMutation } from '@organization/src/graphql/updateService.generated';
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
import { OneTimeServiceForm } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/Services/modals/OneTimeServiceForm';
import {
  ServiceDTO,
  ServiceForm,
} from '@organization/src/components/Tabs/panels/AccountPanel/Contract/Services/modals/Service.dto';
import { SubscriptionServiceFrom } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/Services/modals/SubscriptionServiceForm';

interface SubscriptionServiceModalProps {
  isOpen: boolean;
  onClose: () => void;
  data?: ServiceLineItem;
}

export const UpdateServiceModal = ({
  data,
  isOpen,
  onClose,
}: SubscriptionServiceModalProps) => {
  const id = useParams()?.id as string;
  const initialRef = useRef(null);
  const formId = `update-service-item`;
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const queryKey = useGetContractsQuery.getKey({ id });
  const defaultValues = ServiceDTO.toForm(data);
  const updateService = useUpdateServiceMutation(client, {
    onMutate: ({ input }) => {
      queryClient.cancelQueries({ queryKey });
      queryClient.setQueryData<GetContractsQuery>(queryKey, (currentCache) => {
        return produce(currentCache, (draft) => {
          if (draft?.['organization']?.['contracts']) {
            draft['organization']['contracts']?.map((contractData, index) => {
              return (contractData.serviceLineItems ?? []).map(
                (serviceItem) => {
                  if (serviceItem.id === input.serviceLineItemId) {
                    const { serviceLineItemId, ...rest } = input;

                    return {
                      ...serviceItem,
                      ...rest,
                    };
                  }

                  return serviceItem;
                },
              );
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
      toastError('Failed to update service', 'update-service-error');
    },
    onSuccess: () => {
      onClose();
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
  }, [
    defaultValues.renewalCycle,
    defaultValues.billed,
    defaultValues.appSource,
    defaultValues.quantity,
    defaultValues.serviceStartedAt,
    defaultValues.externalReference,
  ]);

  const updateServiceData = () => {
    if (!data?.id) return;
    updateService.mutate({
      input: {
        ...ServiceDTO.toUpdatePayload(state.values),
        serviceLineItemId: data.id,
      },
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
            {data?.billed === BilledType.Once
              ? 'Update one-time service'
              : 'Update subscription service'}
          </Heading>
        </ModalHeader>
        <ModalBody pb='0'>
          {data?.billed === BilledType.Once ? (
            <OneTimeServiceForm formId={formId} />
          ) : (
            <SubscriptionServiceFrom formId={formId} />
          )}
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
            onClick={updateServiceData}
          >
            Update
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
