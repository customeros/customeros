'use client';
import dynamic from 'next/dynamic';

import { ChartCard } from '@dashboard/components/ChartCard';
import ParentSize from '@visx/responsive/lib/components/ParentSize';
// import { useMrrPerCustomerQuery } from '@dashboard/graphql/.generated';

// import { getGraphQLClient } from '@shared/util/getGraphQLClient';
// import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';

import { PercentageTrend } from '../../PercentageTrend';
// import {} from './NewCustomers.chart';

const NewCustomersChart = dynamic(() => import('./NewCustomers.chart'), {
  ssr: false,
});

export const NewCustomers = () => {
  // const client = getGraphQLClient();
  // const { data } = useMrrPerCustomerQuery(client, {
  //   year: 2023,
  // });

  // const chartData = (data?.dashboard_MRRPerCustomer?.perMonth ?? []).map(
  //   (d) => ({
  //     month: d?.month,
  //     value: d?.value,
  //   }),
  // );

  // const stat = formatCurrency(
  //   data?.dashboard_MRRPerCustomer?.mrrPerCustomer ?? 0,
  // );
  // const percentage = data?.dashboard_MRRPerCustomer?.increasePercentage ?? 0;

  return (
    <ChartCard
      flex='1'
      stat={'127'}
      title='New Customers'
      renderSubStat={() => <PercentageTrend percentage={2} />}
    >
      <ParentSize>
        {({ width }) => <NewCustomersChart width={width} />}
      </ParentSize>
    </ChartCard>
  );
};
