import type { Channel } from 'phoenix';

import { match } from 'ts-pattern';
import { when, action, observable, runInAction } from 'mobx';

import type { RootStore } from './root';
import type { Transport } from './transport';

import { GroupOperation, GroupSyncPacket } from './types';
import { Persister, PersisterInstance } from './persister';

type StoreOptions = {
  name: string;
};

export class Store<T> {
  @observable accessor isLoading = false;
  @observable accessor isHydrated = false;
  @observable accessor isBootstrapped = false;
  @observable accessor isBootstrapping = false;
  @observable accessor error: string | null = null;
  @observable accessor value: Map<string, T> = new Map();
  channel?: Channel;
  persister?: PersisterInstance;
  private options: StoreOptions;

  constructor(
    public root: RootStore,
    public transport: Transport,
    opts: StoreOptions,
  ) {
    this.options = opts;
    when(
      () => !!this.root.session.sessionToken && !this.root.demoMode,
      () => {
        this.persister = Persister.getInstance(opts.name);
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

    window.addEventListener('focus', async () => {
      await this.getRecentChanges();
    });
  }

  public async getRecentChanges() {}

  public async bootstrap() {}

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
    // getId: (data: T) => string;
  }) => {
    return;
    await this.drop(options?.idsToDrop ?? []);

    try {
      // const stores: [string, T][] = [];

      const persistedData = await this.persister?.getItem<Map<string, T>>(
        'data',
      );

      // persistedData?.forEach((data) => {
      //   const id = options.getId(data);
      //   const syncableItem = new this.SyncableStore(
      //     this.root,
      //     this.transport,
      //     data,
      //     this.channel,
      //   );
      //
      //   if (!removedIdsMap.has(id)) {
      //     stores.push([id, syncableItem as TSyncable]);
      //   }
      // });

      runInAction(() => {
        if (!persistedData) return;
        this.value = persistedData;
      });
    } catch (e) {
      console.error('Failed to hydrate group', e);
    }
    runInAction(() => {
      this.isHydrated = true;
    });
  };

  private async initChannelConnection(tenant: string) {
    try {
      const connection = await this.transport.join(
        this.options.name,
        tenant,
        0,
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
    });
  }

  private applyGroupOperation(operation: GroupOperation) {
    match(operation.action)
      .with('APPEND', () => {
        operation.ids.forEach((_id) => {
          // const newSyncableItem = new this.SyncableStore(
          //   this.root,
          //   this.transport,
          //   // eslint-disable-next-line @typescript-eslint/no-explicit-any
          //   this.SyncableStore.getDefaultValue() as any,
          // );

          runInAction(() => {
            // newSyncableItem.setId(id);
            // this.value.set(id, newSyncableItem as TSyncable);
          });

          setTimeout(() => {
            // this.value.get(id)?.invalidate();
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

          // item.invalidate();
        });
      })
      .otherwise(() => {});
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

  public toActiveArray() {
    return this.toArray().map((v) => new ActiveRecord(this, v));
  }

  public toArray() {
    const arr = Array.from(this.value.values());

    return arr;
  }

  public toComputedArray(
    compute: (arr: ActiveRecord<T>[]) => ActiveRecord<T>[],
  ) {
    const arr = this.toActiveArray();

    return compute(arr);
  }

  public getById(id: string) {
    const data = this.value.get(id);

    return new ActiveRecord(this, data);
  }
}

export class ActiveRecord<T> {
  value: T;

  constructor(private store: Store<T>, data: T) {
    this.value = data;
  }

  commit() {}
}
