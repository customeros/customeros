import { Transport } from '@store/transport';

import SlackChannelsDocument from './slackChannels.graphql';
import {
  SlackChannelsQuery,
  SlackChannelsQueryVariables,
} from './slackChannels.generated';

export class SettingsService {
  private static instance: SettingsService;
  private transport: Transport;

  private constructor(transport: Transport) {
    this.transport = transport;
  }

  static getInstance(transport: Transport) {
    if (!SettingsService.instance) {
      SettingsService.instance = new SettingsService(transport);
    }

    return SettingsService.instance;
  }

  async getSlackChannels(payload: SlackChannelsQueryVariables) {
    return this.transport.graphql.request<
      SlackChannelsQuery,
      SlackChannelsQueryVariables
    >(SlackChannelsDocument, payload);
  }
}
