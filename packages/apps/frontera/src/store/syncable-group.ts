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
  canBypassBootstrap = false;
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
      hydrate: action,
      subscribe: action,
      error: observable,
      value: observable,
      version: observable,
      channel: observable,
      history: observable,
      isHydrated: observable,
      isLoading: observable,
      channelName: computed,
      initPersister: action,
      checkIfCanHydrate: action,
      isBootstrapped: observable,
      applyGroupOperation: action,
      initChannelConnection: action,
      canBypassBootstrap: observable,
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

    this.isBootstrapped = true;
  }

  public async hydrate() {
    try {
      const stores: [string, TSyncable][] = [];

      await this.persister?.iterate<T, void>((data, id) => {
        const syncableItem = new this.SyncableStore(
          this.root,
          this.transport,
          data,
          this.channel,
        );

        stores.push([id, syncableItem as TSyncable]);
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
  }

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
      const idsCount = await this.persister?.length();
      const canBypass = typeof idsCount !== 'undefined' && idsCount > 0;

      runInAction(() => {
        this.canBypassBootstrap = canBypass;
        this.isBootstrapped = true;
      });

      return canBypass;
    } catch (e) {
      console.error('Failed to get persisted ids length', e);
    }
  }

  static SyncableStore = Syncable;
}
