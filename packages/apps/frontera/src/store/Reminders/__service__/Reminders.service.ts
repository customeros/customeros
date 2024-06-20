import { Transport } from '@store/transport';

import RemindersDocument from './reminders.graphql';
import UpdateReminderDocument from './updateReminder.graphql';
import CreateReminderDocument from './createReminder.graphql';
import { RemindersQuery, RemindersQueryVariables } from './reminders.generated';
import {
  UpdateReminderMutation,
  UpdateReminderMutationVariables,
} from './updateReminder.generated';
import {
  CreateReminderMutation,
  CreateReminderMutationVariables,
} from './createReminder.generated';

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
