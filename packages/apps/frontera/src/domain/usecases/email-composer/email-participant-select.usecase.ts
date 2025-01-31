import { RootStore } from '@store/root';
import { action, computed, reaction, observable } from 'mobx';

import { validateEmail } from '@utils/email.ts';
import { SelectOption } from '@ui/utils/types.ts';

export class EmailParticipantSelectUsecase {
  @observable public accessor searchTerm = '';
  @observable private accessor newEmailOptions = new Set();
  @observable public accessor selectedEmails: {
    label: string;
    value: string;
  }[] = [];
  @observable public accessor emailOptions: { label: string; value: string }[] =
    [];

  private root = RootStore.getInstance();
  private organizationId: string;

  constructor(organizationId: string) {
    this.organizationId = organizationId;
    this.select = this.select.bind(this);
    this.create = this.create.bind(this);
    this.removeSelectedEmail = this.removeSelectedEmail.bind(this);
    this.setSearchTerm = this.setSearchTerm.bind(this);
    this.computeInitialEmailOptions =
      this.computeInitialEmailOptions.bind(this);

    reaction(
      () => this.newEmailOptions.size || this.organization?.contacts.length,

      this.computeInitialEmailOptions,
    );
    this.computeInitialEmailOptions();
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
    this.emailOptions = (this.organization?.contacts || [])
      .flatMap((e) => e.emails.map((d) => ({ ...d, contactName: e.name })))
      .filter(Boolean)
      .sort((a, b) => {
        const aInOrg = this.orgContactEmails.has(a.email);
        const bInOrg = this.orgContactEmails.has(b.email);

        if (aInOrg && !bInOrg) return -1;
        if (!aInOrg && bInOrg) return 1;

        return 0;
      })
      .map((contact) => {
        return {
          value: contact?.email || '',
          label: contact?.contactName || '',
        };
      });
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

    return sorted.filter((tag) =>
      tag.label.toLowerCase().includes(this.searchTerm.toLowerCase()),
    );
  }

  @action
  public removeSelectedEmail(val: string) {
    this.selectedEmails = this.selectedEmails.filter((e) => e.value !== val);
  }

  @action
  public setSearchTerm(searchTerm: string) {
    this.searchTerm = searchTerm;
  }

  @action
  public reset() {
    this.setSearchTerm('');
    this.newEmailOptions.clear();
    this.selectedEmails = [];
  }

  @action
  public close() {
    this.reset();
    this.root.ui.commandMenu.setOpen(false);
  }

  @action
  public select(values?: SelectOption[] | SelectOption) {
    if (!values) {
      this.selectedEmails = [];

      return;
    }
    this.selectedEmails = Array.isArray(values)
      ? values
      : [...this.selectedEmails, values];
  }

  @action
  public addOption() {
    if (!this.searchTerm) return;
    const emailValidation = validateEmail(this.searchTerm);

    if (emailValidation) {
      this.root.ui.toastError('Invalid email address', 'invalid-email');

      return;
    }

    this.emailOptions.push({ label: '', value: this.searchTerm });

    this.select({ label: '', value: this.searchTerm });
    this.searchTerm = '';
  }

  @action
  public create() {
    if (!this.organization) return;

    this.root.contacts.createWithEmail(
      this.organizationId,
      {
        onSuccess: () =>
          this.root.ui.toastSuccess('Contact created', 'contact-email-created'),
      },
      { email: { email: this.searchTerm, primary: true } },
    );
  }
}
