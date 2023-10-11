import {
  Action,
  InteractionEvent,
  Issue,
  LogEntry,
  Meeting,
  User,
} from '@graphql/types';

export type InteractionEventWithDate = InteractionEvent & { date: string };
export type LogEntryWithAliases = LogEntry & {
  logEntryStartedAt: string;
  logEntryCreatedBy: User;
};

export type IssueWithAliases = Issue & {
  issueStatus: string;
};

export type TimelineEvent =
  | InteractionEventWithDate
  | Meeting
  | Action
  | IssueWithAliases
  | LogEntryWithAliases;
