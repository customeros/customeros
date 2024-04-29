import { Channel } from 'phoenix';
import { toJS, makeAutoObservable } from 'mobx';
import { getDiff, applyDiff } from 'recursive-diff';

import { RootStore } from './root';
import { TransportLayer } from './transport';
import { Operation, SyncPacket } from './types';

export interface AbstractStore<T extends { id: string }> {
  value: T;
  meta: Meta<T>;
  subscribe?(): void;
  load(data: T): Promise<void>;
  update(updater: (prev: T) => T): void;
}

export type AbstractStoreClass<T extends { id: string }> = new (
  rootStore: RootStore,
  transportLayer: TransportLayer,
) => AbstractStore<T>;

export class Meta<T extends { id: string }> {
  version: number = 0;
  history: Operation[] = [];
  channel?: Channel;
  channelName: string = '';
  isLoading: boolean = false;
  error?: string | null = null;

  constructor(
    private store: AbstractStore<T>,
    private transportLayer: TransportLayer,
    { channelName }: { channelName: string },
  ) {
    this.channelName = channelName;
    this.store.subscribe = this.subscribe.bind(this);

    makeAutoObservable(this);
  }

  async load(data: T) {
    this.store.value = data;

    try {
      const connection = await this.transportLayer.join(
        this.channelName,
        data.id,
        this.version,
      );

      if (!connection) return;

      this.channel = connection.channel;
      this.subscribe();
    } catch (e) {
      console.error(e);
    }
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

  update(updater: (prev: T) => T, mutator?: () => Promise<void>) {
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
