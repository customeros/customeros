import { Channel } from 'phoenix';
import { match } from 'ts-pattern';
import {
  when,
  action,
  computed,
  observable,
  runInAction,
  makeObservable,
} from 'mobx';

import { RootStore } from './root';
import { Syncable } from './syncable';
import { Transport } from './transport';
import { GroupOperation, GroupSyncPacket } from './types';

export class SyncableGroup<T extends object, TSyncable extends Syncable<T>> {
  version = 0;
  channel?: Channel;
  isLoading = false;
  isBootstrapped = false;
  error: string | null = null;
  history: GroupOperation[] = [];
  value: Map<string, TSyncable> = new Map();

  constructor(
    public root: RootStore,
    public transport: Transport,
    private SyncableStore: typeof Syncable<T>,
  ) {
    makeObservable<
      SyncableGroup<T, TSyncable>,
      'initChannelConnection' | 'subscribe' | 'applyGroupOperation'
    >(this, {
      load: action,
      sync: action,
      subscribe: action,
      error: observable,
      value: observable,
      version: observable,
      channel: observable,
      history: observable,
      isLoading: observable,
      channelName: computed,
      isBootstrapped: observable,
      applyGroupOperation: action,
      initChannelConnection: action,
    });

    when(
      () => !!this.root.session.value.tenant && !this.root.demoMode,
      async () => {
        const tenant = this.root.session.value.tenant;

        try {
          await this.initChannelConnection(tenant);
        } catch (e) {
          console.error(e);
        }
      },
    );
  }

  get channelName() {
    return '';
  }

  load(data: T[], options: { getId: (data: T) => string }) {
    data.forEach((item) => {
      const id = options.getId(item);

      if (this.value.has(id)) {
        this.value.get(id)?.load(item);

        return;
      }

      const syncableItem = new this.SyncableStore(
        this.root,
        this.transport,
        item,
      );

      syncableItem.load(item);
      this.value.set(id, syncableItem as TSyncable);
    });

    this.isBootstrapped = true;
  }

  sync(operation: GroupOperation) {
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

  private async initChannelConnection(tenant: string) {
    try {
      const connection = await this.transport.join(
        this.channelName,
        tenant,
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
    if (!this.channel || this.root.demoMode) return;

    this.channel.on('sync_group_packet', (packet: GroupSyncPacket) => {
      if (packet.ref === this.transport.refId) return;
      this.applyGroupOperation(packet);
      this.history.push(packet);
    });
  }

  private applyGroupOperation(operation: GroupOperation) {
    match(operation.action)
      .with('APPEND', () => {
        operation.ids.forEach((id) => {
          const newSyncableItem = new this.SyncableStore(
            this.root,
            this.transport,
            this.SyncableStore.getDefaultValue(),
          );

          runInAction(() => {
            newSyncableItem.setId(id);
            this.value.set(id, newSyncableItem as TSyncable);
          });

          setTimeout(() => {
            this.value.get(id)?.invalidate();
          }, 1000);
        });
      })
      .with('DELETE', () => {
        operation.ids.forEach((id) => {
          runInAction(() => {
            this.value.delete(id);
          });
        });
      })
      .with('INVALIDATE', () => {
        operation.ids.forEach((id) => {
          const item = this.value.get(id);

          if (!item) return;

          item.invalidate();
        });
      })
      .otherwise(() => {});
  }

  static SyncableStore = Syncable;
}
