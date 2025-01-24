import { Channel } from 'phoenix';
import { RootStore } from '@store/root';
import { Transport } from '@infra/transport';
import { GroupOperation } from '@store/types';
import { UserService } from '@store/Users/User.service.ts';
import { computed, runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

import { User } from '@graphql/types';

import { UserStore } from './User.store';

export class UsersStore implements GroupStore<User> {
  channel?: Channel | undefined;
  error: string | null = null;
  history: GroupOperation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  totalElements: number = 0;
  value: Map<string, UserStore> = new Map();
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncableGroup.load<User>();
  subscribe = makeAutoSyncableGroup.subscribe;
  private service: UserService;

  constructor(public root: RootStore, public transport: Transport) {
    this.service = UserService.getInstance(transport);
    makeAutoSyncableGroup(this, {
      channelName: 'Users',
      ItemStore: UserStore,
      getItemId: (user) => user.id,
    });
    makeAutoObservable(this, {
      usersWithoutBotsAndInternal: computed,
      tenantUsers: computed,
    });
  }

  get usersWithoutBotsAndInternal() {
    return this.toComputedArray((users) =>
      users.filter((user) => !user.value.bot && !user.value.internal),
    );
  }

  get tenantUsers() {
    return this.toComputedArray((users) =>
      users.filter(
        (user) => !user.value.bot && !user.value.test && !user.value.internal,
      ),
    );
  }

  async bootstrap() {
    if (this.isBootstrapped || this.isLoading) return;

    try {
      this.isLoading = true;

      const { users } = await this.service.getUsers({
        pagination: {
          limit: 1000,
          page: 0,
        },
      });

      this.load(users.content);

      runInAction(() => {
        this.isBootstrapped = true;
        this.totalElements = users.totalElements;
      });
    } catch (error) {
      runInAction(() => {
        this.error = (error as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  toArray() {
    return Array.from(this.value.values());
  }

  toComputedArray(compute: (arr: UserStore[]) => UserStore[]) {
    return compute(this.toArray());
  }
}
