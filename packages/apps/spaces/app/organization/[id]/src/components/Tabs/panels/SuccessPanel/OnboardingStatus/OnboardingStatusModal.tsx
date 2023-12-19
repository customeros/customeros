'use client';

import { useParams } from 'next/navigation';
import { useForm } from 'react-inverted-form';

import set from 'lodash/set';
import { produce } from 'immer';
import { match } from 'ts-pattern';
import { useQueryClient } from '@tanstack/react-query';
import { OptionProps, chakraComponents } from 'chakra-react-select';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { SelectOption } from '@ui/utils/types';
import { Flag04 } from '@ui/media/icons/Flag04';
import { FormSelect } from '@ui/form/SyncSelect';
import { Heading } from '@ui/typography/Heading';
import { Trophy01 } from '@ui/media/icons/Trophy01';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { toastError, toastSuccess } from '@ui/presentation/Toast';
import { OnboardingStatus, OnboardingDetails } from '@graphql/types';
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
}

const formId = 'onboarding-status-update-form';

const getIcon = (status: OnboardingStatus) => {
  const color = match(status)
    .returnType<string>()
    .with(
      OnboardingStatus.Successful,
      OnboardingStatus.OnTrack,
      OnboardingStatus.Done,
      () => 'success.500',
    )
    .with(OnboardingStatus.Late, OnboardingStatus.Stuck, () => 'warning.500')
    .otherwise(() => 'gray.500');

  return match(status)
    .with(OnboardingStatus.Successful, () => <Trophy01 color={color} mr='3' />)
    .otherwise(() => <Flag04 color={color} mr='3' />);
};

export const OnboardingStatusModal = ({
  data,
  isOpen,
  onClose,
}: OnboardingStatusModalProps) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const id = useParams()?.id as string;
  const queryKey = useOrganizationQuery.getKey({ id });
  const updateOnboardingStatus = useUpdateOnboardingStatusMutation(client, {
    onMutate: ({ input }) => {
      queryClient.cancelQueries({ queryKey });

      const previousEntries =
        queryClient.getQueryData<OrganizationQuery>(queryKey);
      queryClient.setQueryData<OrganizationQuery>(queryKey, (currentCache) => {
        return produce(currentCache, (draft) => {
          if (!draft) return;
          set<OrganizationQuery>(
            draft,
            'organization.accountDetails.onboarding',
            { ...input, updatedAt: new Date().toISOString() },
          );
        });
      });

      return { previousEntries };
    },
    onSuccess: () => {
      toastSuccess(
        'Onboarding status updated',
        `${id}-onboarding-status-update`,
      );
    },
    onError: (_, __, context) => {
      queryClient.setQueryData<OrganizationQuery>(
        queryKey,
        context?.previousEntries,
      );
      toastError(
        'Failed to update onboarding status',
        `${id}-onboarding-status-update-error`,
      );
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey });
      onClose();
    },
  });

  const defaultValues = OnboardingStatusDto.toForm(data);
  const { handleSubmit } = useForm<OnboardingStatusForm>({
    formId,
    defaultValues,
    onSubmit: async (values) => {
      updateOnboardingStatus.mutate({
        input: OnboardingStatusDto.toPayload({ id, ...values }),
      });
    },
    stateReducer: (_, action, next) => {
      if (action.type === 'HAS_SUBMITTED') {
        return { ...next, values: { ...next.values, comments: '' } };
      }

      return next;
    },
  });

  return (
    <Modal isOpen={isOpen} onClose={onClose} closeOnEsc>
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
          <FeaturedIcon size='lg' colorScheme='success'>
            <Flag04 />
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
            leftElement={<Flag04 color='gray.500' mr='3' />}
            isDisabled={updateOnboardingStatus.isLoading}
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
          <div>
            <Text as='label' htmlFor='reason' fontSize='sm'>
              <b>Reason for change</b> (optional)
            </Text>
            <FormAutoresizeTextarea
              pt='0'
              formId={formId}
              name='comments'
              spellCheck='false'
              isDisabled={updateOnboardingStatus.isLoading}
              placeholder={`What is the reason for changing the onboarding status?`}
            />
          </div>
        </ModalBody>
        <ModalFooter p='6'>
          <Button
            w='full'
            variant='outline'
            onClick={onClose}
            isDisabled={updateOnboardingStatus.isLoading}
          >
            Cancel
          </Button>
          <Button
            ml='3'
            w='full'
            type='submit'
            variant='outline'
            colorScheme='primary'
            isLoading={updateOnboardingStatus.isLoading}
          >
            Update status
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
