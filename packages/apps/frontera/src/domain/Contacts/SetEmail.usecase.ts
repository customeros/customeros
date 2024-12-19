import type { Contact } from '@store/Contacts/Contact.dto';

import { set } from 'lodash';
import { action, observable } from 'mobx';
import { ContactService } from '@store/Contacts/__service__/Contacts.service';

import { Email } from '@shared/types/__generated__/graphql.types';

export class SetEmailCase {
  @observable accessor email: string = '';
  @observable accessor oldEmail: string = '';
  @observable accessor entity: Contact | null = null;
  @observable accessor emailError: boolean = false;
  private service: ContactService;

  constructor() {
    this.service = new ContactService();
  }

  @action
  setEntity(entity: Contact) {
    this.entity = entity;
  }

  @action
  setEmail(email: string) {
    this.email = email;
  }

  @action
  getEmail() {
    return this.email;
  }

  @action
  setEmailFromContact() {
    const primaryEmail = this.entity?.value.emails.find(
      (e) => e.primary,
    )?.email;

    this.email = primaryEmail ?? '';
  }

  @action
  setPreviousEmail(index: number) {
    this.oldEmail = this.entity?.value.emails[index]?.email || '';
  }

  public setEmailForContact() {
    this.entity?.draft();
    this.entity?.value.emails.push({ email: this.email } as Email);
    this.entity?.commit();
  }

  @action
  public setPrimaryEmailForContact(syncOnly = false) {
    if (this.entity) {
      this.entity.draft();
      set(this.entity.value, 'primaryEmail.email', this.email);
      this.entity.commit({ syncOnly: syncOnly });

      const foundIndex = this.entity.value.emails.findIndex(
        (e) => e.email === this.email,
      );
      const foundOldIndex = this.entity.value.emails.findIndex(
        (e) => e.primary === true,
      );

      if (foundIndex !== -1 && foundOldIndex !== -1) {
        this.entity.draft();
        this.entity.value.emails[foundOldIndex].primary = false;
        this.entity.value.emails[foundIndex].primary = true;
        this.entity.commit({ syncOnly: true });
      }
    }
  }

  public updateEmailForContact(emailIdx: number) {
    if (this.entity) {
      this.entity.draft();

      set(this.entity.value.emails[emailIdx], 'email', this.email);
      // this.entity.value.emails[emailIdx].email = this.email;
      this.entity.commit();
    }
  }

  @action
  public updatePrimaryEmailForContact() {
    const primaryEmail = this.entity?.value.emails.find((e) => e.primary);

    if (this.entity) {
      this.entity.draft();
      primaryEmail!.email = this.email;
      this.entity.commit({ syncOnly: true });
    }
  }
}
