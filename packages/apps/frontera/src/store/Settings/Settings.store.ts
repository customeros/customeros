import type { RootStore } from '@store/root';

import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';
import { OauthTokenStore } from '@store/Settings/OauthTokenStore.store';

import { Slack } from './Slack.store';
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
  oauthToken: OauthTokenStore;
  tenant: TenantStore;
  features: FeaturesStore;
  integrations: IntegrationsStore;
  isLoading = false;
  error: string | null = null;

  constructor(private root: RootStore, private transport: Transport) {
    this.slack = new Slack(this.root, this.transport);
    this.oauthToken = new OauthTokenStore(this.root, this.transport);
    this.features = new FeaturesStore(this.root, this.transport);
    this.tenant = new TenantStore(this.root, this.transport);
    this.integrations = new IntegrationsStore(this.root, this.transport);
    makeAutoObservable(this);
  }

  get isBootstrapping() {
    return (
      this.slack.isLoading ||
      this.oauthToken.isLoading ||
      this.tenant.isLoading ||
      this.features.isLoading ||
      this.integrations.isLoading
    );
  }
  get bootstrapError() {
    return (
      this.slack.error ||
      this.oauthToken.error ||
      this.tenant.error ||
      this.features.error ||
      this.integrations.error
    );
  }
  get isBootstrapped() {
    return (
      this.slack.isBootstrapped &&
      this.oauthToken.isBootstrapped &&
      this.tenant.isBootstrapped &&
      this.features.isBootstrapped &&
      this.integrations.isBootstrapped
    );
  }

  async bootstrap() {
    if (this.isBootstrapped) return;

    await Promise.all([
      await this.slack.load(),
      await this.oauthToken.load(),
      await this.tenant.bootstrap(),
      await this.features.load(),
      await this.integrations.load(),
    ]);
  }

  async updateUser(
    payload: unknown,
    options?: {
      onError?: (err: Error) => void;
      onSuccess?: (res: unknown) => void;
    },
  ) {
    try {
      this.isLoading = true;

      const res = this.transport.http.post('/ua/updateUser', payload);
      options?.onSuccess?.(res);
    } catch (err) {
      this.error = (err as Error)?.message;
      options?.onError?.(err as Error);
    } finally {
      this.isLoading = false;
    }
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
