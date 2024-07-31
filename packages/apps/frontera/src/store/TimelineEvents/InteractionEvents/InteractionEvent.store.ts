import { Channel } from 'phoenix';
import { RootStore } from '@store/root';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';
import { AnalysisStore } from '@store/TimelineEvents/Analyses/Analysis.store';

import { DataSource, InteractionEvent } from '@graphql/types';

export class InteractionEventStore implements Store<InteractionEvent> {
  value: InteractionEvent = defaultValue;
  channel?: Channel | undefined;
  error: string | null = null;
  history: Operation[] = [];
  isBootstrapped: boolean = false;
  isLoading: boolean = false;
  version: number = 0;
  load = makeAutoSyncable.load<InteractionEvent>();
  subscribe = makeAutoSyncable.subscribe;
  update = makeAutoSyncable.update<InteractionEvent>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, {
      channelName: 'InteractionEvent',
      mutator: this.save,
      getId: (item) => item?.id,
    });
    makeAutoObservable(this);
  }

  async bootstrap() {}

  async invalidate() {}

  async save() {}

  get id() {
    return this.value.id;
  }

  set id(id: string) {
    this.value.id = id;
  }
}

const defaultValue: InteractionEvent = {
  id: crypto.randomUUID(),
  appSource: DataSource.Openline,
  createdAt: new Date().toISOString(),
  source: DataSource.Openline,
  __typename: 'InteractionEvent',
  content: '',
  externalLinks: [],
  includes: [],
  sentBy: [],
  sentTo: [],
  sourceOfTruth: DataSource.Openline,
  actionItems: [],
  actions: [],
  channel: '',
  channelData: '',
  contentType: '',
  customerOSInternalIdentifier: '',
  eventIdentifier: '',
  eventType: '',
  interactionSession: undefined,
  issue: undefined,
  meeting: undefined,
  repliesTo: undefined,
  summary: AnalysisStore.getDefaultValue(),
};
