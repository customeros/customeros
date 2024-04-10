import React from 'react';

import { utcToZonedTime } from 'date-fns-tz';

import { Contract } from '@graphql/types';
import { Button } from '@ui/form/Button/Button';
import { DateTimeUtils } from '@spaces/utils/date';
import { billingFrequencyOptions } from '@organization/src/components/Tabs/panels/AccountPanel/utils';

export const ContractSubtitle = ({
  data,
  onOpenEditModal,
}: {
  data: Contract;
  onOpenEditModal: () => void;
}) => {
  const serviceStarted = data?.serviceStarted
    ? utcToZonedTime(data?.serviceStarted, 'UTC').toUTCString()
    : null;

  const contractEnded = data?.contractEnded
    ? utcToZonedTime(data?.contractEnded, 'UTC').toUTCString()
    : null;
  const renewalDate = data?.opportunities?.[0]?.renewedAt
    ? DateTimeUtils.format(
        utcToZonedTime(data.opportunities[0].renewedAt, 'UTC').toUTCString(),
        DateTimeUtils.dateWithAbreviatedMonth,
      )
    : null;
  const hasStartedService =
    serviceStarted && !DateTimeUtils.isFuture(serviceStarted);

  const serviceStartDate =
    serviceStarted && DateTimeUtils.isFuture(serviceStarted)
      ? DateTimeUtils.format(
          serviceStarted,
          DateTimeUtils.dateWithAbreviatedMonth,
        )
      : null;

  const endDate = contractEnded
    ? DateTimeUtils.format(contractEnded, DateTimeUtils.dateWithAbreviatedMonth)
    : null;

  const renewalPeriod = billingFrequencyOptions.find(
    (e) => e.value === data?.contractRenewalCycle,
  )?.label;

  if (
    !renewalPeriod &&
    !hasStartedService &&
    !serviceStartDate &&
    data?.contractRenewalCycle
  ) {
    return (
      <p className='font-normal shadow-none text-sm  text-gray-500 focus:text-gray-500 hover:text-gray-500 hover:no-underline focus:no-underline'>
        Contract starting...
        <Button
          className='underline ml-1 p-0 font-normal text-sm text-gray-500 hover:text-gray-500 focus:text-gray-500'
          variant='link'
          size='xs'
          onClick={onOpenEditModal}
        >
          Edit contract
        </Button>
      </p>
    );
  }
  if (!hasStartedService && !serviceStartDate && data?.contractRenewalCycle) {
    return (
      <p className='font-normal shadow-none text-sm  text-gray-500 focus:text-gray-500 hover:text-gray-500 hover:no-underline focus:no-underline'>
        {renewalPeriod} contract starting ...
        <Button
          className='underline ml-1 p-0 font-normal text-sm text-gray-500 hover:text-gray-500 focus:text-gray-500'
          variant='link'
          size='xs'
          onClick={onOpenEditModal}
        >
          Edit contract
        </Button>
      </p>
    );
  }
  if (!hasStartedService && serviceStartDate && data?.contractRenewalCycle) {
    return (
      <p className='font-normal shadow-none text-sm  text-gray-500 focus:text-gray-500 hover:text-gray-500 hover:no-underline focus:no-underline'>
        {renewalPeriod} contract starting {serviceStartDate}
      </p>
    );
  }
  if (hasStartedService && endDate) {
    return (
      <p className='font-normal shadow-none text-sm  text-gray-500 focus:text-gray-500 hover:text-gray-500 hover:no-underline focus:no-underline'>
        {renewalPeriod} contract{' '}
        {contractEnded && DateTimeUtils.isFuture(contractEnded)
          ? 'ending'
          : 'ended on'}{' '}
        {endDate}
      </p>
    );
  }
  if (hasStartedService && renewalDate && data?.autoRenew) {
    return (
      <p className='font-normal shadow-none text-sm  text-gray-500 focus:text-gray-500 hover:text-gray-500 hover:no-underline focus:no-underline'>
        {renewalPeriod} contract auto-renewing {renewalDate}
      </p>
    );
  }

  if (hasStartedService && !data?.autoRenew) {
    return (
      <p className='font-normal shadow-none text-sm  text-gray-500 focus:text-gray-500 hover:text-gray-500 hover:no-underline focus:no-underline'>
        {renewalPeriod} contract until {renewalDate}, not auto-renewing
      </p>
    );
  }

  return null;
};
