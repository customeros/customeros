import { CapabilityType } from '@graphql/types';

import { WebsiteTrackerCapability } from './WebsiteTrackerCapability';

export const capabilities: Record<CapabilityType, () => JSX.Element> = {
  [CapabilityType.IdentifyWebVisitor]: () => <WebsiteTrackerCapability />,
  [CapabilityType.SendSlackNotification]: () => (
    <div>Send Slack Notification</div>
  ),
  [CapabilityType.AnalyzeWebSessionIntent]: () => <div />,
  [CapabilityType.CreateOrganization]: () => <div />,
  [CapabilityType.WebVisitorSendSlackNotification]: () => <div />,
};
