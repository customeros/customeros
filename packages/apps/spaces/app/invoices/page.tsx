'use client';

import { Flex } from '@ui/layout/Flex';

import { InvoicesTable } from './src/components/InvoicesTable';

export default function InvoicesPage() {
  return (
    <Flex pl='3' pt='4'>
      <InvoicesTable />
      {/*<Invoice {...dummyData} />*/}
    </Flex>
  );
}
