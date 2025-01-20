import { match } from 'ts-pattern';
import { Operation } from '@store/types';
import { Transport } from '@infra/transport';
import { MailboxStore } from '@store/Settings/Mailbox.store';

import BuyDomainsDocument from './buyDomains.graphql';
import GetDomainsDocument from './getDomains.graphql';
import { GetDomainsQuery } from './getDomains.generated';
import GetMailboxesDocument from './getMailboxes.graphql';
import { GetMailboxesQuery } from './getMailboxes.generated';
import ValidateDomainsDocument from './validateDomains.graphql';
import GetPaymentIntentDocument from './getPaymentIntent.graphql';
import MailstackSetUserDocument from './updateMailstackSetUser.graphql';
import GetDomainSuggestionsDocument from './getDomainSuggestions.graphql';
import {
  BuyDomainsMutation,
  BuyDomainsMutationVariables,
} from './buyDomains.generated';
import {
  ValidateDomainsQuery,
  ValidateDomainsQueryVariables,
} from './validateDomains.generated';
import {
  GetPaymentIntentMutation,
  GetPaymentIntentMutationVariables,
} from './getPaymentIntent.generated';
import {
  GetDomainSuggestionsQuery,
  GetDomainSuggestionsQueryVariables,
} from './getDomainSuggestions.generated';
import {
  MailstackSetUserMutation,
  MailstackSetUserMutationVariables,
} from './updateMailstackSetUser.generated';

export class MailboxesService {
  private static instance: MailboxesService;
  private transport: Transport;

  private constructor(transport: Transport) {
    this.transport = transport;
  }

  static getInstance(transport: Transport) {
    if (!MailboxesService.instance) {
      MailboxesService.instance = new MailboxesService(transport);
    }

    return MailboxesService.instance;
  }

  async getDomainSuggestions(payload: GetDomainSuggestionsQueryVariables) {
    return this.transport.graphql.request<
      GetDomainSuggestionsQuery,
      GetDomainSuggestionsQueryVariables
    >(GetDomainSuggestionsDocument, payload);
  }

  async getMailstackDomains() {
    return this.transport.graphql.request<GetDomainsQuery>(GetDomainsDocument);
  }

  async validateDomains(payload: ValidateDomainsQueryVariables) {
    return this.transport.graphql.request<
      ValidateDomainsQuery,
      ValidateDomainsQueryVariables
    >(ValidateDomainsDocument, payload);
  }

  async getMailboxes() {
    return this.transport.graphql.request<GetMailboxesQuery>(
      GetMailboxesDocument,
    );
  }

  async updateUser(payload: MailstackSetUserMutationVariables) {
    return this.transport.graphql.request<
      MailstackSetUserMutation,
      MailstackSetUserMutationVariables
    >(MailstackSetUserDocument, payload);
  }

  async buyDomains(payload: BuyDomainsMutationVariables) {
    return this.transport.graphql.request<
      BuyDomainsMutation,
      BuyDomainsMutationVariables
    >(BuyDomainsDocument, payload);
  }

  async getPaymentIntent(payload: GetPaymentIntentMutationVariables) {
    return this.transport.graphql.request<
      GetPaymentIntentMutation,
      GetPaymentIntentMutationVariables
    >(GetPaymentIntentDocument, payload);
  }

  public async mutateOperation(operation: Operation, store: MailboxStore) {
    const diff = operation.diff?.[0];
    const path = diff?.path;
    const mailboxId = operation.entityId;

    if (!operation.diff.length) {
      return;
    }

    if (!mailboxId) {
      console.error('Missing entityId in Operation! Mutations will not fire.');

      return;
    }
    match(path)
      .with(['userId'], async () => {
        try {
          await this.updateUser({
            mailbox: store.value.mailbox,
            userId: diff.val,
          });

          store.root.ui.toastSuccess(
            'Mailbox user updated',
            'mailbox-user-update',
          );
        } catch (err) {
          store.root.ui.toastError(
            'Could not update mailbox user',
            'mailbox-user-update',
          );
        }
      })

      .otherwise(() => {});
  }
}
