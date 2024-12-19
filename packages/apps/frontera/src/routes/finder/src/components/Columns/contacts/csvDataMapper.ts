import type { Contact } from '@store/Contacts/Contact.dto';

import { DateTimeUtils } from '@utils/date.ts';
import { ColumnViewType } from '@graphql/types';

export const csvDataMapper = {
  [ColumnViewType.ContactsAvatar]: (d: Contact) => d?.value?.profilePhotoUrl,
  [ColumnViewType.ContactsName]: (d: Contact) => d.name,
  CONTACTS_FIRST_NAME: (d: Contact) => d.value.firstName,
  CONTACTS_LAST_NAME: (d: Contact) => d.value.lastName,
  [ColumnViewType.ContactsCity]: (d: Contact) =>
    d?.value?.locations?.[0]?.locality,
  [ColumnViewType.ContactsCountry]: (d: Contact) => d?.country,
  [ColumnViewType.ContactsPrimaryEmail]: (d: Contact) =>
    d?.value?.emails
      ?.filter((e) => e.primary)
      .map((e) => e.email)
      .join('; '),

  [ColumnViewType.ContactsExperience]: () => null,
  [ColumnViewType.ContactsJobTitle]: (d: Contact) =>
    d.value.primaryOrganizationJobRoleTitle,

  [ColumnViewType.ContactsTimeInCurrentRole]: (d: Contact) => {
    const jobRole = d.value.primaryOrganizationJobRoleEndDate
      ? d.value.primaryOrganizationJobRoleEndDate
      : '';

    if (!jobRole?.startedAt) return '';

    return DateTimeUtils.timeAgo(jobRole.startedAt);
  },
  [ColumnViewType.ContactsLinkedin]: (d: Contact) => {
    return d?.value?.linkedInUrl;
  },
  [ColumnViewType.ContactsFlows]: (d: Contact) =>
    d?.flows?.map((e) => e.value.name).join('; '),
  [ColumnViewType.ContactsLinkedinFollowerCount]: (d: Contact) =>
    d?.value?.linkedInFollowerCount,
  [ColumnViewType.ContactsOrganization]: (d: Contact) =>
    d?.value?.primaryOrganizationName,
  [ColumnViewType.ContactsPersona]: (d: Contact) =>
    `${(d?.value?.tags ?? [])?.map((e) => e.name).join('; ')}`?.trim(),
  [ColumnViewType.ContactsTags]: (d: Contact) =>
    `${(d?.value?.tags ?? [])?.map((e) => e.name).join('; ')}`?.trim(),
  // [ColumnViewType.ContactsPhoneNumbers]: (d: Contact) =>
  //   d.value?.phoneNumbers?.map((e) => e?.e164)?.join('; '),
  [ColumnViewType.ContactsConnections]: (data: Contact) => {
    return data.connectedUsers?.map((e) => e?.name)?.join('; ');
  },
  [ColumnViewType.ContactsRegion]: (data: Contact) => {
    return data.value?.locations?.[0]?.region;
  },
};
