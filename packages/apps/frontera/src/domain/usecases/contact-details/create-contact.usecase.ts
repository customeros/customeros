import { RootStore } from '@store/root';
import { action, computed, observable } from 'mobx';
import { ContactService } from '@store/Contacts/__service__/Contacts.service';

type ContactType = 'linkedin' | 'email' | 'name';

interface ValidationError {
  message: string;
  isEmpty: boolean;
  isInvalid: boolean;
}

export class CreateContactUsecase {
  private service = ContactService.getInstance();
  private root = RootStore.getInstance();

  @observable accessor inputValue: string = '';
  @observable accessor isLoading: boolean = false;
  @observable accessor type: ContactType = 'linkedin';
  @observable accessor organizationId: string = '';
  @observable accessor errors: Record<ContactType, ValidationError> = {
    linkedin: { isEmpty: false, isInvalid: false, message: '' },
    email: { isEmpty: false, isInvalid: false, message: '' },
    name: { isEmpty: false, isInvalid: false, message: '' },
  };

  private readonly PATTERNS = {
    linkedin: /linkedin\.com\/in\/[^/]+$/,
    email: /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/,
  };

  @action
  setType(type: ContactType) {
    this.type = type;
    this.clearErrors();
  }

  @action
  setOrganizationId(organizationId: string) {
    this.organizationId = organizationId;
  }

  @action
  setInputValue(inputValue: string) {
    this.inputValue = inputValue;
    this.clearErrors();
  }

  @computed
  get currentError(): ValidationError {
    return this.errors[this.type];
  }

  @action
  private validateInput() {
    const error = this.errors[this.type];

    if (!this.inputValue) {
      error.isEmpty = true;
    }

    if (this.type === 'name') {
      return;
    }

    const pattern = this.PATTERNS[this.type];

    if (!pattern.test(this.inputValue)) {
      error.isInvalid = true;
    }
  }

  @action
  private async checkExisting(): Promise<boolean> {
    const error = this.errors[this.type];

    try {
      if (this.type === 'linkedin') {
        const { contact_ByLinkedIn } =
          await this.service.contactExistsByLinkedIn({
            linkedIn: this.inputValue,
          });

        if (contact_ByLinkedIn?.metadata.id) {
          error.message = 'A contact with this LinkedIn already exists';

          return false;
        }
      }

      if (this.type === 'email') {
        const { contact_ByEmail } = await this.service.contactExistsByEmail({
          email: this.inputValue,
        });

        if ((contact_ByEmail?.emails?.length ?? 0) > 0) {
          error.message = 'A contact with this email already exists';

          return false;
        }
      }

      return true;
    } catch (err) {
      error.message = 'Error checking existing contact';

      return false;
    }
  }

  @action
  private async createContact() {
    const options = {
      onSuccess: () => {
        this.isLoading = false;
        this.root.ui.commandMenu.toggle('AddSingleContact');
        this.root.ui.commandMenu.clearContext();
        this.root.ui.toastSuccess('Contact created', 'contact-email-created');
      },
      onError: () => {
        this.isLoading = false;
      },
    };

    switch (this.type) {
      case 'linkedin':
        await this.root.contacts.createWithSocial({
          organizationId: this.organizationId,
          socialUrl: this.inputValue,
        });
        break;

      case 'name':
        await this.root.contacts.create(this.organizationId, options, {
          name: this.inputValue,
        });
        break;

      case 'email':
        await this.root.contacts.createWithEmail(this.organizationId, options, {
          email: { email: this.inputValue, primary: true },
        });
        break;
    }
  }

  @action
  clearErrors() {
    Object.values(this.errors).forEach((error) => {
      error.isEmpty = false;
      error.isInvalid = false;
      error.message = '';
    });
  }

  @action
  clearState() {
    this.clearErrors();
    this.inputValue = '';
    this.isLoading = false;
    this.type = 'linkedin';
  }

  @action
  async submit() {
    this.clearErrors();
    this.isLoading = true;
    this.validateInput();

    if (this.errors[this.type].isEmpty || this.errors[this.type].isInvalid) {
      this.isLoading = false;

      return;
    }

    if (this.type !== 'name' && !(await this.checkExisting())) {
      this.isLoading = false;

      return;
    }

    await this.createContact();
    this.clearState();
  }
}
