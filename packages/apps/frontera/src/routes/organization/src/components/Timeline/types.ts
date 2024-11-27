import { MarkdownEventType } from '@store/TimelineEvents/MarkdownEvent/types';

import {
  User,
  Issue,
  Action,
  Meeting,
  Invoice,
  LogEntry,
  InteractionEvent,
} from '@graphql/types';

export type InteractionEventWithDate = InteractionEvent & { date: string };
export type LogEntryWithAliases = LogEntry & {
  logEntryCreatedBy: User;
  logEntryStartedAt: string;
};

export type IssueWithAliases = Issue & {
  issueStatus: string;
};
export type InvoiceWithId = Invoice & {
  id: string;
};

export type TimelineEvent =
  | InteractionEventWithDate
  | Meeting
  | Action
  | IssueWithAliases
  | MarkdownEventType
  | Pick<InvoiceWithId, 'id' | '__typename'>
  | LogEntryWithAliases;
