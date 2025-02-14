import type { Channel } from 'phoenix';
import type { Transport } from '@infra/transport';

import set from 'lodash/set';
import { match } from 'ts-pattern';
import { getDiff, applyDiff } from 'recursive-diff';
import { when, action, observable, runInAction } from 'mobx';

import type { RootStore } from './root';
import type { Entity, EntityFactoryClass } from './record';

import { Persister, PersisterInstance } from './persister';
import {
  Operation,
  SyncPacket,
  GroupOperation,
  GroupSyncPacket,
} from './types';

type StoreOptions<T extends object> = {
  name: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  factory: any;
  getId: (data: T) => string;
};

export class Store<T extends object, E extends Entity<T> = Entity<T>> {
  @observable accessor size = 0;
  @observable accessor version = 0;
  @observable accessor totalElements = 0;
  @observable accessor isLoading = false;
  @observable accessor isHydrated = false;
  @observable accessor isBootstrapped = false;
  @observable accessor isBootstrapping = false;
  @observable accessor error: string | null = null;
  @observable accessor value: Map<string, E> = new Map();
  @observable accessor searchResults: Map<string, string[]> = new Map();

  channel?: Channel;
  options: StoreOptions<T>;
  persister?: PersisterInstance;

  private snapshots: Map<string, T> = new Map();
  @observable public accessor views: Map<string, E[]> = new Map();
  @observable private accessor searchTerms: Map<string, string> = new Map();

  constructor(
    public root: RootStore,
    public transport: Transport,
    opts: StoreOptions<T>,
  ) {
    this.options = opts;
    when(
      () => !!this.root.session?.sessionToken,
      () => {
        this.persister = Persister.getInstance(opts.name);
      },
    );

    when(
      () => !!this.root.session?.value?.tenant && !this.root.demoMode,
      async () => {
        const tenant = this.root.session?.value?.tenant;

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

    // temporary commented out -> will be used in the future;
    // window.addEventListener('focus', async () => {
    //   await this.getRecentChanges();
    // });
  }

  public getViewById(id: string) {
    const view = this.views.get(id);

    return view ?? [];
  }

  @action
  public async getRecentChanges() {}

  @action
  public async bootstrap() {}

  @action
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
      this.size = this.value.size;

      return removedIdsMap;
    } catch (e) {
      this.error = (e as Error)?.message;
    }

    return removedIdsMap;
  };

  @action
  public hydrate = async (options?: { idsToDrop: string[] }) => {
    await this.drop(options?.idsToDrop ?? []);

    try {
      const factory = this.options.factory as EntityFactoryClass<T, E>;
      const persisted = await this.persister?.getItem<Map<string, T>>('data');

      if (!persisted) {
        runInAction(() => {
          this.isHydrated = true;
        });

        return;
      }

      const initialized = new Map<string, E>();

      persisted.forEach((v, k) => {
        initialized.set(k, new factory(this, v));
      });

      runInAction(() => {
        this.value = initialized;
        this.totalElements = initialized.size;
      });
    } catch (e) {
      console.error('Failed to hydrate group', e);
    } finally {
      runInAction(() => {
        this.isHydrated = true;
      });
    }
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

    this.channel.on('sync_packet', (packet: SyncPacket) => {
      if (packet.operation.ref === this.transport.refId) return;
      this.applyOperation(packet);
    });
  }

  public sync = (operation: GroupOperation) => {
    const op = {
      ...operation,
      ref: this.transport.refId,
    };

    this?.channel
      ?.push('sync_group_packet', { payload: { operation: op } })
      ?.receive('ok', () => {});
  };

  @action
  private applyOperation(packet: SyncPacket) {
    const targetId = packet.operation.entityId;

    if (!targetId) return;

    const target = this.value.get(targetId);

    if (!target) return;

    const diff = packet.operation.diff;

    applyDiff(target.value, diff);

    this.version++;
  }

