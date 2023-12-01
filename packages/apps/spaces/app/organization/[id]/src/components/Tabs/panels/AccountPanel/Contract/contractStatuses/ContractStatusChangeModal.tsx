'use client';
import React from 'react';

import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { ContractStatus } from '@graphql/types';
import { Heading } from '@ui/typography/Heading';
import { Icons, FeaturedIcon } from '@ui/media/Icon';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
import { DatePicker } from '@ui/form/DatePicker/DatePicker';
import {
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal';

import { contractOptionIcon, confirmationModalDataByStatus } from './utils';

interface ContractStatusChangeModalProps {
  isOpen: boolean;
  onClose: () => void;
  organizationName: string;
  mode?: ContractStatus.Live | ContractStatus.Draft | ContractStatus.Ended;
}

export const ContractStatusChangeModal = ({
  isOpen,
  onClose,
  organizationName,
  mode = ContractStatus.Draft,
}: ContractStatusChangeModalProps) => {
  const status = {
    isLive: mode === ContractStatus.Live,
    isDraft: mode === ContractStatus.Draft,
    isEnded: mode === ContractStatus.Ended,
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
        <ModalHeader>
          <FeaturedIcon
            size='lg'
            colorScheme={confirmationModalDataByStatus[mode].colorScheme}
          >
            {contractOptionIcon[mode]}
          </FeaturedIcon>
          <Heading fontSize='lg' mt='4'>
            {confirmationModalDataByStatus[mode].title}
          </Heading>
        </ModalHeader>
        <ModalBody pb='0'>
          <Text fontSize='sm'>
            {confirmationModalDataByStatus[mode].description(
              organizationName,
              new Date().toISOString(),
            )}
          </Text>

          {(status.isLive || status.isEnded) && (
            <DatePicker
              label={
                ContractStatus.Live === mode
                  ? 'Service start date'
                  : 'Contract end date'
              }
              placeholder='Signed date'
              formId='tbd'
              name='date'
              calendarIconHidden
              inset='120% auto auto 0px'
            />
          )}

          {status.isEnded && (
            <>
              <Text fontSize='xs' mb={4}>
                Your current contract end date{' '}
                {/* todo: add condition when data will be available for fixed contract */}
                Auto-calculated to match the next{' '}
                {/* todo: add condition when data will be available for contract with renewal*/}
                billing cycle.
              </Text>
              <FormAutoresizeTextarea
                label='Reason for change'
                placeholder='What is the reason for ending this contract?'
                name='reason'
                formId='tbd'
                leftElement={<Icons.Target5 color='gray.500' />}
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
            colorScheme={status.isEnded ? 'error' : 'primary'}
          >
            {confirmationModalDataByStatus[mode].submit}
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
