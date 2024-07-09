import { merge } from 'lodash';
import { Channel } from 'phoenix';
import { Store } from '@store/store';
import { gql } from 'graphql-request';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { GroupOperation } from '@store/types';
import { runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

import { Tag, TagInput } from '@shared/types/__generated__/graphql.types';

import { TagStore } from './Tag.store';

export class TagsStore implements GroupStore<Tag> {
  channel?: Channel | undefined;
  error: string | null = null;
  history: GroupOperation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  totalElements: number = 0;
  value: Map<string, TagStore> = new Map();
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncableGroup.load<Tag>();
  subscribe = makeAutoSyncableGroup.subscribe;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: 'Tags',
      getItemId: (item) => item?.id,
      ItemStore: TagStore,
    });
  }

  async invalidate() {
    try {
      this.isLoading = true;
      const { tags } = await this.transport.graphql.request<TAG_QUERY_RESPONSE>(
        TAG_QUERY,
      );
      runInAction(() => {
        this.load(tags);
        this.isBootstrapped = true;
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  async bootstrap() {
    if (this.isBootstrapped || this.isLoading) return;

    this.invalidate();
  }

  create = async (
    payload?: TagInput,
    options?: { onSucces?: (serverId: string) => void },
  ) => {
    const newTag = new TagStore(this.root, this.transport);
    const tempId = newTag.value.id;
    let serverId = '';

    if (payload) {
      merge(newTag.value, payload);
    }

    this.value.set(tempId, newTag);

    try {
      const { tag_Create } = await this.transport.graphql.request<
        CREATE_TAG_RESPONSE,
        CREATE_TAG_PAYLOAD
      >(CREATE_TAG_MUTATION, {
        input: {
          name: payload?.name || '',
        },
      });

      runInAction(() => {
        serverId = tag_Create.id;

        newTag.value.id = serverId;

        this.value.set(serverId, newTag);
        this.value.delete(tempId);

        this.sync({
          action: 'APPEND',
          ids: [serverId],
        });
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    } finally {
      if (serverId) {
        options?.onSucces?.(serverId);
      }
    }
  };

  toArray() {
    return Array.from(this.value.values());
  }

  toComputedArray<T extends Store<Tag>>(compute: (arr: Store<Tag>[]) => T[]) {
    const arr = this.toArray();

    return compute(arr);
  }
}

type TAG_QUERY_RESPONSE = {
  tags: Tag[];
};

const TAG_QUERY = gql`
  query getTags {
    tags {
      metadata {
        id
        source
        sourceOfTruth
        appSource
        created
        lastUpdated
      }
      id
      name
      source
      updatedAt
      createdAt
      appSource
    }
  }
`;

type CREATE_TAG_PAYLOAD = {
  input: TagInput;
};

type CREATE_TAG_RESPONSE = {
  tag_Create: Tag;
};

const CREATE_TAG_MUTATION = gql`
  mutation createTag($input: TagInput!) {
    tag_Create(input: $input) {
      name
      id
    }
  }
`;
