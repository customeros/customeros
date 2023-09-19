import {
  Action,
  InteractionEvent,
  LogEntry,
  Meeting,
  User,
} from '@graphql/types';

export type InteractionEventWithDate = InteractionEvent & { date: string };
export type LogEntryWithAliases = LogEntry & {
  logEntryStartedAt: string;
  logEntryCreatedBy: User;
};

export type TimelineEvent =
  | InteractionEventWithDate
  | Meeting
  | Action
  | LogEntryWithAliases;
