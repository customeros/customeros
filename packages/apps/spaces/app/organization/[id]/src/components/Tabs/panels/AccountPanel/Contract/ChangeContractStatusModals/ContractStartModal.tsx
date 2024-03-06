'use client';
import { useForm } from 'react-inverted-form';
import React, { useRef, useEffect } from 'react';

import { UseMutationResult } from '@tanstack/react-query';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { ModalBody } from '@ui/overlay/Modal';
import { FeaturedIcon } from '@ui/media/Icon';
import { Heading } from '@ui/typography/Heading';
import { DotLive } from '@ui/media/icons/DotLive';
import { Exact, ContractUpdateInput } from '@graphql/types';
import { DatePickerUnderline } from '@ui/form/DatePicker/DatePickerUnderline';
import { GetContractsQuery } from '@organization/src/graphql/getContracts.generated';
import { UpdateContractMutation } from '@organization/src/graphql/updateContract.generated';
import {
  Modal,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal';

interface ContractStartModalProps {
  isOpen: boolean;
  contractId: string;
  onClose: () => void;
  organizationName: string;
  serviceStartedAt?: string;
  onUpdateContract: UseMutationResult<
    UpdateContractMutation,
    unknown,
    Exact<{ input: ContractUpdateInput }>,
    { previousEntries: GetContractsQuery | undefined }
  >;
}

export const ContractStartModal = ({
  isOpen,
  onClose,
  contractId,
  organizationName,
  serviceStartedAt,
  onUpdateContract,
}: ContractStartModalProps) => {
  const initialRef = useRef(null);
  const formId = `contract-starts-on-form-${contractId}`;
  const { state, setDefaultValues } = useForm<{
    serviceStartedAt?: string | Date | null;
  }>({
    formId,
    defaultValues: {
      serviceStartedAt: serviceStartedAt
        ? new Date(serviceStartedAt)
        : new Date(),
    },
    stateReducer: (_, action, next) => {
      return next;
    },
  });

  useEffect(() => {
    if (serviceStartedAt) {
      setDefaultValues({
        serviceStartedAt: new Date(serviceStartedAt),
      });
    }
  }, [serviceStartedAt]);

  const handleApplyChanges = () => {
    onUpdateContract.mutate(
      {
        input: {
          contractId,
          patch: true,
          serviceStartedAt: state.values.serviceStartedAt,
          endedAt: '0001-01-01T00:00:00.000000Z',
        },
      },
      {
        onSuccess: () => {
          onClose();
        },
      },
    );
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      initialFocusRef={initialRef}
      size='md'
    >
      <ModalOverlay />
      <ModalContent borderRadius='2xl'>
        <ModalHeader>
          <FeaturedIcon size='lg' colorScheme='primary'>
            <DotLive color='primary.600' />
          </FeaturedIcon>
          <Heading fontSize='lg' mt='4'>
            Make this contract live?
          </Heading>
        </ModalHeader>
        <ModalBody as={Flex} flexDir='column' gap={4}>
          <Text>
            Congrats! Let’s make {organizationName}
            ’s contract live starting on
            <Box ml={1}>
              <DatePickerUnderline
                placeholder='Start date'
                formId={formId}
                name='serviceStartedAt'
                calendarIconHidden
                value={state.values.serviceStartedAt}
              />
            </Box>
          </Text>
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
            loadingText='Applying changes...'
            onClick={handleApplyChanges}
          >
            Make live
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
