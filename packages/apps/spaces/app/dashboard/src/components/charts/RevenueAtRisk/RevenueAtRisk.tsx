'use client';
import dynamic from 'next/dynamic';

import { ChartCard } from '@dashboard/components/ChartCard';
import ParentSize from '@visx/responsive/lib/components/ParentSize';
import { useRevenueAtRiskQuery } from '@dashboard/graphql/revenueAtRisk.generated';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';

import { RevenueAtRiskDatum } from './RevenueAtRisk.chart';

const RevenueAtRiskChart = dynamic(() => import('./RevenueAtRisk.chart'), {
  ssr: false,
});

export const RevenueAtRisk = () => {
  const client = getGraphQLClient();
  const { data } = useRevenueAtRiskQuery(client, {
    year: 2023,
  });

  const chartData: RevenueAtRiskDatum = {
    atRisk: data?.dashboard_RevenueAtRisk?.atRisk ?? 0,
    highConfidence: data?.dashboard_RevenueAtRisk?.highConfidence ?? 0,
  };

  return (
    <ChartCard
      flex='1'
      title='Revenue at Risk'
      renderSubStat={() => (
        <Flex mt='4' justify='space-between'>
          <Flex flexDir='column'>
            <Flex gap='3' align='center'>
              <Flex w='3' h='3' bg='#66C61C' borderRadius='full' />
              <Text>High Confidence</Text>
            </Flex>
            <Text fontSize='sm' fontWeight='medium'>
              {formatCurrency(chartData.highConfidence)}
            </Text>
          </Flex>

          <Flex flexDir='column'>
            <Flex gap='3' align='center'>
              <Flex w='3' h='3' bg='yellow.400' borderRadius='full' />
              <Text>At Risk</Text>
            </Flex>
            <Text fontSize='sm'>{formatCurrency(chartData.atRisk)}</Text>
          </Flex>
        </Flex>
      )}
    >
      <ParentSize>
        {({ width, height }) => (
          <RevenueAtRiskChart width={width} height={height} data={chartData} />
        )}
      </ParentSize>
    </ChartCard>
  );
};
