import type { RootStore } from '@store/root';

import { Transport } from '@store/transport';
import { runInAction, makeAutoObservable } from 'mobx';

export type OauthToken = {
  type: string;
  email: string;

  provider: string;
  needsManualRefresh: boolean;
};

export class OauthTokenStore {
  tokens: Array<OauthToken> = [];
  isLoading = false;
  error: string | null = null;
  isBootstrapped = false;

  constructor(private root: RootStore, private transport: Transport) {
    makeAutoObservable(this);
  }

  async load() {
    if (this.root.demoMode) {
      return;
    }

    try {
      this.isLoading = true;

      const { data } = await this.transport.http.get<OauthToken[]>(
        `/sa/user/settings/oauth/${this.root.session.value.tenant}`,
      );

      runInAction(() => {
        this.tokens = data;
        this.isBootstrapped = true;
      });
    } catch (error) {
      runInAction(() => {
        this.error = (error as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  async enableSync(tokenType: string, provider: string) {
    if (this.root.demoMode) {
      return;
    }

    try {
      const { data } = await this.transport.http.get<{ url: string }>(
        `/enable/${provider}-sync?origin=${window.location.pathname}${window.location.search}&type=${tokenType}`,
      );

      window.location.href = data.url;
    } catch (err) {
      console.error(err);
    }
  }

  async disableSync(email: string, provider: string) {
    if (this.root.demoMode) {
      return;
    }
    this.isLoading = true;

    this.root.settings.revokeAccess(
      {
        tenant: this.root.session.value.tenant,
        provider: provider,
        email: email,
      },
      {
        onSuccess: () => this.onDisableSuccess(email),
        onError: (err) => this.onDisableError(err, provider),
      },
    );
  }

  private onDisableSuccess(email: string) {
    this.isLoading = false;
    this.root.ui.toastSuccess(
      `We have unlinked ${email}`,
      'disable-google-sync',
    );
    setTimeout(() => this.load(), 500);
  }

  private onDisableError(err: Error, provider: string) {
    const providerLabel = provider === 'google' ? 'Google' : 'Microsoft 365';

    this.error = err.message;
    this.isLoading = false;
    this.root.ui.toastError(
      `An error occurred while disabling the ${providerLabel} sync!`,
      'disable-google-sync',
    );
  }
}
