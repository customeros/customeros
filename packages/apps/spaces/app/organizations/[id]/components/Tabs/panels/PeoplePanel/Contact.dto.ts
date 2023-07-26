import { Contact } from '@graphql/types';

import { UpdateContactMutationVariables } from '@organization/graphql/updateContact.generated';

export interface ContactForm {
  id: string;
  name: string;
  title: string;
  role: string;
  roleId: string;
  note: string;
  email: string;
  phone: string;
  phoneId: string;
  timezone: string;
  companyName?: string;
}

export class ContactFormDto implements ContactForm {
  id: string; // auxiliary field
  name: string;
  role: string;
  roleId: string; // auxiliary field
  title: string;
  note: string;
  email: string;
  phone: string;
  phoneId: string; // auxiliary field
  timezone: string;
  companyName?: string | undefined;

  constructor(data?: Partial<Contact> | null) {
    this.id = data?.id || ''; // auxiliary field
    this.name = data?.firstName || '';
    this.role = data?.jobRoles?.[0]?.jobTitle || '';
    this.roleId = data?.jobRoles?.[0]?.id || ''; // auxiliary field
    this.title = data?.prefix || '';
    this.note = data?.description || '';
    this.email = data?.emails?.[0]?.email || '';
    this.phone = data?.phoneNumbers?.[0]?.rawPhoneNumber || '';
    this.phoneId = data?.phoneNumbers?.[0]?.id || ''; // auxiliary field
    this.timezone = '';
    this.companyName = '';
  }

  static toForm(data: Contact) {
    return new ContactFormDto(data);
  }

  static toDto(payload: ContactForm): UpdateContactMutationVariables {
    return {
      input: {
        id: payload.id,
        firstName: payload.name,
        description: payload.note,
        prefix: payload.title,
      },
    };
  }
}
