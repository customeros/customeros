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
    this.isLoading = true;

    this.root.settings.revokeAccess(
      {
        tenant: this.root.session.value.tenant,
        provider: provider,
        email: email,
      },
      {
        onSuccess: this.onDisableSuccess.bind(this),
        onError: this.onDisableError.bind(this),
      },
    );
  }

  private onDisableSuccess() {
    //todo
    this.isLoading = false;
    this.root.ui.toastSuccess(
      'We have successfully disabled the google sync!',
      'disable-google-sync',
    );
    setTimeout(() => this.load(), 500);
  }

  private onDisableError(err: Error) {
    this.error = err.message;
    this.isLoading = false;
    this.root.ui.toastError(
      'An error occurred while disabling the google sync!',
      'disable-google-sync',
    );
  }

  private onUserChangeSuccess() {
    this.isLoading = false;
    this.root.ui.toastSuccess(
      'We have successfully changed the user!',
      'change-user-token',
    );
    setTimeout(() => this.load(), 500);
  }

  private onUserChangeError(err: Error) {
    this.error = err.message;
    this.isLoading = false;
    this.root.ui.toastError(
      'An error occurred while changing the owner!',
      'change-user-token',
    );
  }
}
