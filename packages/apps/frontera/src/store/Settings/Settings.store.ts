import type { RootStore } from '@store/root';

import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { BankAccountsStore } from '@store/BankAccounts/BankAccounts.store.ts';
import { TenantBillingProfilesStore } from '@store/TenantBillingProfiles/TenantBillingProfiles.store.ts';

import { Slack } from './Slack.store';
import { Google } from './Google.store';
import { TenantStore } from './Tenant.store';
import { FeaturesStore } from './Features.store';
import { IntegrationsStore } from './Integrations.store';

export interface OAuthToken {
  scope: string;
  expiresAt: Date;
  idToken: string;
  accessToken: string;
  refreshToken: string;
  providerAccountId: string;
}

export class SettingsStore {
  slack: Slack;
  google: Google;
  tenant: TenantStore;
  features: FeaturesStore;
  integrations: IntegrationsStore;
  bankAccounts: BankAccountsStore;
  tenantBillingProfiles: TenantBillingProfilesStore;
  isLoading = false;
  error: string | null = null;

  constructor(private root: RootStore, private transport: Transport) {
    this.slack = new Slack(this.root, this.transport);
    this.google = new Google(this.root, this.transport);
    this.features = new FeaturesStore(this.root, this.transport);
    this.tenant = new TenantStore(this.root, this.transport);
    this.integrations = new IntegrationsStore(this.root, this.transport);
    this.bankAccounts = new BankAccountsStore(this.root, this.transport);
    this.tenantBillingProfiles = new TenantBillingProfilesStore(
      this.root,
      this.transport,
    );
    makeAutoObservable(this);
  }

  get isBootstrapping() {
    return (
      this.slack.isLoading ||
      this.google.isLoading ||
      this.tenant.isLoading ||
      this.features.isLoading ||
      this.bankAccounts.isLoading ||
      this.tenantBillingProfiles.isLoading ||
      this.integrations.isLoading
    );
  }
  get bootstrapError() {
    return (
      this.slack.error ||
      this.google.error ||
      this.tenant.error ||
      this.features.error ||
      this.tenantBillingProfiles.error ||
      this.integrations.error
    );
  }
  get isBootstrapped() {
    return (
      this.slack.isBootstrapped &&
      this.google.isBootstrapped &&
      this.tenant.isBootstrapped &&
      this.features.isBootstrapped &&
      this.bankAccounts.isBootstrapped &&
      this.tenantBillingProfiles.isBootstrapped &&
      this.integrations.isBootstrapped
    );
  }

  async bootstrap() {
    if (this.isBootstrapped) return;

    await Promise.all([
      await this.slack.load(),
      await this.google.load(),
      await this.tenant.bootstrap(),
      await this.features.load(),
      await this.bankAccounts.bootstrap(),
      await this.tenantBillingProfiles.bootstrap(),
      await this.integrations.load(),
    ]);
  }

  async revokeAccess(
    payload: unknown,
    options?: {
      onError?: (err: Error) => void;
      onSuccess?: (res: unknown) => void;
    },
  ) {
    try {
      this.isLoading = true;

      const res = this.transport.http.post('/ua/revoke', payload);
      options?.onSuccess?.(res);
    } catch (err) {
      this.error = (err as Error)?.message;
      options?.onError?.(err as Error);
    } finally {
      this.isLoading = false;
    }
  }
}
