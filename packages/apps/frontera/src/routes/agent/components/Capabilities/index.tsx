import { CapabilityType } from '@graphql/types';

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
  [CapabilityType.IcpQualify]: () => <></>,
  [CapabilityType.CreateMarkdownTimelineEvent]: () => <></>,
  [CapabilityType.ApplyTag]: () => <></>,
};
