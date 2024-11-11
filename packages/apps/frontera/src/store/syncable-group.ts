import type { Channel } from 'phoenix';

import { match } from 'ts-pattern';
import { Persister, type PersisterInstance } from '@store/persister';
import {
  when,
  action,
  computed,
  observable,
  runInAction,
  makeObservable,
} from 'mobx';

import type { RootStore } from './root';
import type { Transport } from './transport';
import type { GroupOperation, GroupSyncPacket } from './types';

import { Syncable } from './syncable';

export class SyncableGroup<T extends object, TSyncable extends Syncable<T>> {
  version = 0;
  channel?: Channel;
  isLoading = false;
  isHydrated = false;
  isBootstrapped = false;
  isBootstrapping = false;
  error: string | null = null;
  persister?: PersisterInstance;
  history: GroupOperation[] = [];
  value: Map<string, TSyncable> = new Map();

  constructor(
    public root: RootStore,
    public transport: Transport,
    private SyncableStore: typeof Syncable<T>,
  ) {
    makeObservable<
      SyncableGroup<T, TSyncable>,
      | 'initChannelConnection'
      | 'subscribe'
      | 'applyGroupOperation'
      | 'initPersister'
      | 'checkIfCanHydrate'
    >(this, {
      load: action,
      sync: action,
      drop: action,
      hydrate: action,
      subscribe: action,
      error: observable,
      value: observable,
      version: observable,
      channel: observable,
      history: observable,
      isLoading: observable,
      channelName: computed,
      initPersister: action,
      isHydrated: observable,
      getRecentChanges: action,
      checkIfCanHydrate: action,
      isBootstrapped: observable,
      isBootstrapping: observable,
      applyGroupOperation: action,
      initChannelConnection: action,
    });

    when(
      () => !!this.root.session.sessionToken && !this.root.demoMode,
      async () => {
        try {
          await this.initPersister();
        } catch (e) {
          console.error(e);
        }
      },
    );

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

    when(
      () => this.isBootstrapped,
      () => {
        this.persister?.setItem('isBootstrapped', true);
      },
    );

    if (this.channelName) {
      window.addEventListener('focus', async () => {
        await this.getRecentChanges();
      });
    }
  }

  get channelName() {
    return '';
  }

  get persisterKey() {
    return '';
  }

  public load(data: T[], options: { getId: (data: T) => string }) {
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
        this.channel,
      );

      syncableItem.load(item);
      this.value.set(id, syncableItem as TSyncable);
    });

    this.persister?.getItem<Map<string, T>>('data', (err, value) => {
      if (err) {
        console.error('Failed to get persisted data', err);

        return;
      }

      const persistedMap = value ?? new Map();

      for (let i = 0; i < data.length; i++) {
        persistedMap.set(options.getId(data[i]), data[i]);
      }

      this.persister?.setItem('data', persistedMap);
    });
  }

  public drop = async (ids: string[]) => {
    const removedIdsMap = new Map();

    if (!ids.length) return removedIdsMap;

    try {
      const items = await this.persister?.getItem<Map<string, T>>('data');

      ids.forEach((id) => {
        removedIdsMap.set(id, true);

        this.value?.delete(id);
        items?.delete(id);
      });

      await this.persister?.setItem('data', items);

      return removedIdsMap;
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    }

    return removedIdsMap;
  };

  public hydrate = async (options: {
    idsToDrop: string[];
    getId: (data: T) => string;
  }) => {
    const removedIdsMap = await this.drop(options?.idsToDrop ?? []);

    try {
      const stores: [string, TSyncable][] = [];

      const persistedData = await this.persister?.getItem<T[]>('data');

      persistedData?.forEach((data) => {
        const id = options.getId(data);
        const syncableItem = new this.SyncableStore(
          this.root,
          this.transport,
          data,
          this.channel,
        );

        if (!removedIdsMap.has(id)) {
          stores.push([id, syncableItem as TSyncable]);
        }
      });

      runInAction(() => {
        this.value = new Map<string, TSyncable>(stores);
      });
    } catch (e) {
      console.error('Failed to hydrate group', e);
    }
    runInAction(() => {
      this.isHydrated = true;
    });
  };

  public sync(operation: GroupOperation) {
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
        true,
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
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            this.SyncableStore.getDefaultValue() as any,
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

  private async initPersister() {
    this.persister = Persister.getInstance(this.persisterKey);
  }

  public async checkIfCanHydrate() {
    try {
      const isBootstrapped = await this.persister?.getItem<boolean>(
        'isBootstrapped',
      );

      runInAction(() => {
        this.isBootstrapped = isBootstrapped ?? false;
      });

      return isBootstrapped;
    } catch (e) {
      console.error('Failed to get persisted ids length', e);
    }
  }

  public async getRecentChanges() {}

  static SyncableStore = Syncable;
}
