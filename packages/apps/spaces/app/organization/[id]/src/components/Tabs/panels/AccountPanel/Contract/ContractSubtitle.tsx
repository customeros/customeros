import React, { useMemo } from 'react';

import { Text } from '@ui/typography/Text';
import { DateTimeUtils } from '@spaces/utils/date';
import { Contract, ContractRenewalCycle } from '@graphql/types';
import { billingFrequencyOptions } from '@organization/src/components/Tabs/panels/AccountPanel/utils';

export const ContractSubtitle = ({ data }: { data: Contract }) => {
  const hasStartedService =
    data?.serviceStarted && !DateTimeUtils.isFuture(data.serviceStarted);

  const serviceStartDate =
    data?.serviceStarted && DateTimeUtils.isFuture(data.serviceStarted)
      ? DateTimeUtils.format(
          data.serviceStarted,
          DateTimeUtils.dateWithAbreviatedMonth,
        )
      : null;
  const renewalDate = data?.opportunities?.[0]?.renewedAt
    ? DateTimeUtils.format(
        data?.opportunities?.[0]?.renewedAt,
        DateTimeUtils.dateWithAbreviatedMonth,
      )
    : null;
  const calcContractEndDate = useMemo(() => {
    switch (data?.contractRenewalCycle) {
      case ContractRenewalCycle.AnnualRenewal:
        return DateTimeUtils.format(
          DateTimeUtils.addYears(
            data.serviceStarted,
            data?.committedPeriods ?? 1,
          ).toISOString(),
          DateTimeUtils.dateWithAbreviatedMonth,
        );
      case ContractRenewalCycle.MonthlyRenewal:
        return DateTimeUtils.format(
          DateTimeUtils.addMonth(
            data.serviceStarted,
            data.committedPeriods ?? 1,
          ).toISOString(),
          DateTimeUtils.dateWithAbreviatedMonth,
        );
      case ContractRenewalCycle.QuarterlyRenewal:
        return DateTimeUtils.format(
          DateTimeUtils.addMonth(data.serviceStarted, 3).toISOString(),
          DateTimeUtils.dateWithAbreviatedMonth,
        );

      default:
        return null;
    }
  }, [
    data?.contractRenewalCycle,
    data?.serviceStarted,
    data?.committedPeriods,
  ]);

  const endDate = data?.contractEnded
    ? DateTimeUtils.format(
        data.contractEnded,
        DateTimeUtils.dateWithAbreviatedMonth,
      )
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
        {DateTimeUtils.isFuture(data.contractEnded) ? 'ending' : 'ended on'}{' '}
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
