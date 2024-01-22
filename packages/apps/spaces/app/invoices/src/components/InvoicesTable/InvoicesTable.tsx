'use client';

import { useRef, useState } from 'react';

// import { useIsRestoring } from '@tanstack/react-query';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { Flex } from '@ui/layout/Flex';
import { Invoice } from '@graphql/types';
import { Heading } from '@ui/typography/Heading';
import { Table, SortingState } from '@ui/presentation/Table';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { EmptyState } from '@shared/components/Invoice/EmptyState/EmptyState';

import { columns } from './Columns/Columns';
import { useGetInvoicesQuery } from '../../graphql/getInvoices.generated';

export function InvoicesTable() {
  // const isRestoring = useIsRestoring();
  const enableFeature = useFeatureIsOn('gp-dedicated-1');
  const [sorting, setSorting] = useState<SortingState>([
    { id: 'LAST_TOUCHPOINT', desc: true },
  ]);
  const client = getGraphQLClient();
  const tableRef = useRef(null);

  const { data } = useGetInvoicesQuery(client, {
    pagination: {
      page: 0,
      limit: 40,
    },
  });

  if (data?.invoices.totalElements === 0) {
    return <EmptyState />;
  }

  return (
    <Flex flexDir='column' as='article'>
      <Heading fontSize='lg' pb={3}>
        Invoices
      </Heading>
      <Table<Invoice>
        data={(data?.invoices?.content as Invoice[]) ?? []}
        columns={columns}
        sorting={sorting}
        enableTableActions={enableFeature !== null ? enableFeature : true}
        enableRowSelection={true}
        fullRowSelection={true}
        canFetchMore={false}
        onSortingChange={setSorting}
        // onFetchMore={handleFetchMore}
        isLoading={false}
        totalItems={4}
        tableRef={tableRef}
        borderColor='gray.100'
        rowHeight={48}
      />
    </Flex>
  );
}
