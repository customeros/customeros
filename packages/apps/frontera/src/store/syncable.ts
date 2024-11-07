import type { Channel } from 'phoenix';

import { getDiff, applyDiff } from 'recursive-diff';
import { Persister, type PersisterInstance } from '@store/persister';
import {
  toJS,
  when,
  action,
  computed,
  observable,
  runInAction,
  makeObservable,
} from 'mobx';

import type { RootStore } from './root';
import type { Transport } from './transport';
import type { Operation, SyncPacket } from './types';

type SyncableUpdateOptions = {
  mutate?: boolean;
  syncMutate?: boolean;
};

export class Syncable<T extends object> {
  value: T;
  snapshot: T;
  isLoading = false;
  error: string | null = null;
  channel: Channel | null = null;
  private persister?: PersisterInstance;

  constructor(
    public root: RootStore,
    public transport: Transport,
    data: T,
    _channel?: Channel,
  ) {
    this.value = data;
    this.snapshot = Object.assign({}, data);
    this.persister = Persister.getInstance(this.getChannelName().split(':')[0]);
    this.persist = this.persist.bind(this);
    this.subscribe = this.subscribe.bind(this);

    makeObservable<Syncable<T>, 'subscribe'>(this, {
      id: computed,
      load: action,
      save: action,
      getId: action,
      setId: action,
      update: action,
      subscribe: action,
      error: observable,
      value: observable,
      isLoading: observable,
      getChannelName: action,
    });

    when(
      () => !!this.root.session.value.tenant && !this.root.demoMode,
      () => {
        const channel = this.transport.channels.get(this.getChannelName());

        if (channel) {
          this.channel = channel;
          this.subscribe();
        }
      },
    );
  }

  get id() {
    if (!this.value || !('id' in this.value)) return '';

    return this.value?.id as string;
  }

  getId() {
    if (!this.value || !('id' in this.value)) return '';

    return this.value?.id as string;
  }

  setId(id: string) {
    if (!this.value || !('id' in this.value)) return;

    this.value.id = id;
  }

  getChannelName() {
    return '';
  }

  public async load(data: T) {
    requestIdleCallback(() => {
      runInAction(() => {
        Object.assign(this.value, data);
        Object.assign(this.snapshot, data);
      });
    });
  }

  private subscribe() {
    if (!this.channel || this.root.demoMode) return;

    this.channel.on('sync_packet', (packet: SyncPacket) => {
      if (
        packet.operation?.entityId !== this.getId() ||
        packet.operation.ref === this.transport.refId
      )
        return;

      const prev = toJS(this.value);
      const diff = packet.operation.diff;

      const next = applyDiff(prev, diff);

      runInAction(() => {
        this.value = next;
      });
    });
  }

  /**
   * @deprecated
   * use Syncable.commit instead.
   */
  public update(updater: (prev: T) => T, options?: SyncableUpdateOptions) {
    const lhs = toJS(this.value);
    const next = updater(this.value);
    const rhs = toJS(next);
    const diff = getDiff(lhs, rhs, true);

    const operation: Operation = {
      id: 0,
      diff,
      entityId: this.getId(),
      ref: this.transport.refId,
      entity: this.getChannelName().split(':')[0],
    };

    this.value = next;

    if (this?.save) {
      (async () => {
        try {
          this.error = null;

          if (options?.mutate && !this.root.demoMode) {
            await this.save(operation);
          }

          this?.channel
            ?.push('sync_packet', { payload: { operation } })
            ?.receive('ok', ({ version: _version }: { version: number }) => {});
        } catch (e) {
          console.error(e);
          this.value = lhs;
        }
      })();
    }
  }

  public commit(
    opts: {
      syncOnly?: boolean;
      onFailled?: () => void;
      onCompleted?: () => void;
    } = { syncOnly: false },
  ) {
    const operation = this.makeChangesetOperation();

    Object.assign(this.snapshot, toJS(this.value));

    this.root.transactions.commit(operation, {
      ...opts,
      persist: this.persist,
    });
  }

  /**
   * @deprecated
   * will be removed as soon as we migrate all stores to tx queues
   */
  public save(_operation: Operation) {
    /* Placeholder: should be overwritten by sub-classes with the apropiate mutation logic */
  }

  public async invalidate() {
    /* Placeholder: should be overwritten by sub-classes with the apropiate invalidation logic */
  }

  private makeChangesetOperation() {
    const lhs = toJS(this.snapshot);
    const rhs = toJS(this.value);
    const diff = getDiff(lhs, rhs, true);

    const operation: Operation = {
      id: 0,
      diff,
      entityId: this.id || this.getId(),
      ref: this.transport.refId,
      tenant: this.root.session.value.tenant,
      entity: this.getChannelName().split(':')[0],
    };

    return operation;
  }

  private async persist() {
    try {
      const persistedData = await this.persister?.getItem<Map<string, T>>(
        'data',
      );

      persistedData?.set(this.getId(), toJS(this.snapshot));
    } catch (e) {
      console.error('Failed to persist', e);
    }
  }

  static getDefaultValue() {
    return {};
  }
}
