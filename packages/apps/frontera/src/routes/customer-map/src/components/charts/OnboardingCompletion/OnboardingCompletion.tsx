import ParentSize from '@visx/responsive/lib/components/ParentSize';

import { cn } from '@ui/utils/cn';
import { Skeleton } from '@ui/feedback/Skeleton';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';

import { HelpContent } from './HelpContent';
import { ChartCard } from '../../ChartCard';
import { PercentageTrend } from '../../PercentageTrend';
import { useOnboardingCompletionQuery } from '../../../graphql/onboardingCompletion.generated';
import OnboardingCompletionChart, {
  OnboardingCompletionDatum,
} from './OnboardingCompletion.chart';

export const OnboardingCompletion = () => {
  const client = getGraphQLClient();
  const { data: globalCache } = useGlobalCacheQuery(client);
  const { data, isLoading } = useOnboardingCompletionQuery(client);

  const hasContracts = globalCache?.global_Cache?.contractsExist;
  const chartData = (data?.dashboard_OnboardingCompletion?.perMonth ?? []).map(
    (d) => ({
      month: d?.month,
      value: d?.value,
    }),
  ) as OnboardingCompletionDatum[];

  const stat = `${
    data?.dashboard_OnboardingCompletion?.completionPercentage ?? 0
  }%`;

  const percentage = `${
    data?.dashboard_OnboardingCompletion?.increasePercentage ?? 0
  }`;

  return (
    <ChartCard
      stat={stat}
      className='flex-1'
      hasData={hasContracts}
      title='Onboarding completion'
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
              <OnboardingCompletionChart
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
