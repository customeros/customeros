import type { RootStore } from '@store/root';

import merge from 'lodash/merge';
import { Channel } from 'phoenix';
import { P, match } from 'ts-pattern';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';
import { makeAutoSyncableGroup } from '@store/group-store';
import { FlowSequenceService } from '@store/Sequences/__service__';

import {
  DataSource,
  FlowStatus,
  FlowSequence,
  FlowSequenceStatus,
} from '@graphql/types';

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

  private async save(operation: Operation) {
    const diff = operation.diff?.[0];
    const path = diff?.path;

    match(path)
      .with(['status', ...P.array()], () => {
        this.service.updateSequenceStatus({
          id: this.id,
          stage: this.value.status as FlowSequenceStatus,
        });
      })
      .with(['name', ...P.array()], () => {
        this.service.updateSequence({
          input: {
            id: this.id,
            name: this.value.name,
            description: this.value.description,
          },
        });
      });
  }

  linkContact(contactId: string, emailId: string) {
    return this.service.linkContact({
      sequenceId: this.id,
      contactId,
      emailId,
    });
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
  flow: {
    metadata: {
      source: DataSource.Openline,
      appSource: DataSource.Openline,
      id: crypto.randomUUID(),
      created: new Date().toISOString(),
      lastUpdated: new Date().toISOString(),
      sourceOfTruth: DataSource.Openline,
    },
    description: '',
    name: '',
    status: FlowStatus.Inactive,
    sequences: [],
  },
  senders: [],
  name: '',
  status: FlowSequenceStatus.Inactive,
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
