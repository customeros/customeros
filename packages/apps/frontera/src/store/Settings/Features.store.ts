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
    if (this.root.demoMode) {
      this.values = mock.features;
      this.isBootstrapped = true;

      return;
    }

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

const mock = {
  status: 200,
  features: {
    'gp-dedicated-1': {
      defaultValue: true,
      rules: [
        {
          condition: {
            tenant: 'gasposco',
          },
          force: false,
        },
      ],
    },
    'show-parent-relationship-selector': {
      defaultValue: false,
      rules: [
        {
          condition: {
            tenant: {
              $in: [
                'openlineai',
                'PureThePhantom',
                'gasposco',
                'DistinctMiracleman',
              ],
            },
          },
          force: true,
        },
      ],
    },
    'parent-relationship-selector-read-only': {
      defaultValue: false,
      rules: [
        {
          condition: {
            tenant: 'gasposco',
          },
          force: true,
        },
      ],
    },
    'taller-customer-map-chart': {
      defaultValue: false,
      rules: [
        {
          condition: {
            tenant: {
              $in: ['govlycom'],
            },
          },
          force: true,
        },
      ],
    },
    'decrease-circle-top-radius': {
      defaultValue: false,
      rules: [
        {
          condition: {
            tenant: {
              $in: ['gasposco', 'encordcom', 'govlycom'],
            },
          },
          force: true,
        },
      ],
    },
    'settings-master-plans-view': {
      defaultValue: false,
      rules: [
        {
          condition: {
            tenant: {
              $in: ['PureThePhantom', 'openlineai', 'DistinctMiracleman'],
            },
          },
          force: true,
        },
      ],
    },
    'my-views-nav-item': {
      defaultValue: false,
      rules: [
        {
          condition: {
            tenant: {
              $in: [
                'openlineai',
                'PureThePhantom',
                'govlycom',
                'SolidNightveil',
                'DistinctMiracleman',
              ],
            },
          },
          force: true,
        },
      ],
    },
    kmenu: {
      defaultValue: false,
      rules: [
        {
          condition: {
            tenant: {
              $in: ['openlineai', 'PureThePhantom'],
            },
          },
          force: true,
        },
      ],
    },
    'onboarding-plans': {
      defaultValue: false,
      rules: [
        {
          condition: {
            tenant: {
              $in: ['PureThePhantom', 'openlineai', 'DistinctMiracleman'],
            },
          },
          force: true,
        },
      ],
    },
    'org-name-readonly': {
      defaultValue: false,
      rules: [
        {
          condition: {
            tenant: 'gasposco',
          },
          force: true,
        },
      ],
    },
    presence: {
      defaultValue: false,
      rules: [
        {
          condition: {
            tenant: {
              $in: [
                'openlineai',
                'PureThePhantom',
                'DistinctMiracleman',
                'SolidNightveil',
              ],
            },
          },
          force: true,
        },
      ],
    },
    'contract-new': {
      defaultValue: false,
      rules: [
        {
          condition: {
            tenant: {
              $in: [
                'PureThePhantom',
                'openlineai',
                'customerosai',
                'SettledSpiderHam',
              ],
            },
          },
          force: true,
        },
      ],
    },
    'edit-columns': {
      defaultValue: false,
      rules: [
        {
          condition: {
            tenant: {
              $in: [
                'PureThePhantom',
                'customerosai',
                'openlineai',
                'SolidNightveil',
                'DistinctMiracleman',
              ],
            },
          },
          force: true,
        },
      ],
    },
    'invoice-sim': {
      defaultValue: true,
      rules: [
        {
          condition: {
            tenant: {
              $in: [
                'PureThePhantom',
                'openlineai',
                'customerosai',
                'SolidNightveil',
                'DistinctMiracleman',
                'SettledSpiderHam',
              ],
            },
          },
          force: true,
        },
      ],
    },
    prospects: {
      defaultValue: false,
      rules: [
        {
          condition: {
            tenant: {
              $in: [
                'PureThePhantom',
                'openlineai',
                'customerosai',
                'SolidNightveil',
                'DistinctMiracleman',
                'ellipsisdev',
              ],
            },
          },
          force: true,
        },
      ],
    },
    'invoice-simulation': {
      defaultValue: false,
      rules: [
        {
          condition: {
            tenant: {
              $in: [
                'PureThePhantom',
                'openlineai',
                'customerosai',
                'SolidNightveil',
                'DistinctMiracleman',
                'SettledSpiderHam',
              ],
            },
          },
          force: true,
        },
      ],
    },
  },
  dateUpdated: '2024-06-06T15:33:13.096Z',
};
