import type { RootStore } from '@store/root';

import { Channel } from 'phoenix';
import { P, match } from 'ts-pattern';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';
import { makeAutoSyncableGroup } from '@store/group-store';

import { Tag, DataSource } from '@shared/types/__generated__/graphql.types';

import { TagService } from './Tag.service';

export class TagStore implements Store<Tag> {
  value: Tag = defaultValue;
  version: number = 0;
  isLoading = false;
  history: Operation[] = [];
  isBootstrapped: boolean = false;
  error: string | null = null;
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncable.load<Tag>();
  update = makeAutoSyncable.update<Tag>();
  private service: TagService;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncable(this, {
      channelName: 'Tag',
      mutator: this.save,
      getId: (item) => item?.id,
    });
    this.service = new TagService(transport);
  }

  init(data: Tag): Tag {
    this.value = data;

    return data;
  }

  get tagName() {
    return this.value.name;
  }

  set id(id: string) {
    this.value.id = id;
  }

  async bootstrap() {}

  async invalidate() {}

  async updateTag() {
    try {
      await this.service.updateTag({
        input: {
          id: this.value.id,
          name: this.value.name,
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

    match(path).with(['name', ...P.array()], () => {
      if (type === 'update') {
        this.updateTag();
      }
    });
  }
}

const defaultValue: Tag = {
  id: crypto.randomUUID(),
  name: '',
  source: DataSource.Na,
  createdAt: '',
  appSource: '',
  updatedAt: '',
  metadata: {
    id: crypto.randomUUID(),
    source: DataSource.Openline,
    sourceOfTruth: DataSource.Openline,
    appSource: 'organization',
    created: new Date().toISOString(),
    lastUpdated: new Date().toISOString(),
  },
};
