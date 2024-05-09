import type { RootStore } from '@store/root';

import { makeAutoObservable } from 'mobx';
import { TransportLayer } from '@store/transport';

export class Slack {
  enabled = false;
  isLoading = false;
  error: string | null = null;
  isBootstrapped = false;

  constructor(
    private rootStore: RootStore,
    private transportLayer: TransportLayer,
  ) {
    makeAutoObservable(this);
  }

  async load() {
    try {
      this.isLoading = true;
      const { data } = await this.transportLayer.http.get(
        '/sa/user/settings/slack',
      );
      this.enabled = data.slackEnabled;
      this.isBootstrapped = true;
    } catch (err) {
      this.error = (err as Error)?.message;
    } finally {
      this.isLoading = false;
    }
  }

  async oauthCallback(code: string) {
    try {
      this.isLoading = true;
      await this.transportLayer.http.post(
        `/ua/slack/oauth/callback?code=${code}`,
      );
      this.load();
    } catch (err) {
      this.error = (err as Error)?.message;
    } finally {
      this.isLoading = false;
    }
  }

  async enableSync() {
    try {
      this.isLoading = true;
      const { data } = await this.transportLayer.http.get(
        `/ua/slack/requestAccess`,
      );
      window.location.href = data.url;
    } catch (err) {
      this.error = (err as Error)?.message;
    } finally {
      this.isLoading = false;
    }
  }

  async disableSync() {
    this.isLoading = true;
    this.rootStore.settingsStore.revokeAccess('slack', {
      onSuccess: () => {
        this.enabled = false;
        this.isLoading = false;
        this.rootStore.uiStore.toastSuccess(
          `We have successfully revoked the access to your Slack account!`,
          'revoke-slack-access',
        );
        this.load();
      },
      onError: (err) => {
        this.error = err.message;
        this.isLoading = false;
        this.rootStore.uiStore.toastError(
          'An error occurred while revoking access to your Slack account!',
          'revoke-slack-access',
        );
      },
    });
  }
}
