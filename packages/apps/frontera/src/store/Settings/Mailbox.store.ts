import { RootStore } from '@store/root';
import { Syncable } from '@store/syncable';
import { Transport } from '@store/transport';
import { action, override, makeObservable } from 'mobx';

import {
  MailstackBuyRequest,
  MailstackBuyRequestStatus,
} from '@shared/types/__generated__/graphql.types';

import { MailboxesService } from './__service__/Mailboxes/Mailboxes.service';

export class MailboxStore extends Syncable<MailstackBuyRequest> {
  private service: MailboxesService;

  constructor(
    public root: RootStore,
    public transport: Transport,
    data: MailstackBuyRequest,
  ) {
    super(root, transport, data ?? getDefaultValue());
    this.service = MailboxesService.getInstance(transport);

    makeObservable<MailboxStore>(this, {
      id: override,
      save: override,
      setId: override,
      getId: override,
      invalidate: action,
      getMailbox: action,
      getChannelName: override,
    });
  }

  getMailbox() {
    return this.service.getMailstackMailboxes();
  }

  getId() {
    return this.value.id;
  }

  setId(id: string) {
    this.value.id = id;
  }

  getChannelName(): string {
    return 'Mailbox';
  }

  static getDefaultValue(): MailstackBuyRequest {
    return {
      id: crypto.randomUUID(),
      domains: [],
      mailboxes: [],
      createdAt: new Date().toISOString(),
      status: MailstackBuyRequestStatus.Pending,
    };
  }
}

export const getDefaultValue = (): MailstackBuyRequest => ({
  id: crypto.randomUUID(),
  domains: [],
  createdAt: new Date().toISOString(),
  mailboxes: [],
  status: MailstackBuyRequestStatus.Pending,
});
