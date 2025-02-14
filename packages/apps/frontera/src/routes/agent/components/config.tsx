import { ReactElement, ComponentType } from 'react';

import { AgentType, CapabilityType, AgentListenerEvent } from '@graphql/types';

import { NewMeetingRecording, NewWebSessionListener } from './Listeners';
import {
  ApplyTag,
  EvaluateCompanyIcpFit,
  DetectSupportWebVisit,
  SendSlackNotificationCapability,
} from './Capabilities';

type ConfigComponent = ComponentType | (() => ReactElement);

type ConfigMap = {
  [K in CapabilityType | AgentListenerEvent]?: ConfigComponent;
};
type GoalMap = {
  [K in AgentType]?: string[];
};

export const configs: ConfigMap = {
  //////////////////
  // CAPABILITIES //
  //////////////////
  [CapabilityType.IcpQualify]: EvaluateCompanyIcpFit,
  [CapabilityType.WebVisitorSendSlackNotification]:
    SendSlackNotificationCapability,
  [CapabilityType.ApplyTagToCompany]: ApplyTag,
  [CapabilityType.DetectSupportWebvisit]: DetectSupportWebVisit,
  //////////////////
  //  LISTENERS   //
  //////////////////
  [AgentListenerEvent.NewWebSession]: NewWebSessionListener,
  [AgentListenerEvent.NewMeetingRecording]: NewMeetingRecording,
};

export const goals: GoalMap = {
  [AgentType.WebVisitIdentifier]: ['Identify companies that visit my website'],
  [AgentType.IcpQualifier]: ['Qualify companies'],
  [AgentType.SupportSpotter]: ['Spot interactions where help might be needed'],
  [AgentType.MeetingKeeper]: ['Capture external meetings'],
};
