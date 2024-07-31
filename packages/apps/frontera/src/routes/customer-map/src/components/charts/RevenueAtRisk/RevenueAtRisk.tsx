import ParentSize from '@visx/responsive/lib/components/ParentSize';

import { Skeleton } from '@ui/feedback/Skeleton';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';

import { HelpContent } from './HelpContent';
import { ChartCard } from '../../ChartCard';
import RevenueAtRiskChart, { RevenueAtRiskDatum } from './RevenueAtRisk.chart';
import { useRevenueAtRiskQuery } from '../../../graphql/revenueAtRisk.generated';

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
      hasData={hasContracts}
      title='Revenue at risk'
      renderHelpContent={HelpContent}
      stat={!hasContracts ? 'N/A' : undefined}
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
              {isLoading ? (
                <Skeleton className='h-[200px]' />
              ) : (
                <RevenueAtRiskChart
                  width={width}
                  height={height}
                  data={chartData}
                  hasContracts={hasContracts}
                />
              )}
            </>
          );
        }}
      </ParentSize>
    </ChartCard>
  );
};
