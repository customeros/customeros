'use client';
import { useState } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Heading } from '@ui/typography/Heading';
import { Button, ButtonGroup } from '@ui/form/Button';
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
import { Dot } from '@ui/media/Dot';

export type Likelihood = 'HIGH' | 'MEDIUM' | 'LOW' | 'ZERO' | 'NOT_SET';
export type Value = { likelihood: Likelihood; reason: string };

interface RenewalLikelihoodModalProps {
  isOpen: boolean;
  onClose: () => void;
  value: Value;
  onChange: (value: Value) => void;
}

export const RenewalLikelihoodModal = ({
  value,
  isOpen,
  onClose,
  onChange,
}: RenewalLikelihoodModalProps) => {
  const [likelihood, setLikelihood] = useState<Likelihood>(value.likelihood);
  const [reason, setReason] = useState<string>(value.reason);

  const handleSet = () => {
    onChange({ likelihood, reason });
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
            {`${
              value.likelihood === 'NOT_SET' ? 'Set' : 'Update'
            } renewal likelihood`}
          </Heading>
          <Text mt='1' fontSize='sm' fontWeight='normal'>
            {value.likelihood === 'NOT_SET' ? 'Setting' : 'Updating'}{' '}
            <b>Acme Corp’s</b> renewal likelihood will change how its renewal
            estimates are calculated and actions are prioritised.
          </Text>
        </ModalHeader>
        <ModalBody as={Flex} flexDir='column' pb='0'>
          <ButtonGroup w='full' isAttached>
            <Button
              w='full'
              variant='outline'
              leftIcon={<Dot colorScheme='success' />}
              onClick={() => setLikelihood('HIGH')}
              bg={likelihood === 'HIGH' ? 'gray.100' : 'white'}
            >
              High
            </Button>
            <Button
              w='full'
              variant='outline'
              leftIcon={<Dot colorScheme='warning' />}
              onClick={() => setLikelihood('MEDIUM')}
              bg={likelihood === 'MEDIUM' ? 'gray.100' : 'white'}
            >
              Medium
            </Button>
            <Button
              w='full'
              variant='outline'
              leftIcon={<Dot colorScheme='error' />}
              onClick={() => setLikelihood('LOW')}
              bg={likelihood === 'LOW' ? 'gray.100' : 'white'}
            >
              Low
            </Button>
            <Button
              variant='outline'
              w='full'
              leftIcon={<Dot />}
              onClick={() => setLikelihood('ZERO')}
              bg={likelihood === 'ZERO' ? 'gray.100' : 'white'}
            >
              Zero
            </Button>
          </ButtonGroup>

          {likelihood !== 'NOT_SET' && (
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
                  value.likelihood === 'NOT_SET' ? 'setting' : 'updating'
                } the renewal likelihood?`}
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
            {value.likelihood === 'NOT_SET' ? 'Set' : 'Update'}
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
