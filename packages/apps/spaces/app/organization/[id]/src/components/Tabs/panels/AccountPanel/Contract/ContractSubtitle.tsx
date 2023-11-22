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
  const endDate = data?.endedAt
    ? DateTimeUtils.format(data.endedAt, DateTimeUtils.dateWithAbreviatedMonth)
    : null;

  if (!data?.signedAt) {
    return <Text>No start date or services yet</Text>;
  }

  return (
    <Flex flexDir='column' alignItems='flex-start' justifyContent='center'>
      {serviceStartDate && <Text>Service starts on {serviceStartDate}</Text>}
      {hasStartedService && data.status !== ContractStatus.Ended && (
        <Text>
          Renews {getLabelFromValue(data.renewalCycle)} on {renewalDate}
        </Text>
      )}
      {data?.endedAt && DateTimeUtils.isFuture(data.endedAt) && (
        <Text>Ends on {endDate}</Text>
      )}
      {!DateTimeUtils.isFuture(data.endedAt) && <Text>Ended on {endDate}</Text>}
    </Flex>
  );
};
