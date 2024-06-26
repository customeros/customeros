import { Transport } from '@store/transport';

import GetLogEntryDocument from './logEntry.graphql';
import CreateLogEntryDocument from './createLogEntry.graphql';
import {
  GetLogEntryQuery,
  GetLogEntryQueryVariables,
} from './logEntry.generated';
import {
  CreateLogEntryMutation,
  CreateLogEntryMutationVariables,
} from './createLogEntry.generated';

export class LogEntriesService {
  private static instance: LogEntriesService | null = null;
  private transport: Transport;

  constructor(transport: Transport) {
    this.transport = transport;
  }

  static getInstance(transport: Transport): LogEntriesService {
    if (!LogEntriesService.instance) {
      LogEntriesService.instance = new LogEntriesService(transport);
    }

    return LogEntriesService.instance;
  }

  async createLogEntry(
    payload: CreateLogEntryMutationVariables,
  ): Promise<CreateLogEntryMutation> {
    return this.transport.graphql.request<
      CreateLogEntryMutation,
      CreateLogEntryMutationVariables
    >(CreateLogEntryDocument, payload);
  }

  async getLogEntry(id: string) {
    return this.transport.graphql.request<
      GetLogEntryQuery,
      GetLogEntryQueryVariables
    >(GetLogEntryDocument, { id });
  }
}
