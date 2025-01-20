import type { RootStore } from '@store/root';

import { Channel } from 'phoenix';
import { P, match } from 'ts-pattern';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@infra/transport';
import { Store, makeAutoSyncable } from '@store/store';
import { makeAutoSyncableGroup } from '@store/group-store';

import { GetTagsQuery } from './__service__/getTags.generated';

export type TagDatum = NonNullable<GetTagsQuery['tags'][number]>;

import { EntityType } from '@graphql/types';

import { TagService } from './__service__/Tag.service';

export class TagStore implements Store<TagDatum> {
  value: TagDatum = defaultValue;
  version: number = 0;
  isLoading = false;
  history: Operation[] = [];
  isBootstrapped: boolean = false;
  error: string | null = null;
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncable.load<TagDatum>();
  update = makeAutoSyncable.update<TagDatum>();
  private service = TagService.getInstance();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncable(this, {
      channelName: 'Tag',
      mutator: this.save,
      getId: (item) => item?.metadata.id,
    });
  }

  init(data: TagDatum): TagDatum {
    this.value = data;

    return data;
  }

  get tagName() {
    return this.value.name;
  }

  get colorScheme() {
    return this.value.colorCode;
  }

  get id() {
    return this.value.metadata.id;
  }

  set id(id: string) {
    this.value.metadata.id = id;
  }

  async bootstrap() {}

  async invalidate() {}

  async updateTag() {
    try {
      await this.service.updateTag({
        input: {
          id: this.value.metadata.id,
          name: this.value.name,
          colorCode: this.value.colorCode,
        },
      });
    } catch (error) {
      console.error(error);
    }
    this.root.ui.toastSuccess('Tag updated', 'tag-updated');
  }

  private async save(operation: Operation) {
    const diff = operation.diff?.[0];
    const type = diff?.op;
    const path = diff?.path;

    match(path)
      .with(['name', ...P.array()], () => {
        if (type === 'update') {
          this.updateTag();
        }
      })
      .with(['colorCode', ...P.array()], () => {
        if (type === 'update') {
          this.updateTag();
        }
      });
  }
}

const defaultValue: TagDatum = {
  name: '',
  entityType: EntityType.Organization,
  colorCode: 'grayModern',
  metadata: {
    id: crypto.randomUUID(),
  },
};
