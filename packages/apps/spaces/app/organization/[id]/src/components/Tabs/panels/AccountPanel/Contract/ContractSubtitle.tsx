import React from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { DateTimeUtils } from '@spaces/utils/date';
import { Contract, ContractStatus, ContractRenewalCycle } from '@graphql/types';

function getLabelFromValue(value: string): string | undefined {
  if (ContractRenewalCycle.AnnualRenewal === value) {
    return 'annually';
  }
  if (ContractRenewalCycle.MonthlyRenewal === value) {
    return 'monthly';
  }
}
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
  if (
    !renewalDate &&
    hasStartedService &&
    data?.contractStatus !== ContractStatus.Ended
  ) {
    const serviceStarted = hasStartedService
      ? DateTimeUtils.format(
          data.serviceStarted,
          DateTimeUtils.dateWithAbreviatedMonth,
        )
      : null;

    return <Text>Service started {serviceStarted}</Text>;
  }

  const endDate = data?.contractEnded
    ? DateTimeUtils.format(
        data.contractEnded,
        DateTimeUtils.dateWithAbreviatedMonth,
      )
    : null;

  const isActiveAndRenewable =
    hasStartedService &&
    data.contractStatus !== ContractStatus.Ended &&
    !!data.contractRenewalCycle &&
    data.contractRenewalCycle !== ContractRenewalCycle.None;

  return (
    <Flex flexDir='column' alignItems='flex-start' justifyContent='center'>
      {serviceStartDate && <Text>Service starts {serviceStartDate}</Text>}
      {isActiveAndRenewable && (
        <Text>
          Renews {getLabelFromValue(data.contractRenewalCycle)} on {renewalDate}
        </Text>
      )}
      {data?.contractEnded && DateTimeUtils.isFuture(data.contractEnded) && (
        <Text>Ends {endDate}</Text>
      )}

      {data.contractStatus === ContractStatus.Ended && (
        <Text>Ended on {endDate}</Text>
      )}
    </Flex>
  );
};
