import { RootStore } from '@store/root';
import { Syncable } from '@store/syncable';
import { Transport } from '@infra/transport';
import { action, override, computed, makeObservable } from 'mobx';

import { MailboxesService } from './__service__/Mailboxes/Mailboxes.service';
import { GetMailboxesQuery } from './__service__/Mailboxes/getMailboxes.generated';

export type Mailbox = NonNullable<
  GetMailboxesQuery['mailstack_Mailboxes'][number]
>;

export class MailboxStore extends Syncable<Mailbox> {
  private service: MailboxesService;

  constructor(
    public root: RootStore,
    public transport: Transport,
    data: Mailbox,
  ) {
    super(root, transport, data ?? MailboxStore.getDefaultValue());
    this.service = MailboxesService.getInstance(transport);

    makeObservable<MailboxStore>(this, {
      id: override,
      save: override,
      setId: override,
      getId: override,
      invalidate: action,
      getChannelName: override,
      user: computed,
    });
  }

  get user() {
    if (!this.value.userId) return null;

    return this.root.users.value.get(this.value.userId) ?? null;
  }

  getId() {
    return this.value.mailbox;
  }

  setId(id: string) {
    this.value.mailbox = id;
  }

  getChannelName(): string {
    return 'Mailboxes';
  }

  static getDefaultValue(): Mailbox {
    return {
      domain: '',
      mailbox: crypto.randomUUID(),
      userId: '',
      rampUpMax: 0,
      rampUpRate: 0,
      rampUpCurrent: 0,
      scheduledEmails: 0,
      currentFlowIds: [],
    };
  }
}
