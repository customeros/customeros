import { merge } from 'lodash';
import { Channel } from 'phoenix';
import { Store } from '@store/store';
import { RootStore } from '@store/root';
import { Transport } from '@infra/transport';
import { GroupOperation } from '@store/types';
import { runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store';

import { TagInput, EntityType } from '@graphql/types';

import { TagStore, type TagDatum } from './Tag.store';
import { TagService } from './__service__/Tag.service';

export class TagsStore implements GroupStore<TagDatum> {
  channel?: Channel | undefined;
  error: string | null = null;
  history: GroupOperation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  totalElements: number = 0;
  value: Map<string, TagStore> = new Map();
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncableGroup.load<TagDatum>();
  subscribe = makeAutoSyncableGroup.subscribe;
  private service = TagService.getInstance();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: 'Tags',
      getItemId: (item) => item?.metadata.id,
      ItemStore: TagStore,
    });
  }

  async invalidate() {
    try {
      this.isLoading = true;

      const { tags } = await this.service.getTags();

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

  async getTagsByEntityType(entityType: EntityType) {
    try {
      this.isLoading = true;

      const { tags_ByEntityType } = await this.service.getTagsByEntityType(
        entityType as EntityType,
      );

      runInAction(() => {
        this.load(tags_ByEntityType);
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

  createNew = (payload?: Partial<TagDatum>) => {
    const newTag = new TagStore(this.root, this.transport);

    if (payload) {
      merge(newTag.value, payload);
    }

    this.value.set(newTag.value.metadata.id, newTag);

    return newTag;
  };

  create = async (
    payload?: TagInput,
    options?: { onSucces?: (serverId: string) => void },
  ) => {
    const newTag = new TagStore(this.root, this.transport);
    const tempId = newTag.value.metadata.id;
    let serverId = '';

    if (payload) {
      merge(newTag.value, payload);
    }

    this.value.set(tempId, newTag);

    try {
      const { tag_Create } = await this.service.createTag({
        input: {
          name: payload?.name || '',
          entityType: payload?.entityType,
          colorCode: payload?.colorCode,
        },
      });

      runInAction(() => {
        serverId = tag_Create.metadata.id;

        newTag.value.metadata.id = serverId;

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
    this.root.ui.toastSuccess('Tag created', 'tag-created');
  };

  async deleteTag(tagId: string) {
    this.value.delete(tagId);

    try {
      this.isLoading = true;
      await this.service.deleteTag({
        id: tagId,
      });
    } catch (error) {
      console.error(error);
    } finally {
      this.isLoading = false;
    }
    this.sync({
      action: 'DELETE',
      ids: [tagId],
    });
    this.root.ui.toastSuccess('Tag deleted', 'tag-deleted');
  }

  toArray() {
    return Array.from(this.value.values());
  }

  toComputedArray(compute: (arr: Store<TagDatum>[]) => Store<TagDatum>[]) {
    const arr = this.toArray();

    return compute(arr);
  }

  getById(id: string) {
    return this.value.get(id);
  }

  getByEntityType(entityType: EntityType) {
    const tags = this.toArray();

    return tags.filter((tag) => tag.value.entityType === entityType);
  }
}
