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
  if (!data.serviceStartedAt) {
    return <Text>No start date or services yet</Text>;
  }
  if (!data?.serviceLineItems?.length) {
    return <Text>No services added yet</Text>;
  }

  const hasStartedService =
    data?.serviceStartedAt && !DateTimeUtils.isFuture(data.serviceStartedAt);

  const serviceStartDate =
    data?.serviceStartedAt && DateTimeUtils.isFuture(data.serviceStartedAt)
      ? DateTimeUtils.format(
          data.serviceStartedAt,
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
    data?.status !== ContractStatus.Ended
  ) {
    const serviceStarted = hasStartedService
      ? DateTimeUtils.format(
          data.serviceStartedAt,
          DateTimeUtils.dateWithAbreviatedMonth,
        )
      : null;

    return <Text>Service started {serviceStarted}</Text>;
  }

  const endDate = data?.endedAt
    ? DateTimeUtils.format(data.endedAt, DateTimeUtils.dateWithAbreviatedMonth)
    : null;

  const isActiveAndRenewable =
    hasStartedService &&
    data.status !== ContractStatus.Ended &&
    !!data.renewalCycle &&
    data.renewalCycle !== ContractRenewalCycle.None;

  return (
    <Flex flexDir='column' alignItems='flex-start' justifyContent='center'>
      {serviceStartDate && <Text>Service starts {serviceStartDate}</Text>}
      {isActiveAndRenewable && (
        <Text>
          Renews {getLabelFromValue(data.renewalCycle)} on {renewalDate}
        </Text>
      )}
      {data?.endedAt && DateTimeUtils.isFuture(data.endedAt) && (
        <Text>Ends {endDate}</Text>
      )}

      {data.status === ContractStatus.Ended && <Text>Ended on {endDate}</Text>}
    </Flex>
  );
};
