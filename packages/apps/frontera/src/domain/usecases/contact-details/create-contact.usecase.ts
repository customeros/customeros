import { RootStore } from '@store/root';
import { action, computed, observable } from 'mobx';
import { ContactService } from '@store/Contacts/__service__/Contacts.service';

export class CreateContact {
  private service = ContactService.getInstance();
  private root = RootStore.getInstance();

  @observable accessor inputValue: string = '';
  @observable accessor type: 'linkedin' | 'email' | 'name' = 'linkedin';
  @observable accessor organizationId: string = '';
  @observable accessor invalidName: boolean = false;
  @observable accessor invalidLinkedInUrl: boolean = false;
  @observable accessor emptyLinkedInUrl: boolean = false;
  @observable accessor errorLinkedIn: string = '';
  @observable accessor errorEmail: string = '';
  @observable accessor emptyEmail: boolean = false;
  @observable accessor invalidEmail: boolean = false;

  constructor() {
    this.setInputValue = this.setInputValue.bind(this);
    this.setType = this.setType.bind(this);
    this.setErrorEmail = this.setErrorEmail.bind(this);
    this.setErrorLinkedin = this.setErrorLinkedin.bind(this);
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

  @computed
  get getType() {
    return this.type;
  }

  @action
  getInputValue() {
    return this.inputValue;
  }

  @action
  setErrorEmail(error: string) {
    this.errorEmail = error;
  }

  @action
  setErrorLinkedin(error: string) {
    this.errorLinkedIn = error;
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
  async checkIfLinkedInUrlExists(linkedInUrl: string) {
    const { contact_ByLinkedIn } = await this.service.contactExistsByLinkedIn({
      linkedIn: linkedInUrl,
    });

    if (contact_ByLinkedIn?.metadata.id) {
      this.setErrorLinkedin('A contact with this LinkedIn already exists');
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
  validateEmail() {
    const email = this.inputValue;

    if (email.length === 0) {
      this.emptyEmail = true;
    } else {
      this.emptyEmail = false;
    }

    const emailPatern = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

    if (!emailPatern.test(email)) {
      this.invalidEmail = true;
    } else {
      this.invalidEmail = false;
    }
  }

  @action
  async checkIfEmailExists(email: string) {
    const { contact_ByEmail } = await this.service.contactExistsByEmail({
      email: email,
    });

    if ((contact_ByEmail?.emails?.length ?? 0) > 0) {
      this.setErrorEmail('A contact with this email already exists');
    }
  }

  @action
  clearState() {
    this.errorEmail = '';
    this.errorLinkedIn = '';
    this.emptyLinkedInUrl = false;
    this.invalidLinkedInUrl = false;
    this.emptyEmail = false;
    this.invalidEmail = false;
  }

  @action
  async submit() {
    if (this.type === 'linkedin') {
      this.validateLinkedInUrl();

      if (this.emptyLinkedInUrl) return;
      if (this.invalidLinkedInUrl) return;
      await this.checkIfLinkedInUrlExists(this.inputValue);

      if (this.errorLinkedIn) return;

      this.root.contacts.createWithSocial({
        organizationId: this.organizationId,
        socialUrl: this.inputValue,
      });
    }

    if (this.type === 'name') {
      this.validateName();

      if (this.invalidName) return;
      this.root.contacts.create(
        this.organizationId,
        {
          onSuccess: () =>
            this.root.ui.toastSuccess(
              'Contact created',
              'contact-email-created',
            ),
        },
        { name: this.inputValue },
      );
    }

    if (this.type === 'email') {
      this.validateEmail();

      if (this.emptyEmail) return;
      if (this.invalidEmail) return;
      await this.checkIfEmailExists(this.inputValue);

      if (this.errorEmail) return;

      this.root.contacts.createWithEmail(
        this.organizationId,
        {
          onSuccess: () =>
            this.root.ui.toastSuccess(
              'Contact created',
              'contact-email-created',
            ),
        },
        { email: { email: this.inputValue, primary: true } },
      );
    }
    this.clearState();
  }
}
