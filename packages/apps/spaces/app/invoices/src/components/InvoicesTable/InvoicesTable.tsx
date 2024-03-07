'use client';

import { useRef, useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';

import { useIsRestoring } from '@tanstack/react-query';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Heading } from '@ui/typography/Heading';
import { Table, TableInstance } from '@ui/presentation/Table';
import { EmptyState } from '@shared/components/Invoice/EmptyState/EmptyState';
import {
  InvoiceTableData,
  useInfiniteInvoices,
} from '@shared/components/Invoice/hooks/useInfiniteInvoices';

import { columns } from './Columns/Columns';

export function InvoicesTable() {
  const isRestoring = useIsRestoring();
  const searchParams = useSearchParams();
  const selectedInvoiceId = searchParams?.get('invoice');
  const router = useRouter();

  const tableRef = useRef<TableInstance<InvoiceTableData>>(null);
  const {
    invoiceFlattenData,
    totalInvoicesCount,
    isFetching,
    isFetched,
    fetchNextPage,
    hasNextPage,
  } = useInfiniteInvoices();

  useEffect(() => {
    if (!selectedInvoiceId && tableRef.current) {
      const newParams = new URLSearchParams(searchParams ?? '');
      const firstId = invoiceFlattenData?.[0]?.id;
      if (!firstId) return;
      newParams.set('invoice', firstId);
      router.replace(`/invoices?${newParams.toString()}`);
    }
  }, [selectedInvoiceId, isFetched]);

  useEffect(() => {
    if (tableRef.current && isFetched) {
      tableRef.current
        ?.getRowModel()
        ?.rows?.find((e) => e.original.id === selectedInvoiceId)
        ?.toggleSelected(true);
    }
  }, [tableRef, isFetched, selectedInvoiceId]);

  if (totalInvoicesCount === 0) {
    return <EmptyState maxW={500} isDashboard />;
  }

  const handleOpenInvoice = (id: string) => {
    const params = new URLSearchParams(searchParams?.toString() ?? '');
    params.set('invoice', id);
    router.push(`?${params}`);
  };

  return (
    <Flex flexDir='column' as='article'>
      <Heading fontSize='lg' pb={3}>
        Invoices
      </Heading>
      <Box ml={-3}>
        <Table<InvoiceTableData>
          data={invoiceFlattenData ?? []}
          columns={columns}
          onFullRowSelection={(id) => id && handleOpenInvoice(id)}
          enableRowSelection={false}
          fullRowSelection={true}
          canFetchMore={hasNextPage}
          onFetchMore={fetchNextPage}
          isLoading={isRestoring ? false : isFetching}
          totalItems={isRestoring ? 10 : totalInvoicesCount || 0}
          tableRef={tableRef}
          borderColor='gray.100'
          rowHeight={48}
        />
      </Box>
    </Flex>
  );
}
