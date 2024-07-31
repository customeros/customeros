import type { RootStore } from '@store/root';

import { Transport } from '@store/transport';
import { runInAction, makeAutoObservable } from 'mobx';

import { SlackChannel } from '@graphql/types';

import { SettingsService } from './__service__/Settings.service';

export class Slack {
  enabled = false;
  isLoading = false;
  error: string | null = null;
  isBootstrapped = false;
  channels: SlackChannel[] = [];

  private service: SettingsService;

  constructor(private root: RootStore, private transportLayer: Transport) {
    this.service = SettingsService.getInstance(transportLayer);
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
      const { slack_Channels } = await this.service.getSlackChannels({
        pagination: { page: 0, limit: 1000 },
      });

      runInAction(() => {
        this.enabled = data.slackEnabled;
        this.channels = slack_Channels.content as SlackChannel[];
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
          `Access to your Slack account has been revoked`,
          'revoke-slack-access',
        );
        this.load();
      },
      onError: (err) => {
        this.error = err.message;
        this.isLoading = false;
        this.root.ui.toastError(
          'An error occurred while revoking access to your Slack account',
          'revoke-slack-access',
        );
      },
    });
  }
}

const mock = { slackEnabled: true };
