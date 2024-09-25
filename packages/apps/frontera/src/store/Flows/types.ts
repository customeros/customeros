export enum FlowActionType {
  EMAIL_NEW = 'EMAIL_NEW',
  EMAIL_REPLY = 'EMAIL_REPLY',
  FLOW_START = 'FLOW_START',
  FLOW_END = 'FLOW_END',
}

export enum FlowNodeType {
  Trigger = 'trigger',
  Action = 'action',
  Control = 'control',
  Wait = 'wait',
}
export interface StartNodeData {
  entity: 'CONTACT';
  action: FlowActionType.FLOW_START;
}
export interface EmailNodeData {
  subject: string;
  waitBefore: number;
  bodyTemplate: string;
  action: FlowActionType.EMAIL_NEW;
}
export interface EmailReplyNodeData {
  subject: string;
  waitBefore: number;
  bodyTemplate: string;
  action: FlowActionType.EMAIL_REPLY;
}
