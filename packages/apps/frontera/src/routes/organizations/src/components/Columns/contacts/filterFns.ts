import { match } from 'ts-pattern';
import { FilterItem } from '@store/types.ts';
import { ContactStore } from '@store/Contacts/Contact.store.ts';

import {
  Tag,
  Filter,
  Social,
  ColumnViewType,
  ComparisonOperator,
} from '@graphql/types';

const getFilterFn = (filter: FilterItem | undefined | null) => {
  const noop = (_row: ContactStore) => true;
  if (!filter) return noop;

  return match(filter)
    .with({ property: 'STAGE' }, (filter) => (row: ContactStore) => {
      const filterValues = filter?.value;

      if (!filterValues || !row.value?.organizations.content.length) {
        return false;
      }
      const hasOrgWithMatchingStage = row.value?.organizations.content.every(
        (o) => {
          const stage = row.root?.organizations?.value.get(o.metadata.id)?.value
            ?.stage;

          return filterValues.includes(stage);
        },
      );

      return hasOrgWithMatchingStage;
    })
    .with({ property: 'RELATIONSHIP' }, (filter) => (row: ContactStore) => {
      const filterValues = filter?.value;
      if (!filterValues || !row.value?.organizations.content.length) {
        return false;
      }
      const hasOrgWithMatchingRelationship =
        row.value?.organizations.content.every((o) => {
          const stage = row.root?.organizations?.value.get(o.metadata.id)?.value
            ?.relationship;

          return filterValues.includes(stage);
        });

      return hasOrgWithMatchingRelationship;
    })
    .with(
      { property: ColumnViewType.ContactsName },
      (filter) => (row: ContactStore) => {
        const filterValue = filter?.value;
        if (!filter.active) return true;

        if (!row.value?.name?.length && filter.includeEmpty) return true;
        if (!filterValue || !row.value?.name?.length) return false;

        return filterValue.toLowerCase().includes(row.value.name.toLowerCase());
      },
    )
    .with(
      { property: ColumnViewType.ContactsOrganization },
      (filter) => (row: ContactStore) => {
        const filterValues = filter?.value;

        if (!filter.active) return true;

        const orgs = row.value?.organizations?.content?.map((o) =>
          o.name.toLowerCase().trim(),
        );

        return orgs?.some((e) => e.includes(filterValues));
      },
    )
    .with(
      { property: ColumnViewType.ContactsEmails },
      (filter) => (row: ContactStore) => {
        const filterValues = filter?.value;
        if (!filter.active) return true;

        if (!filterValues) return false;

        return row.value?.emails?.some(
          (e) => e.email && filterValues.includes(e.email.toLowerCase().trim()),
        );
      },
    )
    .with({ property: 'EMAIL_VERIFIED' }, (filter) => (row: ContactStore) => {
      const filterValue = filter?.value;

      if (!filter.active) return true;
      if (row.value?.emails?.length === 0) return false;

      return row.value?.emails?.every((e) => {
        const { validated, isReachable, isValidSyntax } =
          e.emailValidationDetails;

        if (filterValue === 'verified') {
          return validated && isReachable !== 'invalid' && isValidSyntax;
        }

        return !validated || isReachable === 'invalid' || !isValidSyntax;
      });
    })
    .with(
      { property: ColumnViewType.ContactsPhoneNumbers },
      (filter) => (row: ContactStore) => {
        const filterValue = filter?.value;
        if (!filter.active) return true;

        if (!row.value?.phoneNumbers?.length && filter.includeEmpty)
          return true;
        if (!filterValue) return false;

        return row.value?.phoneNumbers?.[0]?.e164?.includes(filterValue);
      },
    )

    .with(
      { property: ColumnViewType.ContactsLinkedin },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        // specific logic for linkedin
        const linkedInUrl = row.value.socials?.find(
          (v: { id: string; url: string }) => v.url.includes('linkedin'),
        )?.url;
        const linkedinAlias = row.value.socials?.find(
          (v: { id: string; url: string }) => v.url.includes('linkedin'),
        )?.alias;

        if (!linkedInUrl && filter.includeEmpty) return true;

        return (
          (linkedInUrl && linkedInUrl.includes(filterValue)) ||
          (linkedinAlias && linkedinAlias.includes(filterValue))
        );
      },
    )
    .with(
      { property: ColumnViewType.ContactsCity },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;
        const cities = row.value.locations?.map((l) => l.locality);

        if (!cities.length && filter.includeEmpty) return true;

        return row.value.locations?.some((l) =>
          l?.locality?.includes(filterValue),
        );
      },
    )
    .with(
      { property: ColumnViewType.ContactsPersona },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;
        const tags = row.value.tags?.map((l: Tag) => l.name);
        if (!filter.value?.length) return true;

        if (!tags?.length) return false;

        return filter.value.some((f: string) => tags.includes(f));
      },
    )
    .with(
      { property: ColumnViewType.ContactsLinkedinFollowerCount },
      (filter) => (row: ContactStore) => {
        if (!filter.active) return true;
        const filterValue = filter?.value;

        const operator = filter.operation;
        const followers = row.value?.socials?.find((e: Social) =>
          e?.url?.includes('linkedin'),
        )?.followersCount;

        if (operator === ComparisonOperator.Lte) {
          return followers <= filterValue[0];
        }
        if (operator === ComparisonOperator.Gte) {
          return followers >= filterValue[0];
        }

        if (operator === ComparisonOperator.Between) {
          const filterValue = filter?.value?.map(Number) as number[];

          return followers >= filterValue[0] && followers <= filterValue[1];
        }

        return true;
      },
    )

    .otherwise(() => noop);
};

export const getContactFilterFns = (filters: Filter | null) => {
  if (!filters || !filters.AND) return [];

  const data = filters?.AND;

  return data.map(({ filter }) => getFilterFn(filter));
};
