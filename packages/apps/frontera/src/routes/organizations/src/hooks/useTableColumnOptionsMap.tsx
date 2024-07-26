import { useMemo } from 'react';

import { match } from 'ts-pattern';

import { TableViewType } from '@graphql/types';
import {
  contactsOptionsMap,
  invoicesOptionsMap,
  renewalsOptionsMap,
  contactsHelperTextMap,
  invoicesHelperTextMap,
  renewalsHelperTextMap,
  organizationsOptionsMap,
  organizationsHelperTextMap,
} from '@shared/components/ViewSettings/EditColumns/columnOptions.ts';

export const useTableColumnOptionsMap = (type?: TableViewType) => {
  return useMemo(
    () =>
      match(type)
        .with(TableViewType.Contacts, () => [
          contactsOptionsMap,
          contactsHelperTextMap,
        ])
        .with(TableViewType.Invoices, () => [
          invoicesOptionsMap,
          invoicesHelperTextMap,
        ])
        .with(TableViewType.Renewals, () => [
          renewalsOptionsMap,
          renewalsHelperTextMap,
        ])
        .otherwise(() => [organizationsOptionsMap, organizationsHelperTextMap]),
    [type],
  );
};
