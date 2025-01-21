import { RootStore } from '@store/root';
import { ContactService } from '@domain/services';
import { action, computed, observable } from 'mobx';

export class EditContactNameUseCase {
  private root = RootStore.getInstance();
  private service = new ContactService();
  @observable accessor name: string | undefined = undefined;
  private contactId: string;

  constructor(contactId: string) {
    this.contactId = contactId;
    this.setName = this.setName.bind(this);
  }

  @action
  setName(name: string) {
    this.name = name;
  }

  @computed
  get contactName() {
    return this.name ?? this.root.contacts.getById(this.contactId)?.value.name;
  }

  execute() {
    const contact = this.root.contacts.getById(this.contactId);

    if (!contact || !this.name) return;
    this.service.changeContactName(contact, this.name || '');
  }
}
