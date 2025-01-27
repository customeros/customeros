import { RootStore } from '@store/root';
import { action, computed, observable } from 'mobx';

import { SelectOption } from '@ui/utils/types.ts';

function extractEventUUID(url: string) {
  const matches = url.match(
    /events=([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})/,
  );

  return matches ? matches[1] : null;
}

type EmailOption = {
  label: string;
  value: string;
  active?: boolean;
  provider?: string;
};
export class EmailSenderSelectUsecase {
  @observable public accessor searchTerm = '';
  @observable private accessor newEmailOptions = new Set();
  @observable public accessor selectedEmail: EmailOption | undefined =
    undefined;
  @observable public accessor emailOptions: {
    label: string;
    value: string;
    active?: boolean;
    provider?: string;
  }[] = [];

  private root = RootStore.getInstance();
  private organizationId: string;
  private attendees: string[] = [];
  private currentUserId?: string;

  constructor(
    organizationId: string,
    attendees: string[],
    currentUserId?: string,
  ) {
    this.organizationId = organizationId;
    this.attendees = attendees;
    this.currentUserId = currentUserId;
    this.computeInitialEmailOptions();
    this.select = this.select.bind(this);
  }

  @computed
  get organization() {
    if (!this.organizationId) return;

    return this.root.organizations.getById(this.organizationId);
  }

  @computed
  get orgContactEmails() {
    return new Set(
      (this.organization?.contacts ?? []).flatMap((e) =>
        e.emails.map((d) => d.email),
      ),
    );
  }

  @action
  private computeInitialEmailOptions() {
    const id = extractEventUUID(window.location.href);

    const userMailboxes = this.root.globalCache.value?.user?.mailboxes;

    const activeEmailOptions = this.root.globalCache.value?.activeEmailTokens
      .filter((a) => (id && this.attendees?.includes(a.email)) || !id)
      .map((v) => ({
        label: v.email,
        value: v.email,
        provider: v.provider,
        active: true,
      }));

    const inactiveEmailOptions =
      this.root.globalCache.value?.inactiveEmailTokens
        .filter((a) => (id && this.attendees.includes(a.email)) || !id)
        .map((v) => ({
          label: v.email,
          value: v.email,
          provider: v.provider,
          active: false,
        }));

    const mailboxOptions = this.root.globalCache.value?.mailboxes
      .filter((v) => this.attendees.includes(v))
      .map((v) => ({
        label: v,
        value: v,
        provider: 'mailbox',
        active: true,
      }));

    const userMailboxOptions = userMailboxes?.map((v) => ({
      label: v,
      value: v,
      provider: 'mailbox',
      active: true,
    }));

    this.emailOptions = [
      ...(activeEmailOptions || []),
      ...(inactiveEmailOptions || []),
      ...(mailboxOptions || []),
      ...(userMailboxOptions || []),
    ];
  }

  @computed
  get emailOptionsList() {
    const sorted = this.emailOptions.slice().sort((a, b) => {
      const aInOrg = this.newEmailOptions.has(a.label);
      const bInOrg = this.newEmailOptions.has(b.label);

      if (aInOrg && !bInOrg) return -1;
      if (!aInOrg && bInOrg) return 1;

      return 0;
    });

    return sorted.filter(
      (email) =>
        email.label.toLowerCase().includes(this.searchTerm.toLowerCase()) ||
        email.value.toLowerCase().includes(this.searchTerm.toLowerCase()),
    );
  }

  @action
  public setSearchTerm(searchTerm: string) {
    this.searchTerm = searchTerm;
  }

  @action
  public reset() {
    this.setSearchTerm('');
    this.newEmailOptions.clear();
  }

  @action
  public close() {
    this.reset();
    this.root.ui.commandMenu.setOpen(false);
  }

  @action
  public select(value: SelectOption) {
    if (!value) {
      this.selectedEmail = undefined;
    }
    this.selectedEmail = value;
  }
}
