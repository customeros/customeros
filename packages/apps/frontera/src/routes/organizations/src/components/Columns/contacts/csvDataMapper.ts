import { ContactStore } from '@store/Contacts/Contact.store';

import { DateTimeUtils } from '@utils/date.ts';
import { JobRole, ColumnViewType } from '@graphql/types';

export const csvDataMapper = {
  [ColumnViewType.ContactsAvatar]: (d: ContactStore) =>
    d?.value?.profilePhotoUrl,
  [ColumnViewType.ContactsName]: (d: ContactStore) => d.name,
  CONTACTS_FIRST_NAME: (d: ContactStore) => d.value.firstName,
  CONTACTS_LAST_NAME: (d: ContactStore) => d.value.lastName,
  [ColumnViewType.ContactsCity]: (d: ContactStore) =>
    d?.value?.locations?.[0]?.locality,
  [ColumnViewType.ContactsCountry]: (d: ContactStore) => d?.country,
  [ColumnViewType.ContactsEmails]: (d: ContactStore) =>
    d?.value?.emails?.map((e) => e.email).join('; '),
  [ColumnViewType.ContactsExperience]: () => null,
  [ColumnViewType.ContactsJobTitle]: (d: ContactStore) =>
    d?.value?.jobRoles?.[0]?.jobTitle,
  [ColumnViewType.ContactsLanguages]: (_d: ContactStore) => '',
  [ColumnViewType.ContactsLastInteraction]: (_d: ContactStore) => '',
  [ColumnViewType.ContactsSchools]: (_d: ContactStore) => '',
  [ColumnViewType.ContactsSkills]: (_d: ContactStore) => '',
  [ColumnViewType.ContactsTimeInCurrentRole]: (d: ContactStore) => {
    const jobRole = d.value.jobRoles?.find((role: JobRole) => {
      return role?.endedAt !== null;
    });

    if (!jobRole?.startedAt) return '';

    return DateTimeUtils.timeAgo(jobRole.startedAt);
  },
  [ColumnViewType.ContactsLinkedin]: (d: ContactStore) => {
    return d?.value?.socials.find((e) => e?.url?.includes('linkedin'))?.url;
  },
  [ColumnViewType.ContactsLinkedinFollowerCount]: (d: ContactStore) =>
    d?.value?.socials.find((e) => e?.url?.includes('linkedin'))?.followersCount,
  [ColumnViewType.ContactsOrganization]: (d: ContactStore) =>
    d?.value?.organizations?.content?.[0]?.name,
  [ColumnViewType.ContactsPersona]: (d: ContactStore) =>
    `${(d?.value?.tags ?? [])?.map((e) => e.name).join('; ')}`?.trim(),
  [ColumnViewType.ContactsTags]: (d: ContactStore) =>
    `${(d?.value?.tags ?? [])?.map((e) => e.name).join('; ')}`?.trim(),
  [ColumnViewType.ContactsPhoneNumbers]: (d: ContactStore) =>
    d.value?.phoneNumbers?.map((e) => e?.e164)?.join('; '),
  [ColumnViewType.ContactsConnections]: (data: ContactStore) => {
    return data.connectedUsers?.map((e) => e?.name)?.join('; ');
  },
  [ColumnViewType.ContactsRegion]: (data: ContactStore) => {
    return data.value?.locations?.[0]?.region;
  },
};
