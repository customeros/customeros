import React, { useState } from 'react';

import { cn } from '@ui/utils/cn';
import { DateTimeUtils } from '@utils/date';
import { Button } from '@ui/form/Button/Button';
import { ContractStatus } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { Radio, RadioGroup } from '@ui/form/Radio/Radio';
import { RefreshCw05 } from '@ui/media/icons/RefreshCw05';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { DatePickerUnderline2 } from '@ui/form/DatePicker/DatePickerUnderline2.tsx';

interface ContractEndModalProps {
  contractId: string;
  onClose: () => void;
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
}: ContractEndModalProps) => {
  const store = useStore();
  const contractStore = store.contracts.value.get(contractId);
  const renewsAt = contractStore?.value?.opportunities?.find(
    (e) => e.internalStage === 'OPEN',
  )?.renewedAt;
  const [value, setValue] = useState(RenewContract.Now);
  const timeToRenewal = renewsAt
    ? DateTimeUtils.format(renewsAt, DateTimeUtils.dateWithAbreviatedMonth)
    : null;
  const renewsToday = renewsAt && DateTimeUtils.isToday(renewsAt);
  const renewsTomorrow = renewsAt && DateTimeUtils.isTomorrow(renewsAt);
  const [renewsAtData, setRenewsAt] = useState<
    string | Date | null | undefined
  >(renewsAt || new Date().toString());

  const handleApplyChanges = () => {
    contractStore?.update((prev) => ({
      ...prev,
      renewalDate: renewsAtData,
      opportunities: prev?.opportunities?.map((opportunity) => {
        if (opportunity.internalStage === 'OPEN') {
          return {
            ...opportunity,
            renewedAt: renewsAtData,
          };
        }

        return opportunity;
      }),
    }));
    onClose();
  };

  const handleChangeEndsOnOption = (nextValue: string | null) => {
    if (nextValue === RenewContract.Now) {
      setRenewsAt(new Date());
      setValue(RenewContract.Now);

      return;
    }
    if (nextValue === RenewContract.EndOfCurrentBillingPeriod) {
      setRenewsAt(contractStore?.value?.upcomingInvoices?.[0]?.issued);
      setValue(RenewContract.EndOfCurrentBillingPeriod);

      return;
    }
    if (nextValue === RenewContract.CustomDate) {
      setRenewsAt(new Date());
      setValue(RenewContract.CustomDate);

      return;
    }
    if (nextValue === RenewContract.EndOfCurrentRenewalPeriod) {
      setRenewsAt(renewsAt);
      setValue(RenewContract.EndOfCurrentRenewalPeriod);

      return;
    }
  };

  return (
    <>
      <div>
        <div>
          {!contractStore?.value?.upcomingInvoices?.length && (
            <FeaturedIcon size='lg' colorScheme='primary'>
              <RefreshCw05 className='text-primary-600' />
            </FeaturedIcon>
          )}

          <h1
            className={cn('text-lg font-semibold  mb-1', {
              'mt-4': !contractStore?.value?.upcomingInvoices?.length,
            })}
          >
            {status === ContractStatus.OutOfContract
              ? 'Renew this contract?'
              : 'When should this contract renew?'}
          </h1>
        </div>

        <p className='flex flex-col mb-3 text-base'>
          Renewing this contract will extend it with another{' '}
          {getCommittedPeriodLabel(
            contractStore?.value.committedPeriodInMonths,
          )}{' '}
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
                      <DatePickerUnderline2
                        value={renewsAtData || new Date().toString()}
                        onChange={(e) => setRenewsAt(e)}
                      />
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
        >
          Renew{' '}
          {RenewContract.Now === value || renewsToday
            ? 'now'
            : DateTimeUtils.format(
                renewsAtData as string,
                DateTimeUtils.defaultFormatShortString,
              )}
        </Button>
      </div>
    </>
  );
};
