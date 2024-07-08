import { Contact, ColumnViewType } from '@graphql/types';

export const csvDataMapper = {
  [ColumnViewType.ContactsAvatar]: (d: Contact) => d?.profilePhotoUrl,
  [ColumnViewType.ContactsName]: (d: Contact) =>
    d.name ?? `${d?.firstName ?? ''} ${d?.lastName ?? ''}`.trim(),
  [ColumnViewType.ContactsCity]: (d: Contact) => d?.locations?.[0]?.locality,
  [ColumnViewType.ContactsCountry]: (d: Contact) => d?.locations?.[0]?.country,

  [ColumnViewType.ContactsEmails]: (d: Contact) =>
    d?.emails?.map((e) => e.email).join('; '),
  [ColumnViewType.ContactsExperience]: () => null,
  [ColumnViewType.ContactsJobTitle]: (d: Contact) => d?.jobRoles?.[0]?.jobTitle,
  [ColumnViewType.ContactsLanguages]: (_d: Contact) => '',
  [ColumnViewType.ContactsLastInteraction]: (_d: Contact) => '',
  [ColumnViewType.ContactsSchools]: (_d: Contact) => '',
  [ColumnViewType.ContactsSkills]: (_d: Contact) => '',
  [ColumnViewType.ContactsTimeInCurrentRole]: (_d: Contact) => '',
  [ColumnViewType.ContactsLinkedin]: (d: Contact) => {
    return d.socials.find((e) => e?.url?.includes('linkedin'))?.url;
  },
  [ColumnViewType.ContactsLinkedinFollowerCount]: (d: Contact) =>
    d.socials.find((e) => e?.url?.includes('linkedin'))?.followersCount,
  [ColumnViewType.ContactsOrganization]: (d: Contact) =>
    d?.organizations?.content?.[0]?.name,
  [ColumnViewType.ContactsPersona]: (d: Contact) =>
    `${(d.tags ?? [])?.map((e) => e.name).join('; ')}`?.trim(),
  [ColumnViewType.ContactsTags]: (d: Contact) =>
    `${(d.tags ?? [])?.map((e) => e.name).join('; ')}`?.trim(),
  [ColumnViewType.ContactsPhoneNumbers]: (data: Contact) =>
    data.phoneNumbers?.map((e) => e?.e164)?.join('; '),
};
