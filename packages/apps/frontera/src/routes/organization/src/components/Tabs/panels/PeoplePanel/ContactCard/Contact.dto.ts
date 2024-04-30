import { Social, Contact } from '@graphql/types';
import { SelectOption } from '@shared/types/SelectOptions';
import { UpdateContactMutationVariables } from '@organization/graphql/updateContact.generated';

export interface ContactForm {
  id: string;
  name: string;
  note: string;
  title: string;
  email: string;
  phone: string;
  roleId: string;
  phoneId: string;
  startedAt: string;
  role: SelectOption<string>[];
  timezone: SelectOption<string> | null;
  socials: Pick<Social, 'id' | 'url'>[];
}

const getNameFromEmail = (email: string) => {
  const name = email?.split('@')[0];

  return `${name}`?.trim()?.replace(/\b\w/g, (char) => char.toUpperCase());
};
export class ContactFormDto implements ContactForm {
  id: string; // auxiliary field
  name: string;
  role: SelectOption<string>[];
  roleId: string; // auxiliary field
  title: string;
  note: string;
  email: string;
  phone: string;
  phoneId: string; // auxiliary field
  timezone: SelectOption<string> | null;
  socials: Pick<Social, 'id' | 'url'>[];
  startedAt: string;

  constructor(data?: Partial<Contact> | null) {
    const nameFromEmail = data?.emails?.[0]?.email
      ? getNameFromEmail(data?.emails?.[0]?.email ?? '')
      : '';

    this.id = data?.id || ''; // auxiliary field
    this.name =
      data?.name ||
      `${data?.firstName} ${data?.lastName}`.trim() ||
      nameFromEmail;

    this.title = data?.jobRoles?.[0]?.jobTitle || '';
    this.roleId = data?.jobRoles?.[0]?.id || ''; // auxiliary field
    this.role = (() => {
      const _role = data?.jobRoles?.[0]?.description;
      if (!_role?.length) return [];

      return _role?.split(',').map((v) => ({ value: v, label: v })) || [];
    })();
    this.note = data?.description || '';
    this.email = data?.emails?.[0]?.email || '';
    this.phone = data?.phoneNumbers?.[0]?.rawPhoneNumber || '';
    this.phoneId = data?.phoneNumbers?.[0]?.id || ''; // auxiliary field
    this.timezone = data?.timezone
      ? { label: data?.timezone, value: data?.timezone }
      : null;
    this.socials = data?.socials || [];
    this.startedAt = data?.jobRoles?.[0]?.startedAt || '';
  }

  static toForm(data: Contact) {
    return new ContactFormDto(data);
  }

  static toDto(
    payload: Partial<
      Omit<
        ContactForm & { timezone?: string | null; description?: string | null },
        'id' | 'role'
      >
    >,
    id: string,
  ): UpdateContactMutationVariables {
    return {
      input: {
        ...payload,
        id,
        patch: true,
      },
    };
  }
}
