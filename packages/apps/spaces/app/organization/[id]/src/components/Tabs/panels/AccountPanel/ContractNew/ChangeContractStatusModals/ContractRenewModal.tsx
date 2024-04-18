'use client';
import React, { useRef } from 'react';
import { useParams } from 'next/navigation';

import { useQueryClient } from '@tanstack/react-query';

import { FeaturedIcon } from '@ui/media/Icon';
import { Button } from '@ui/form/Button/Button';
import { ModalBody } from '@ui/overlay/Modal/Modal';
import { RefreshCw05 } from '@ui/media/icons/RefreshCw05';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGetContractsQuery } from '@organization/src/graphql/getContracts.generated';
import { useRenewContractMutation } from '@organization/src/graphql/renewContract.generated';
import {
  Modal,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal';

interface ContractEndModalProps {
  isOpen: boolean;
  contractId: string;
  onClose: () => void;
  organizationName: string;
}

export const ContractRenewsModal = ({
  isOpen,
  onClose,
  contractId,
  organizationName,
}: ContractEndModalProps) => {
  const initialRef = useRef(null);
  const client = getGraphQLClient();
  const id = useParams()?.id as string;

  const queryKey = useGetContractsQuery.getKey({ id });
  const queryClient = useQueryClient();
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  const { mutate } = useRenewContractMutation(client, {
    onSuccess: () => {
      onClose();
    },
    onSettled: () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }

      timeoutRef.current = setTimeout(() => {
        queryClient.invalidateQueries({ queryKey });
        queryClient.invalidateQueries({
          queryKey: ['GetTimeline.infinite'],
        });
      }, 1000);
    },
  });

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      initialFocusRef={initialRef}
      size='md'
    >
      <ModalOverlay />
      <ModalContent borderRadius='2xl'>
        <ModalHeader className='pb-2'>
          <FeaturedIcon size='lg' colorScheme='primary'>
            <RefreshCw05 className='text-primary-600' />
          </FeaturedIcon>
          <h2 className='text-lg mt-4'>Renew contract</h2>
        </ModalHeader>
        <ModalBody className='flex flex-col gap-3'>
          Let’s renew {organizationName}’s contract from today
        </ModalBody>
        <ModalFooter p='6'>
          <Button
            variant='outline'
            colorScheme='gray'
            className='w-full'
            onClick={onClose}
          >
            Cancel
          </Button>
          <Button
            className='ml-3 w-full'
            variant='outline'
            colorScheme='primary'
            onClick={() => mutate({ contractId })}
          >
            Renew now
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
