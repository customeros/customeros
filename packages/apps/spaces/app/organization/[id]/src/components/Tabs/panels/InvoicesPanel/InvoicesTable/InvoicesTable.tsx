'use client';

import { FC, useRef } from 'react';

import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { Invoice } from '@graphql/types';
import { Table } from '@ui/presentation/Table';
import { EmptyState } from '@shared/components/Invoice/EmptyState/EmptyState';

import { columns } from './Columns/Columns';

export const InvoicesTable: FC<{
  totalElements: number;
  invoices: Array<Invoice>;
}> = ({ invoices, totalElements }) => {
  const enableFeature = useFeatureIsOn('invoices');

  const tableRef = useRef(null);

  if (totalElements === 0) {
    return <EmptyState maxW={440} />;
  }

  return (
    <Table<Invoice>
      data={invoices}
      columns={columns}
      enableTableActions={enableFeature !== null ? enableFeature : true}
      enableRowSelection={false}
      fullRowSelection={true}
      canFetchMore={false}
      // onFetchMore={handleFetchMore}
      isLoading={false}
      totalItems={totalElements}
      tableRef={tableRef}
      rowHeight={4}
      borderColor='gray.100'
      contentHeight={'80vh'}
    />
  );
};
