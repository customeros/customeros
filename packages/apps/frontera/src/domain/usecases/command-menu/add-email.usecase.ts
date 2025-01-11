import { action, observable } from 'mobx';
import { Contact } from '@store/Contacts/Contact.dto';
import { ContactService } from '@store/Contacts/__service__/Contacts.service';

import { Email } from '@shared/types/__generated__/graphql.types';

export class AddEmailCase {
  private service = ContactService.getInstance();
  @observable accessor inputValue: string = '';
  @observable accessor entity: Contact | null = null;
  @observable accessor error: string = '';

  constructor() {
    this.setInputValue = this.setInputValue.bind(this);
    this.setErrors = this.setErrors.bind(this);
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
  async checkIfEmailExists(email: string) {
    const { contact_ByEmail } = await this.service.contactExistsByEmail({
      email: email,
    });

    if ((contact_ByEmail?.emails?.length ?? 0) > 0) {
      this.setErrors('A contact with this email already exists');
    }
  }

  @action
  setErrors(error: string) {
    this.error = error;
  }

  @action
  getErrors() {
    return this.error;
  }

  @action
  async submit() {
    const noEmails = this.entity?.value.emails.length === 0;

    await this.checkIfEmailExists(this.inputValue);

    if (this.error) return;
    if (!this.inputValue) return;
    this.entity?.draft();
    this.entity?.value.emails.push({
      email: this.inputValue,
      primary: noEmails,
    } as Email);
    this.entity?.commit();
    this.inputValue = '';
  }
}
