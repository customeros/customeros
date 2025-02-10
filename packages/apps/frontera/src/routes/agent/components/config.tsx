import { AgentType, CapabilityType, AgentListenerEvent } from '@graphql/types';

import {
  EvaluateCompanyIcpFit,
  WebsiteTrackerCapability,
  SendSlackNotificationCapability,
} from './Capabilities';

export const configs: Record<
  CapabilityType | AgentListenerEvent,
  () => JSX.Element
> = {
  //////////////////
  // CAPABILITIES //
  //////////////////
  [CapabilityType.IdentifyWebVisitor]: () => <WebsiteTrackerCapability />,
  [CapabilityType.CreateOrganization]: () => <></>,
  [CapabilityType.AnalyzeWebSessionIntent]: () => <></>,
  [CapabilityType.SendSlackNotification]: () => <></>,
  [CapabilityType.WebVisitorSendSlackNotification]: () => (
    <SendSlackNotificationCapability />
  ),
  [CapabilityType.ApplyTag]: () => <></>,
  [CapabilityType.GatherCompanyIntelligence]: () => <></>,
  [CapabilityType.UpdateCompanyStatus]: () => <></>,
  [CapabilityType.CreateMarkdownTimelineEvent]: () => <></>,
  [CapabilityType.IcpQualify]: () => <EvaluateCompanyIcpFit />,
  //////////////////
  //  LISTENERS   //
  //////////////////
  [AgentListenerEvent.NewLead]: () => <></>,
  [AgentListenerEvent.NewWebSession]: () => <></>,
  [AgentListenerEvent.IcpFit]: () => <></>,
  [AgentListenerEvent.IcpNotAFit]: () => <></>,
  [AgentListenerEvent.RunIcpQualifierAgent]: () => <></>,
  [AgentListenerEvent.WebVisitorIdentified]: () => <></>,
  [AgentListenerEvent.WebVisitorNotIdentified]: () => <></>,
};

export const goals: Record<AgentType, string[]> = {
  [AgentType.WebVisitIdentifier]: ['Identify companies that visit my website'],
  [AgentType.IcpQualifier]: ['Qualify companies'],
  [AgentType.TagSupport]: ['Spot interactions where help might be needed'],
};
