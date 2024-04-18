'use client';
import dynamic from 'next/dynamic';

import ParentSize from '@visx/responsive/lib/components/ParentSize';

import { cn } from '@ui/utils/cn';
import { Skeleton } from '@ui/feedback/Skeleton';
import { ChartCard } from '@customerMap/components/ChartCard';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import { useNewCustomersQuery } from '@customerMap/graphql/newCustomers.generated';

import { HelpContent } from './HelpContent';
import { PercentageTrend } from '../../PercentageTrend';
import { NewCustomersDatum } from './NewCustomers.chart';

const NewCustomersChart = dynamic(() => import('./NewCustomers.chart'), {
  ssr: false,
});

export const NewCustomers = () => {
  const client = getGraphQLClient();
  const { data: globalCache } = useGlobalCacheQuery(client);
  const { data, isLoading } = useNewCustomersQuery(client);

  const hasContracts = globalCache?.global_Cache?.contractsExist;
  const chartData = (data?.dashboard_NewCustomers?.perMonth ?? []).map((d) => ({
    month: d?.month,
    value: d?.count,
  })) as NewCustomersDatum[];

  const stat = `${data?.dashboard_NewCustomers?.thisMonthCount ?? 0}`;

  const percentage =
    data?.dashboard_NewCustomers?.thisMonthIncreasePercentage ?? '0';

  return (
    <ChartCard
      className='flex-1'
      stat={stat}
      title='New customers'
      hasData={hasContracts}
      renderHelpContent={HelpContent}
      renderSubStat={() => <PercentageTrend percentage={percentage} />}
    >
      <ParentSize>
        {({ width }) => {
          return (
            <>
              {isLoading && (
                <Skeleton
                  className={cn(isLoading ? 'h-[200px]' : 'h-full', 'w-full')}
                />
              )}
              <NewCustomersChart
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
