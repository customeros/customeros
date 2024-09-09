import { match } from 'ts-pattern';
import { FlowSequenceStore } from '@store/Sequences/FlowSequence.store';

import { ColumnViewType } from '@graphql/types';
import { flowSequencesOptions } from '@organizations/components/Columns/sequences/utils.ts';

export const getSequenceColumnSortFn = (columnId: string) =>
  match(columnId)
    .with(
      ColumnViewType.FlowSequenceStatus,
      () => (row: FlowSequenceStore) =>
        row?.value?.status
          ? flowSequencesOptions.find((e) => e.value === row.value.status)
              ?.label
          : null,
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
