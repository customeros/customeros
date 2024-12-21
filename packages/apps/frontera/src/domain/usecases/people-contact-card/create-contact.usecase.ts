import { action, computed, observable } from 'mobx';
import { Organization } from '@store/Organizations/Organization.dto';

export class CreateContact {
  @observable accessor inputValue: string = '';
  @observable accessor type: 'linkedin' | 'email' | 'name' = 'linkedin';
  @observable accessor organizationId: string = '';
  @observable accessor entity: Organization | null = null;

  constructor() {
    this.setInputValue = this.setInputValue.bind(this);
    this.setType = this.setType.bind(this);
    this.setEntity = this.setEntity.bind(this);
  }

  @action
  setType(type: 'linkedin' | 'email' | 'name') {
    this.type = type;
  }

  @action
  setOrganizationId(organizationId: string) {
    this.organizationId = organizationId;
  }

  @action
  setInputValue(inputValue: string) {
    this.inputValue = inputValue;
  }

  @action
  setEntity(entity: Organization) {
    this.entity = entity;
  }

  @computed
  get getType() {
    return this.type;
  }

  @action
  getInputValue() {
    return this.inputValue;
  }

  @action
  submit() {
    if (this.type === 'linkedin') {
      this.entity?.store.root.contacts.createWithSocial({
        organizationId: this.organizationId,
        socialUrl: this.inputValue,
      });
    }

    if (this.type === 'name') {
      this.entity?.store.root.contacts.create(
        this.organizationId,
        {},
        { name: this.inputValue },
      );
    }
  }
}
