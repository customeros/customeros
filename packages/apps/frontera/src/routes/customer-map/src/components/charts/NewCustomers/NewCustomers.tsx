import { observer } from 'mobx-react-lite';
import ParentSize from '@visx/responsive/lib/components/ParentSize';

import { cn } from '@ui/utils/cn';
import { Skeleton } from '@ui/feedback/Skeleton';
import { useStore } from '@shared/hooks/useStore';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

import { HelpContent } from './HelpContent';
import { ChartCard } from '../../ChartCard';
import { PercentageTrend } from '../../PercentageTrend';
import NewCustomersChart, { NewCustomersDatum } from './NewCustomers.chart';
import { useNewCustomersQuery } from '../../../graphql/newCustomers.generated';

export const NewCustomers = observer(() => {
  const store = useStore();
  const client = getGraphQLClient();
  const { data, isLoading } = useNewCustomersQuery(client);

  const hasContracts = store.globalCache?.value?.contractsExist ?? false;
  const chartData = (data?.dashboard_NewCustomers?.perMonth ?? []).map((d) => ({
    month: d?.month,
    value: d?.count,
  })) as NewCustomersDatum[];

  const stat = `${data?.dashboard_NewCustomers?.thisMonthCount ?? 0}`;

  const percentage =
    data?.dashboard_NewCustomers?.thisMonthIncreasePercentage ?? '0';

  return (
    <ChartCard
      stat={stat}
      className='flex-1'
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
});
