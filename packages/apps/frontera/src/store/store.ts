import { toJS } from 'mobx';
import { Channel } from 'phoenix';
import { getDiff, applyDiff } from 'recursive-diff';

import { RootStore } from './root';
import { Transport } from './transport';
import { Operation, SyncPacket } from './types';

type StoreNames = keyof Pick<
  RootStore,
  'contracts' | 'tableViewDefs' | 'organizations'
>;
export interface Store<T> {
  value: T;
  version: number;
  root: RootStore;
  channel?: Channel;
  subscribe(): void;
  isLoading: boolean;
  set id(id: string);
  error: string | null;
  history: Operation[];
  transport: Transport;
  load(data: T): Promise<void>;
  /** Method that handles loading asynchronous data */
  invalidate: () => Promise<void>;
  update(updater: (prev: T) => T): void;
}

export type StoreConstructor<T> = new (
  root: RootStore,
  transport: Transport,
) => Store<T>;

export function makeAutoSyncable<T extends Record<string, unknown>>(
  instance: InstanceType<StoreConstructor<T>>,
  options: {
    channelName: string;
    getId?: (data: T) => string;
    mutator?: (operation: Operation) => Promise<void>;
    storeMapper?: Partial<
      Record<keyof T, { storeName: StoreNames; getItemId: (data: T) => string }>
    >;
  },
) {
  const {
    channelName,
    mutator = null,
    getId = (data) => data?.id as string,
    storeMapper,
  } = options;

  function subscribe(this: Store<typeof instance.value>) {
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

  async function load(
    this: Store<typeof instance.value>,
    data: typeof instance.value,
  ) {
    this.value = data;

    if (storeMapper) {
      Object.entries(storeMapper).forEach(([key, options]) => {
        if (!options) return;

        const { storeName, getItemId } = options;
        const value = data[key];

        if (Array.isArray(value)) {
          (this.value as Record<string, unknown>)[key] = value.map((item) => {
            const id = getItemId(item);

            if (this.root[storeName].value.has(id)) {
              return this.root[storeName].value.get(id)?.value;
            }

            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            this.root[storeName].load([value as any]);

            return this.root[storeName].value.get(id)?.value;
          });
        } else {
          const id = getItemId(value as T);

          if (this.root[storeName].value.has(id)) {
            return this.root[storeName].value.get(id)?.value;
          }

          // eslint-disable-next-line @typescript-eslint/no-explicit-any
          this.root[storeName].load([value as any]);

          (this.value as Record<string, unknown>)[key] =
            this.root[storeName].value.get(id)?.value;
        }
      });
    }

    try {
      const id = getId(data);
      const connection = await this.transport.join(
        channelName,
        id,
        this.version,
      );

      if (!connection) return;

      this.channel = connection.channel;
      this.subscribe();
    } catch (e) {
      console.error(e);
    }
  }

  function update(
    this: Store<typeof instance.value>,
    updater: (prev: typeof instance.value) => typeof instance.value,
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
          await mutator.bind(this)(operation);

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

  instance.subscribe = subscribe.bind(instance);
  instance.load = load.bind(instance);
  instance.update = update.bind(instance);
}

makeAutoSyncable.subscribe = function () {};
makeAutoSyncable.load = function <T>() {
  return async function (
    // @ts-expect-error - we don't want to prefix parameters with `_`
    data: T,
  ): Promise<void> {};
};
makeAutoSyncable.update = function <T>() {
  // @ts-expect-error - we don't want to prefix parameters with `_`
  return function (updater: (prev: T) => T, mutator?: () => Promise<void>) {};
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
