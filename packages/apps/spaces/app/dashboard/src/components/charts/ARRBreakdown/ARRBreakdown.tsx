'use client';
import dynamic from 'next/dynamic';

import { ChartCard } from '@dashboard/components/ChartCard';
import ParentSize from '@visx/responsive/lib/components/ParentSize';
import { useArrBreakdownQuery } from '@dashboard/graphql/arrBreakdown.generated';

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
  const { data } = useArrBreakdownQuery(client, {
    year: 2023,
  });

  const chartData = (data?.dashboard_ARRBreakdown?.perMonth ?? []).map((d) => ({
    month: d?.month,
    values: {
      upsells: d?.upsells,
      churned: d?.churned,
      renewals: d?.renewals,
      downgrades: d?.downgrades,
      cancellations: d?.cancellations,
      newlyContracted: d?.newlyContracted,
    },
  })) as ARRBreakdownDatum[];

  const stat = formatCurrency(data?.dashboard_ARRBreakdown?.arrBreakdown ?? 0);
  const percentage = data?.dashboard_ARRBreakdown?.increasePercentage ?? 0;

  return (
    <ChartCard
      flex='3'
      stat={stat}
      title='ARR Breakdown'
      renderSubStat={() => <PercentageTrend percentage={percentage} />}
    >
      <ParentSize>
        {({ width }) => (
          <RevenueRetentionRateChart width={width} data={chartData} />
        )}
      </ParentSize>
    </ChartCard>
  );
};
