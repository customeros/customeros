'use client';
import dynamic from 'next/dynamic';

import ParentSize from '@visx/responsive/lib/components/ParentSize';

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
    (d) => ({
      month: d?.month,
      value: d?.percentage,
    }),
  ) as GrossRevenueRetentionDatum[];

  const stat = `${
    data?.dashboard_GrossRevenueRetention?.grossRevenueRetention ?? 0
  }%`;

  const percentage =
    data?.dashboard_GrossRevenueRetention?.increasePercentage ?? 0;

  return (
    <ChartCard
      flex='2'
      stat={stat}
      hasData={hasContracts}
      title='Gross Revenue Retention'
      renderHelpContent={HelpContent}
      renderSubStat={() => <PercentageTrend percentage={percentage} />}
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
              hasContracts={hasContracts}
            />
          </Skeleton>
        )}
      </ParentSize>
    </ChartCard>
  );
};
