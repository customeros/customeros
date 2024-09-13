import { match } from 'ts-pattern';
import { FlowStore } from '@store/Flows/Flow.store';

import { ColumnViewType } from '@graphql/types';
import { flowOptions } from '@organizations/components/Columns/flows/utils.ts';

export const getFlowsColumnSortFn = (columnId: string) =>
  match(columnId)
    .with(
      ColumnViewType.FlowSequenceStatus,
      () => (row: FlowStore) =>
        row?.value?.status
          ? flowOptions.find((e) => e.value === row.value.status)?.label
          : null,
    )

    .with(ColumnViewType.FlowName, () => (row: FlowStore) => {
      const value = row.value?.name?.toLowerCase();

      return value || null;
    })
    .with(ColumnViewType.FlowSequenceContactCount, () => (row: FlowStore) => {
      return row.value?.contacts?.length;
    })
    .otherwise(() => (_row: FlowStore) => null);
