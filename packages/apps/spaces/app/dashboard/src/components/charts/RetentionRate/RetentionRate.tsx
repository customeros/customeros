'use client';
import dynamic from 'next/dynamic';

import { ChartCard } from '@dashboard/components/ChartCard';
import ParentSize from '@visx/responsive/lib/components/ParentSize';
import { useRetentionRateQuery } from '@dashboard/graphql/retentionRate.generated';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';

import { PercentageTrend } from '../../PercentageTrend';
import { RetentionRateDatum } from './RetentionRate.chart';

const RetentionRateChart = dynamic(() => import('./RetentionRate.chart'), {
  ssr: false,
});

export const RetentionRate = () => {
  const client = getGraphQLClient();
  const { data } = useRetentionRateQuery(client, {
    year: 2023,
  });

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
      title='Retention Rate'
      renderSubStat={() => <PercentageTrend percentage={percentage} />}
    >
      <ParentSize>
        {({ width }) => <RetentionRateChart width={width} data={chartData} />}
      </ParentSize>
    </ChartCard>
  );
};
