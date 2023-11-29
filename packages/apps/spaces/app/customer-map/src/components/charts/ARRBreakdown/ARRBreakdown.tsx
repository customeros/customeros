'use client';
import dynamic from 'next/dynamic';

import { ChartCard } from '@customerMap/components/ChartCard';
import ParentSize from '@visx/responsive/lib/components/ParentSize';
import { useArrBreakdownQuery } from '@customerMap/graphql/arrBreakdown.generated';

import { Skeleton } from '@ui/presentation/Skeleton';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';

import { PercentageTrend } from '../../PercentageTrend';
import { ARRBreakdownDatum } from './ARRBreakdown.chart';

const RevenueRetentionRateChart = dynamic(
  () => import('./ARRBreakdown.chart'),
  {
    ssr: false,
  },
);

export const ARRBreakdown = () => {
  const client = getGraphQLClient();
  const { data, isLoading } = useArrBreakdownQuery(client);

  const chartData = (data?.dashboard_ARRBreakdown?.perMonth ?? []).map((d) => ({
    month: d?.month,
    upsells: d?.upsells,
    renewals: d?.renewals,
    newlyContracted: d?.newlyContracted,
    churned: d?.churned,
    cancellations: d?.cancellations,
    downgrades: d?.downgrades,
  })) as ARRBreakdownDatum[];

  const stat = formatCurrency(data?.dashboard_ARRBreakdown?.arrBreakdown ?? 0);
  const percentage = data?.dashboard_ARRBreakdown?.increasePercentage ?? 0;

  return (
    <ChartCard
      flex='3'
      stat={stat}
      title='ARR breakdown'
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
            <RevenueRetentionRateChart width={width} data={chartData} />
          </Skeleton>
        )}
      </ParentSize>
    </ChartCard>
  );
};
