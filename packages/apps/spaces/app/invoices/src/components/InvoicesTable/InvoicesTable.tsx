'use client';

import { useRef, useState } from 'react';

import { useIsRestoring } from '@tanstack/react-query';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { Flex } from '@ui/layout/Flex';
import { Invoice } from '@graphql/types';
import { Heading } from '@ui/typography/Heading';
import { Table, SortingState } from '@ui/presentation/Table';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGetInvoicesQuery } from '@shared/graphql/getInvoices.generated';
import { EmptyState } from '@shared/components/Invoice/EmptyState/EmptyState';

import { columns } from './Columns/Columns';

export function InvoicesTable() {
  const isRestoring = useIsRestoring();
  const enableFeature = useFeatureIsOn('gp-dedicated-1');
  const [sorting, setSorting] = useState<SortingState>([
    { id: 'DUE_DATE', desc: true },
  ]);
  const client = getGraphQLClient();
  const tableRef = useRef(null);

  const { data, isFetching } = useGetInvoicesQuery(client, {
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
        isLoading={isRestoring ? false : isFetching}
        totalItems={isRestoring ? 40 : data?.invoices?.totalElements || 0}
        tableRef={tableRef}
        borderColor='gray.100'
        rowHeight={48}
      />
    </Flex>
  );
}
