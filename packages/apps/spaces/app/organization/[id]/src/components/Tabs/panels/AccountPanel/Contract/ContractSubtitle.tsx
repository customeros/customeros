import React, { useMemo } from 'react';

import { utcToZonedTime } from 'date-fns-tz';

import { Text } from '@ui/typography/Text';
import { DateTimeUtils } from '@spaces/utils/date';
import { Contract, ContractRenewalCycle } from '@graphql/types';
import { billingFrequencyOptions } from '@organization/src/components/Tabs/panels/AccountPanel/utils';

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

  const calcContractEndDate = useMemo(() => {
    if (!serviceStarted) return null;
    switch (data?.contractRenewalCycle) {
      case ContractRenewalCycle.AnnualRenewal:
        return DateTimeUtils.format(
          DateTimeUtils.addYears(
            serviceStarted,
            data?.committedPeriods ?? 1,
          ).toISOString(),
          DateTimeUtils.dateWithAbreviatedMonth,
        );
      case ContractRenewalCycle.MonthlyRenewal:
        return DateTimeUtils.format(
          DateTimeUtils.addMonth(
            serviceStarted,
            data.committedPeriods ?? 1,
          ).toISOString(),
          DateTimeUtils.dateWithAbreviatedMonth,
        );
      case ContractRenewalCycle.QuarterlyRenewal:
        return DateTimeUtils.format(
          DateTimeUtils.addMonth(serviceStarted, 3).toISOString(),
          DateTimeUtils.dateWithAbreviatedMonth,
        );

      default:
        return null;
    }
  }, [data?.contractRenewalCycle, serviceStarted, data?.committedPeriods]);

  const endDate = contractEnded
    ? DateTimeUtils.format(contractEnded, DateTimeUtils.dateWithAbreviatedMonth)
    : null;

  const renewalPeriod = billingFrequencyOptions.find(
    (e) => e.value === data?.contractRenewalCycle,
  )?.label;

  if (!hasStartedService && !serviceStartDate && data?.contractRenewalCycle) {
    return <Text>{renewalPeriod} contract starting ... Edit contract</Text>;
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
        {renewalPeriod} contract until {calcContractEndDate}, not auto-renewing
      </Text>
    );
  }

  return null;
};
