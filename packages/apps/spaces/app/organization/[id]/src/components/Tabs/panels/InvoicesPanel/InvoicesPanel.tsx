'use client';
import React, { useRef } from 'react';
import { useParams, useRouter } from 'next/navigation';

import { motion } from 'framer-motion';
import { useIsRestoring } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Invoice } from '@graphql/types';
import { Text } from '@ui/typography/Text';
import { Table } from '@ui/presentation/Table';
import { IconButton } from '@ui/form/IconButton';
import { ChevronDown } from '@ui/media/icons/ChevronDown';
import { EmptyState } from '@shared/components/Invoice/EmptyState/EmptyState';
import { useInfiniteInvoices } from '@shared/components/Invoice/hooks/useInfiniteInvoices';
import { columns } from '@organization/src/components/Tabs/panels/InvoicesPanel/Columns/Columns';
import { OrganizationPanel } from '@organization/src/components/Tabs/panels/OrganizationPanel/OrganizationPanel';
import { useTimelineEventPreviewMethodsContext } from '@organization/src/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

const slideUpVariants = {
  initial: { y: '100%', opacity: 0 },
  animate: {
    y: 0,
    opacity: 1,
    transition: { type: 'interia', stiffness: 100 },
  },
  exit: { y: '100%', opacity: 0, transition: { duration: 3 } },
};
export const InvoicesPanel = () => {
  const id = useParams()?.id as string;
  const isRestoring = useIsRestoring();
  const router = useRouter();
  const tableRef = useRef(null);
  const { handleOpenInvoice } = useTimelineEventPreviewMethodsContext();

  const {
    invoiceFlattenData,
    totalInvoicesCount,
    isFetching,
    fetchNextPage,
    hasNextPage,
  } = useInfiniteInvoices(id);
  if (totalInvoicesCount === 0) {
    return (
      <Flex justifyContent='center'>
        <EmptyState maxW={448} />
      </Flex>
    );
  }

  return (
    <OrganizationPanel title='Account'>
      <motion.div
        key='invoices'
        variants={slideUpVariants}
        initial='initial'
        animate='animate'
        exit={{ x: -500, opacity: 0 }}
        style={{ width: '100%' }}
      >
        <Flex justifyContent='space-between' mb={2}>
          <Text fontSize='sm' fontWeight='semibold'>
            Invoices
          </Text>
          <IconButton
            aria-label='Go back'
            variant='ghost'
            size='xs'
            icon={<ChevronDown color='gray.400' />}
            onClick={() => router.push(`?tab=account`)}
          />
        </Flex>
        <Flex mx={-5}>
          <Table<Invoice>
            data={invoiceFlattenData ?? []}
            columns={columns}
            enableRowSelection={false}
            fullRowSelection={true}
            onFullRowSelection={(id) => id && handleOpenInvoice(id)}
            canFetchMore={hasNextPage}
            onFetchMore={fetchNextPage}
            tableRef={tableRef}
            isLoading={isRestoring ? false : isFetching}
            totalItems={isRestoring ? 10 : totalInvoicesCount}
            rowHeight={4}
            borderColor='gray.100'
            contentHeight={'80vh'}
          />
        </Flex>
      </motion.div>
    </OrganizationPanel>
  );
};
