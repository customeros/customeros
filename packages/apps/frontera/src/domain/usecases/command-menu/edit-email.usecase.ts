import { action, observable } from 'mobx';
import { Contact } from '@store/Contacts/Contact.dto';

export class EditEmailCase {
  private static instance: EditEmailCase;
  @observable accessor email: string = '';
  @observable accessor oldEmail: string = '';
  @observable accessor entity: Contact | null = null;

  private constructor() {
    this.setEmail = this.setEmail.bind(this);
  }

  public static getInstance(): EditEmailCase {
    if (!EditEmailCase.instance) {
      EditEmailCase.instance = new EditEmailCase();
    }

    return EditEmailCase.instance;
  }

  @action
  public setEmail(email: string) {
    this.email = email;
  }

  @action
  public setOldEmail(oldEmail: string) {
    this.oldEmail = oldEmail;
  }

  @action
  public setEntity(entity: Contact) {
    this.entity = entity;
  }

  @action
  submit() {
    if (!this.entity) return;
    const foundIdx = this.entity.value.emails.findIndex(
      (e) => e.email === this.oldEmail,
    );

    this.entity.draft();
    this.entity.value.emails[foundIdx].email = this.email;
    this.entity.commit();
  }
}
