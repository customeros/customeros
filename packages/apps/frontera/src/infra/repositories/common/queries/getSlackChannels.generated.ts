import * as Types from '../../../../routes/src/types/__generated__/graphql.types';

export type GetSlackChannelsQueryVariables = Types.Exact<{
  [key: string]: never;
}>;

export type GetSlackChannelsQuery = {
  __typename?: 'Query';
  slackChannelsWithBot: Array<{
    __typename?: 'AgentSlackChannel';
    channelId: string;
    name: string;
  }>;
};
