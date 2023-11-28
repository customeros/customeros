'use client';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';

import { CustomerMap } from './src/components/charts/CustomerMap';
import { ARRBreakdown } from './src/components/charts/ARRBreakdown';
import { NewCustomers } from './src/components/charts/NewCustomers';
import { RevenueAtRisk } from './src/components/charts/RevenueAtRisk';
import { RetentionRate } from './src/components/charts/RetentionRate';
import { MrrPerCustomer } from './src/components/charts/MrrPerCustomer';
import { GrossRevenueRetention } from './src/components/charts/GrossRevenueRetention';

export default function DashboardPage() {
  return (
    <Flex flexDir='column' gap='3' pl='1' pt='4'>
      <Text fontWeight='semibold' fontSize='xl'>
        Customer map
      </Text>
      <CustomerMap />

      <Flex gap='3'>
        <MrrPerCustomer />
        <GrossRevenueRetention />
      </Flex>

      <Flex gap='3'>
        <ARRBreakdown />
        <RevenueAtRisk />
      </Flex>

      <Flex gap='3'>
        <NewCustomers />
        <RetentionRate />
      </Flex>
    </Flex>
  );
}