  @action
  private applyGroupOperation(operation: GroupOperation) {
    match(operation.action)
      .with('APPEND', () => {
        operation.ids.forEach((id) => {
          const factory = this.options.factory as EntityFactoryClass<T, E>;
          const record = new factory(this, factory.default!());

          set(record, 'id', id);
          this.value.set(id, record);

          this.searchResults.forEach((v) => {
            if (v.includes(id)) return;
            v.unshift(id);
          });
          this.size++;
          this.version++;

          setTimeout(() => {
            this.invalidate(id);
          }, 1000);
        });
      })
      .with('DELETE', () => {
        operation.ids.forEach((id) => {
          this.value.delete(id);
          this.size--;
        });
      })
      .with('INVALIDATE', () => {
        operation.ids.forEach((id) => {
          this.invalidate(id);
        });
      })
      .otherwise(() => {});
  }

  @action
  public async checkIfCanHydrate() {
    try {
      const isBootstrapped = await this.persister?.getItem<boolean>(
        'isBootstrapped',
      );

      this.isBootstrapped = isBootstrapped ?? false;

      return isBootstrapped;
    } catch (e) {
      console.error('Failed to get persisted ids length', e);
    }
  }

  public toArray() {
    const arr = Array.from(this.value.values());

    return arr;
  }

  public toComputedArray(compute: (arr: E[]) => E[]) {
    const arr = compute(this.toArray());

    return arr;
  }

  public getById(id: string): E | null {
    if (!this.value || typeof id !== 'string') return null;
    const data = this.value.get(id);

    return data as E;
  }

  public snapshot(id: string, current: T) {
    if (this.hasSnapshot(id)) return;

    this.snapshots.set(id, current);
  }

  public getSnapshot(id: string) {
    return this.snapshots.get(id);
  }

  public hasSnapshot(id: string) {
    return this.snapshots.has(id);
  }

  public clearSnapshot(id: string) {
    this.snapshots.delete(id);
  }

  @action
  public setSearchTerm(term: string, viewId: string) {
    this.searchTerms.set(viewId, term);
  }

  public getSearchTermByView(viewId: string) {
    return this.searchTerms.get(viewId);
  }

  @action
  public invalidate(_id: string) {}

  public commit(
    id: string,
    opts: {
      syncOnly?: boolean;
      onFailled?: () => void;
      onCompleted?: () => void;
    } = { syncOnly: false },
  ) {
    const operation = this.makeChangesetOperation(id);

    this.clearSnapshot(id);

    this.root.transactions.commit(operation, {
      ...opts,
      // persist: () => this.persist(id),
    });

    // this.reconcile(id);
    this.version++;
  }

  // private async persist(id: string) {
  //   try {
  //     const record = this.value.get(id);

  //     if (!record) return;
  //     const data = record.toRaw();

  //     const persisted = await this.persister?.getItem<Map<string, T>>('data');

  //     persisted?.set(id, data);
  //     await this.persister?.setItem('data', persisted);
  //   } catch (e) {
  //     console.error('Failed to persist', e);
  //   }
  // }

  // public persistGroup() {
  //   this.persister?.getItem<Map<string, T>>('data', (err) => {
  //     if (err) {
  //       console.error('Failed to get persisted data', err);

  //       return;
  //     }

  //     const persisted = new Map<string, T>();

  //     this.value.forEach((v, k) => persisted.set(k, v.toRaw()));

  //     this.persister?.setItem('data', persisted, (err) => {
  //       if (err) {
  //         console.error('Failed to persist store data', err);
  //       }
  //     });
  //   });
  // }

  private makeChangesetOperation(id: string) {
    const lhs = this.snapshots.get(id)!;
    const rhs = this.value.get(id)!.toRaw();

    const diff = getDiff(lhs, rhs, true);

    const operation: Operation = {
      id: 0,
      diff,
      entityId: id,
      ref: this.transport.refId,
      tenant: this.root.session.value.tenant,
      entity: this.options.name,
    };

    return operation;
  }

  @action
  public setView = (key: string, filterFn: (records: E[]) => E[]) => {
    this.views.set(key, filterFn(this.toArray()));
  };

  public findOne(selector: (object: E, idx: number, arr: E[]) => boolean) {
    const record = this.toArray().find(selector);

    if (!record) return null;

    return record;
  }

  public findMany(selector: (object: E, idx: number, arr: E[]) => boolean) {
    return this.toArray().filter(selector);
  }
}
