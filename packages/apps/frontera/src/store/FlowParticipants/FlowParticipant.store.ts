import type { RootStore } from '@store/root';

import { Channel } from 'phoenix';
import { Operation } from '@store/types';
import { Transport } from '@infra/transport';
import { Store, makeAutoSyncable } from '@store/store';
import { runInAction, makeAutoObservable } from 'mobx';
import { makeAutoSyncableGroup } from '@store/group-store';

import {
  DataSource,
  FlowEntityType,
  FlowParticipant,
  FlowParticipantStatus,
} from '@graphql/types';

import { FlowParticipantsService } from './__service__';

export class FlowParticipantStore implements Store<FlowParticipant> {
  value: FlowParticipant = getDefaultValue();
  version = 0;
  isLoading = false;
  history: Operation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  sync = makeAutoSyncableGroup.sync;
  load = makeAutoSyncable.load<FlowParticipant>();
  update = makeAutoSyncable.update<FlowParticipant>();
  private service: FlowParticipantsService;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, {
      channelName: 'FlowParticipant',
      getId: (d: FlowParticipant) => d?.metadata?.id,
    });
    makeAutoObservable(this);

    this.service = FlowParticipantsService.getInstance(transport);
  }

  get id() {
    return this.value.metadata?.id;
  }

  get contactId() {
    return this.value.entityId;
  }

  get contact() {
    return this.root.contacts.value.get(this.value.entityId);
  }

  setId(id: string) {
    this.value.metadata.id = id;
  }

  async invalidate() {
    try {
      const { flowParticipant } = await this.service.getFlowParticipant({
        id: this.id,
      });

      runInAction(() => {
        this.value = flowParticipant as FlowParticipant;
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    }
  }
}

const getDefaultValue = (): FlowParticipant => ({
  entityId: '',
  entityType: FlowEntityType.Contact,
  executions: [],
  requirementsUnmeet: [],
  metadata: {
    source: DataSource.Openline,
    appSource: DataSource.Openline,
    id: crypto.randomUUID(),
    created: new Date().toISOString(),
    lastUpdated: new Date().toISOString(),
    sourceOfTruth: DataSource.Openline,
  },
  status: FlowParticipantStatus.Scheduled,
});
