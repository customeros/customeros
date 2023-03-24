import { Contact, User } from '../graphQL/__generated__/generated';

export const getContactDisplayName = (
  contact?: Partial<Contact | User> | null,
) => {
  if (!contact) return 'Unnamed';

  if (contact.__typename === 'Contact' && contact?.name) {
    return contact.name;
  }

  const name = `${contact?.firstName} ${contact?.lastName}`;
  return name.trim().length ? name : 'Unnamed';
};
