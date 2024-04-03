'use client';

import { useRef, useCallback } from 'react';
import { useParams } from 'next/navigation';
import { useForm } from 'react-inverted-form';

import set from 'lodash/set';
import { produce } from 'immer';
import { match } from 'ts-pattern';
import { useQueryClient } from '@tanstack/react-query';
import { OptionProps, chakraComponents } from 'chakra-react-select';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { SelectOption } from '@ui/utils/types';
import { Flag04 } from '@ui/media/icons/Flag04';
import { Heading } from '@ui/typography/Heading';
import { toastError } from '@ui/presentation/Toast';
import { Trophy01 } from '@ui/media/icons/Trophy01';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
import { FormSelect, SelectInstance } from '@ui/form/SyncSelect';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { OnboardingStatus, OnboardingDetails } from '@graphql/types';
import { useTimelineMeta } from '@organization/src/components/Timeline/state';
import { useInfiniteGetTimelineQuery } from '@organization/src/graphql/getTimeline.generated';
import { useUpdateOnboardingStatusMutation } from '@organization/src/graphql/updateOnboardingStatus.generated';
import {
  OrganizationQuery,
  useOrganizationQuery,
} from '@organization/src/graphql/organization.generated';
import {
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
  ModalCloseButton,
} from '@ui/overlay/Modal';

import { options } from './util';
import {
  OnboardingStatusDto,
  OnboardingStatusForm,
} from './OboardingStatus.dto';

interface OnboardingStatusModalProps {
  isOpen: boolean;
  onClose: () => void;
  data?: OnboardingDetails | null;
  onFetching?: (status: boolean) => void;
}

const formId = 'onboarding-status-update-form';

const getIconcolorScheme = (status: OnboardingStatus) =>
  match(status)
    .returnType<string>()
    .with(
      OnboardingStatus.Successful,
      OnboardingStatus.OnTrack,
      OnboardingStatus.Done,
      () => 'success',
    )
    .with(OnboardingStatus.Late, OnboardingStatus.Stuck, () => 'warning')
    .otherwise(() => 'gray');

const getIcon = (status: OnboardingStatus) => {
  const color = `${getIconcolorScheme(status)}.500`;

  return match(status)
    .with(OnboardingStatus.Successful, () => <Trophy01 color={color} mr='3' />)
    .otherwise(() => <Flag04 color={color} mr='3' />);
};

export const OnboardingStatusModal = ({
  data,
  isOpen,
  onClose,
  onFetching,
}: OnboardingStatusModalProps) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const id = useParams()?.id as string;
  const timeout = useRef<NodeJS.Timeout>();
  const queryKey = useOrganizationQuery.getKey({ id });
  const initialFocusRef = useRef<SelectInstance>(null);

  const [timelineMeta] = useTimelineMeta();
  const timelineQueryKey = useInfiniteGetTimelineQuery.getKey(
    timelineMeta.getTimelineVariables,
  );

  const updateOnboardingStatus = useUpdateOnboardingStatusMutation(client, {
    onMutate: ({ input }) => {
      onFetching?.(true);
      queryClient.cancelQueries({ queryKey });

      const { previousEntries } = useOrganizationQuery.mutateCacheEntry(
        queryClient,
        { id },
      )((currentCache) => {
        return produce(currentCache, (draft) => {
          if (!draft) return;
          set(draft, 'organization.accountDetails.onboarding', {
            ...input,
            updatedAt: new Date().toISOString(),
          });
        });
      });

      return { previousEntries };
    },
    onError: (_, __, context) => {
      queryClient.setQueryData<OrganizationQuery>(
        queryKey,
        context?.previousEntries,
      );
      toastError(
        `We couldn't update the onboarding status`,
        `${id}-onboarding-status-update-error`,
      );
      onFetching?.(false);
    },
    onSettled: () => {
      onClose();
      if (timeout.current) clearTimeout(timeout.current);

      timeout.current = setTimeout(() => {
        queryClient.invalidateQueries({ queryKey });
        queryClient.invalidateQueries({ queryKey: timelineQueryKey });
        onFetching?.(false);
      }, 500);
    },
  });

  const onSubmit = useCallback(
    async (values: OnboardingStatusForm) => {
      updateOnboardingStatus.mutate({
        input: OnboardingStatusDto.toPayload({ id, ...values }),
      });
    },
    [updateOnboardingStatus.mutate],
  );

  const defaultValues = OnboardingStatusDto.toForm(data);
  const { state, handleSubmit } = useForm<OnboardingStatusForm>({
    formId,
    defaultValues,
    onSubmit,
    stateReducer: (_, action, next) => {
      if (action.type === 'HAS_SUBMITTED') {
        return { ...next, values: { ...next.values, comments: '' } };
      }

      return next;
    },
  });

  return (
    <Modal
      closeOnEsc
      isOpen={isOpen}
      onClose={onClose}
      initialFocusRef={initialFocusRef}
    >
      <ModalOverlay />
      <ModalContent
        as='form'
        onSubmit={handleSubmit}
        borderRadius='2xl'
        backgroundImage='/backgrounds/organization/circular-bg-pattern.png'
        backgroundRepeat='no-repeat'
        sx={{
          backgroundPositionX: '1px',
          backgroundPositionY: '-7px',
        }}
      >
        <ModalCloseButton />
        <ModalHeader>
          <FeaturedIcon
            size='lg'
            colorScheme={getIconcolorScheme(
              state?.values?.status?.value ?? OnboardingStatus.NotApplicable,
            )}
          >
            {state?.values?.status?.value === OnboardingStatus.Successful ? (
              <Trophy01 />
            ) : (
              <Flag04 />
            )}
          </FeaturedIcon>
          <Heading fontSize='lg' mt='4'>
            Update onboarding status
          </Heading>
        </ModalHeader>
        <ModalBody pb='0' gap={4} as={Flex} flexDir='column'>
          <FormSelect
            name='status'
            label='Status'
            isLabelVisible
            formId={formId}
            options={options}
            openMenuOnFocus
            ref={initialFocusRef}
            isDisabled={updateOnboardingStatus.isPending}
            components={{
              Option: ({
                data,
                children,
                ...rest
              }: OptionProps<SelectOption<OnboardingStatus>>) => {
                const icon = getIcon(data.value);

                return (
                  <chakraComponents.Option data={data} {...rest}>
                    {icon}
                    {children}
                  </chakraComponents.Option>
                );
              },
            }}
          />
          {defaultValues.status.value !== state?.values?.status?.value && (
            <Box>
              <Text as='label' htmlFor='reason' fontSize='sm'>
                <b>Reason for change</b> (optional)
              </Text>
              <FormAutoresizeTextarea
                pt='0'
                formId={formId}
                name='comments'
                spellCheck='false'
                isDisabled={updateOnboardingStatus.isPending}
                placeholder={`What is the reason for changing the onboarding status?`}
              />
            </Box>
          )}
        </ModalBody>
        <ModalFooter p='6'>
          <Button
            w='full'
            variant='outline'
            onClick={onClose}
            isDisabled={updateOnboardingStatus.isPending}
          >
            Cancel
          </Button>
          <Button
            ml='3'
            w='full'
            type='submit'
            variant='outline'
            colorScheme='primary'
            isLoading={updateOnboardingStatus.isPending}
          >
            Update status
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
