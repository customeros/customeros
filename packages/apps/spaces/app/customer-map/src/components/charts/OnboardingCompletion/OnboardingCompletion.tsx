'use client';
import dynamic from 'next/dynamic';

import ParentSize from '@visx/responsive/lib/components/ParentSize';

import { Skeleton } from '@ui/presentation/Skeleton';
import { ChartCard } from '@customerMap/components/ChartCard';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import { useOnboardingCompletionQuery } from '@customerMap/graphql/onboardingCompletion.generated';

import { HelpContent } from './HelpContent';
import { PercentageTrend } from '../../PercentageTrend';
import { OnboardingCompletionDatum } from './OnboardingCompletion.chart';

const OnboardingCompletionChart = dynamic(
  () => import('./OnboardingCompletion.chart'),
  {
    ssr: false,
  },
);

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
  }`;

  const percentage = `${
    data?.dashboard_OnboardingCompletion?.increasePercentage ?? 0
  }`;

  return (
    <ChartCard
      flex='1'
      stat={stat}
      title='Onboarding completion'
      hasData={hasContracts}
      renderHelpContent={HelpContent}
      renderSubStat={() => <PercentageTrend percentage={percentage} />}
    >
      <ParentSize>
        {({ width }) => (
          <Skeleton
            w='full'
            h={isLoading ? '200px' : 'full'}
            endColor='gray.300'
            startColor='gray.300'
            isLoaded={!isLoading}
          >
            <OnboardingCompletionChart
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
