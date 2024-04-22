'use client';
import { useParams } from 'next/navigation';
import { useForm } from 'react-inverted-form';
import React, { useRef, useEffect } from 'react';

import { useQueryClient } from '@tanstack/react-query';
import { UseMutationResult } from '@tanstack/react-query';

<<<<<<< HEAD
import { cn } from '@ui/utils/cn';
import { FeaturedIcon } from '@ui/media/Icon';
import { Button } from '@ui/form/Button/Button';
import { DotLive } from '@ui/media/icons/DotLive';
import { DateTimeUtils } from '@spaces/utils/date';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
=======
import { Button } from '@ui/form/Button/Button';
import { DotLive } from '@ui/media/icons/DotLive';
import { DateTimeUtils } from '@spaces/utils/date';
import { ModalBody } from '@ui/overlay/Modal/Modal';
import { RefreshCw05 } from '@ui/media/icons/RefreshCw05';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon2';
>>>>>>> 1281f1234 (refactor account)
import { Exact, ContractStatus, ContractUpdateInput } from '@graphql/types';
import { DatePickerUnderline } from '@ui/form/DatePicker/DatePickerUnderline';
import { UpdateContractMutation } from '@organization/src/graphql/updateContract.generated';
import {
  GetContractsQuery,
  useGetContractsQuery,
} from '@organization/src/graphql/getContracts.generated';
import { useContractModalStatusContext } from '@organization/src/components/Tabs/panels/AccountPanel/context/ContractStatusModalsContext';

interface ContractStartModalProps {
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
  onClose,
  contractId,
  organizationName,
  serviceStarted,
  onUpdateContract,
  status,
}: ContractStartModalProps) => {
  const queryClient = useQueryClient();
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const id = useParams()?.id as string;
  const queryKey = useGetContractsQuery.getKey({ id });

  const { nextInvoice } = useContractModalStatusContext();
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
      },
    );
  };

  return (
<<<<<<< HEAD
    <>
      <div
        className={
          'rounded-2xl max-w-[500px] h-full flex flex-col justify-between'
        }
      >
        <div>
          <div>
            {!nextInvoice && (
              <FeaturedIcon size='lg' colorScheme='primary'>
                <DotLive className='text-primary-600' />
              </FeaturedIcon>
=======
    <Modal open={isOpen} onOpenChange={onClose}>
      <ModalOverlay />
      <ModalContent className='rounded-2xl max-w-[500px]'>
        <ModalHeader>
          <FeaturedIcon
            size='lg'
            colorScheme='primary'
            className='ml-[14px] mb-[25px] mt-[12px]'
          >
            {status === ContractStatus.OutOfContract ? (
              <RefreshCw05 className='text-primary-600' />
            ) : (
              <DotLive className='text-primary-600' />
>>>>>>> 1281f1234 (refactor account)
            )}

            <h1
              className={cn('text-lg font-semibold  mb-1', {
                'mt-4': !nextInvoice,
              })}
            >
              {status === ContractStatus.OutOfContract
                ? 'Renew contract'
                : 'Make this contract live?'}
            </h1>
          </div>
          <div className='flex flex-col'>
            <p className='text-sm'>
              Congrats! Let’s make
              <span className='font-medium mx-1'>{organizationName}’s </span>
              contract live starting on
              <div className='ml-1 inline'>
                <DatePickerUnderline formId={formId} name='serviceStarted' />
              </div>
            </p>
            <p className='text-sm mt-3'>
              Once the contract goes live, we’ll start sending invoices.
            </p>
            {nextInvoice && (
              <p className='text-sm'>
                The first one will be for
                <span className='text-sm ml-1 font-medium'>
                  {formatCurrency(
                    nextInvoice.amountDue,
                    2,
                    nextInvoice.currency,
                  )}{' '}
                  on{' '}
                  {DateTimeUtils.format(
                    nextInvoice.issued,
                    DateTimeUtils.defaultFormatShortString,
                  )}{' '}
                  (
                  {DateTimeUtils.format(
                    nextInvoice.invoicePeriodStart,
                    DateTimeUtils.dateDayAndMonth,
                  )}{' '}
                  -{' '}
                  {DateTimeUtils.format(
                    nextInvoice.invoicePeriodEnd,
                    DateTimeUtils.dateDayAndMonth,
                  )}
                  )
                </span>
              </p>
            )}
          </div>
        </div>

        <div className='mt-6 flex'>
          <Button
            variant='outline'
            size='lg'
            className='w-full'
            onClick={onClose}
          >
            Not now
          </Button>
          <Button
            className='ml-3 w-full'
            variant='outline'
            size='lg'
            colorScheme='primary'
            onClick={handleApplyChanges}
          >
            Go live{' '}
            {DateTimeUtils.format(
              state.values.serviceStarted as string,
              DateTimeUtils.defaultFormatShortString,
            )}
          </Button>
        </div>
      </div>
    </>
  );
};
