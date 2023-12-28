'use client';
import dynamic from 'next/dynamic';

import ParentSize from '@visx/responsive/lib/components/ParentSize';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Skeleton } from '@ui/presentation/Skeleton';
import { ChartCard } from '@customerMap/components/ChartCard';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import { useRevenueAtRiskQuery } from '@customerMap/graphql/revenueAtRisk.generated';

import { HelpContent } from './HelpContent';
import { RevenueAtRiskDatum } from './RevenueAtRisk.chart';

const RevenueAtRiskChart = dynamic(() => import('./RevenueAtRisk.chart'), {
  ssr: false,
});

export const RevenueAtRisk = () => {
  const client = getGraphQLClient();
  const { data: globalCache } = useGlobalCacheQuery(client);
  const { data, isLoading } = useRevenueAtRiskQuery(client);

  const hasContracts = globalCache?.global_Cache?.contractsExist;
  const chartData: RevenueAtRiskDatum = {
    atRisk: data?.dashboard_RevenueAtRisk?.atRisk ?? 0,
    highConfidence: data?.dashboard_RevenueAtRisk?.highConfidence ?? 0,
  };

  return (
    <ChartCard
      flex='1'
      stat={!hasContracts ? 'N/A' : undefined}
      hasData={hasContracts}
      title='Revenue at risk'
      renderHelpContent={HelpContent}
      renderSubStat={() => (
        <Flex mt='4' justify='space-between'>
          <Flex flexDir='column'>
            <Flex gap='2' align='center'>
              <Flex w='2' h='2' bg='greenLight.500' borderRadius='full' />
              <Text fontSize='sm'>High Confidence</Text>
            </Flex>
            <Text fontSize='sm'>
              {formatCurrency(chartData.highConfidence)}
            </Text>
          </Flex>

          <Flex flexDir='column'>
            <Flex gap='2' align='center'>
              <Flex w='2' h='2' bg='warning.300' borderRadius='full' />
              <Text fontSize='sm' color='gray.500'>
                At Risk
              </Text>
            </Flex>
            <Text fontSize='sm'>{formatCurrency(chartData.atRisk)}</Text>
          </Flex>
        </Flex>
      )}
    >
      <ParentSize>
        {({ width, height }) => (
          <Skeleton
            w='full'
            h={isLoading ? '200px' : 'full'}
            endColor='gray.300'
            startColor='gray.300'
            isLoaded={!isLoading}
          >
            <RevenueAtRiskChart
              width={width}
              height={height}
              data={chartData}
              hasContracts={hasContracts}
            />
          </Skeleton>
        )}
      </ParentSize>
    </ChartCard>
  );
};
