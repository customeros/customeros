import { GetSlackChannelsQuery } from './queries/getSlackChannels.generated';
import { TenantImpersonateListQuery } from '../common/queries/impersonateList.generated.ts';

export type SlackChannelDatum =
  GetSlackChannelsQuery['slackChannelsWithBot'][0];

export type TenantImpersonateDatum = NonNullable<
  TenantImpersonateListQuery['tenant_impersonateList'][0]
>;

export type BrowserAutomationConfigResponse = {
  success: boolean;
  data: {
    id: number;
    userId: string;
    tenant: string;
    cookies: string;
    userAgent: string;
    createdAt: string;
    updatedAt: string;
    sessionStatus: string;
  };
};
