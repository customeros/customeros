import pick from 'lodash/pick';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { runInAction, makeAutoObservable } from 'mobx';

import { TenantSettings, TenantSettingsInput } from '@graphql/types';

import { SettingsService } from './__service__/Settings.service';

// TODO: Refactor this store to use the new syncable store
export class TenantStore {
  value: TenantSettings | null = null;
  isLoading = false;
  isBootstrapped = false;
  error: string | null = null;
  private service: SettingsService;

  constructor(public root: RootStore, public transportLayer: Transport) {
    this.service = SettingsService.getInstance(transportLayer);
    makeAutoObservable(this);
  }

  async bootstrap() {
    if (this.root.demoMode) {
      this.value = mock.data.tenantSettings as TenantSettings;
      this.isBootstrapped = true;

      return;
    }

    if (this.isBootstrapped || this.isLoading) return;

    this.load();
  }

  async load() {
    try {
      this.isLoading = true;

      const { tenantSettings } = await this.service.getTenantSettings();

      runInAction(() => {
        this.value = tenantSettings;
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

  update(
    updated: (value: TenantSettings) => TenantSettings,
    options: { mutate: boolean } = { mutate: true },
  ) {
    this.value = updated(this.value as TenantSettings);

    if (options?.mutate) this.save();
  }

  // Temporary - This whole store needs to be refactored to use the new syncable store
  // at which point this method will be removed
  async saveOpportunityStage(stage: string) {
    try {
      const stageIndex = this.value?.opportunityStages.findIndex(
        (s) => s.value === stage,
      );

      if (!stageIndex) return;

      const payload = pick(
        this.value?.opportunityStages[stageIndex],
        'id',
        'label',
        'visible',
        'likelihoodRate',
      );

      await this.service.updateOpportunityStage({
        input: {
          ...payload,
          id: payload.id as string,
        },
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error).message;
      });
    }
  }

  async save() {
    const { opportunityStages, ...rest } = this.value as TenantSettings;

    try {
      this.isLoading = true;
      await this.service.updateTenantSettings({
        input: {
          ...(rest as TenantSettingsInput),
          patch: true,
        },
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
  data: {
    tenantSettings: {
      logoUrl: '59e1ad09-49fe-40b1-9e9a-e1f94682d12d',
      logoRepositoryFileId: '59e1ad09-49fe-40b1-9e9a-e1f94682d12d',
      baseCurrency: 'USD',
      billingEnabled: true,
      opportunityStages: [
        {
          id: '1',
          value: 'STAGE1',
          order: 1,
          label: 'Identified',
          visible: true,
        },
        {
          id: '2',
          value: 'STAGE2',
          order: 2,
          label: 'Qualified',
          visible: true,
        },
        {
          id: '3',
          value: 'STAGE3',
          order: 3,
          label: 'Committed',
          visible: true,
        },
      ],
    },
  },
};
