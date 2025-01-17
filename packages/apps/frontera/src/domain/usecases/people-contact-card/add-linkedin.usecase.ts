import type { Contact } from '@store/Contacts/Contact.dto';

import { action, observable } from 'mobx';
import { ContactService } from '@store/Contacts/__service__/Contacts.service';

export class LinkedIn {
  private service = ContactService.getInstance();
  @observable accessor linkedInUrl: string = '';
  @observable accessor inputValue: string = '';
  @observable accessor entity: Contact | null = null;
  @observable accessor emptyLinkedInUrl: boolean = false;
  @observable accessor invalidLinkedInUrl: boolean = false;
  @observable accessor error: string = '';
  @observable accessor errorsTriggered: boolean = false;

  constructor() {
    this.setInputValue = this.setInputValue.bind(this);
    this.setErrors = this.setErrors.bind(this);
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
  setErrors(error: string) {
    this.error = error;
  }

  @action
  setLinkedInUrl() {
    if (!this.entity) return;

    this.entity.draft();
    this.entity.value.linkedInUrl = this.linkedInUrl;
    this.entity.commit();
  }

  @action
  async checkIfLinkedInUrlExists(linkedInUrl: string) {
    if (!this.entity) return;

    const { contact_ByLinkedIn } = await this.service.contactExistsByLinkedIn({
      linkedIn: linkedInUrl,
    });

    if (contact_ByLinkedIn?.metadata.id) {
      this.setErrors('A contact with this LinkedIn already exists');
    }
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
    this.errorsTriggered = false;
  }

  @action
  resetErrors() {
    this.error = '';
    this.emptyLinkedInUrl = false;
    this.invalidLinkedInUrl = false;
  }

  @action
  async getErrors() {
    if (!this.errorsTriggered) return false;

    return !!(this.validateLinkedInUrl() || (await this.error));
  }

  @action
  async submitLinkedInUrl() {
    this.errorsTriggered = true;

    if (this.validateLinkedInUrl()) return;
    await this.checkIfLinkedInUrlExists(this.linkedInUrl);

    if (this.error) return;

    this.setLinkedInUrl();
    this.clearState();
  }
}
