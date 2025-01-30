import { RootStore } from '@store/root';
import { TagStore } from '@store/Tags/Tag.store';
import { Contact } from '@store/Contacts/Contact.dto';
import { ContactService as ContactRepo } from '@store/Contacts/__service__/Contacts.service';

import { unwrap } from '@utils/unwrap';

export class ContactService {
  private store = RootStore.getInstance();
  private contactRepo = ContactRepo.getInstance();

  constructor() {}

  public async addTag(contact: Contact, tag: TagStore) {
    contact.addTag(tag.id);

    const [res, err] = await unwrap(
      this.contactRepo.addTagsToContact({
        input: { contactId: contact.id, tag: { name: tag.tagName } },
      }),
    );

    if (err) {
      console.error(err);

      this.store.ui.toastError(
        'Failed to add tag to contact',
        'tag-add-failed',
      );

      return [null, err];
    }

    return [res, err];
  }

  public async removeTag(contact: Contact, tag: TagStore) {
    contact.removeTag(tag.id);

    const [res, err] = await unwrap(
      this.contactRepo.removeTagsFromContact({
        input: { contactId: contact.id, tag: { name: tag.tagName } },
      }),
    );

    if (err) {
      console.error(err);

      this.store.ui.toastError(
        'Failed to remove tag from contact',
        'tag-remove-failed',
      );

      return [null, err];
    }

    return [res, err];
  }

  public async changeContactName(contact: Contact, name: string) {
    contact.changeName(name);

    const [res, err] = await unwrap(
      this.contactRepo.updateContact({
        input: { id: contact.id, name },
      }),
    );

    if (err) {
      console.error(err);

      this.store.ui.toastError(
        'Failed to change contact name',
        'contact-name-change-failed',
      );

      return [null, err];
    }

    return [res, err];
  }
}
