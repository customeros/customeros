import { action, observable } from 'mobx';
import { Organization } from '@store/Organizations/Organization.dto.ts';

export class AddOrganizationDomainCase {
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
  setInputValue(inputValue: string) {
    this.inputValue = inputValue;
  }

  @action
  getInputValue() {
    return this.inputValue;
  }

  @action
  submit() {
    this.entity?.draft();
    this.entity?.value.domains.push(this.inputValue);
    this.entity?.commit();
    this.inputValue = '';
  }
}
