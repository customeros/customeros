import type { Channel } from 'phoenix';

import set from 'lodash/set';
import { match } from 'ts-pattern';
import { getDiff, applyDiff } from 'recursive-diff';
import { when, action, reaction, observable } from 'mobx';

import type { RootStore } from './root';
import type { Transport } from './transport';
import type { Record, RecordFactoryClass } from './record';

import { Persister, PersisterInstance } from './persister';
import {
  Operation,
  SyncPacket,
  GroupOperation,
  GroupSyncPacket,
} from './types';

type ValueOf<R extends Record> = R['value'];

type StoreOptions<R extends Record> = {
  name: string;
  factory: RecordFactoryClass<R>;
  getId: (data: ValueOf<R>) => string;
};

export class Store<R extends Record> {
  @observable accessor size = 0;
  @observable accessor version = 0;
  @observable accessor totalElements = 0;
  @observable accessor isLoading = false;
  @observable accessor isHydrated = false;
  @observable accessor isBootstrapped = false;
  @observable accessor isBootstrapping = false;
  @observable accessor error: string | null = null;
  @observable accessor value: Map<string, R> = new Map();
  @observable accessor range: [startIndex: number, endIndex: number] = [0, 0];

  channel?: Channel;
  options: StoreOptions<R>;
  persister?: PersisterInstance;

  private snapshots: Map<string, ValueOf<R>> = new Map();
  @observable private accessor views: Map<string, R[]> = new Map();

  constructor(
    public root: RootStore,
    public transport: Transport,
    opts: StoreOptions<R>,
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

    reaction(
      () => this.size,
      () => {
        if (this.value.size === 0) {
          this.setActiveRange(0, 99);
        }
        this.persistGroup();
      },
    );

    window.addEventListener('focus', async () => {
      await this.getRecentChanges();
    });
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
      const items = await this.persister?.getItem<Map<string, ValueOf<R>>>(
        'data',
      );

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
  public hydrate = async (options: { idsToDrop: string[] }) => {
    await this.drop(options?.idsToDrop ?? []);

    try {
      const { factory } = this.options;
      const persisted = await this.persister?.getItem<Map<string, ValueOf<R>>>(
        'data',
      );

      if (!persisted) return;

      const initialized = new Map<string, R>();

      persisted.forEach((v, k) => {
        initialized.set(k, new factory(this, v));
      });

      this.value = initialized;
      this.totalElements = initialized.size;
    } catch (e) {
      console.error('Failed to hydrate group', e);
    }
    this.isHydrated = true;
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

  private applyOperation(packet: SyncPacket) {
    const targetId = packet.operation.entityId;

    if (!targetId) return;

    const target = this.value.get(targetId);

    if (!target) return;

    const diff = packet.operation.diff;

    applyDiff(target, diff);

    this.version++;
  }

  @action
  private applyGroupOperation(operation: GroupOperation) {
    match(operation.action)
      .with('APPEND', () => {
        operation.ids.forEach((id) => {
          const { factory } = this.options;
          const record = new factory(this, factory.default!());

          set(record, 'id', id);
          this.value.set(id, record);

          this.size++;

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

  public toComputedArray(compute: (arr: R[]) => R[]) {
    const arr = compute(this.toArray());

    return arr;
  }

  public getById(id: string) {
    const data = this.value.get(id);

    return data as R;
  }

  public snapshot(id: string, current: R) {
    if (this.hasSnapshot(id)) return;

    this.snapshots.set(id, current);
  }

  public hasSnapshot(id: string) {
    return this.snapshots.has(id);
  }

  public clearSnapshot(id: string) {
    this.snapshots.delete(id);
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
      persist: () => this.persist(id),
    });
  }

  private async persist(id: string) {
    try {
      const record = this.value.get(id);

      if (!record) return;
      const data = record.toRaw();

      const persisted = await this.persister?.getItem<Map<string, ValueOf<R>>>(
        'data',
      );

      persisted?.set(id, data);

      await this.persister?.setItem('data', persisted);
    } catch (e) {
      console.error('Failed to persist', e);
    }
  }

  public persistGroup() {
    this.persister?.getItem<Map<string, ValueOf<R>>>('data', (err) => {
      if (err) {
        console.error('Failed to get persisted data', err);

        return;
      }

      const persisted = new Map<string, ValueOf<R>>();

      this.value.forEach((v, k) => persisted.set(k, v.toRaw()));

      this.persister?.setItem('data', persisted, (err) => {
        if (err) {
          console.error('Failed to persist store data', err);
        }
      });
    });
  }

  @action
  public setActiveRange = (startIndex: number, endIndex: number) => {
    this.range[0] = startIndex;
    this.range[1] = endIndex;
  };

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
  public setView = (key: string, filterFn: (records: R[]) => R[]) => {
    this.views.set(key, filterFn(this.toArray()));
  };

  public findOne(selector: (object: R, idx: number, arr: R[]) => boolean) {
    const record = this.toArray().find(selector);

    if (!record) return null;

    return record;
  }

  public findMany(selector: (object: R, idx: number, arr: R[]) => boolean) {
    return this.toArray().filter(selector);
  }
}
