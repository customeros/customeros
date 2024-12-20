import { action, observable } from 'mobx';
import { Contact } from '@store/Contacts/Contact.dto';

import { Email } from '@shared/types/__generated__/graphql.types';

export class AddEmailCase {
  @observable accessor inputValue: string = '';
  @observable accessor entity: Contact | null = null;

  constructor() {
    this.setInputValue = this.setInputValue.bind(this);
  }

  @action
  setEntity(entity: Contact) {
    this.entity = entity;
  }

  @action
  setInputValue(inputValue: string) {
    this.inputValue = inputValue;
  }

  @action
  submit() {
    const noEmails = this.entity?.value.emails.length === 0;

    this.entity?.draft();
    this.entity?.value.emails.push({ email: this.inputValue } as Email);
    this.entity?.commit();

    if (noEmails) {
      this.entity?.draft();
      this.entity!.value.emails[0].primary = true;
      this.entity?.commit({ syncOnly: true });
    }

    this.inputValue = '';
  }
}
