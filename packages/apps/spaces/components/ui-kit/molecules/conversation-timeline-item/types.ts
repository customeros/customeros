export type FeedItem = {
  id: string;
  initiatorFirstName: string;
  initiatorLastName: string;
  initiatorUsername: Participant;
  initiatorType: string;
  lastSenderFirstName: string;
  lastSenderLastName: string;
  lastContentPreview: string;
  lastTimestamp: Time;
};
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
  destination: string[];
  subject?: string;
  replyTo?: string;
};

export interface Content {
  type?: string;
  mimetype: string;
  body: string;
}
