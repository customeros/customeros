import type { Contact } from '@store/Contacts/Contact.dto';

import { action, observable } from 'mobx';

export class AddLinkedin {
  @observable accessor linkedInUrl: string = '';
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
  getLinkdInUrl() {
    return this.linkedInUrl;
  }

  @action
  setInputValue(inputValue: string) {
    this.inputValue = inputValue;
    this.linkedInUrl = inputValue;
  }

  @action
  getInputValue() {
    return this.inputValue;
  }

  @action
  setLinkedInUrl() {
    if (!this.entity) return;

    this.entity.draft();
    this.entity.value.linkedInUrl = this.linkedInUrl;
    this.entity.commit();
    this.inputValue = '';
  }
}
