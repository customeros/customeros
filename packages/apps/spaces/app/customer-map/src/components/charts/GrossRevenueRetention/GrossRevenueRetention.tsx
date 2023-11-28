'use client';
import dynamic from 'next/dynamic';

import { ChartCard } from '@customerMap/components/ChartCard';
import ParentSize from '@visx/responsive/lib/components/ParentSize';
import { useGrossRevenueRetentionQuery } from '@customerMap/graphql/grossRevenueRetention.generated';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';

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
  const { data } = useGrossRevenueRetentionQuery(client);

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
      title='Gross Revenue Retention'
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
