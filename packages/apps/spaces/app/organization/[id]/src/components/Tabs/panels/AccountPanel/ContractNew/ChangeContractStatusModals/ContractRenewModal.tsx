'use client';
import React, { useRef } from 'react';
import { useParams } from 'next/navigation';

import { useQueryClient } from '@tanstack/react-query';

import { cn } from '@ui/utils/cn';
import { FeaturedIcon } from '@ui/media/Icon';
import { ContractStatus } from '@graphql/types';
import { Button } from '@ui/form/Button/Button';
import { DateTimeUtils } from '@spaces/utils/date';
import { Radio, RadioGroup } from '@ui/form/Radio/Radio2';
import { RefreshCw05 } from '@ui/media/icons/RefreshCw05';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { DatePickerUnderline } from '@ui/form/DatePicker/DatePickerUnderline';
import { useGetContractsQuery } from '@organization/src/graphql/getContracts.generated';
import { useRenewContractMutation } from '@organization/src/graphql/renewContract.generated';
import { useContractModalStatusContext } from '@organization/src/components/Tabs/panels/AccountPanel/context/ContractStatusModalsContext';

interface ContractEndModalProps {
  renewsAt?: string;
  contractId: string;
  onClose: () => void;
  organizationName: string;

  status?: ContractStatus | null;
}
export enum EndContract {
  Now = 'Now',
  EndOfCurrentBillingPeriod = 'EndOfCurrentBillingPeriod',
  EndOfCurrentRenewalPeriod = 'EndOfCurrentRenewalPeriod',
  CustomDate = 'CustomDate',
}
export const ContractRenewsModal = ({
  onClose,
  contractId,
  organizationName,
  status,
  renewsAt,
}: ContractEndModalProps) => {
  const client = getGraphQLClient();
  const id = useParams()?.id as string;

  const queryKey = useGetContractsQuery.getKey({ id });
  const queryClient = useQueryClient();
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const { nextInvoice } = useContractModalStatusContext();
  const [value, setValue] = React.useState(EndContract.Now);
  const formId = `contract-ends-on-form-${contractId}`;
  const timeToRenewal = renewsAt
    ? DateTimeUtils.format(renewsAt, DateTimeUtils.dateWithAbreviatedMonth)
    : null;
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

  const getText = () => {
    if (
      status === ContractStatus.OutOfContract &&
      renewsAt &&
      DateTimeUtils.isSameDay(renewsAt, new Date().toISOString())
    ) {
      return 'Renewing this contract will extend it with another year until 2 Aug 2025';
    }
  };
  const { state, setDefaultValues } = useForm<{
    renewsAt?: string | Date | null;
  }>({
    formId,
    defaultValues: { renewsAt: null || new Date() },
    stateReducer: (_, action, next) => {
      return next;
    },
  });

  const handleApplyChanges = () => {};

  const handleChangeEndsOnOption = (nextValue: string | null) => {
    if (nextValue === EndContract.Now) {
      setDefaultValues({ renewsAt: today });
      setValue(EndContract.Now);

      return;
    }
    if (nextValue === EndContract.EndOfCurrentBillingPeriod) {
      setDefaultValues({ renewsAt: nextInvoice?.issued });
      setValue(EndContract.EndOfCurrentBillingPeriod);

      return;
    }
    if (nextValue === EndContract.CustomDate) {
      setDefaultValues({ renewsAt: new Date(today) });
      setValue(EndContract.CustomDate);

      return;
    }
    if (nextValue === EndContract.EndOfCurrentRenewalPeriod) {
      setDefaultValues({ renewsAt: renewsAt });
      setValue(EndContract.EndOfCurrentRenewalPeriod);

      return;
    }
  };

  return (
    <>
      <div>
        {!nextInvoice && (
          <FeaturedIcon size='lg' colorScheme='primary'>
            <RefreshCw05 className='text-primary-600' />
          </FeaturedIcon>
        )}

        <h1
          className={cn('text-lg font-semibold  mb-1', {
            'mt-4': !nextInvoice,
          })}
        >
          {status === ContractStatus.OutOfContract
            ? 'Renew this contract?'
            : 'When should this contract renew?'}
        </h1>
      </div>

      <p className='flex flex-col gap-3'>{getText()}</p>

      <RadioGroup
        value={value}
        onValueChange={handleChangeEndsOnOption}
        className='flex flex-col gap-1 text-base'
      >
        <Radio value={EndContract.Now}>
          <span className='mr-1'>Now</span>
        </Radio>

        {timeToRenewal && (
          <Radio value={EndContract.EndOfCurrentRenewalPeriod}>
            <span className='ml-1'>
              End of current renewal period, {timeToRenewal}
            </span>
          </Radio>
        )}

        <Radio value={EndContract.CustomDate}>
          <div className='flex items-center max-h-6'>
            On{' '}
            {value === EndContract.CustomDate ? (
              <div className='ml-1'>
                <DatePickerUnderline
                  placeholder='Renewal date'
                  defaultOpen={true}
                  // minDate={state.values.serviceStarted}
                  formId={formId}
                  name='renewsAt'
                  calendarIconHidden
                  value={state.values.renewsAt}
                />
              </div>
            ) : (
              'custom date'
            )}
          </div>
        </Radio>
      </RadioGroup>
      <div className='flex'>
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
      </div>
    </>
  );
};
