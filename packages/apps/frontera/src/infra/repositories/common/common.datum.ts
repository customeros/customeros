import { GetSlackChannelsQuery } from './queries/getSlackChannels.generated';

export type SlackChannelDatum =
  GetSlackChannelsQuery['slackChannelsWithBot'][0];
