import { Transport } from '@store/transport';

import {
  RemindersQuery,
  RemindersQueryVariables,
} from '@organization/graphql/reminders.generated';
import {
  CreateReminderMutation,
  CreateReminderMutationVariables,
} from '@organization/graphql/createReminder.generated';
import {
  UpdateReminderMutation,
  UpdateReminderMutationVariables,
} from '@organization/graphql/updateReminder.generated';

import RemindersDocument from './reminders.graphql';
import UpdateReminderDocument from './updateReminder.graphql';
import CreateReminderDocument from './createReminder.graphql';

export class RemindersService {
  private static instance: RemindersService | null = null;
  private transport: Transport;

  constructor(transport: Transport) {
    this.transport = transport;
  }

  static getInstance(transport: Transport): RemindersService {
    if (!RemindersService.instance) {
      RemindersService.instance = new RemindersService(transport);
    }

    return RemindersService.instance;
  }

  async getRemindersByOrganizationId(
    payload: RemindersQueryVariables,
  ): Promise<RemindersQuery> {
    return this.transport.graphql.request<
      RemindersQuery,
      RemindersQueryVariables
    >(RemindersDocument, payload);
  }

  async updateReminder(payload: UpdateReminderMutationVariables) {
    return this.transport.graphql.request<
      UpdateReminderMutation,
      UpdateReminderMutationVariables
    >(UpdateReminderDocument, payload);
  }

  async createReminder(payload: CreateReminderMutationVariables) {
    return this.transport.graphql.request<
      CreateReminderMutation,
      CreateReminderMutationVariables
    >(CreateReminderDocument, payload);
  }
}
