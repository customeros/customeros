'use client';
import dynamic from 'next/dynamic';

import ParentSize from '@visx/responsive/lib/components/ParentSize';
import { useCustomerMapQuery } from '@customerMap/graphql/customerMap.generated';

import { Skeleton } from '@ui/presentation/Skeleton';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

import { CustomerMapDatum } from './CustomerMap.chart';

const CustomerMapChart = dynamic(() => import('./CustomerMap.chart'), {
  ssr: false,
});

export const CustomerMap = () => {
  const client = getGraphQLClient();
  const { data, isLoading } = useCustomerMapQuery(client);

  const chartData = (data?.dashboard_CustomerMap ?? []).map((d) => ({
    x: d?.contractSignedDate ? new Date(d?.contractSignedDate) : new Date(),
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
        <Skeleton
          w='full'
          h='350px'
          endColor='gray.300'
          startColor='gray.300'
          isLoaded={!isLoading}
        >
          <CustomerMapChart width={width} height={350} data={chartData} />
        </Skeleton>
      )}
    </ParentSize>
  );
};
