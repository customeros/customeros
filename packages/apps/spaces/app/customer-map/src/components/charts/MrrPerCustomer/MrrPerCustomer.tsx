'use client';
import dynamic from 'next/dynamic';

import ParentSize from '@visx/responsive/lib/components/ParentSize';

import { Skeleton } from '@ui/presentation/Skeleton';
import { ChartCard } from '@customerMap/components/ChartCard';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import { useMrrPerCustomerQuery } from '@customerMap/graphql/mrrPerCustomer.generated';

import { HelpContent } from './HelpContent';
import { PercentageTrend } from '../../PercentageTrend';
import { MrrPerCustomerDatum } from './MrrPerCustomer.chart';

const MrrPerCustomerChart = dynamic(() => import('./MrrPerCustomer.chart'), {
  ssr: false,
});

export const MrrPerCustomer = () => {
  const client = getGraphQLClient();
  const { data: globalCacheData } = useGlobalCacheQuery(client);
  const { data, isLoading } = useMrrPerCustomerQuery(client);

  const hasContracts = globalCacheData?.global_Cache?.contractsExist;
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
      hasData={hasContracts}
      title='MRR per customer'
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
            <MrrPerCustomerChart
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
