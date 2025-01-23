import type { Channel } from 'phoenix';
import type { Transport } from '@infra/transport';

import set from 'lodash/set';
import { match } from 'ts-pattern';
import { getDiff, applyDiff } from 'recursive-diff';
import {
  when,
  toJS,
  action,
  reaction,
  computed,
  intercept,
  observable,
  isObservable,
  makeAutoObservable,
} from 'mobx';

import type { RootStore } from './root';

import { Persister, PersisterInstance } from './persister';
import {
  Operation,
  SyncPacket,
  GroupOperation,
  DTOFactoryClass,
  GroupSyncPacket,
} from './types';

type StoreOptions<T extends object> = {
  name: string;
  getId: (data: T) => string;
  factory: DTOFactoryClass<T & object>;
};

export class Store<T extends object> {
  @observable accessor size = 0;
  @observable accessor version = 0;
  @observable accessor totalElements = 0;
  @observable accessor isLoading = false;
  @observable accessor isHydrated = false;
  @observable accessor isBootstrapped = false;
  @observable accessor isBootstrapping = false;
  @observable accessor error: string | null = null;
  @observable accessor active: Map<string, T> = new Map();
  @observable accessor range: [startIndex: number, endIndex: number] = [0, 0];

  channel?: Channel;
  options: StoreOptions<T>;
  persister?: PersisterInstance;
  value: Map<string, T> = new Map();

  private snapshots: Map<string, T> = new Map();
  @observable private accessor views: Map<string, string[]> = new Map();

  @computed
  get activeArray() {
    return this.toArray()
      .slice(this.range[0], this.range[1] + 1)
      .map((obj) => new ActiveRecord(this, obj));
  }

  constructor(
    public root: RootStore,
    public transport: Transport,
    opts: StoreOptions<T>,
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
        if (this.active.size === 0) {
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

    return (view ?? []).map(
      (id) => new ActiveRecord(this, this.value.get(id) as T),
    );
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
  public hydrate = async (options: { idsToDrop: string[] }) => {
    await this.drop(options?.idsToDrop ?? []);

    try {
      const persisted = await this.persister?.getItem<Map<string, T>>('data');

      if (!persisted) return;

      const initialized = new Map<string, T>();

      persisted.forEach((v, k) => {
        initialized.set(k, this.options.factory.of(this.root, v));
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

    if (this.active.has(targetId)) {
      const activeRecord = this.active.get(targetId);

      applyDiff(activeRecord, diff);
    }

    this.version++;
  }

  @action
  private applyGroupOperation(operation: GroupOperation) {
    match(operation.action)
      .with('APPEND', () => {
        operation.ids.forEach((id) => {
          const { factory } = this.options;
          const dto = factory.of(this.root, factory.default());

          set(dto, 'id', id);
          this.value.set(id, dto);

          this.size++;

          setTimeout(() => {
            this.invalidate(id);
          }, 1000);
        });
      })
      .with('DELETE', () => {
        operation.ids.forEach((id) => {
          this.active.delete(id);
          this.value.delete(id);
          this.size--;
        });
        // this.version++;
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

  public toComputedArray(compute: (arr: T[]) => T[]): T[] {
    const arr = compute(this.toArray());

    return arr;
  }

  public getById(id: string) {
    const data = this.value.get(id);

    return new ActiveRecord(this, data!);
  }

  public snapshot(id: string, current: T) {
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
      const { factory } = this.options;
      const instance = toJS(this.value.get(id));

      if (!instance) return;
      const current = factory.toPersistable(instance);

      const persistedData = await this.persister?.getItem<Map<string, T>>(
        'data',
      );

      persistedData?.set(id, current);

      await this.persister?.setItem('data', persistedData);
    } catch (e) {
      console.error('Failed to persist', e);
    }
  }

  public persistGroup() {
    this.persister?.getItem<Map<string, T>>('data', (err) => {
      if (err) {
        console.error('Failed to get persisted data', err);

        return;
      }

      const payload = new Map<string, T>();

      this.value.forEach((v, k) =>
        payload.set(k, this.options.factory.toPersistable(v)),
      );

      this.persister?.setItem('data', payload, (err) => {
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
    const { factory } = this.options;
    const lhs = factory.toPersistable(this.snapshots.get(id)!);
    const rhs = factory.toPersistable(toJS(this.value.get(id)!));

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
  public setView = (key: string, filterFn: (data: T[]) => string[]) => {
    this.views.set(key, filterFn(this.toArray()));
  };

  public findOne(selector: (object: T, idx: number, arr: T[]) => boolean) {
    const record = this.toArray().find(selector);

    if (!record) return null;

    return new ActiveRecord(this, record);
  }

  public findMany(selector: (object: T, idx: number, arr: T[]) => boolean) {
    return this.toArray().filter(selector);
  }

  public findManyActive(
    selector: (object: T, idx: number, arr: T[]) => boolean,
  ) {
    return this.toArray()
      .filter(selector)
      .map((record) => new ActiveRecord(this, record));
  }
}

export class ActiveRecord<T extends object> {
  private _value: T;

  constructor(private store: Store<T>, data: T) {
    const id = this.store.options.getId(data);

    if (!isObservable(data)) {
      if (!this.store.active.has(id)) {
        this.store.active.set(id, makeAutoObservable(data));
      }
    }

    this._value = data;

    if (isObservable(this._value)) {
      intercept(this._value!, (change) => {
        if (!this.store.hasSnapshot(id)) {
          const current = toJS(this._value);

          this.store.snapshot(id, current);
        }

        return change;
      });
    }
  }

  get value() {
    return this._value;
  }

  get id() {
    return this.store.options.getId(this.value);
  }

  public commit(
    opts: {
      syncOnly?: boolean;
      onFailled?: () => void;
      onCompleted?: () => void;
    } = { syncOnly: false },
  ) {
    this.store.commit(this.id, opts);
  }
}
