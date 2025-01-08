import { action, computed, observable } from 'mobx';
import { Organization } from '@store/Organizations/Organization.dto';

export class CreateContact {
  @observable accessor inputValue: string = '';
  @observable accessor type: 'linkedin' | 'email' | 'name' = 'linkedin';
  @observable accessor organizationId: string = '';
  @observable accessor entity: Organization | null = null;
  @observable accessor invalidName: boolean = false;
  @observable accessor invalidLinkedInUrl: boolean = false;
  @observable accessor emptyLinkedInUrl: boolean = false;

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
  validateName() {
    if (this.inputValue.length === 0) {
      this.invalidName = true;
    } else {
      this.invalidName = false;
    }
  }

  @action
  validateLinkedInUrl() {
    const url = this.inputValue;

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
  }

  @action
  clearState() {
    this.inputValue = '';
    this.emptyLinkedInUrl = false;
    this.invalidLinkedInUrl = false;
  }

  @action
  submit() {
    if (this.type === 'linkedin') {
      this.validateLinkedInUrl();

      if (this.emptyLinkedInUrl || this.invalidLinkedInUrl) return;
      this.entity?.store.root.contacts.createWithSocial({
        organizationId: this.organizationId,
        socialUrl: this.inputValue,
      });
    }

    if (this.type === 'name') {
      this.validateName();

      if (this.invalidName) return;
      this.entity?.store.root.contacts.create(
        this.organizationId,
        {},
        { name: this.inputValue },
      );
    }
    this.clearState();
  }
}
