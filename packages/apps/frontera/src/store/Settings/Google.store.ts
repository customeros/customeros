import type { RootStore } from '@store/root';

import { makeAutoObservable } from 'mobx';
import { Transport } from '@store/transport';

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
      this.gmailEnabled = data.gmailSyncEnabled;
      this.calendarEnabled = data.googleCalendarSyncEnabled;
      this.isBootstrapped = true;
    } catch (error) {
      this.error = (error as Error)?.message;
    } finally {
      this.isLoading = false;
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
        onSuccess: () => {
          this.gmailEnabled = false;
          this.calendarEnabled = false;
          this.isLoading = false;
          this.root.ui.toastSuccess(
            'We have successfully revoked the access to your google account!',
            'revoke-google-access',
          );
          this.load();
        },
        onError: (err) => {
          this.error = err.message;
          this.isLoading = false;
          this.root.ui.toastError(
            'An error occurred while revoking access to your google account!',
            'revoke-google-access',
          );
        },
      },
    );
  }
}
