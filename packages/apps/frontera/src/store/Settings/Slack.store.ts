import type { RootStore } from '@store/root';

import { Transport } from '@store/transport';
import { runInAction, makeAutoObservable } from 'mobx';

export class Slack {
  enabled = false;
  isLoading = false;
  error: string | null = null;
  isBootstrapped = false;

  constructor(private root: RootStore, private transportLayer: Transport) {
    makeAutoObservable(this);
  }

  async load() {
    if (this.root.demoMode) {
      this.enabled = mock.slackEnabled;
      this.isBootstrapped = true;

      return;
    }

    try {
      this.isLoading = true;
      const { data } = await this.transportLayer.http.get(
        '/sa/user/settings/slack',
      );
      runInAction(() => {
        this.enabled = data.slackEnabled;
        this.isBootstrapped = true;
      });
    } catch (err) {
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
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
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
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
      runInAction(() => {
        this.error = (err as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  async disableSync() {
    this.isLoading = true;
    this.root.settings.revokeAccess('slack', {
      onSuccess: () => {
        this.enabled = false;
        this.isLoading = false;
        this.root.ui.toastSuccess(
          `We have successfully revoked the access to your Slack account!`,
          'revoke-slack-access',
        );
        this.load();
      },
      onError: (err) => {
        this.error = err.message;
        this.isLoading = false;
        this.root.ui.toastError(
          'An error occurred while revoking access to your Slack account!',
          'revoke-slack-access',
        );
      },
    });
  }
}

const mock = { slackEnabled: true };
