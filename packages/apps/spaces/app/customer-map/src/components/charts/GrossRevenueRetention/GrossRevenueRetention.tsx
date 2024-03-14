'use client';
import dynamic from 'next/dynamic';

import ParentSize from '@visx/responsive/lib/components/ParentSize';

import { Text } from '@ui/typography/Text';
import { Skeleton } from '@ui/presentation/Skeleton';
import { ChartCard } from '@customerMap/components/ChartCard';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import { useGrossRevenueRetentionQuery } from '@customerMap/graphql/grossRevenueRetention.generated';

import { HelpContent } from './HelpContent';
import { PercentageTrend } from '../../PercentageTrend';
import { GrossRevenueRetentionDatum } from './GrossRevenueRetention.chart';

const RevenueRetentionRateChart = dynamic(
  () => import('./GrossRevenueRetention.chart'),
  {
    ssr: false,
  },
);

export const GrossRevenueRetention = () => {
  const client = getGraphQLClient();
  const { data: globalCacheData } = useGlobalCacheQuery(client);
  const { data, isLoading } = useGrossRevenueRetentionQuery(client);

  const hasContracts = globalCacheData?.global_Cache?.contractsExist;
  const chartData = (data?.dashboard_GrossRevenueRetention?.perMonth ?? []).map(
    (d, index, arr) => {
      const decIndex = arr.findIndex((d) => d?.month === 12);

      return {
        month: d?.month,
        value: d?.percentage,
        index: decIndex > index - 1 ? 1 : 2,
      };
    },
  ) as GrossRevenueRetentionDatum[];
  const hasMissingData = chartData.every(
    (d) => d.value === 0 || d.value === 100,
  );

  const stat = `${
    data?.dashboard_GrossRevenueRetention?.grossRevenueRetention ?? 0
  }%`;

  const percentage =
    data?.dashboard_GrossRevenueRetention?.increasePercentage ?? '0';

  return (
    <ChartCard
      flex='2'
      stat={hasMissingData ? undefined : stat}
      hasData={hasContracts}
      title='Gross Revenue Retention'
      renderHelpContent={HelpContent}
      renderSubStat={
        hasMissingData
          ? () => (
              <Text fontWeight='semibold' color='gray.400' mb='10'>
                Key data missing.
              </Text>
            )
          : () => <PercentageTrend percentage={percentage} />
      }
    >
      <ParentSize>
        {({ width }) => (
          <Skeleton
            w='full'
            h='200px'
            endColor='gray.300'
            startColor='gray.300'
            isLoaded={!isLoading}
          >
            <RevenueRetentionRateChart
              width={width}
              data={chartData}
              hasContracts={hasMissingData ? false : hasContracts}
            />
          </Skeleton>
        )}
      </ParentSize>
    </ChartCard>
  );
};
