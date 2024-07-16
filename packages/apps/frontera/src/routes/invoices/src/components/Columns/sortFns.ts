import { match } from 'ts-pattern';
import { InvoiceStore } from '@store/Invoices/Invoice.store.ts';

import { InvoiceStatus, ColumnViewType } from '@graphql/types';

export const getColumnSortFn = (columnId: string) =>
  match(columnId)
    .with(
      ColumnViewType.InvoicesInvoiceStatus,
      () => (row: InvoiceStore) =>
        match(row.value?.status)
          .with(InvoiceStatus.Empty, () => null)
          .with(InvoiceStatus.Initialized, () => 1)
          .with(InvoiceStatus.OnHold, () => 2)
          .with(InvoiceStatus.Scheduled, () => 3)
          .with(InvoiceStatus.Void, () => 4)
          .with(InvoiceStatus.Paid, () => 5)
          .with(InvoiceStatus.Due, () => 6)
          .with(InvoiceStatus.Overdue, () => 7)
          .otherwise(() => null),
    )

    .with(ColumnViewType.InvoicesDueDate, () => (row: InvoiceStore) => {
      const value = row.value?.due;

      return value ? new Date(value) : null;
    })
    .with(ColumnViewType.InvoicesIssueDate, () => (row: InvoiceStore) => {
      const value = row.value?.due;

      return value ? new Date(value) : null;
    })
    .with(ColumnViewType.InvoicesIssueDatePast, () => (row: InvoiceStore) => {
      const value = row.value?.metadata.created;

      return value ? new Date(value) : null;
    })
    .otherwise(() => (_row: InvoiceStore) => null);
