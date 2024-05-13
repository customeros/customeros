import { toJS } from 'mobx';
import { Channel } from 'phoenix';
import { getDiff, applyDiff } from 'recursive-diff';

import { RootStore } from './root';
import { Transport } from './transport';
import { Operation, SyncPacket } from './types';
import { Store, StoreConstructor } from './store';

export interface GroupStore<T> {
  version: number;
  root: RootStore;
  channel?: Channel;
  subscribe(): void;
  isLoading: boolean;
  history: Operation[];
  error: string | null;
  transport: Transport;
  load(data: T[]): void;
  isBootstrapped: boolean;
  value: Map<string, Store<T>>;
  update(update: (prev: Map<string, Store<T>>) => Map<string, Store<T>>): void;
}

type GroupStoreConstructor<T> = new (
  root: RootStore,
  transport: Transport,
) => GroupStore<T>;

export function makeAutoSyncableGroup<T extends Record<string, unknown>>(
  instance: InstanceType<GroupStoreConstructor<T>>,
  options: {
    channelName: string;
    mutator?: () => Promise<void>;
    ItemStore: StoreConstructor<T>;
    getItemId: (data: T) => string;
  },
) {
  const {
    ItemStore,
    // channelName,
    mutator = null,
    getItemId = (data) => data?.id as string,
  } = options;

  function load(this: GroupStore<T>, data: T[]) {
    data.forEach((item) => {
      const id = getItemId(item);
      if (this.value.has(id)) return;

      const itemStore = new ItemStore(this.root, this.transport);
      itemStore.load(item);
      this.value.set(id, itemStore);
    });

    // channel join logic needs to be implemented here
    // after channel join, subscribe to the channel

    this.isBootstrapped = true;
  }

  function subscribe(this: GroupStore<T>) {
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

  function update(
    this: GroupStore<T>,
    updater: (prev: Map<string, Store<T>>) => Map<string, Store<T>>,
  ) {
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

    if (mutator) {
      (async () => {
        try {
          this.error = null;
          await mutator.bind(this)();

          this?.channel
            ?.push('sync_packet', { payload: { operation } })
            ?.receive('ok', ({ version }: { version: number }) => {
              this.version = version;
            });
        } catch (e) {
          console.error(e);
          this.value = lhs;
          this.history.pop();
        }
      })();
    }
  }

  instance.load = load.bind(instance);
  instance.subscribe = subscribe.bind(instance);
  instance.update = update.bind(instance);
}

makeAutoSyncableGroup.subscribe = function () {};
makeAutoSyncableGroup.load = function <T>() {
  return function (data: T[]): void {};
};
makeAutoSyncableGroup.update = function <T>() {
  return function (
    updater: (prev: Map<string, Store<T>>) => Map<string, Store<T>>,
  ): void {};
};

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
