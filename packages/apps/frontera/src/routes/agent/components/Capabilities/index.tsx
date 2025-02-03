import { CapabilityType } from '@graphql/types';

import { EvaluateCompanyIcpFit } from './EvaluateCompanyIcpFit';
import { WebsiteTrackerCapability } from './WebsiteTrackerCapability';
import { SendSlackNotificationCapability } from './SendSlackNotificationCapability';

export const capabilities: Record<CapabilityType, () => JSX.Element> = {
  [CapabilityType.IdentifyWebVisitor]: () => <WebsiteTrackerCapability />,
  [CapabilityType.CreateOrganization]: () => <></>,
  [CapabilityType.AnalyzeWebSessionIntent]: () => <></>,
  [CapabilityType.SendSlackNotification]: () => <></>,
  [CapabilityType.WebVisitorSendSlackNotification]: () => (
    <SendSlackNotificationCapability />
  ),
  [CapabilityType.ApplyTag]: () => <></>,
  [CapabilityType.CreateMarkdownTimelineEvent]: () => <></>,
  [CapabilityType.IcpQualify]: () => <EvaluateCompanyIcpFit />,
};
