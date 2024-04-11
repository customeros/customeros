'use client';
import React, { useEffect } from 'react';
import { useForm } from 'react-inverted-form';

import { UseMutationResult } from '@tanstack/react-query';

import { FeaturedIcon } from '@ui/media/Icon';
import { Button } from '@ui/form/Button/Button';
import { DotLive } from '@ui/media/icons/DotLive';
import { DateTimeUtils } from '@spaces/utils/date';
import { ModalBody } from '@ui/overlay/Modal/Modal';
import { RefreshCw05 } from '@ui/media/icons/RefreshCw05';
import { Exact, ContractStatus, ContractUpdateInput } from '@graphql/types';
import { DatePickerUnderline } from '@ui/form/DatePicker/DatePickerUnderline';
import { GetContractsQuery } from '@organization/src/graphql/getContracts.generated';
import { UpdateContractMutation } from '@organization/src/graphql/updateContract.generated';
import {
  Modal,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal/Modal';

interface ContractStartModalProps {
  isOpen: boolean;
  contractId: string;
  onClose: () => void;
  serviceStarted?: string;
  organizationName: string;
  status?: ContractStatus | null;
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
  serviceStarted,
  onUpdateContract,
  status,
}: ContractStartModalProps) => {
  const formId = `contract-starts-on-form-${contractId}`;
  const { state, setDefaultValues } = useForm<{
    serviceStarted?: string | Date | null;
  }>({
    formId,
    defaultValues: {
      serviceStarted: serviceStarted ? new Date(serviceStarted) : new Date(),
    },
    stateReducer: (_, action, next) => {
      return next;
    },
  });

  useEffect(() => {
    if (serviceStarted) {
      setDefaultValues({
        serviceStarted: new Date(serviceStarted),
      });
    }
  }, [serviceStarted]);

  const handleApplyChanges = () => {
    onUpdateContract.mutate(
      {
        input: {
          contractId,
          patch: true,
          serviceStarted: state.values.serviceStarted,
          approved: true,
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
    <Modal open={isOpen} onOpenChange={onClose}>
      <ModalOverlay />
      <ModalContent className='rounded-2xl'>
        <ModalHeader>
          <FeaturedIcon size='lg' colorScheme='primary'>
            {status === ContractStatus.OutOfContract ? (
              <RefreshCw05 className='text-primary-600' />
            ) : (
              <DotLive className='text-primary-600' />
            )}
          </FeaturedIcon>
          <h1 className='text-lg mt-4'>
            {status === ContractStatus.OutOfContract
              ? 'Renew contract'
              : 'Make this contract live?'}
          </h1>
        </ModalHeader>
        <ModalBody className='flex flex-col gap-4'>
          {status === ContractStatus.OutOfContract ? (
            <div>a</div>
          ) : (
            <p className='text-base'>
              Congrats! Let’s make {organizationName}
              ’s contract live starting on
              <div className='ml-1 inline'>
                <DatePickerUnderline
                  placeholder='Start date'
                  formId={formId}
                  name='serviceStarted'
                  calendarIconHidden
                  value={state.values.serviceStarted}
                />
              </div>
            </p>
          )}
        </ModalBody>
        <ModalFooter className='p-6'>
          <Button variant='outline' className='w-full' onClick={onClose}>
            Cancel
          </Button>
          <Button
            className='ml-3 w-full'
            variant='outline'
            colorScheme='primary'
            onClick={handleApplyChanges}
          >
            {status === ContractStatus.OutOfContract ? 'Renew' : 'Start'}
            {DateTimeUtils.format(
              state.values.serviceStarted as string,
              DateTimeUtils.defaultFormatShortString,
            )}
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
