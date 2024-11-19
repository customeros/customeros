import { User, Contact } from '@graphql/types';

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
