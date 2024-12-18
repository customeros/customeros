import { observer } from 'mobx-react-lite';
import ParentSize from '@visx/responsive/lib/components/ParentSize';

import { cn } from '@ui/utils/cn';
import { Skeleton } from '@ui/feedback/Skeleton';
import { useStore } from '@shared/hooks/useStore';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

import { HelpContent } from './HelpContent';
import { ChartCard } from '../../ChartCard';
import { PercentageTrend } from '../../PercentageTrend';
import { useOnboardingCompletionQuery } from '../../../graphql/onboardingCompletion.generated';
import OnboardingCompletionChart, {
  OnboardingCompletionDatum,
} from './OnboardingCompletion.chart';

export const OnboardingCompletion = observer(() => {
  const store = useStore();
  const client = getGraphQLClient();
  const { data, isLoading } = useOnboardingCompletionQuery(client);

  const hasContracts = store.globalCache?.value?.contractsExist ?? false;
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
});
