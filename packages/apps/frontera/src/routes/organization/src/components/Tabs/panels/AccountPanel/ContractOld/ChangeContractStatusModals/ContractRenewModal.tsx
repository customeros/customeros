import React, { useRef } from 'react';
import { useParams } from 'react-router-dom';
import { useForm } from 'react-inverted-form';

import { useQueryClient } from '@tanstack/react-query';

import { cn } from '@ui/utils/cn';
import { ContractStatus } from '@graphql/types';
import { Button } from '@ui/form/Button/Button';
import { DateTimeUtils } from '@spaces/utils/date';
import { Radio, RadioGroup } from '@ui/form/Radio/Radio';
import { RefreshCw05 } from '@ui/media/icons/RefreshCw05';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { DatePickerUnderline } from '@ui/form/DatePicker/DatePickerUnderline';
import { useGetContractsQuery } from '@organization/graphql/getContracts.generated';
import { useRenewContractMutation } from '@organization/graphql/renewContract.generated';
import { useContractModalStatusContext } from '@organization/components/Tabs/panels/AccountPanel/context/ContractStatusModalsContext';

interface ContractEndModalProps {
  renewsAt?: string;
  contractId: string;
  onClose: () => void;
  status?: ContractStatus | null;
}
export enum RenewContract {
  Now = 'Now',
  EndOfCurrentBillingPeriod = 'EndOfCurrentBillingPeriod',
  EndOfCurrentRenewalPeriod = 'EndOfCurrentRenewalPeriod',
  CustomDate = 'CustomDate',
}

export function getCommittedPeriodLabel(months: string | number) {
  if (`${months}` === '1') {
    return 'month';
  }
  if (`${months}` === '3') {
    return 'quarter';
  }

  if (`${months}` === '12') {
    return 'year';
  }

  return `${months} months`;
}
export const ContractRenewsModal = ({
  onClose,
  contractId,
  status,
  renewsAt,
}: ContractEndModalProps) => {
  const client = getGraphQLClient();
  const id = useParams()?.id as string;

  const queryKey = useGetContractsQuery.getKey({ id });
  const queryClient = useQueryClient();
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const { nextInvoice, committedPeriodInMonths } =
    useContractModalStatusContext();
  const [value, setValue] = React.useState(RenewContract.Now);
  const formId = `contract-ends-on-form-${contractId}`;
  const timeToRenewal = renewsAt
    ? DateTimeUtils.format(renewsAt, DateTimeUtils.dateWithAbreviatedMonth)
    : null;
  const renewsToday = renewsAt && DateTimeUtils.isToday(renewsAt);
  const renewsTomorrow = renewsAt && DateTimeUtils.isTomorrow(renewsAt);

  const { mutate, isPending } = useRenewContractMutation(client, {
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

  const { state, setDefaultValues } = useForm<{
    renewsAt?: string | Date | null;
  }>({
    formId,
    defaultValues: { renewsAt: renewsAt },
    stateReducer: (_, _action, next) => {
      return next;
    },
  });

  const handleApplyChanges = () => {
    mutate(
      {
        input: {
          contractId,
          renewalDate: state.values.renewsAt,
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

  const handleChangeEndsOnOption = (nextValue: string | null) => {
    if (nextValue === RenewContract.Now) {
      setDefaultValues({ renewsAt: new Date() });
      setValue(RenewContract.Now);

      return;
    }
    if (nextValue === RenewContract.EndOfCurrentBillingPeriod) {
      setDefaultValues({ renewsAt: nextInvoice?.issued });
      setValue(RenewContract.EndOfCurrentBillingPeriod);

      return;
    }
    if (nextValue === RenewContract.CustomDate) {
      setDefaultValues({ renewsAt: new Date() });
      setValue(RenewContract.CustomDate);

      return;
    }
    if (nextValue === RenewContract.EndOfCurrentRenewalPeriod) {
      setDefaultValues({ renewsAt: renewsAt });
      setValue(RenewContract.EndOfCurrentRenewalPeriod);

      return;
    }
  };

  return (
    <>
      <div>
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

        <p className='flex flex-col mb-3 text-base'>
          Renewing this contract will extend it with another{' '}
          {getCommittedPeriodLabel(committedPeriodInMonths)}{' '}
        </p>

        {!renewsToday && (
          <RadioGroup
            value={value}
            onValueChange={handleChangeEndsOnOption}
            className='flex flex-col gap-1 text-base'
          >
            <Radio value={RenewContract.Now}>
              <span className='mr-1'>Now</span>
            </Radio>

            {timeToRenewal && (
              <Radio value={RenewContract.EndOfCurrentRenewalPeriod}>
                <span className='ml-1'>
                  End of current renewal period, {timeToRenewal}
                </span>
              </Radio>
            )}

            {!renewsTomorrow && (
              <Radio value={RenewContract.CustomDate}>
                <div className='flex items-center max-h-6'>
                  On{' '}
                  {value === RenewContract.CustomDate ? (
                    <div className='ml-1'>
                      <DatePickerUnderline formId={formId} name='renewsAt' />
                    </div>
                  ) : (
                    'custom date'
                  )}
                </div>
              </Radio>
            )}
          </RadioGroup>
        )}
      </div>

      <div className='flex'>
        <Button
          size='lg'
          variant='outline'
          colorScheme='gray'
          className='w-full'
          onClick={onClose}
        >
          Cancel
        </Button>
        <Button
          size='lg'
          className='ml-3 w-full'
          variant='outline'
          colorScheme='primary'
          onClick={handleApplyChanges}
          loadingText='Renewing...'
          isLoading={isPending}
        >
          Renew{' '}
          {RenewContract.Now === value || renewsToday
            ? 'now'
            : DateTimeUtils.format(
                state.values.renewsAt as string,
                DateTimeUtils.defaultFormatShortString,
              )}
        </Button>
      </div>
    </>
  );
};
