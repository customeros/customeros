import { Transport } from '@store/transport';

import {
  CreateLogEntryMutation,
  CreateLogEntryMutationVariables,
} from '@organization/graphql/createLogEntry.generated';
import {
  UpdateLogEntryMutation,
  UpdateLogEntryMutationVariables,
} from '@organization/graphql/updateLogEntry.generated';

import GetLogEntryDocument from './logEntry.graphql';
import UpdateLogEntryDocument from './updateLogEntry.graphql';
import CreateLogEntryDocument from './createLogEntry.graphql';
import AddTagToLogEntryDocument from './addTagToLogEntry.graphql';
import RemoveTagFromLogEntryDocument from './removeTagFromLogEntry.graphql';
import {
  GetLogEntryQuery,
  GetLogEntryQueryVariables,
} from './logEntry.generated';
import {
  AddTagToLogEntryMutation,
  AddTagToLogEntryMutationVariables,
} from './addTagToLogEntry.generated';
import {
  RemoveTagFromLogEntryMutation,
  RemoveTagFromLogEntryMutationVariables,
} from './removeTagFromLogEntry.generated';

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

  async updateLogEntry(
    payload: UpdateLogEntryMutationVariables,
  ): Promise<UpdateLogEntryMutation> {
    return this.transport.graphql.request<
      UpdateLogEntryMutation,
      UpdateLogEntryMutationVariables
    >(UpdateLogEntryDocument, payload);
  }

  async addTagToLogEntry(
    payload: AddTagToLogEntryMutationVariables,
  ): Promise<AddTagToLogEntryMutation> {
    return this.transport.graphql.request<
      AddTagToLogEntryMutation,
      AddTagToLogEntryMutationVariables
    >(AddTagToLogEntryDocument, payload);
  }

  async removeTagFromLogEntry(
    payload: RemoveTagFromLogEntryMutationVariables,
  ): Promise<RemoveTagFromLogEntryMutation> {
    return this.transport.graphql.request<
      RemoveTagFromLogEntryMutation,
      RemoveTagFromLogEntryMutationVariables
    >(RemoveTagFromLogEntryDocument, payload);
  }
}
