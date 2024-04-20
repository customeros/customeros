import type { Channel } from 'phoenix';
import type { RootStore } from '@store/root';
import type { TransportLayer } from '@store/transport';
import type { Operation, SyncPacket } from '@store/types';

import { gql } from 'graphql-request';
import { toJS, makeAutoObservable } from 'mobx';
import { getDiff, applyDiff } from 'recursive-diff';

import {
  TableViewDef,
  TableViewType,
  TableViewDefUpdateInput,
} from '@graphql/types';

export class TableViewDefStore {
  value: TableViewDef = {
    columns: [],
    createdAt: '',
    filters: '',
    icon: '',
    id: '',
    name: '',
    order: 0,
    sorting: '',
    updatedAt: '',
    tableType: TableViewType.Organizations,
  };
  version: number = 0;
  history: Operation[] = [];
  channel?: Channel;
  isLoading = false;
  error?: string | null = null;

  constructor(
    private rootStore: RootStore,
    private transportLayer: TransportLayer,
  ) {
    makeAutoObservable(this);
  }

  async load(tableViewDef: TableViewDef) {
    this.value = tableViewDef;

    try {
      const connection = await this.transportLayer.join(
        'TableViewDef',
        tableViewDef.id,
        this.version,
      );

      if (!connection) return;

      this.channel = connection.channel;
      this.subscribe();
    } catch (e) {
      console.error(e);
    }
  }

  subscribe() {
    if (!this.channel) return;

    this.channel.on('sync_packet', (packet: SyncPacket) => {
      const prev = toJS(this.value);
      const diff = packet.operation.diff;
      const next = applyDiff(prev, diff);

      this.value = next;
      this.version = packet.version;
      this.history.push(packet.operation);
    });
  }

  update(updater: (prev: TableViewDef) => TableViewDef) {
    const lhs = toJS(this.value);
    const next = updater(this.value);
    const rhs = toJS(next);
    const diff = getDiff(lhs, rhs);

    const operation: Operation = {
      id: this.version,
      diff,
    };

    this.history.push(operation);
    this.value = next;

    if (!this.channel) return;

    this.channel
      .push('sync_packet', { payload: { operation } })
      .receive('ok', ({ version }: { version: number }) => {
        this.version = version;
        this.save();
      });
  }

  reorderColumn(fromIndex: number, toIndex: number) {
    this.update((value) => {
      const column = value.columns[fromIndex];

      value.columns.splice(fromIndex, 1);
      value.columns.splice(toIndex, 0, column);

      return value;
    });
  }

  orderColumnsByVisibility() {
    this.update((value) => {
      value.columns.sort((a, b) => {
        if (a.visible === b.visible) return 0;
        if (a.visible) return -1;

        return 1;
      });

      return value;
    });
  }

  async save() {
    try {
      this.isLoading = true;
      await this.transportLayer.client.request<
        { id: string },
        TableViewDefUpdateInput
      >(UPDATE_TABLE_VIEW_DEF, this.value);
    } catch (e) {
      this.error = (e as Error)?.message;
    } finally {
      this.isLoading = false;
    }
  }
}

const UPDATE_TABLE_VIEW_DEF = gql`
  mutation updateTableViewDef($input: TableViewDefUpdateInput!) {
    tableViewDef_Update(input: $input) {
      id
    }
  }
`;

// function _transformChangesets(
//   changeset1: Operation[],
//   changeset2: Operation[],
// ): Operation[] {
//   // Merge the changesets
//   const mergedChangeset = [...changeset1, ...changeset2];

//   // Sort the merged changeset by the order of occurrence
//   mergedChangeset.sort((a, b) => a.id - b.id);

//   // Apply the LWW strategy
//   const transformedChangeset = mergedChangeset.reduce((result, operation) => {
//     // Check if the operation conflicts with any previous operation
//     const conflictIndex = result.findIndex((prevOperation) =>
//       conflicts(prevOperation, operation),
//     );
//     if (conflictIndex !== -1) {
//       // Resolve conflict using Last Write Wins strategy
//       const prevOperation = result[conflictIndex];
//       if (prevOperation.op === 'delete') {
//         // If previous operation was delete, discard current operation
//         return result;
//       } else {
//         // If previous operation was add or update, replace it with current operation
//         result.splice(conflictIndex, 1, operation);

//         return result;
//       }
//     } else {
//       // No conflict, add current operation to the result
//       result.push(operation);

//       return result;
//     }
//   }, [] as Operation[]);

//   return transformedChangeset;
// }

// function conflicts(operation1: Operation, operation2: Operation): boolean {
//   // Check if operation1 and operation2 modify the same field
//   return (
//     JSON.stringify(operation1.path) === JSON.stringify(operation2.path) &&
//     operation1.id === operation2.id
//   );
// }
