import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { computed, reaction, observable, runInAction } from 'mobx';
import {
  CommonRepository,
  SlackChannelDatum,
  TenantImpersonateDatum,
} from '@infra/repositories/common';

import { unwrap } from '@utils/unwrap';

export class CommonStore {
  private repository = CommonRepository.getInstance();
  @observable accessor slackChannels: Map<string, SlackChannelDatum> =
    new Map();
  @observable accessor impersonateAccounts: TenantImpersonateDatum[] = [];

  constructor(private root: RootStore) {
    reaction(
      () => this.root?.globalCache?.value?.isPlatformOwner,
      () => {
        if (this.root?.globalCache?.value?.isPlatformOwner) {
          this.fetchImpersonateAccounts.bind(this)();
        }
      },
    );
  }

  @computed
  get slackChannelsArray() {
    return Array.from(this.slackChannels.values());
  }

  async fetchSlackChannels() {
    const span = Tracer.span('CommonStore.fetchSlackChannels');
    const [data, err] = await unwrap(this.repository.getSlackChannels());

    if (err) {
      console.error(
        'CommonStore.fetchSlackChannels',
        'Failed to fetch slack channels',
        err,
      );

      return;
    }

    runInAction(() => {
      data?.slackChannelsWithBot?.forEach((channel) => {
        this.slackChannels.set(channel.channelId, channel);
      });
    });

    span.end();
  }

  async fetchImpersonateAccounts() {
    const span = Tracer.span('CommonStore.fetchImpersonateAccounts');
    const [data, err] = await unwrap(
      this.repository.getTenantImpersonateList(),
    );

    if (err) {
      console.error(
        'CommonStore.fetchImpersonateAccounts',
        'Failed to fetch impersonate accounts',
        err,
      );

      return;
    }

    runInAction(() => {
      this.impersonateAccounts = data?.tenant_impersonateList ?? [];
    });

    span.end();
  }

  async bootstrap() {
    await this.fetchSlackChannels();
  }
}
