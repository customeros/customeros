'use client';
import dynamic from 'next/dynamic';

import ParentSize from '@visx/responsive/lib/components/ParentSize';
import { useCustomerMapQuery } from '@customerMap/graphql/customerMap.generated';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { useDisclosure } from '@ui/utils';
import { Text } from '@ui/typography/Text';
import { Skeleton } from '@ui/presentation/Skeleton';
import { InfoDialog } from '@ui/overlay/AlertDialog/InfoDialog';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

import { HelpContent } from './HelpContent';
import { HelpButton } from '../../HelpButton';
import { CustomerMapDatum } from './CustomerMap.chart';

const CustomerMapChart = dynamic(() => import('./CustomerMap.chart'), {
  ssr: false,
});

export const CustomerMap = () => {
  const client = getGraphQLClient();
  const { data, isLoading } = useCustomerMapQuery(client);
  const { isOpen, onOpen, onClose } = useDisclosure();

  const chartData = (data?.dashboard_CustomerMap ?? []).map((d) => ({
    x: d?.contractSignedDate ? new Date(d?.contractSignedDate) : new Date(),
    r: d?.arr,
    values: {
      id: d?.organization?.id,
      name: d?.organization?.name || 'Unnamed',
      status: d?.state,
    },
  })) as CustomerMapDatum[];

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
            <Flex gap='2' align='center'>
              <Text fontWeight='semibold' fontSize='xl'>
                Customer map
              </Text>
              <HelpButton isOpen={isOpen} onOpen={onOpen} />
            </Flex>
            <Skeleton
              w='full'
              h='350px'
              endColor='gray.300'
              startColor='gray.300'
              isLoaded={!isLoading}
            >
              <CustomerMapChart width={width} height={350} data={chartData} />
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
