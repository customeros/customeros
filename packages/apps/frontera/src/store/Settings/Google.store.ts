import type { RootStore } from '@store/root';

import { Transport } from '@store/transport';
import { runInAction, makeAutoObservable } from 'mobx';

type GoogleToken = {
  email: string;
  userId: string;
  needsManualRefresh: boolean;
};

export class Google {
  tokens: Array<GoogleToken> = [];
  isLoading = false;
  error: string | null = null;
  isBootstrapped = false;

  constructor(private root: RootStore, private transport: Transport) {
    makeAutoObservable(this);
  }

  async load() {
    try {
      this.isLoading = true;
      const { data } = await this.transport.http.get<GoogleToken[]>(
        `/sa/user/settings/google/${this.root.session.value.tenant}`,
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

  async enableSync() {
    try {
      const { data } = await this.transport.http.get<{ url: string }>(
        `/enable/google-sync?origin=${window.location.pathname}${window.location.search}`,
      );

      window.location.href = data.url;
    } catch (err) {
      console.error(err);
    }
  }

  async updateUser(email: string, userId: string) {
    this.isLoading = true;

    this.root.settings.updateUser(
      {
        tenant: this.root.session.value.tenant,
        email: email,
        userId: userId,
      },
      {
        onSuccess: this.onUserChangeSuccess.bind(this),
        onError: this.onUserChangeError.bind(this),
      },
    );
  }

  async disableSync(email: string) {
    this.isLoading = true;

    this.root.settings.revokeAccess(
      {
        tenant: this.root.session.value.tenant,
        provider: 'google',
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
