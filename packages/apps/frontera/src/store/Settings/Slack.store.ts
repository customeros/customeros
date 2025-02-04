import type { RootStore } from '@store/root';

import { Tracer } from '@infra/tracer';
import { Transport } from '@infra/transport';
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
    const span = Tracer.span('Slack.load');

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
      span.end();
    }
  }

  async oauthCallback(code: string, redirect_uri?: string) {
    try {
      this.isLoading = true;
      await this.transportLayer.http.post(
        `/sa/slack/oauth/callback?code=${code}&redirect_url=${redirect_uri}`,
      );
      this.load();
      this.root.common.fetchSlackChannels();
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

  async enableSync(redirect_uri?: string, state?: string) {
    Tracer.span('Slack.enableSync');

    try {
      this.isLoading = true;

      const redirectUri = redirect_uri ?? 'https://app.customeros.ai/settings';

      const { data } = await this.transportLayer.http.get(
        `/sa/slack/requestAccess?redirect_url=${redirectUri}&state=${state}`,
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
    const span = Tracer.span('Slack.disableSync');

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
        span.end();
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
