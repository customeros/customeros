import * as Types from '../../../routes/src/types/__generated__/graphql.types';

export type SlackChannelsQueryVariables = Types.Exact<{
  pagination?: Types.InputMaybe<Types.Pagination>;
}>;

export type SlackChannelsQuery = {
  __typename?: 'Query';
  slack_Channels: {
    __typename?: 'SlackChannelPage';
    totalElements: any;
    content: Array<{
      __typename?: 'SlackChannel';
      channelId: string;
      channelName: string;
      metadata: {
        __typename?: 'Metadata';
        id: string;
        appSource: string;
        source: Types.DataSource;
        sourceOfTruth: Types.DataSource;
      };
      organization?: {
        __typename?: 'Organization';
        metadata: { __typename?: 'Metadata'; id: string };
      } | null;
    }>;
  };
};
