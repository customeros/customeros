import { Contact } from '@graphql/types';

import { getContactDisplayName } from './getContactName';

export const getContactPageTitle = (contact: Partial<Contact>): string => {
  if (!contact) {
    return 'Unnamed';
  }
  const contactName = getContactDisplayName(contact);

  if (contactName === 'Unnamed') {
    if (contact?.emails?.[0]?.email?.length) {
      return contact.emails[0].email;
    } else if (contact?.phoneNumbers?.length) {
      const phoneNumber =
        contact.phoneNumbers[0]?.e164 ||
        contact.phoneNumbers[0]?.rawPhoneNumber;

      return `${phoneNumber} (Unnamed contact)`;
    } else {
      return 'Unnamed';
    }
  }
  const organizationName = contact?.jobRoles?.[0]?.organization?.name;
  if (!!organizationName?.length && contact?.jobRoles?.length === 1) {
    // from requirements, do not show if there are multiple
    return `${contactName} • ${organizationName}`;
  }

  return contactName;
};
