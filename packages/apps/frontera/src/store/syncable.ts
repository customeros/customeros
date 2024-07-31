import { Channel } from 'phoenix';
import { getDiff, applyDiff } from 'recursive-diff';
import {
  toJS,
  when,
  action,
  computed,
  observable,
  runInAction,
  makeObservable,
} from 'mobx';

import { RootStore } from './root';
import { Transport } from './transport';
import { Operation, SyncPacket } from './types';

type SyncableUpdateOptions = {
  mutate?: boolean;
  syncMutate?: boolean;
};

export class Syncable<T extends object> {
  value: T;
  version = 0;
  isLoading = false;
  history: Operation[] = [];
  error: string | null = null;
  channel: Channel | null = null;

  constructor(public root: RootStore, public transport: Transport, data: T) {
    this.value = data;
    makeObservable<Syncable<T>, 'initChannelConnection' | 'subscribe'>(this, {
      id: computed,
      load: action,
      save: action,
      getId: action,
      setId: action,
      update: action,
      subscribe: action,
      error: observable,
      value: observable,
      history: observable,
      channel: observable,
      version: observable,
      isLoading: observable,
      getChannelName: action,
      initChannelConnection: action,
    });

    when(
      () => !!this.root.session.value.tenant && !this.root.demoMode,
      async () => {
        try {
          // await this.initChannelConnection();
        } catch (e) {
          console.error(e);
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

  async load(data: T) {
    requestIdleCallback(() => {
      runInAction(() => {
        Object.assign(this.value, data);
        this.initChannelConnection();
      });
    });
  }

  private async initChannelConnection() {
    try {
      const connection = await this.transport.join(
        this.getChannelName(),
        this.id,
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

    this.channel.on('sync_packet', (packet: SyncPacket) => {
      if (packet.operation.ref === this.transport.refId) return;

      const prev = toJS(this.value);
      const diff = packet.operation.diff;

      const next = applyDiff(prev, diff);

      runInAction(() => {
        this.value = next;
        this.version = packet.version;
        this.history.push(packet.operation);
      });
    });
  }

  public update(
    updater: (prev: T) => T,
    options: SyncableUpdateOptions = {
      mutate: true,
      syncMutate: false,
    },
  ) {
    const lhs = toJS(this.value);
    const next = updater(this.value);
    const rhs = toJS(next);
    const diff = getDiff(lhs, rhs, true);

    const operation: Operation = {
      id: this.version,
      diff,
      ref: this.transport.refId,
    };

    this.history.push(operation);
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

  public save(_operation: Operation) {
    /* Placeholder: should be overwritten by sub-classes with the apropiate mutation logic */
  }

  public async invalidate() {
    /* Placeholder: should be overwritten by sub-classes with the apropiate invalidation logic */
  }

  static getDefaultValue<T extends object>(_data?: T): T {
    return {} as T;
  }
}
