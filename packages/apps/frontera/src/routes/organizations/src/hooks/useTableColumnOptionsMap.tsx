import { useMemo } from 'react';

import { match } from 'ts-pattern';

import { TableViewType } from '@graphql/types';
import {
  contractsMap,
  contactsOptionsMap,
  invoicesOptionsMap,
  renewalsOptionsMap,
  contactsHelperTextMap,
  invoicesHelperTextMap,
  renewalsHelperTextMap,
  contractsHelperTextMap,
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

        .with(TableViewType.Contracts, () => [
          contractsMap,
          contractsHelperTextMap,
        ])
        .otherwise(() => [organizationsOptionsMap, organizationsHelperTextMap]),
    [type],
  );
};
