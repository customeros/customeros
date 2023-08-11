'use client';
import { useState } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Heading } from '@ui/typography/Heading';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { AutoresizeTextarea } from '@ui/form/Textarea';
import { Icons, FeaturedIcon } from '@ui/media/Icon';
import {
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
  ModalCloseButton,
} from '@ui/overlay/Modal';
import { CurrencyInput } from '@ui/form/CurrencyInput';

export type Value = { forecast: string; reason: string };

interface RenewalForecastModalProps {
  isOpen: boolean;
  onClose: () => void;
  value: Value;
  onChange: (value: Value) => void;
}

export const RenewalForecastModal = ({
  value,
  isOpen,
  onClose,
  onChange,
}: RenewalForecastModalProps) => {
  const [forecast, setForecast] = useState<string>(value.forecast);
  const [reason, setReason] = useState<string>(value.reason);

  const handleSet = () => {
    onChange({ forecast, reason });
    onClose();
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
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
          <FeaturedIcon size='lg' colorScheme='warning'>
            <Icons.AlertTriangle />
          </FeaturedIcon>
          <Heading fontSize='lg' mt='4'>
            {`${!value.forecast ? 'Set' : 'Update'} renewal forecast`}
          </Heading>
          <Text mt='1' fontSize='sm' fontWeight='normal'>
            {!value.forecast ? 'Setting' : 'Updating'} <b>Acme Corp’s</b>{' '}
            renewal forecast will change how expected revenue is reported.
          </Text>
        </ModalHeader>
        <ModalBody as={Flex} flexDir='column' pb='0'>
          <CurrencyInput onChange={setForecast} value={forecast} w='full' />

          {forecast && (
            <>
              <Text as='label' htmlFor='reason' mt='5' fontSize='sm'>
                <b>Reason for change</b> (optional)
              </Text>
              <AutoresizeTextarea
                pt='0'
                id='reason'
                value={reason}
                spellCheck='false'
                onChange={(e) => setReason(e.target.value)}
                placeholder={`What is the reason for ${
                  !value.forecast ? 'setting' : 'updating'
                } the renewal forecast?`}
              />
            </>
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
            onClick={handleSet}
          >
            {!value.forecast ? 'Set' : 'Update'}
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
