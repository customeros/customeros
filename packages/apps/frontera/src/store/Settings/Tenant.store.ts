import { set } from 'lodash';
import pick from 'lodash/pick';
import { RootStore } from '@store/root';
import { Transport } from '@infra/transport';
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
    if (this.isBootstrapped || this.isLoading) return;

    this.load();
  }

  async load() {
    try {
      runInAction(() => {
        this.isLoading = true;
      });

      const { tenantSettings } = await this.service.getTenantSettings();

      runInAction(() => {
        this.value = tenantSettings;
        set(
          this.root.session.value.profile,
          'workspaceName',
          this.value.workspaceName ?? undefined,
        );
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

  updateBillingStatus(newStatus: boolean) {
    if (!this.value) {
      console.error('TenantStore - value is not defined');

      return;
    }
    this.value.billingEnabled = newStatus;
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
