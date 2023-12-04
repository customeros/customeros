'use client';
import dynamic from 'next/dynamic';

import ParentSize from '@visx/responsive/lib/components/ParentSize';

import { Skeleton } from '@ui/presentation/Skeleton';
import { ChartCard } from '@customerMap/components/ChartCard';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import { useRetentionRateQuery } from '@customerMap/graphql/retentionRate.generated';

import { HelpContent } from './HelpContent';
import { PercentageTrend } from '../../PercentageTrend';
import { RetentionRateDatum } from './RetentionRate.chart';

const RetentionRateChart = dynamic(() => import('./RetentionRate.chart'), {
  ssr: false,
});

export const RetentionRate = () => {
  const client = getGraphQLClient();
  const { data: globalCache } = useGlobalCacheQuery(client);
  const { data, isLoading } = useRetentionRateQuery(client);

  const hasContracts = globalCache?.global_Cache?.contractsExist;
  const chartData = (data?.dashboard_RetentionRate?.perMonth ?? []).map(
    (d) => ({
      month: d?.month,
      values: {
        churned: d?.churnCount ?? 0,
        renewed: d?.renewCount ?? 0,
      },
    }),
  ) as RetentionRateDatum[];

  const stat = `${data?.dashboard_RetentionRate?.retentionRate ?? 0}%`;
  const percentage = data?.dashboard_RetentionRate?.increasePercentage ?? 0;

  return (
    <ChartCard
      flex='1'
      stat={stat}
      title='Retention rate'
      hasData={hasContracts}
      renderHelpContent={HelpContent}
      renderSubStat={() => <PercentageTrend percentage={percentage} />}
    >
      <ParentSize>
        {({ width }) => (
          <Skeleton
            w='full'
            h={isLoading ? '200px' : 'full'}
            endColor='gray.300'
            startColor='gray.300'
            isLoaded={!isLoading}
          >
            <RetentionRateChart
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
