import { Contact } from '../graphQL/__generated__/generated';

export const getContactDisplayName = (contact?: Partial<Contact> | null) => {
  if (!contact) return 'Unnamed';
  const name = `${contact?.firstName} ${contact?.lastName} ${
    contact?.name || ''
  }`;
  return name.trim().length ? name : 'Unnamed';
};
