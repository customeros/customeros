import React, { useState } from 'react';

import { ContractStore } from '@store/Contracts/Contract.store.ts';

import { cn } from '@ui/utils/cn';
import { DateTimeUtils } from '@utils/date';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { DotLive } from '@ui/media/icons/DotLive';
import { Invoice, ContractStatus } from '@graphql/types';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber';
import { DatePickerUnderline2 } from '@ui/form/DatePicker/DatePickerUnderline2.tsx';

interface ContractStartModalProps {
  contractId: string;
  onClose: () => void;
  serviceStarted?: string;
  organizationName: string;
  status?: ContractStatus | null;
}

export const ContractStartModal = ({
  onClose,
  contractId,
  organizationName,
  serviceStarted,
  status,
}: ContractStartModalProps) => {
  const store = useStore();
  const contractStore = store.contracts.value.get(contractId) as ContractStore;
  const nextInvoice: Invoice | undefined =
    contractStore?.value?.upcomingInvoices?.find(
      (invoice: Invoice) => invoice.issued === nextInvoice,
    );

  const [serviceStartedData, setServiceStarted] = useState<
    string | Date | null | undefined
  >(serviceStarted ? new Date(serviceStarted) : new Date());

  const handleApplyChanges = () => {
    contractStore?.update((prev) => ({
      ...prev,
      serviceStarted: serviceStartedData,
      approved: true,
    }));
    if (
      DateTimeUtils.isPast(serviceStartedData as string) ||
      DateTimeUtils.isToday(serviceStartedData as string)
    ) {
      contractStore?.update(
        (prev) => ({
          ...prev,
          status: ContractStatus.Live,
        }),
        { mutate: false },
      );
    } else {
      contractStore?.update(
        (prev) => ({
          ...prev,
          status: ContractStatus.Scheduled,
        }),
        { mutate: false },
      );
    }
    onClose();
  };

  return (
    <>
      <div
        className={
          'rounded-2xl max-w-[600px] h-full flex flex-col justify-between'
        }
      >
        <div>
          <div>
            {!nextInvoice && (
              <FeaturedIcon size='lg' colorScheme='primary'>
                <DotLive className='text-primary-600' />
              </FeaturedIcon>
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
              Congrats! Let’s make{' '}
              <span className='font-medium '>{organizationName}’s </span>
              contract live starting on
              <div className='ml-1 inline-flex text-sm'>
                <DatePickerUnderline2
                  value={serviceStartedData as string}
                  onChange={(e) => setServiceStarted(e)}
                />
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
            loadingText='Saving...'
          >
            Go live{' '}
            {DateTimeUtils.format(
              serviceStartedData as string,
              DateTimeUtils.defaultFormatShortString,
            )}
          </Button>
        </div>
      </div>
    </>
  );
};
