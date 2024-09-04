import type { RootStore } from '@store/root';

import merge from 'lodash/merge';
import { Channel } from 'phoenix';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';
import { makeAutoSyncableGroup } from '@store/group-store';

import { DataSource, FlowSequence } from '@graphql/types';

import { FlowSequenceService } from './FlowSequence.service.ts';

export class FlowSequenceStore implements Store<FlowSequence> {
  value: FlowSequence = getDefaultValue();
  version = 0;
  isLoading = false;
  history: Operation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncable.load<FlowSequence>();
  update = makeAutoSyncable.update<FlowSequence>();
  private service: FlowSequenceService;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncable(this, {
      channelName: 'FlowSequence',
      mutator: this.save,
      getId: (d) => d?.metadata?.id,
    });
    this.service = FlowSequenceService.getInstance(transport);
  }

  get id() {
    return this.value.metadata?.id;
  }

  setId(id: string) {
    this.value.metadata.id = id;
  }

  private async save() {
    // todo
  }

  invalidate() {
    // todo
    return Promise.resolve();
  }

  init(data: FlowSequence) {
    return merge(this.value, data);
  }
}

const getDefaultValue = (): FlowSequence => ({
  contacts: [],
  description: '',
  flow: [],
  mailboxes: [],
  name: '',
  status: undefined,
  steps: [],
  metadata: {
    source: DataSource.Openline,
    appSource: DataSource.Openline,
    id: crypto.randomUUID(),
    created: new Date().toISOString(),
    lastUpdated: new Date().toISOString(),
    sourceOfTruth: DataSource.Openline,
  },
});
