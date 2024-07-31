import { useFeatureIsOn } from '@growthbook/growthbook-react';
import ParentSize from '@visx/responsive/lib/components/ParentSize';

import { cn } from '@ui/utils/cn';
import { Skeleton } from '@ui/feedback/Skeleton';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { InfoDialog } from '@ui/overlay/AlertDialog/InfoDialog/InfoDialog';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';

import { HelpContent } from './HelpContent';
import { HelpButton } from '../../HelpButton';
import CustomerMapChart, { CustomerMapDatum } from './CustomerMap.chart';
import { useCustomerMapQuery } from '../../../graphql/customerMap.generated';

export const CustomerMap = () => {
  const client = getGraphQLClient();
  const { data, isLoading } = useCustomerMapQuery(client);
  const { data: globalCacheData } = useGlobalCacheQuery(client);
  const { open: isOpen, onOpen, onClose } = useDisclosure();
  const isTaller = useFeatureIsOn('taller-customer-map-chart');

  const chartData = (data?.dashboard_CustomerMap ?? []).map((d) => ({
    x: d?.contractSignedDate ? new Date(d?.contractSignedDate) : new Date(),
    r: d?.arr,
    values: {
      id: d?.organization?.id,
      name: d?.organization?.name || 'Unnamed',
      status: d?.state,
    },
  })) as CustomerMapDatum[];

  const hasContracts = globalCacheData?.global_Cache?.contractsExist;

  return (
    <div className='w-full group'>
      <ParentSize>
        {({ width }) => (
          <>
            <div className='flex flex-col relative'>
              <div className='flex gap-2 items-center mb-1'>
                <p className='font-semibold text-xl'>Customer map</p>
                <HelpButton isOpen={isOpen} onOpen={onOpen} />
              </div>

              {!hasContracts && (
                <p className='bottom-1 text-gray-400 font-semibold text-lg absolute transform translate-y-full'>
                  No data yet
                </p>
              )}
            </div>
            {isLoading && (
              <Skeleton
                className={cn(isTaller ? 'h-[700px]' : 'h-[350px]', 'w-full')}
              />
            )}

            <CustomerMapChart
              width={width}
              data={chartData}
              hasContracts={hasContracts}
              height={isTaller ? 700 : 350}
            />
            <InfoDialog
              isOpen={isOpen}
              onClose={onClose}
              onConfirm={onClose}
              label='Customer map'
              confirmButtonLabel='Got it'
            >
              <HelpContent />
            </InfoDialog>
          </>
        )}
      </ParentSize>
    </div>
  );
};
