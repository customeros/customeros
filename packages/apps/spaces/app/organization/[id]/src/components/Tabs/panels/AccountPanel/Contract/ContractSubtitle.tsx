import React from 'react';

import { utcToZonedTime } from 'date-fns-tz';

import { Contract } from '@graphql/types';
import { Text } from '@ui/typography/Text';
import { DateTimeUtils } from '@spaces/utils/date';
import { contractRenewalCycle } from '@organization/src/components/Tabs/panels/AccountPanel/utils';

export const ContractSubtitle = ({ data }: { data: Contract }) => {
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

  const renewalPeriod = contractRenewalCycle.find(
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
        <span className='underline ml-1'> Edit contract</span>
      </p>
    );
  }
  if (!hasStartedService && !serviceStartDate && data?.contractRenewalCycle) {
    return (
      <p className='font-normal shadow-none text-sm  text-gray-500 focus:text-gray-500 hover:text-gray-500 hover:no-underline focus:no-underline'>
        {renewalPeriod} contract starting ...
        <span className='underline ml-1'> Edit contract</span>
      </p>
    );
  }
  if (!hasStartedService && serviceStartDate && data?.contractRenewalCycle) {
    return (
      <Text>
        {renewalPeriod} contract starting {serviceStartDate}
      </Text>
    );
  }
  if (hasStartedService && endDate) {
    return (
      <Text>
        {renewalPeriod} contract{' '}
        {contractEnded && DateTimeUtils.isFuture(contractEnded)
          ? 'ending'
          : 'ended on'}{' '}
        {endDate}
      </Text>
    );
  }
  if (hasStartedService && renewalDate && data?.autoRenew) {
    return (
      <Text>
        {renewalPeriod} contract auto-renewing {renewalDate}
      </Text>
    );
  }

  if (hasStartedService && !data?.autoRenew) {
    return (
      <Text>
        {renewalPeriod} contract until {renewalDate}, not auto-renewing
      </Text>
    );
  }

  return null;
};
