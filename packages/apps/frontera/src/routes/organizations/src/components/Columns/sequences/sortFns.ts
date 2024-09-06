import { match } from 'ts-pattern';
import { FlowSequenceStore } from '@store/Sequences/FlowSequence.store';

import { ColumnViewType, FlowSequenceStatus } from '@graphql/types';

export const getSequenceColumnSortFn = (columnId: string) =>
  match(columnId)
    .with(
      ColumnViewType.FlowSequenceStatus,
      () => (row: FlowSequenceStore) =>
        match(row.value?.status)
          .with(FlowSequenceStatus.Inactive, () => 4)
          .with(FlowSequenceStatus.Active, () => 3)
          .with(FlowSequenceStatus.Paused, () => 2)
          .with(FlowSequenceStatus.Archived, () => 1)
          .otherwise(() => null),
    )

    .with(ColumnViewType.FlowName, () => (row: FlowSequenceStore) => {
      const value = row.value?.flow?.name?.toLowerCase();

      return value || null;
    })
    .with(ColumnViewType.FlowSequenceName, () => (row: FlowSequenceStore) => {
      const value = row.value?.name?.toLowerCase();

      return value || null;
    })
    .with(
      ColumnViewType.FlowSequenceContactCount,
      () => (row: FlowSequenceStore) => {
        return row.value?.contacts?.length;
      },
    )
    .otherwise(() => (_row: FlowSequenceStore) => null);
