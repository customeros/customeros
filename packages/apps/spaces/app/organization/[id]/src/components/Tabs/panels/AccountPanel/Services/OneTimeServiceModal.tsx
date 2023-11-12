'use client';
import { useRef, useState } from 'react';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { FeaturedIcon } from '@ui/media/Icon';
import { FormSelect } from '@ui/form/SyncSelect';
import { Heading } from '@ui/typography/Heading';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { AutoresizeTextarea } from '@ui/form/Textarea';
import { CurrencyInput } from '@ui/form/CurrencyInput';
import { ClockCheck } from '@ui/media/icons/ClockCheck';
import { RenewalLikelihoodProbability } from '@graphql/types';
import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar';
import {
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
  ModalCloseButton,
} from '@ui/overlay/Modal';

export type OneTimeServiceValue = {
  amount?: string | null;
  comment?: string | null;
};

interface OneTimeServiceModalProps {
  name: string;
  isOpen: boolean;
  onClose: () => void;
  OneTimeService: OneTimeServiceValue;
  renewalProbability?: RenewalLikelihoodProbability | null;
}

export const OneTimeServiceModal = ({
  OneTimeService,
  isOpen,
  onClose,
}: OneTimeServiceModalProps) => {
  const initialRef = useRef(null);

  const [amount, setAmount] = useState<string>(OneTimeService?.amount || '');
  const [reason, setReason] = useState<string>(OneTimeService?.comment || '');

  const handleSetOneTimeServiceData = () => {
    // todo COS-857
    onClose();
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
        <ModalCloseButton />
        <ModalHeader>
          <FeaturedIcon size='lg' colorScheme='primary'>
            <DotSingle color='primary.600' />
          </FeaturedIcon>
          <Heading fontSize='lg' mt='4'>
            Add a new one-time service
          </Heading>
        </ModalHeader>
        <ModalBody pb='0'>
          <Flex>
            <FormSelect
              isReadOnly
              label='Billed'
              isLabelVisible
              name='billed'
              formId='tbd'
              options={[{ value: 'once', label: 'Once' }]}
              value={{ value: 'once', label: 'Once' }}
              leftElement={<ClockCheck mr='3' color='gray.500' />}
            />
            <CurrencyInput
              onChange={setAmount}
              value={`${amount}`}
              w='full'
              placeholder='Price'
              isLabelVisible
              label='Price'
              min={0}
              ref={initialRef}
              leftElement={
                <Box color='gray.500'>
                  <CurrencyDollar height='16px' />
                </Box>
              }
            />
          </Flex>

          <AutoresizeTextarea
            pt='0'
            id='description'
            value={reason}
            label='Description (Optional)'
            isLabelVisible
            spellCheck='false'
            onChange={(e) => setReason(e.target.value)}
            placeholder='What is this service about?'
          />
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
            onClick={handleSetOneTimeServiceData}
          >
            Add
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
