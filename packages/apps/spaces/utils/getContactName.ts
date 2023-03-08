import { Contact } from '../graphQL/__generated__/generated';

export const getContactDisplayName = (contact?: Partial<Contact> | null) => {
  if (!contact) return 'Unnamed';

  if (contact?.name) {
    return contact.name;
  }

  const name = `${contact?.firstName} ${contact?.lastName}`;
  return name.trim().length ? name : 'Unnamed';
};
