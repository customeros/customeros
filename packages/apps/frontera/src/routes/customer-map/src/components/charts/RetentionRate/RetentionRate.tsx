import ParentSize from '@visx/responsive/lib/components/ParentSize';

import { cn } from '@ui/utils/cn';
import { Skeleton } from '@ui/feedback/Skeleton';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';

import { HelpContent } from './HelpContent';
import { ChartCard } from '../../ChartCard';
import { PercentageTrend } from '../../PercentageTrend';
import RetentionRateChart, { RetentionRateDatum } from './RetentionRate.chart';
import { useRetentionRateQuery } from '../../../graphql/retentionRate.generated';

export const RetentionRate = () => {
  const client = getGraphQLClient();
  const { data: globalCache } = useGlobalCacheQuery(client);
  const { data, isLoading } = useRetentionRateQuery(client);

  const hasContracts = globalCache?.global_Cache?.contractsExist;
  const chartData = (data?.dashboard_RetentionRate?.perMonth ?? []).map(
    (d) => ({
      month: d?.month,
      values: {
        churned: d?.churnCount ?? 0,
        renewed: d?.renewCount ?? 0,
      },
    }),
  ) as RetentionRateDatum[];

  const stat = `${data?.dashboard_RetentionRate?.retentionRate ?? 0}%`;
  const percentage = data?.dashboard_RetentionRate?.increasePercentage ?? '0';

  return (
    <ChartCard
      stat={stat}
      className='flex-1'
      title='Retention rate'
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

              <RetentionRateChart
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
