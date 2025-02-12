import { ReactElement, FunctionComponent } from 'react';

import { AgentType, CapabilityType, AgentListenerEvent } from '@graphql/types';

import { NewWebSessionListener } from './Listeners';
import {
  EvaluateCompanyIcpFit,
  SendSlackNotificationCapability,
} from './Capabilities';

type ConfigComponent = FunctionComponent | (() => ReactElement);

export const configs: Record<
  CapabilityType | AgentListenerEvent,
  ConfigComponent
> = {
  //////////////////
  // CAPABILITIES //
  //////////////////
  [CapabilityType.AddMeetingNotesToCompany]: () => <></>,
  [CapabilityType.AnalyzeWebSessionIntent]: () => <></>,
  [CapabilityType.ApplyTagToCompany]: () => <></>,
  [CapabilityType.CreateContacts]: () => <></>,
  [CapabilityType.CreateMarkdownTimelineEvent]: () => <></>,
  [CapabilityType.CreateOrganization]: () => <></>,
  [CapabilityType.DetectSupportWebvisit]: () => <></>,
  [CapabilityType.EnrichEmailAddress]: () => <></>,
  [CapabilityType.ExtractMeetingHighlights]: () => <></>,
  [CapabilityType.ExtractSupportSignalsFromMeeting]: () => <></>,
  [CapabilityType.ForwardEmailReply]: () => <></>,
  [CapabilityType.GatherCompanyIntelligence]: () => <></>,
  [CapabilityType.GenerateInvoice]: () => <></>,
  [CapabilityType.IcpQualify]: EvaluateCompanyIcpFit,
  [CapabilityType.IdentifyMeetingParticipants]: () => <></>,
  [CapabilityType.IdentifyWebVisitor]: () => <></>, // deprecated, use NewWebSessionListener instead
  [CapabilityType.ManageCampaignExecution]: () => <></>,
  [CapabilityType.ManageEmailDeliveryFailure]: () => <></>,
  [CapabilityType.SelectOptimalSendingMailbox]: () => <></>,
  [CapabilityType.SendInvoiceViaEmail]: () => <></>,
  [CapabilityType.SendSlackNotification]: () => <></>,
  [CapabilityType.UpdateCompanyStatus]: () => <></>,
  [CapabilityType.ValidateEmailDeliverability]: () => <></>,
  [CapabilityType.WebVisitorSendSlackNotification]:
    SendSlackNotificationCapability,

  //////////////////
  //  LISTENERS   //
  //////////////////
  [AgentListenerEvent.CompanyIdentified]: () => <></>,
  [AgentListenerEvent.CompanyNeedsHelp]: () => <></>,
  [AgentListenerEvent.NewLead]: () => <></>,
  [AgentListenerEvent.NewWebSession]: NewWebSessionListener,
  [AgentListenerEvent.IcpFit]: () => <></>,
  [AgentListenerEvent.IcpNotAFit]: () => <></>,
  [AgentListenerEvent.RunIcpQualifierAgent]: () => <></>,
  [AgentListenerEvent.WebVisitorIdentified]: () => <></>,
  [AgentListenerEvent.WebVisitorNotIdentified]: () => <></>,
  [AgentListenerEvent.StartInvoiceRun]: () => <></>,
  [AgentListenerEvent.StartInvoiceRunWithAutopayment]: () => <></>,

  [AgentListenerEvent.ContactAddedToCampaign]: () => <></>,
  [AgentListenerEvent.EmailBounced]: () => <></>,
  [AgentListenerEvent.EmailReplyReceived]: () => <></>,
  [AgentListenerEvent.NewMeetingRecording]: () => <></>,
};

export const goals: Record<AgentType, string[]> = {
  [AgentType.WebVisitIdentifier]: ['Identify companies that visit my website'],
  [AgentType.IcpQualifier]: ['Qualify companies'],
  [AgentType.SupportSpotter]: ['Spot interactions where help might be needed'],
  [AgentType.CampaignManager]: [''],
  [AgentType.CashflowGuardian]: [''],
  [AgentType.MeetingKeeper]: [''],
};
