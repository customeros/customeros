'use client';
import dynamic from 'next/dynamic';

import ParentSize from '@visx/responsive/lib/components/ParentSize';

import { cn } from '@ui/utils/cn';
import { Skeleton } from '@ui/feedback/Skeleton';
import { ChartCard } from '@customerMap/components/ChartCard';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import { useTimeToOnboardQuery } from '@customerMap/graphql/timeToOnboard.generated';

import { HelpContent } from './HelpContent';
import { PercentageTrend } from '../../PercentageTrend';
import { TimeToOnboardDatum } from './TimeToOnboard.chart';

const TimeToOnboardChart = dynamic(() => import('./TimeToOnboard.chart'), {
  ssr: false,
});

export const TimeToOnboard = () => {
  const client = getGraphQLClient();
  const { data: globalCacheData } = useGlobalCacheQuery(client);
  const { data, isLoading } = useTimeToOnboardQuery(client);

  const hasContracts = globalCacheData?.global_Cache?.contractsExist;
  const chartData = (data?.dashboard_TimeToOnboard?.perMonth ?? []).map(
    (d, index, arr) => {
      const decIndex = arr.findIndex((d) => d.month === 12);

      return {
        month: d?.month,
        value: d?.value,
        index: decIndex > index - 1 ? 1 : 2,
      };
    },
  ) as TimeToOnboardDatum[];

  const statValue = data?.dashboard_TimeToOnboard?.timeToOnboard ?? 0;
  const stat = `${statValue} ${statValue === 1 ? 'day' : 'days'}`;
  const percentage = `${
    data?.dashboard_TimeToOnboard?.increasePercentage ?? 0
  }%`;

  return (
    <ChartCard
      className='flex-1'
      stat={stat}
      hasData={hasContracts}
      title='Time to onboard'
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
              <TimeToOnboardChart
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
