'use client';
import dynamic from 'next/dynamic';

import { useFeatureIsOn } from '@growthbook/growthbook-react';
import ParentSize from '@visx/responsive/lib/components/ParentSize';

import { cn } from '@ui/utils/cn';
import { useDisclosure } from '@ui/utils';
import { Skeleton } from '@ui/feedback/Skeleton';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { InfoDialog } from '@ui/overlay/AlertDialog/InfoDialog/InfoDialog2';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import { useCustomerMapQuery } from '@customerMap/graphql/customerMap.generated';

import { HelpContent } from './HelpContent';
import { HelpButton } from '../../HelpButton';
import { CustomerMapDatum } from './CustomerMap.chart';

const CustomerMapChart = dynamic(() => import('./CustomerMap.chart'), {
  ssr: false,
});

export const CustomerMap = () => {
  const client = getGraphQLClient();
  const { data, isLoading } = useCustomerMapQuery(client);
  const { data: globalCacheData } = useGlobalCacheQuery(client);
  const { isOpen, onOpen, onClose } = useDisclosure();
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
              <div className='flex gap-2 items-center'>
                <p className='font-semibold text-xl'>Customer map</p>
                <HelpButton isOpen={isOpen} onOpen={onOpen} />
              </div>
              {!hasContracts && (
                <p className='bottom-0 text-gray-400 font-semibold text-lg absolute transform translate-y-full'>
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
              height={isTaller ? 700 : 350}
              data={chartData}
              hasContracts={hasContracts}
            />
            <InfoDialog
              label='Customer map'
              isOpen={isOpen}
              onClose={onClose}
              onConfirm={onClose}
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
