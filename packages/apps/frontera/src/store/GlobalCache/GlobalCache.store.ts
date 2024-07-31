import { gql } from 'graphql-request';
import { RootStore } from '@store/root';
import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';

import { GlobalCache } from '@graphql/types';

import mock from './mock.json';

export class GlobalCacheStore {
  value: GlobalCache | null = null;
  isLoading = false;
  isBootstrapped = false;
  error: string | null = null;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
  }

  async bootstrap() {
    if (this.isBootstrapped || this.isLoading) return;

    await this.load();
  }

  async load() {
    if (this.root.demoMode) {
      this.value = mock.data.global_Cache as unknown as GlobalCache;
      this.isBootstrapped = true;

      return;
    }

    try {
      this.isLoading = true;

      const response =
        await this.transport.graphql.request<GLOBAL_CACHE_QUERY_RESULT>(
          GLOBAL_CACHE_QUERY,
        );

      this.value = response.global_Cache;
      this.isBootstrapped = true;
    } catch (error) {
      this.error = (error as Error)?.message;
    } finally {
      this.isLoading = false;
    }
  }
}

type GLOBAL_CACHE_QUERY_RESULT = { global_Cache: GlobalCache };
const GLOBAL_CACHE_QUERY = gql`
  query global_Cache {
    global_Cache {
      cdnLogoUrl
      user {
        id
        emails {
          email
          rawEmail
          primary
        }
        firstName
        lastName
      }
      inactiveEmailTokens {
        email
        provider
      }
      activeEmailTokens {
        email
        provider
      }
      isOwner
      gCliCache {
        id
        type
        display
        data {
          key
          value
          display
        }
      }
      minARRForecastValue
      maxARRForecastValue
      contractsExist
    }
  }
`;
