import {
  User,
  Issue,
  Action,
  Meeting,
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

export type TimelineEvent =
  | InteractionEventWithDate
  | Meeting
  | Action
  | IssueWithAliases
  | LogEntryWithAliases;
