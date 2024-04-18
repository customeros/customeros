'use client';
import dynamic from 'next/dynamic';

import ParentSize from '@visx/responsive/lib/components/ParentSize';

import { cn } from '@ui/utils/cn';
import { Skeleton } from '@ui/feedback/Skeleton';
import { ChartCard } from '@customerMap/components/ChartCard';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import { useArrBreakdownQuery } from '@customerMap/graphql/arrBreakdown.generated';

import { HelpContent } from './HelpContent';
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
  const { data: globalCacheData } = useGlobalCacheQuery(client);
  const { data, isLoading } = useArrBreakdownQuery(client);

  const hasContracts = globalCacheData?.global_Cache?.contractsExist;
  const chartData = (data?.dashboard_ARRBreakdown?.perMonth ?? []).map((d) => ({
    month: d?.month,
    upsells: d?.upsells,
    renewals: d?.renewals,
    newlyContracted: d?.newlyContracted,
    churned: d?.churned,
    cancellations: d?.cancellations,
    downgrades: d?.downgrades,
  })) as ARRBreakdownDatum[];

  const stat = formatCurrency(data?.dashboard_ARRBreakdown?.arrBreakdown ?? 0);
  const percentage = data?.dashboard_ARRBreakdown?.increasePercentage ?? '0';

  return (
    <ChartCard
      className='flex-3'
      stat={stat}
      title='ARR breakdown'
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
              <RevenueRetentionRateChart
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
