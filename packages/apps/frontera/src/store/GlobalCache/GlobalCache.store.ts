import { gql } from 'graphql-request';
import { RootStore } from '@store/root';
import { TransportLayer } from '@store/transport';
import { autorun, makeAutoObservable } from 'mobx';

import { GlobalCache } from '@graphql/types';

export class GlobalCacheStore {
  value: GlobalCache | null = null;
  isLoading = false;
  isBootstrapped = false;
  error: string | null = null;

  constructor(
    private rootStore: RootStore,
    private transportLayer: TransportLayer,
  ) {
    makeAutoObservable(this);

    autorun(() => {
      const sessionStore = this.rootStore.sessionStore;

      if (
        sessionStore.isHydrated &&
        sessionStore.isAuthenticated &&
        this.transportLayer.isAuthenthicated &&
        sessionStore.isBootstrapped
      ) {
        this.load();
      }
    });
  }

  async load() {
    try {
      this.isLoading = true;
      const response =
        await this.transportLayer.client.request<GLOBAL_CACHE_QUERY_RESULT>(
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
      isGoogleActive
      isGoogleTokenExpired
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
