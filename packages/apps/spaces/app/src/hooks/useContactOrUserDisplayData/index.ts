import { useCallback } from 'react';
import { Contact, User } from '@graphql/types';

export const useContactOrUserDisplayName = () => {
  return useCallback((contact?: Partial<Contact | User> | null) => {
    if (!contact) return 'Unnamed';

    if (contact.__typename === 'Contact' && contact?.name) {
      return contact.name;
    }

    const name = `${contact?.firstName} ${contact?.lastName}`;

    if (name.trim()) {
      return name;
    }

    return contact?.emails?.[0]?.email
      ? contact?.emails?.[0]?.email
      : 'Unnamed';
  }, []);
};
