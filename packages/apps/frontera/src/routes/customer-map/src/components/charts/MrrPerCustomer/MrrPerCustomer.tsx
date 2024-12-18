import { observer } from 'mobx-react-lite';
import ParentSize from '@visx/responsive/lib/components/ParentSize';

import { Skeleton } from '@ui/feedback/Skeleton';
import { useStore } from '@shared/hooks/useStore';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber';

import { HelpContent } from './HelpContent';
import { ChartCard } from '../../ChartCard';
import { PercentageTrend } from '../../PercentageTrend';
import { useMrrPerCustomerQuery } from '../../../graphql/mrrPerCustomer.generated';
import MrrPerCustomerChart, {
  MrrPerCustomerDatum,
} from './MrrPerCustomer.chart';

export const MrrPerCustomer = observer(() => {
  const store = useStore();
  const client = getGraphQLClient();
  const { data, isLoading } = useMrrPerCustomerQuery(client);

  const hasContracts = store?.globalCache?.value?.contractsExist ?? false;
  const chartData = (data?.dashboard_MRRPerCustomer?.perMonth ?? []).map(
    (d, index, arr) => {
      const decIndex = arr.findIndex((d) => d?.month === 12);

      return {
        month: d?.month,
        value: d?.value,
        index: decIndex > index - 1 ? 1 : 2,
      };
    },
  ) as MrrPerCustomerDatum[];
  const stat = formatCurrency(
    data?.dashboard_MRRPerCustomer?.mrrPerCustomer ?? 0,
  );
  const percentage = data?.dashboard_MRRPerCustomer?.increasePercentage ?? '0';

  return (
    <ChartCard
      stat={stat}
      className='flex-1'
      hasData={hasContracts}
      title='MRR per customer'
      renderHelpContent={HelpContent}
      renderSubStat={() => <PercentageTrend percentage={percentage} />}
    >
      <ParentSize>
        {({ width }) => {
          return (
            <>
              {isLoading && <Skeleton className='h-[200px] w-full' />}

              <MrrPerCustomerChart
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
