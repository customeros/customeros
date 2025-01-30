import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { computed, observable, runInAction } from 'mobx';
import {
  CommonRepository,
  SlackChannelDatum,
} from '@infra/repositories/common';

import { unwrap } from '@utils/unwrap';

export class CommonStore {
  private repository = CommonRepository.getInstance();
  @observable accessor slackChannels: Map<string, SlackChannelDatum> =
    new Map();

  constructor(private root: RootStore) {}

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

  async bootstrap() {
    await this.fetchSlackChannels();
  }
}
