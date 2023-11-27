'use client';
import dynamic from 'next/dynamic';

import ParentSize from '@visx/responsive/lib/components/ParentSize';
import { useCustomerMapQuery } from '@dashboard/graphql/customerMap.generated';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';

import { CustomerMapDatum } from './CustomerMap.chart';

const CustomerMapChart = dynamic(() => import('./CustomerMap.chart'), {
  ssr: false,
});

export const CustomerMap = () => {
  const client = getGraphQLClient();
  const { data } = useCustomerMapQuery(client);

  const chartData = (data?.dashboard_CustomerMap ?? []).map((d) => ({
    x: new Date(d?.contractSignedDate),
    r: d?.arr,
    values: {
      id: d?.organization?.id,
      name: d?.organization?.name || 'Unnamed',
      status: d?.state,
    },
  })) as CustomerMapDatum[];

  return (
    <ParentSize>
      {({ width }) => (
        <CustomerMapChart width={width} height={350} data={chartData} />
      )}
    </ParentSize>
  );
};
