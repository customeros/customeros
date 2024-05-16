import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { runInAction, makeAutoObservable } from 'mobx';
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

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
  }

  async load() {
    try {
      this.isLoading = true;
      const { data } = await this.transport.http.request<FeaturesResponse>({
        method: 'GET',
        url: 'https://cdn.growthbook.io/api/features/sdk-kU7RLceKTmkcTjxO',
        headers: {
          Authorization: undefined,
        },
      });

      runInAction(() => {
        this.values = data.features;
        this.isBootstrapped = true;
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }
}
