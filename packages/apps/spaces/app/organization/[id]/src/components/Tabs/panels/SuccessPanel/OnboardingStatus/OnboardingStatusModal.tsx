'use client';

import { useParams } from 'next/navigation';
import { useForm } from 'react-inverted-form';

import { match } from 'ts-pattern';
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
import { OnboardingStatus, OnboardingDetails } from '@graphql/types';
import { useUpdateOnboardingStatusMutation } from '@organization/src/graphql/updateOnboardingStatus.generated';
import {
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
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
  const id = useParams()?.id as string;
  const updateOnboardingStatus = useUpdateOnboardingStatusMutation(client);

  const defaultValues = OnboardingStatusDto.toForm(data);
  const { state, handleSubmit } = useForm<OnboardingStatusForm>({
    formId,
    defaultValues,
    onSubmit: async (values) => {
      updateOnboardingStatus.mutate({
        input: OnboardingStatusDto.toPayload({ id, ...values }),
      });
    },
  });

  const icon = getIcon(state?.values?.status?.value);

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
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
            leftElement={icon}
            isDisabled={updateOnboardingStatus.isLoading}
            components={{
              Option: ({
                data,
                children,
                ...props
              }: OptionProps<SelectOption<OnboardingStatus>>) => {
                const icon = getIcon(data.value);

                return (
                  <chakraComponents.Option data={data} {...props}>
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
