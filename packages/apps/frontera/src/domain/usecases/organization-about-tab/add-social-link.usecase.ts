import { action, observable } from 'mobx';
import { Organization } from '@store/Organizations/Organization.dto';

import { Social } from '@graphql/types';

export class AddSocialLinkCase {
  @observable accessor url: string = '';
  @observable accessor inputValue: string = '';
  @observable accessor entity: Organization | null = null;

  constructor() {
    this.setInputValue = this.setInputValue.bind(this);
  }

  @action
  setEntity(entity: Organization) {
    this.entity = entity;
  }

  @action
  getSocialUrl() {
    return this.url;
  }

  @action
  setInputValue(inputValue: string) {
    this.inputValue = inputValue;
    this.url = inputValue;
  }

  @action
  getInputValue() {
    return this.inputValue;
  }

  @action
  submit() {
    if (!this?.entity) return;

    this.entity.draft();
    this.entity.value?.socialMedia.push({
      id: crypto.randomUUID(),
      url: this.url,
    } as Social);
    this.entity.commit();
  }
}
