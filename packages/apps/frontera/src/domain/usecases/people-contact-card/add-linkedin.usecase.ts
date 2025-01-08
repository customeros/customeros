import type { Contact } from '@store/Contacts/Contact.dto';

import { action, observable } from 'mobx';

export class LinkedIn {
  @observable accessor linkedInUrl: string = '';
  @observable accessor inputValue: string = '';
  @observable accessor entity: Contact | null = null;
  @observable accessor emptyLinkedInUrl: boolean = false;
  @observable accessor invalidLinkedInUrl: boolean = false;

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
  }

  @action
  validateLinkedInUrl() {
    const url = this.linkedInUrl;

    if (url.length === 0) {
      this.emptyLinkedInUrl = true;
    } else {
      this.emptyLinkedInUrl = false;

      const linkedInUrlPattern = /^(https?:\/\/)?(www\.)?linkedin\.com\/.*$/;

      if (!linkedInUrlPattern.test(url)) {
        this.invalidLinkedInUrl = true;
      } else {
        this.invalidLinkedInUrl = false;
      }
    }

    return this.emptyLinkedInUrl || this.invalidLinkedInUrl;
  }

  @action
  clearState() {
    this.linkedInUrl = '';
    this.inputValue = '';
    this.emptyLinkedInUrl = false;
    this.invalidLinkedInUrl = false;
  }

  @action
  submitLinkedInUrl() {
    if (this.validateLinkedInUrl()) return;

    this.setLinkedInUrl();
    this.clearState();
  }
}
