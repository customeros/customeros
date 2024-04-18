'use client';
import dynamic from 'next/dynamic';

import ParentSize from '@visx/responsive/lib/components/ParentSize';

import { Skeleton } from '@ui/feedback/Skeleton';
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
    (d, index, arr) => {
      const decIndex = arr.findIndex((d) => d?.month === 12);

      return {
        month: d?.month,
        value: d?.value,
        index: decIndex > index - 1 ? 1 : 2,
      };
    },
  ) as MrrPerCustomerDatum[];
  const stat = formatCurrency(
    data?.dashboard_MRRPerCustomer?.mrrPerCustomer ?? 0,
  );
  const percentage = data?.dashboard_MRRPerCustomer?.increasePercentage ?? '0';

  return (
    <ChartCard
      className='flex-1'
      stat={stat}
      hasData={hasContracts}
      title='MRR per customer'
      renderHelpContent={HelpContent}
      renderSubStat={() => <PercentageTrend percentage={percentage} />}
    >
      <ParentSize>
        {({ width }) => {
          return (
            <>
              {isLoading && <Skeleton className='h-[200px] w-full' />}

              <MrrPerCustomerChart
                width={width}
                data={chartData}
                hasContracts={hasContracts}
              />
            </>
          );
        }}
      </ParentSize>
    </ChartCard>
  );
};
