import { Channel } from 'phoenix';
import { toJS, makeAutoObservable } from 'mobx';
import { getDiff, applyDiff } from 'recursive-diff';

import { RootStore } from './root';
import { TransportLayer } from './transport';
import { Operation, SyncPacket } from './types';
import { AbstractStore, AbstractStoreClass } from './abstract';

export interface AbstractGroupStore<T extends { id: string }> {
  meta: GroupMeta<T>;
  subscribe?(): void;
  load(data: T[]): Promise<void>;
  value: Map<string, AbstractStore<T>>;
  update(
    updater: (
      prev: Map<string, AbstractStore<T>>,
    ) => Map<string, AbstractStore<T>>,
  ): void;
}

export class GroupMeta<T extends { id: string }> {
  version: number = 0;
  channel?: Channel;
  channelName?: string = '';
  history: Operation[] = [];
  isLoading: boolean = false;
  error?: string | null = null;
  isBootstrapped = false;

  constructor(
    private store: AbstractGroupStore<T>,
    private rootStore: RootStore,
    private transportLayer: TransportLayer,
    private abstractStore: AbstractStoreClass<T>,
    options?: { channelName: string },
  ) {
    this.channelName = options?.channelName;
    // this.store.subscribe = this.subscribe.bind(this);

    makeAutoObservable(this);
  }

  async load(data: T[]) {
    data.forEach((value) => {
      if (this.store.value.has(value.id)) return;

      const abstractStore = new this.abstractStore(
        this.rootStore,
        this.transportLayer,
      );
      abstractStore.load(value);
      this.store.value.set(value.id, abstractStore);
    });

    this.isBootstrapped = true;
  }

  private subscribe() {
    if (!this.channel) return;

    this.channel.on('sync_packet', (packet: SyncPacket) => {
      const prev = toJS(this.store.value);
      const diff = packet.operation.diff;
      const next = applyDiff(prev, diff);

      this.store.value = next;
      this.version = packet.version;
      this.history.push(packet.operation);
    });
  }

  update(
    updater: (
      prev: Map<string, AbstractStore<T>>,
    ) => Map<string, AbstractStore<T>>,
    mutator?: () => Promise<void>,
  ) {
    const lhs = toJS(this.store.value);
    const next = updater(this.store.value);
    const rhs = toJS(next);
    const diff = getDiff(lhs, rhs);

    const operation: Operation = {
      id: this.version,
      diff,
    };

    this.history.push(operation);
    this.store.value = next;

    if (mutator) {
      (async () => {
        try {
          this.error = null;
          await mutator();

          this?.channel
            ?.push('sync_packet', { payload: { operation } })
            ?.receive('ok', ({ version }: { version: number }) => {
              this.version = version;
            });
        } catch (e) {
          console.error(e);
          this.store.value = lhs;
          this.history.pop();
        }
      })();
    }
  }
}

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
