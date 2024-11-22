import { match } from 'ts-pattern';
import { Operation } from '@store/types';
import { makePayload } from '@store/util';
import { Transport } from '@store/transport';
import { MailboxStore } from '@store/Settings/Mailbox.store';

import MailstackDomainsDocument from './getMailstackDomains.graphql';
import MailstackSetUserDocument from './updateMailstackSetUser.graphql';
import MailstackMailboxesDocument from './getMailstackMailboxes.graphql';
import MailstackUniqueUsernamesDocument from './getMailstackUniqueUsernames.graphql';
import CheckUnavailableDomainsDocument from './getMailstackCheckUnavailableDomains.graphql';
import RegisterBuyDomainsWithMailboxesDocument from './updateRegisterBuyDomainsWithMailboxes.graphql';
import GetRegisteredBuyDomainsWithMailboxesDocument from './getRegisteredBuyDomainsWithMailboxes.graphql';
import MailstackDomainPurchaseSuggestionsDocument from './getMailstackDomainsPurchaseSuggestions.graphql';
import {
  MailstackDomainsQuery,
  MailstackDomainsQueryVariables,
} from './getMailstackDomains.generated';
import {
  MailstackMailboxesQuery,
  MailstackMailboxesQueryVariables,
} from './getMailstackMailboxes.generated';
import {
  MailstackSetUserMutation,
  MailstackSetUserMutationVariables,
} from './updateMailstackSetUser.generated';
import {
  MailstackUniqueUsernamesQuery,
  MailstackUniqueUsernamesQueryVariables,
} from './getMailstackUniqueUsernames.generated';
import {
  CheckUnavailableDomainsQuery,
  CheckUnavailableDomainsQueryVariables,
} from './getMailstackCheckUnavailableDomains.generated';
import {
  RegisterBuyDomainsWithMailboxesMutation,
  RegisterBuyDomainsWithMailboxesMutationVariables,
} from './updateRegisterBuyDomainsWithMailboxes.generated';
import {
  MailstackDomainPurchaseSuggestionsQuery,
  MailstackDomainPurchaseSuggestionsQueryVariables,
} from './getMailstackDomainsPurchaseSuggestions.generated';
import {
  GetRegisteredBuyDomainsWithMailboxesQuery,
  GetRegisteredBuyDomainsWithMailboxesQueryVariables,
} from './getRegisteredBuyDomainsWithMailboxes.generated';
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

  async getMailstackDomainsSuggestions(
    payload: MailstackDomainPurchaseSuggestionsQueryVariables,
  ) {
    return this.transport.graphql.request<
      MailstackDomainPurchaseSuggestionsQuery,
      MailstackDomainPurchaseSuggestionsQueryVariables
    >(MailstackDomainPurchaseSuggestionsDocument, payload);
  }

  async getMailstackDomains() {
    return this.transport.graphql.request<
      MailstackDomainsQuery,
      MailstackDomainsQueryVariables
    >(MailstackDomainsDocument);
  }

  async getMailstackCheckUnavailableDomains(
    payload: CheckUnavailableDomainsQueryVariables,
  ) {
    return this.transport.graphql.request<
      CheckUnavailableDomainsQuery,
      CheckUnavailableDomainsQueryVariables
    >(CheckUnavailableDomainsDocument, payload);
  }

  async getMailstackMailboxes() {
    return this.transport.graphql.request<
      MailstackMailboxesQuery,
      MailstackMailboxesQueryVariables
    >(MailstackMailboxesDocument);
  }

  async getMailstackUniqueUsernames() {
    return this.transport.graphql.request<
      MailstackUniqueUsernamesQuery,
      MailstackUniqueUsernamesQueryVariables
    >(MailstackUniqueUsernamesDocument);
  }

  async updateMailstackUsernames(payload: MailstackSetUserMutationVariables) {
    return this.transport.graphql.request<
      MailstackSetUserMutation,
      MailstackSetUserMutationVariables
    >(MailstackSetUserDocument, payload);
  }

  async getRegisteredMailboxes() {
    return this.transport.graphql.request<
      GetRegisteredBuyDomainsWithMailboxesQuery,
      GetRegisteredBuyDomainsWithMailboxesQueryVariables
    >(GetRegisteredBuyDomainsWithMailboxesDocument);
  }

  async createMailbox(
    payload: RegisterBuyDomainsWithMailboxesMutationVariables,
  ) {
    return this.transport.graphql.request<
      RegisterBuyDomainsWithMailboxesMutation,
      RegisterBuyDomainsWithMailboxesMutationVariables
    >(RegisterBuyDomainsWithMailboxesDocument, payload);
  }

  public async mutateOperation(operation: Operation, store: MailboxStore) {
    const diff = operation.diff?.[0];
    const path = diff?.path;
    const mailboxNumber = operation.entityId;
    // const mailboxId = store.value.id;

    if (!operation.diff.length) {
      return;
    }

    if (!mailboxNumber) {
      console.error('Missing entityId in Operation! Mutations will not fire.');

      return;
    }
    match(path)
      .with(['domain'], () => {
        const payload =
          makePayload<MailstackDomainPurchaseSuggestionsQueryVariables>(
            operation,
          );

        this.getMailstackDomainsSuggestions({ ...payload });
      })

      .otherwise(() => {});
  }
}
