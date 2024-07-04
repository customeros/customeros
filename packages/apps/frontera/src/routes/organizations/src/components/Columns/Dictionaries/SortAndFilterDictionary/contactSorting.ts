import { match } from 'ts-pattern';
import { Store } from '@store/store.ts';
import countries from '@assets/countries/countries.json';

import { Social, Contact, JobRole, ColumnViewType } from '@graphql/types';

export const getContactColumnSortFn = (columnId: string) =>
  match(columnId)
    .with(ColumnViewType.ContactsName, () => (row: Store<Contact>) => {
      return row.value?.name?.trim().toLowerCase() || null;
    })
    .with(
      ColumnViewType.ContactsOrganization,
      () => (row: Store<Contact>) =>
        row.value?.organizations?.content?.[0]?.name
          ?.trim()
          .toLocaleLowerCase() || null,
    )
    .with(
      ColumnViewType.ContactsCity,
      () => (row: Store<Contact>) =>
        row.value?.locations?.[0]?.locality?.trim().toLowerCase() || null,
    )
    .with(ColumnViewType.ContactsPersona, () => (row: Store<Contact>) => {
      return row.value?.tags?.[0]?.name?.trim().toLowerCase() || null;
    })
    .with(
      ColumnViewType.ContactsLinkedinFollowerCount,
      () => (row: Store<Contact>) => {
        return row.value.socials.find((e: Social) =>
          e?.url?.includes('linkedin'),
        )?.followersCount;
      },
    )
    .with(ColumnViewType.ContactsJobTitle, () => (row: Store<Contact>) => {
      return row.value.jobRoles
        .find((e: JobRole) => e.endedAt === null)
        ?.jobTitle?.toLowerCase();
    })
    .with(ColumnViewType.ContactsCountry, () => (row: Store<Contact>) => {
      const countryName = countries.find(
        (d) =>
          d.alpha2 === row.value.locations?.[0]?.countryCodeA2?.toLowerCase(),
      );

      return countryName?.name?.toLowerCase() || null;
    })
    .otherwise(() => (_row: Store<Contact>) => false);
