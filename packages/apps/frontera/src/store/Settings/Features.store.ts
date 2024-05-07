import { RootStore } from '@store/root';
import { TransportLayer } from '@store/transport';
import { autorun, makeAutoObservable } from 'mobx';
import { FeatureDefinition } from '@growthbook/growthbook-react';

type Features = Record<string, FeatureDefinition>;

type FeaturesResponse = {
  status: number;
  features: Features;
};

export class FeaturesStore {
  isLoading = false;
  values: Features = {};
  error: string | null = null;
  isBootstrapped = false;

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
        this.transportLayer.isAuthenthicated
      ) {
        this.load();
      }
    });
  }

  async load() {
    try {
      this.isLoading = true;
      const { data } = await this.transportLayer.http.request<FeaturesResponse>(
        {
          method: 'GET',
          url: 'https://cdn.growthbook.io/api/features/sdk-kU7RLceKTmkcTjxO',
          headers: {
            Authorization: undefined,
          },
        },
      );

      this.values = data.features;
    } catch (err) {
      this.error = (err as Error).message;
    } finally {
      this.isLoading = false;
      this.isBootstrapped = true;
    }
  }
}
