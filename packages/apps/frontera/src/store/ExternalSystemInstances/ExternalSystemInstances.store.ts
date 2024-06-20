import { gql } from 'graphql-request';
import { RootStore } from '@store/root.ts';
import { Transport } from '@store/transport.ts';
import { runInAction, makeAutoObservable } from 'mobx';

import { ExternalSystemInstance } from '@graphql/types';

import mock from './mock.json';

export class ExternalSystemInstancesStore {
  isLoading = false;
  error: string | null = null;
  isBootstrapped: boolean = false;
  value: Array<ExternalSystemInstance> = [];

  totalElements = 0;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
  }

  async bootstrap() {
    if (this.root.demoMode) {
      this.value = mock.data
        .externalSystemInstances as unknown as ExternalSystemInstance[];
      this.isBootstrapped = true;

      return;
    }

    if (this.isBootstrapped || this.isLoading) return;

    try {
      this.isLoading = true;
      const { externalSystemInstances } =
        await this.transport.graphql.request<EXTERNAL_SYSTEM_INSTANCES_QUERY_RESPONSE>(
          EXTERNAL_INSTANCES_QUERY,
        );

      runInAction(() => {
        this.value = externalSystemInstances;
        this.isBootstrapped = true;
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }
}

type EXTERNAL_SYSTEM_INSTANCES_QUERY_RESPONSE = {
  externalSystemInstances: ExternalSystemInstance[];
};
const EXTERNAL_INSTANCES_QUERY = gql`
  query getExternalSystemInstances {
    externalSystemInstances {
      type
      stripeDetails {
        paymentMethodTypes
      }
    }
  }
`;
