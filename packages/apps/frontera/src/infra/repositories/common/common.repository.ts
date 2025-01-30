import { Transport } from '@infra/transport';

import GetSlackChannelsDocument from './queries/getSlackChannels.graphql';
import { GetSlackChannelsQuery } from './queries/getSlackChannels.generated';

export class CommonRepository {
  static instance: CommonRepository | null = null;
  private transport = Transport.getInstance();

  public static getInstance() {
    if (!CommonRepository.instance) {
      CommonRepository.instance = new CommonRepository();
    }

    return CommonRepository.instance;
  }

  async getSlackChannels() {
    return this.transport.graphql.request<GetSlackChannelsQuery>(
      GetSlackChannelsDocument,
    );
  }
}
