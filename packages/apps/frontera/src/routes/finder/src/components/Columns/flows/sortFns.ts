import { match } from 'ts-pattern';
import { FlowStore } from '@store/Flows/Flow.store';
import { flowOptions } from '@finder/components/Columns/flows/utils.ts';

import { ColumnViewType } from '@graphql/types';

export const getFlowsColumnSortFn = (columnId: string) =>
  match(columnId)
    .with(
      ColumnViewType.FlowActionName,
      () => (row: FlowStore) =>
        row?.value?.status
          ? flowOptions.find((e) => e.value === row.value.status)?.label
          : null,
    )

    .with(ColumnViewType.FlowName, () => (row: FlowStore) => {
      const value = row.value?.name?.toLowerCase();

      return value || null;
    })
    .with(ColumnViewType.FlowOnHoldCount, () => (row: FlowStore) => {
      return row.value?.statistics?.onHold ?? null;
    })
    .with(ColumnViewType.FlowReadyCount, () => (row: FlowStore) => {
      return row.value?.statistics?.ready ?? null;
    })
    .with(ColumnViewType.FlowScheduledCount, () => (row: FlowStore) => {
      return row.value?.statistics?.scheduled ?? null;
    })
    .with(ColumnViewType.FlowCompletedCount, () => (row: FlowStore) => {
      return row.value?.statistics?.completed ?? null;
    })
    .with(ColumnViewType.FlowGoalAchievedCount, () => (row: FlowStore) => {
      return row.value?.statistics?.goalAchieved ?? null;
    })
    .with(ColumnViewType.FlowTotalCount, () => (row: FlowStore) => {
      return row.value?.statistics?.total ?? null;
    })
    .otherwise(() => (_row: FlowStore) => null);
