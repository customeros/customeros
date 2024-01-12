'use client';
import { useRef, useEffect } from 'react';
import { useParams } from 'next/navigation';
import { useForm } from 'react-inverted-form';

import { produce } from 'immer';
import { useSession } from 'next-auth/react';
import { useQueryClient } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { Heading } from '@ui/typography/Heading';
import { toastError } from '@ui/presentation/Toast';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
import { FormCheckbox } from '@ui/form/Checkbox/FormCheckbox';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { Action, BilledType, ServiceLineItem } from '@graphql/types';
import { useTimelineMeta } from '@organization/src/components/Timeline/shared/state';
import { useUpdateServiceMutation } from '@organization/src/graphql/updateService.generated';
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
import { getUpdateServiceEvents } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/Services/modals/utils';
import { OneTimeServiceForm } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/Services/modals/OneTimeServiceForm';
import { RecurringServiceFrom } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/Services/modals/RecurringService';
import {
  ServiceDTO,
  ServiceForm,
} from '@organization/src/components/Tabs/panels/AccountPanel/Contract/Services/modals/Service.dto';

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
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const formId = `update-service-item`;
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const queryKey = useGetContractsQuery.getKey({ id });
  const defaultValues = ServiceDTO.toForm(data);
  const updateTimelineCache = useUpdateCacheWithNewEvent();
  const [timelineMeta] = useTimelineMeta();
  const timelineQueryKey = useInfiniteGetTimelineQuery.getKey(
    timelineMeta.getTimelineVariables,
  );
  const session = useSession();

  const updateService = useUpdateServiceMutation(client, {
    onMutate: ({ input }) => {
      queryClient.cancelQueries({ queryKey });

      queryClient.setQueryData<GetContractsQuery>(queryKey, (currentCache) => {
        return produce(currentCache, (draft) => {
          if (draft?.['organization']?.['contracts']) {
            draft['organization']['contracts']?.map((contractData) => {
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
    onSuccess: (_, variables) => {
      if (data) {
        getUpdateServiceEvents(
          data,
          variables.input,
          session?.data?.user?.name ?? '',
          (event: Action) => updateTimelineCache(event, timelineQueryKey),
        );
      }

      onClose();
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
  const { state, setDefaultValues } = useForm<ServiceForm>({
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
    defaultValues.isRetroactiveCorrection,
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
          <FeaturedIcon size='lg' colorScheme='gray'>
            <DotSingle color='primary.600' />
          </FeaturedIcon>
          <Heading fontSize='lg' mt='4'>
            Modify this
            {data?.billed === BilledType.Once && ' one-time service'}
            {data?.billed === BilledType.Usage && ' service'}
            {data?.billed !== BilledType.Once &&
              data?.billed !== BilledType.Usage &&
              ' recurring service'}
          </Heading>
        </ModalHeader>
        <ModalBody pb='0'>
          {data?.billed === BilledType.Once ||
          data?.billed === BilledType.Usage ? (
            <OneTimeServiceForm formId={formId} billedType={data.billed} />
          ) : (
            <RecurringServiceFrom formId={formId} />
          )}
          <Flex
            p={2}
            my={2}
            border='1px solid'
            borderColor='gray.100'
            bg='gray.25'
            borderRadius='md'
            alignItems='center'
          >
            <FormCheckbox
              formId={formId}
              name='isRetroactiveCorrection'
              fontSize='sm'
            >
              This is a retroactive correction
            </FormCheckbox>
          </Flex>

          <div>
            <Text as='label' htmlFor='reason' fontSize='sm'>
              <b>Reason for change</b> (optional)
            </Text>
            <FormAutoresizeTextarea
              pt='0'
              formId={formId}
              name='reason'
              id='reason'
              spellCheck='false'
              placeholder={`What's this modification about?`}
            />
          </div>
        </ModalBody>
        <ModalFooter p='6'>
          <Button variant='outline' w='full' onClick={onClose}>
            Cancel
          </Button>
          <Button
            ml='3'
            w='full'
            isLoading={updateService.status === 'loading'}
            loadingText='Modifying...'
            variant='outline'
            colorScheme='primary'
            onClick={updateServiceData}
          >
            Modify
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
