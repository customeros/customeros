'use client';
import dynamic from 'next/dynamic';

import { useFeatureIsOn } from '@growthbook/growthbook-react';
import ParentSize from '@visx/responsive/lib/components/ParentSize';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { useDisclosure } from '@ui/utils';
import { Text } from '@ui/typography/Text';
import { Skeleton } from '@ui/presentation/Skeleton';
import { InfoDialog } from '@ui/overlay/AlertDialog/InfoDialog';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
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
    <Box
      w='full'
      _hover={{
        '& #help-button': {
          visibility: 'visible',
        },
      }}
    >
      <ParentSize>
        {({ width }) => (
          <>
            <Flex direction='column' position='relative'>
              <Flex gap='2' align='center'>
                <Text fontWeight='semibold' fontSize='xl'>
                  Customer map
                </Text>
                <HelpButton isOpen={isOpen} onOpen={onOpen} />
              </Flex>
              {!hasContracts && (
                <Text
                  bottom='0'
                  color='gray.400'
                  fontSize='lg'
                  fontWeight='semibold'
                  position='absolute'
                  transform='translateY(100%)'
                >
                  No data yet
                </Text>
              )}
            </Flex>
            <Skeleton
              w='full'
              h={isTaller ? '700px' : '350px'}
              endColor='gray.300'
              startColor='gray.300'
              isLoaded={!isLoading}
            >
              <CustomerMapChart
                width={width}
                height={isTaller ? 700 : 350}
                data={chartData}
                hasContracts={hasContracts}
              />
            </Skeleton>
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
    </Box>
  );
};
