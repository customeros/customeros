import { ContactParticipant, UserParticipant } from '@spaces/graphql';

export interface Props {
  feedId: string;
  source: string;
  first: boolean;
  createdAt: string;
  feedInitiator: {
    firstName: string;
    lastName: string;
    phoneNumber: string;
    lastTimestamp: {
      seconds: number;
    };
  };
  participants?: Array<UserParticipant | ContactParticipant>;
}

export type Time = {
  seconds: number;
};

export type Participant = {
  type: number;
  identifier: string;
};

export type ConversationItem = {
  id: string;
  conversationId: string;
  type: number;
  subtype: number;
  content: string;
  direction: number;
  time: Time;
  senderType: number;
  senderId: string;
  senderUsername: Participant;
};

export type SendMailRequest = {
  username: string;
  content: string;
  channel: string;
  direction: string;
  to: string[];
  cc: string[];
  bcc: string[];
  subject?: string;
  replyTo?: string;
};

export interface Content {
  type?: string;
  mimetype: string;
  body: string;
}
