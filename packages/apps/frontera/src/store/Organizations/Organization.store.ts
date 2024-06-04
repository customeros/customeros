import type { RootStore } from '@store/root';

import omit from 'lodash/omit';
import { Channel } from 'phoenix';
import { gql } from 'graphql-request';
import { Operation } from '@store/types';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { Store, makeAutoSyncable } from '@store/store';

import {
  DataSource,
  Organization,
  OrganizationStage,
  OrganizationUpdateInput,
} from '@graphql/types';

export class OrganizationStore implements Store<Organization> {
  value: Organization = defaultValue;
  version = 0;
  isLoading = false;
  history: Operation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  subscribe = makeAutoSyncable.subscribe;
  load = makeAutoSyncable.load<Organization>();
  update = makeAutoSyncable.update<Organization>();

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoSyncable(this, { channelName: 'Organization', mutator: this.save });
    makeAutoObservable(this);
  }

  async invalidate() {}

  private async save() {
    const payload: PAYLOAD = {
      input: {
        ...omit(this.value, 'metadata', 'owner'),
        id: this.value.metadata?.id,
      },
    };
    try {
      this.isLoading = true;
      await this.transport.graphql.request(UPDATE_ORGANIZATION_DEF, payload);
    } catch (e) {
      this.error = (e as Error)?.message;
    } finally {
      this.isLoading = false;
    }
  }

  public async updateStage(stage: OrganizationStage) {
    this.value.stage = stage;
    await this.save();
  }

  set id(id: string) {
    this.value.metadata.id = id;
  }
}

type PAYLOAD = { input: OrganizationUpdateInput };
const UPDATE_ORGANIZATION_DEF = gql`
  mutation updateOrganization($input: OrganizationUpdateInput!) {
    organization_Update(input: $input) {
      id
    }
  }
`;

const defaultValue: Organization = {
  name: 'Unnamed',
  metadata: {
    id: '1',
    created: new Date().toISOString(),
    lastUpdated: new Date().toISOString(),
    appSource: DataSource.Openline,
    source: DataSource.Openline,
    sourceOfTruth: DataSource.Openline,
  },
  owner: null,
  stage: OrganizationStage.Target,
} as Organization;
