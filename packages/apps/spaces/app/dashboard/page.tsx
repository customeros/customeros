'use client';

import dynamic from 'next/dynamic';

// import ParentSize from '@visx/responsive/lib/components/ParentSize';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';

import { ChartCard } from './src/components/ChartCard';

// const CustomerMap = dynamic(
//   () => import('./src/components/charts/CustomerMap/CustomerMap'),
//   {
//     ssr: false,
//   },
// );
const MrrPerCustomer = dynamic(
  () => import('./src/components/charts/MrrPerCustomer/MrrPerCustomer'),
  {
    ssr: false,
  },
);
const GrossRevenueRetention = dynamic(
  () =>
    import(
      './src/components/charts/GrossRevenueRetention/GrossRevenueRetention'
    ),
  {
    ssr: false,
  },
);
const ARRBreakdown = dynamic(
  () => import('./src/components/charts/ARRBreakdown/ARRBreakdown'),
  {
    ssr: false,
  },
);
const RevenueAtRisk = dynamic(
  () => import('./src/components/charts/RevenueAtRisk/RevenueAtRisk'),
  {
    ssr: false,
  },
);
const NewCustomers = dynamic(
  () => import('./src/components/charts/NewCustomers/NewCustomers'),
  {
    ssr: false,
  },
);
const RetentionRate = dynamic(
  () => import('./src/components/charts/RetentionRate/RetentionRate'),
  {
    ssr: false,
  },
);

export default function DashboardPage() {
  return (
    <Flex flexDir='column' gap='4'>
      <Text fontWeight='medium' fontSize='2xl'>
        Dashboard
      </Text>
      {/* <ParentSize>
        {({ width }) => <CustomerMap width={width} height={350} />}
      </ParentSize> */}

      <Flex gap='4'>
        <ChartCard
          flex='1'
          stat='$4,280'
          title='MRR per Customer'
          renderSubStat={() => <Text fontSize='sm'>+ 5% vs last mth</Text>}
        >
          <MrrPerCustomer />
        </ChartCard>
        <ChartCard
          flex='2'
          stat='95%'
          title='Gross Revenue Retention'
          renderSubStat={() => <Text fontSize='sm'> +5.4% vs last mth</Text>}
        >
          <GrossRevenueRetention />
        </ChartCard>
      </Flex>

      <Flex gap='4'>
        <ChartCard
          flex='3'
          stat='$1,830,990'
          title='ARR Breakdown'
          renderSubStat={() => <Text fontSize='sm'> +2% vs last mth</Text>}
        >
          <ARRBreakdown />
        </ChartCard>
        <ChartCard
          flex='1'
          title='Revenue at Risk'
          renderSubStat={() => (
            <Flex mt='4' justify='space-between'>
              <Flex flexDir='column'>
                <Flex gap='3' align='center'>
                  <Flex w='3' h='3' bg='#66C61C' borderRadius='full' />
                  <Text>High Confidence</Text>
                </Flex>
                <Text fontSize='sm' fontWeight='medium'>
                  $1,830,990
                </Text>
              </Flex>

              <Flex flexDir='column'>
                <Flex gap='3' align='center'>
                  <Flex w='3' h='3' bg='yellow.400' borderRadius='full' />
                  <Text>At Risk</Text>
                </Flex>
                <Text fontSize='sm'>$355,300</Text>
              </Flex>
            </Flex>
          )}
        >
          <RevenueAtRisk />
        </ChartCard>
      </Flex>

      <Flex gap='4'>
        <ChartCard
          flex='1'
          title='New Customers'
          stat='127'
          renderSubStat={() => <Text fontSize='sm'> +2% vs last mth</Text>}
        >
          <NewCustomers />
        </ChartCard>
        <ChartCard
          flex='1'
          title='Retention Rate'
          stat='86%'
          renderSubStat={() => <Text fontSize='sm'> +2% vs last mth</Text>}
        >
          <RetentionRate />
        </ChartCard>
      </Flex>
    </Flex>
  );
}
