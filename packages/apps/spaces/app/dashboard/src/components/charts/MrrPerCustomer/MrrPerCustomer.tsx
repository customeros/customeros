'use client';
import dynamic from 'next/dynamic';

import { ChartCard } from '@dashboard/components/ChartCard';
import ParentSize from '@visx/responsive/lib/components/ParentSize';
import { useMrrPerCustomerQuery } from '@dashboard/graphql/mrrPerCustomer.generated';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';

import { PercentageTrend } from '../../PercentageTrend';
import { MrrPerCustomerDatum } from './MrrPerCustomer.chart';

const MrrPerCustomerChart = dynamic(() => import('./MrrPerCustomer.chart'), {
  ssr: false,
});

export const MrrPerCustomer = () => {
  const client = getGraphQLClient();
  const { data } = useMrrPerCustomerQuery(client, {
    year: 2023,
  });

  const chartData = (data?.dashboard_MRRPerCustomer?.perMonth ?? []).map(
    (d) => ({
      month: d?.month,
      value: d?.value,
    }),
  ) as MrrPerCustomerDatum[];
  const stat = formatCurrency(
    data?.dashboard_MRRPerCustomer?.mrrPerCustomer ?? 0,
  );
  const percentage = data?.dashboard_MRRPerCustomer?.increasePercentage ?? 0;

  return (
    <ChartCard
      flex='1'
      stat={stat}
      title='MRR per Customer'
      renderSubStat={() => <PercentageTrend percentage={percentage} />}
    >
      <ParentSize>
        {({ width }) => <MrrPerCustomerChart width={width} data={chartData} />}
      </ParentSize>
    </ChartCard>
  );
};
