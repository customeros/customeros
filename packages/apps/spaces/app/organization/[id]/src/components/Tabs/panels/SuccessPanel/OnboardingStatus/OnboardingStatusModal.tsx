'use client';

import { useForm } from 'react-inverted-form';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { SelectOption } from '@ui/utils/types';
import { Flag04 } from '@ui/media/icons/Flag04';
import { FormSelect } from '@ui/form/SyncSelect';
import { Heading } from '@ui/typography/Heading';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
import {
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal';

interface OnboardingStatusModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const formId = 'onboarding-status-update-form';
const options: SelectOption[] = [
  { label: 'Not applicable', value: 'NOT_APPLICABLE' },
  { label: 'Not started', value: 'NOT_STARTED' },
  { label: 'On track', value: 'ON_TRACK' },
  { label: 'Late', value: 'LATE' },
  { label: 'Stuck', value: 'STUCK' },
  { label: 'Done', value: 'DONE' },
  { label: 'Success', value: 'SUCCESS' },
];

export const OnboardingStatusModal = ({
  isOpen,
  onClose,
}: OnboardingStatusModalProps) => {
  const { handleSubmit } = useForm({
    formId,
    defaultValues: {
      status: '',
      reason: '',
    },
    onSubmit: async (values) => {
      alert(JSON.stringify(values, null, 2));
    },
  });

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
          />
          <div>
            <Text as='label' htmlFor='reason' fontSize='sm'>
              <b>Reason for change</b> (optional)
            </Text>
            <FormAutoresizeTextarea
              pt='0'
              formId={formId}
              id='reason'
              name='reason'
              spellCheck='false'
              placeholder={`What is the reason for changing the onboarding status?`}
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
            type='submit'
            variant='outline'
            colorScheme='primary'
          >
            Update status
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
