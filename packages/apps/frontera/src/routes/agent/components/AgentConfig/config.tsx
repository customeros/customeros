import { FunctionComponent } from 'react';

import { AgentType, CapabilityType, AgentListenerEvent } from '@graphql/types';

import { NewWebSessionListener } from '../Listeners';
import {
  EvaluateCompanyIcpFit,
  SendSlackNotificationCapability,
} from './Capabilities';

export const configs: Record<
  CapabilityType | AgentListenerEvent,
  FunctionComponent | (() => JSX.Element)
> = {
  //////////////////
  // CAPABILITIES //
  //////////////////
  [CapabilityType.IdentifyWebVisitor]: () => <></>, // deprecated, use NewWebSessionListener instead
  [CapabilityType.CreateOrganization]: () => <></>,
  [CapabilityType.AnalyzeWebSessionIntent]: () => <></>,
  [CapabilityType.SendSlackNotification]: () => <></>,
  [CapabilityType.WebVisitorSendSlackNotification]:
    SendSlackNotificationCapability,
  [CapabilityType.GatherCompanyIntelligence]: () => <></>,
  [CapabilityType.ApplyTagToCompany]: () => <></>,
  [CapabilityType.UpdateCompanyStatus]: () => <></>,
  [CapabilityType.CreateMarkdownTimelineEvent]: () => <></>,
  [CapabilityType.AddMeetingNotesToCompany]: () => <></>,
  [CapabilityType.CreateContacts]: () => <></>,
  [CapabilityType.DetectSupportWebvisit]: () => <></>,
  [CapabilityType.EnrichEmailAddress]: () => <></>,
  [CapabilityType.ExtractMeetingHighlights]: () => <></>,
  [CapabilityType.ExtractSupportSignalsFromMeeting]: () => <></>,
  [CapabilityType.ForwardEmailReply]: () => <></>,
  [CapabilityType.IdentifyMeetingParticipants]: () => <></>,
  [CapabilityType.ManageCampaignExecution]: () => <></>,
  [CapabilityType.ManageEmailDeliveryFailure]: () => <></>,
  [CapabilityType.SelectOptimalSendingMailbox]: () => <></>,
  [CapabilityType.ValidateEmailDeliverability]: () => <></>,
  [CapabilityType.IcpQualify]: EvaluateCompanyIcpFit,
  //////////////////
  //  LISTENERS   //
  //////////////////
  [AgentListenerEvent.NewLead]: () => <></>,
  [AgentListenerEvent.NewWebSession]: NewWebSessionListener,
  [AgentListenerEvent.IcpFit]: () => <></>,
  [AgentListenerEvent.IcpNotAFit]: () => <></>,
  [AgentListenerEvent.RunIcpQualifierAgent]: () => <></>,
  [AgentListenerEvent.WebVisitorIdentified]: () => <></>,
  [AgentListenerEvent.WebVisitorNotIdentified]: () => <></>,
  [AgentListenerEvent.CompanyIdentified]: () => <></>,
  [AgentListenerEvent.CompanyNeedsHelp]: () => <></>,
  [AgentListenerEvent.ContactAddedToCampaign]: () => <></>,
  [AgentListenerEvent.EmailBounced]: () => <></>,
  [AgentListenerEvent.EmailReplyReceived]: () => <></>,
  [AgentListenerEvent.NewMeetingRecording]: () => <></>,
};

export const goals: Record<AgentType, string[]> = {
  [AgentType.WebVisitIdentifier]: ['Identify companies that visit my website'],
  [AgentType.IcpQualifier]: ['Qualify companies'],
  [AgentType.CampaignManager]: ['Someone replies to an email'],
  [AgentType.SupportSpotter]: ['Spot interactions where help might be needed'],
  [AgentType.MeetingKeeper]: [''],
};
