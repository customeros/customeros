import type { Contact } from '@store/Contacts/Contact.dto';

import { action, observable } from 'mobx';

export class CreateContact {
  @observable accessor name: string = '';
  @observable accessor inputValue: string = '';
  @observable accessor entity: Contact | null = null;

  constructor() {}

  @action
  setName(name: string) {
    this.name = name;
  }

  @action
  setEntity(entity: Contact) {
    this.entity = entity;
  }

  @action
  getName() {
    return this.name;
  }

  @action
  setInputValue(inputValue: string) {
    this.inputValue = inputValue;
  }

  @action
  getInputValue() {
    return this.inputValue;
  }
}
