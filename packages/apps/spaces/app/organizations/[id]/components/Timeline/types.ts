import { Action, InteractionEvent, Meeting } from '@graphql/types';

export type InteractionEventWithDate = InteractionEvent & { date: string };

export type TimelineEvent = InteractionEventWithDate | Meeting | Action;
