import { RootStore } from '@store/root';
import { ContactService } from '@domain/services';

export class EditContactNameUseCase {
  private root = RootStore.getInstance();
  private service = new ContactService();
  private contactId: string;

  constructor(contactId: string) {
    this.contactId = contactId;
  }

  execute() {
    const contact = this.root.contacts.getById(this.contactId);

    if (!contact) return;
    this.service.changeContactName(contact, contact.name);
  }
}
