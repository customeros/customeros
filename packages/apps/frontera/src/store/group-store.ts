import { Channel } from 'phoenix';
import { match } from 'ts-pattern';
import { when, runInAction } from 'mobx';

import { RootStore } from './root';
import { Transport } from './transport';
import { Store, StoreConstructor } from './store';
import { GroupOperation, GroupSyncPacket } from './types';

export interface GroupStore<T> {
  version: number;
  root: RootStore;
  channel?: Channel;
  subscribe(): void;
  isLoading: boolean;
  error: string | null;
  transport: Transport;
  load(data: T[]): void;
  isBootstrapped: boolean;
  history: GroupOperation[];
  value: Map<string, Store<T>>;
  sync(operation: GroupOperation): void;
}

type GroupStoreConstructor<T> = new (
  root: RootStore,
  transport: Transport,
) => GroupStore<T>;

export function makeAutoSyncableGroup<T extends Record<string, unknown>>(
  instance: InstanceType<GroupStoreConstructor<T>>,
  options: {
    channelName: string;
    ItemStore: StoreConstructor<T>;
    getItemId: (data: T) => string;
  },
) {
  const {
    ItemStore,
    channelName,
    getItemId = (data) => data?.id as string,
  } = options;

  function load(this: GroupStore<T>, data: T[]) {
    data.forEach((item) => {
      const id = getItemId(item);

      if (this.value.has(id)) {
        this.value.get(id)?.load(item);

        return;
      }

      const itemStore = new ItemStore(this.root, this.transport);

      itemStore.load(item);
      this.value.set(id, itemStore);
    });

    when(
      () => !!this.root.session.value.tenant && !this.root.demoMode,
      async () => {
        const tenant = this.root.session.value.tenant;

        try {
          const connection = await this.transport.join(
            channelName,
            tenant,
            this.version,
            true,
          );

          if (!connection) return;

          this.channel = connection.channel;
          this.subscribe();
        } catch (e) {
          console.error(e);
        }
      },
    );

    this.isBootstrapped = true;
  }

  function subscribe(this: GroupStore<T>) {
    if (!this.channel || this.root.demoMode) return;

    this.channel.on('sync_group_packet', (packet: GroupSyncPacket) => {
      if (packet.ref === this.transport.refId) return;
      applyGroupOperation(this, ItemStore, packet);
      this.history.push(packet);
    });
  }

  function sync(this: GroupStore<T>, operation: GroupOperation) {
    const op = {
      ...operation,
      ref: this.transport.refId,
    };

    this.history.push(op);
    this?.channel
      ?.push('sync_group_packet', { payload: { operation: op } })
      ?.receive('ok', ({ version }: { version: number }) => {
        this.version = version;
      });
  }

  instance.load = load.bind(instance);
  instance.subscribe = subscribe.bind(instance);
  instance.sync = sync.bind(instance);
}

makeAutoSyncableGroup.subscribe = function () {};

makeAutoSyncableGroup.load = function <T>() {
  // @ts-expect-error - we don't want to prefix parameters with `_`
  return function (data: T[]): void {};
};

// @ts-expect-error - we don't want to prefix parameters with `_`
makeAutoSyncableGroup.sync = function (operation: GroupOperation): void {};

function applyGroupOperation<T>(
  instance: GroupStore<T>,
  ItemStore: StoreConstructor<T>,
  operation: GroupOperation,
) {
  match(operation.action)
    .with('APPEND', () => {
      operation.ids.forEach((id) => {
        const newItem = new ItemStore(instance.root, instance.transport);

        runInAction(() => {
          newItem.id = id;
          instance.value.set(id, newItem);
        });

        setTimeout(() => {
          instance.value.get(id)?.invalidate();
        }, 1000);
      });
    })
    .with('DELETE', () => {
      operation.ids.forEach((id) => {
        runInAction(() => {
          instance.value.delete(id);
        });
      });
    })
    .with('INVALIDATE', () => {
      operation.ids.forEach((id) => {
        const item = instance.value.get(id);

        if (!item) return;

        item.invalidate();
      });
    })
    .otherwise(() => {});
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
