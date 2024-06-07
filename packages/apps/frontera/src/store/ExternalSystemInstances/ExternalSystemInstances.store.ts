import { gql } from 'graphql-request';
import { RootStore } from '@store/root.ts';
import { Transport } from '@store/transport.ts';
import { runInAction, makeAutoObservable } from 'mobx';

import { ExternalSystemInstance } from '@graphql/types';

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
    try {
      this.isLoading = true;
      const { externalSystemInstances } =
        await this.transport.graphql.request<ExternalSystemInstanceS_QUERY_RESPONSE>(
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

type ExternalSystemInstanceS_QUERY_RESPONSE = {
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
