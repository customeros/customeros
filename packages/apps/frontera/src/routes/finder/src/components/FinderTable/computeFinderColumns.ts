import { match } from 'ts-pattern';
import { RootStore } from '@store/root';
import { ColumnDef } from '@tanstack/react-table';
import { TableViewDefsStore } from '@store/TableViewDefs/TableViewDefs.store';

import { TableViewDef, TableViewType } from '@graphql/types';

import { getFlowColumnsConfig } from '../Columns/flows';
import { getContactColumnsConfig } from '../Columns/contacts';
import { getInvoiceColumnsConfig } from '../Columns/invoices';
import { getContractColumnsConfig } from '../Columns/contracts';
import { getOpportunityColumnsConfig } from '../Columns/opportunities';
import { getOrganizationColumnsConfig } from '../Columns/organizations';

interface ComputeFinderColumnsOptions {
  tableType: TableViewType;
  currentPreset: string | null;
}

export const computeFinderColumns = (
  store: RootStore,
  options: ComputeFinderColumnsOptions,
) => {
  const { currentPreset, tableType } = options;

  if (!currentPreset) return [];

  const tableViewDefStore = store.tableViewDefs.getById(currentPreset);
  const tableViewDef = tableViewDefStore?.value;
  const parseColumns = makeColumnParser(store);

  const organizationsColumns = parseColumns(
    'organizationsPreset',
    getOrganizationColumnsConfig,
  );
  const contactsColumns = parseColumns(
    'contactsPreset',
    getContactColumnsConfig,
  );
  const contractsColumns = parseColumns(
    'contractsPreset',
    getContractColumnsConfig,
  );
  const opportunitiesColumns = parseColumns(
    'opportunitiesPreset',
    getOpportunityColumnsConfig,
  );
  const pastInvoicesColumns = parseColumns(
    'pastInvoicesPreset',
    getInvoiceColumnsConfig,
  );
  const upcomingInvoicesColumns = parseColumns(
    'upcomingInvoicesPreset',
    getInvoiceColumnsConfig,
  );
  const flowsColumns = parseColumns('flowsPreset', getFlowColumnsConfig);

  return match(tableType)
    .with(TableViewType.Organizations, () =>
      match(currentPreset)
        .with(
          store.tableViewDefs.organizationsPreset ?? '',
          () => organizationsColumns,
        )
        .otherwise(() => getOrganizationColumnsConfig(tableViewDef)),
    )
    .with(TableViewType.Contacts, () =>
      match(currentPreset)
        .with(store.tableViewDefs.contactsPreset ?? '', () => contactsColumns)
        .otherwise(() => getContactColumnsConfig(tableViewDef)),
    )
    .with(TableViewType.Contracts, () =>
      match(currentPreset)
        .with(store.tableViewDefs.contractsPreset ?? '', () => contractsColumns)
        .otherwise(() => getContractColumnsConfig(tableViewDef)),
    )
    .with(TableViewType.Opportunities, () =>
      match(currentPreset)
        .with(
          store.tableViewDefs.opportunitiesPreset ?? '',
          () => opportunitiesColumns,
        )
        .otherwise(() => getOpportunityColumnsConfig(tableViewDef)),
    )
    .with(TableViewType.Invoices, () =>
      match(currentPreset)
        .with(
          store.tableViewDefs.pastInvoicesPreset ?? '',
          () => pastInvoicesColumns,
        )
        .with(
          store.tableViewDefs.upcomingInvoicesPreset ?? '',
          () => upcomingInvoicesColumns,
        )
        .otherwise(() => getInvoiceColumnsConfig(tableViewDef)),
    )
    .with(TableViewType.Flow, () =>
      match(currentPreset)
        .with(store.tableViewDefs.flowsPreset ?? '', () => flowsColumns)
        .otherwise(() => getFlowColumnsConfig(tableViewDef)),
    )
    .otherwise(() => []);
};

function makeColumnParser(store: RootStore) {
  return function (
    presetKey: keyof TableViewDefsStore,
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    getColumnConfig: (viewDef: TableViewDef) => ColumnDef<any, any>[],
  ) {
    const presetId = store.tableViewDefs[presetKey];

    if (!presetId || typeof presetId !== 'string') return [];

    const viewDef = store.tableViewDefs.getById(presetId);

    if (!viewDef) return [];

    return getColumnConfig(viewDef.value);
  };
}
