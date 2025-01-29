import { CapabilityType } from '@graphql/types';

import { SendSlackNotification } from './SendSlackNotification';
import { WebsiteTrackerCapability } from './WebsiteTrackerCapability';

export const capabilities: Record<CapabilityType, () => JSX.Element> = {
  [CapabilityType.IdentifyWebVisitor]: WebsiteTrackerCapability,
  [CapabilityType.CreateOrganization]: () => <></>,
  [CapabilityType.AnalyzeWebSessionIntent]: () => <></>,
  [CapabilityType.SendSlackNotification]: () => <></>,
  [CapabilityType.WebVisitorSendSlackNotification]: SendSlackNotification,
};
