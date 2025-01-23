import { CapabilityType } from '@graphql/types';

import { WebsiteTrackerCapability } from './WebsiteTrackerCapability';

export const capabilities: Record<CapabilityType, () => JSX.Element> = {
  [CapabilityType.WebsiteTracker]: WebsiteTrackerCapability,
  [CapabilityType.SendSlackNotification]: () => (
    <div>Send Slack Notification</div>
  ),
};
