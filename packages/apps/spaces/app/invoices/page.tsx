'use client';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';

import { InvoicesTable } from './src/components/InvoicesTable';
import { PreviewPanel } from './src/components/PreviewPanel/PreviewPanel';

interface InvoicesPageProps {
  searchParams: { invoice?: string };
}
export default function InvoicesPage({ searchParams }: InvoicesPageProps) {
  return (
    <Flex h='100%' w='full'>
      <Box
        maxW={550}
        h='full'
        mr='4'
        pt='4'
        pl='3'
        borderRight='1px solid'
        borderColor='gray.200'
      >
        <InvoicesTable />
      </Box>

      {searchParams?.invoice && (
        <Box
          w='full'
          maxW={575}
          h='full'
          pr={4}
          pt='4'
          borderRight='1px solid'
          borderColor='gray.200'
        >
          <PreviewPanel id={searchParams.invoice} />
        </Box>
      )}
    </Flex>
  );
}
