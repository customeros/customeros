import { Contact, User } from '../graphQL/__generated__/generated';

export const getContactDisplayName = (
  contact?: Partial<Contact | User> | null,
) => {
  if (!contact) return 'Unnamed';

  if (contact.__typename === 'Contact' && contact?.name) {
    return contact.name;
  }

  const name = `${contact?.firstName} ${contact?.lastName}`;
  console.log('🏷️ ----- name: ', name);
  return name.trim().length ? name : 'Unnamed';
};
export const getContactDisplayFirstName = (
  contact?: Partial<Contact | User> | null,
) => {
  if (!contact) return 'Unnamed';

  if (contact.__typename === 'Contact' && contact?.name) {
    return contact.name;
  }

  return contact?.firstName?.length ? contact.firstName : 'Unnamed';
};
export const getContactFirstName = (contact?: Partial<Contact> | null) => {
  if (!contact) return 'Unnamed';

  if (contact?.firstName) {
    return contact.firstName;
  }

  return contact.name || '';
};
