import type { Transport } from '@store/transport';

// import { gql } from 'graphql-request';

import TimelineDocument from './timeline.graphql';
import { TimelineQuery, TimelineQueryVariables } from './timeline.generated';

class TimelineEventsService {
  private static instance: TimelineEventsService | null = null;
  private transport: Transport;

  constructor(transport: Transport) {
    this.transport = transport;
  }

  static getInstance(transport: Transport): TimelineEventsService {
    if (!TimelineEventsService.instance) {
      TimelineEventsService.instance = new TimelineEventsService(transport);
    }

    return TimelineEventsService.instance;
  }

  async getTimeline(payload: TimelineQueryVariables): Promise<TimelineQuery> {
    return this.transport.graphql.request<
      TimelineQuery,
      TimelineQueryVariables
    >(TimelineDocument, payload);
  }
}

export { TimelineEventsService };
