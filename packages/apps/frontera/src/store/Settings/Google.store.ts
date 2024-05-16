import type { RootStore } from '@store/root';

import { Transport } from '@store/transport';
import { runInAction, makeAutoObservable } from 'mobx';

type GoogleSettings = {
  gmailSyncEnabled: boolean;
  googleCalendarSyncEnabled: boolean;
};

export class Google {
  gmailEnabled = false;
  calendarEnabled = false;
  isLoading = false;
  error: string | null = null;
  isBootstrapped = false;

  constructor(private root: RootStore, private transport: Transport) {
    makeAutoObservable(this);
  }

  async load() {
    const playerIdentityId = this.root.session.value.profile.id;

    try {
      this.isLoading = true;
      const { data } = await this.transport.http.get<GoogleSettings>(
        `/sa/user/settings/google/${playerIdentityId}`,
      );
      runInAction(() => {
        this.gmailEnabled = data.gmailSyncEnabled;
        this.calendarEnabled = data.googleCalendarSyncEnabled;
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

  async disableSync() {
    this.isLoading = true;

    this.root.settings.revokeAccess(
      {
        provider: 'google',
        providerAccountId: this.root.session.value.profile.id,
      },
      {
        onSuccess: this.onDisableSuccess.bind(this),
        onError: this.onDisableError.bind(this),
      },
    );
  }

  private onDisableSuccess() {
    this.gmailEnabled = false;
    this.calendarEnabled = false;
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
}
