import { match } from 'ts-pattern';
import { Store } from '@store/store.ts';

import { Contact, ColumnViewType } from '@graphql/types';

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

    .otherwise(() => (_row: Store<Contact>) => false);
