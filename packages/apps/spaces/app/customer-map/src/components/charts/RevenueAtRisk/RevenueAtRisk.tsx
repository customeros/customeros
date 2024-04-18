'use client';

import ParentSize from '@visx/responsive/lib/components/ParentSize';

import { cn } from '@ui/utils/cn';
import { Skeleton } from '@ui/feedback/Skeleton';
import { ChartCard } from '@customerMap/components/ChartCard';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import { useRevenueAtRiskQuery } from '@customerMap/graphql/revenueAtRisk.generated';

import { HelpContent } from './HelpContent';
import RevenueAtRiskChart, { RevenueAtRiskDatum } from './RevenueAtRisk.chart';

export const RevenueAtRisk = () => {
  const client = getGraphQLClient();
  const { data: globalCache } = useGlobalCacheQuery(client);
  const { data, isLoading } = useRevenueAtRiskQuery(client);
  const hasContracts = globalCache?.global_Cache?.contractsExist;
  const chartData: RevenueAtRiskDatum = {
    atRisk: data?.dashboard_RevenueAtRisk?.atRisk ?? 0,
    highConfidence: data?.dashboard_RevenueAtRisk?.highConfidence ?? 0,
  };

  return (
    <ChartCard
      className='flex-1'
      stat={!hasContracts ? 'N/A' : undefined}
      hasData={hasContracts}
      title='Revenue at risk'
      renderHelpContent={HelpContent}
      renderSubStat={() => (
        <div className='flex mt-4 justify-between'>
          <div className='flex flex-col'>
            <div className='flex items-center gap-2'>
              <div className='flex w-2 h-2 bg-greenLight-500 rounded-full' />
              <p className='text-sm'>High Confidence</p>
            </div>
            <p className='text-sm'>
              {formatCurrency(chartData.highConfidence)}
            </p>
          </div>

          <div className='flex flex-col'>
            <div className='flex gap-2 items-center'>
              <div className='flex w-2 h-2 bg-warning-300 rounded-full' />
              <p className='text-sm text-gray-500'>At Risk</p>
            </div>
            <p className='text-sm'>{formatCurrency(chartData.atRisk)}</p>
          </div>
        </div>
      )}
    >
      <ParentSize>
        {({ width, height }) => {
          return (
            <>
              <Skeleton
                className={cn(isLoading ? 'h-[200px]' : 'h-full', 'w-full')}
                isLoaded={!isLoading}
              >
                <RevenueAtRiskChart
                  width={width}
                  height={height}
                  data={chartData}
                  hasContracts={hasContracts}
                />
              </Skeleton>
            </>
          );
        }}
      </ParentSize>
    </ChartCard>
  );
};
