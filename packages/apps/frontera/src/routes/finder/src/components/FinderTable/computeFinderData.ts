import { Params } from 'react-router-dom';

import Fuse from 'fuse.js';
import { match } from 'ts-pattern';
import { RootStore } from '@store/root';
import { inPlaceSort } from 'fast-sort';
import { SortingState } from '@tanstack/table-core';
import { TableViewDef } from '@store/TableViewDefs/TableViewDef.dto';

import { TableViewType } from '@graphql/types';

import { getFlowsFilterFns, getFlowsColumnSortFn } from '../Columns/flows';
import {
  getOpportunitiesSortFn,
  getOpportunityFilterFns,
} from '../Columns/opportunities';
import {
  getInvoicesSortFn,
  getInvoiceFilterFns,
  getInvoiceDefaultFilterFns,
} from '../Columns/invoices';
import {
  getContractSortFn,
  getContractFilterFns,
  getContractDefaultFilters,
} from '../Columns/contracts';

interface ComputeFinderDataOptions {
  searchTerm: string;
  sorting: SortingState;
  tableViewDef: TableViewDef | null;
  urlParams: Readonly<Params<string>>;
}

export const computeFinderData = (
  store: RootStore,
  options: ComputeFinderDataOptions,
) => {
  const { searchTerm, sorting, tableViewDef } = options;

  if (!tableViewDef) return [];

  const preset = tableViewDef.value.id;
  const tableType =
    tableViewDef?.value.tableType || TableViewType.Organizations;

  return match(tableType)
    .with(TableViewType.Organizations, () => {
      return store.organizations.getViewById(preset ?? '');
    })
    .with(TableViewType.Contacts, () => {
      return store.contacts.getViewById(preset ?? '');
    })
    .with(TableViewType.Contracts, () =>
      store.contracts?.toComputedArray((arr) => {
        const defaultFilters = getContractDefaultFilters(
          tableViewDef?.getDefaultFilters(),
        );

        const filters = getContractFilterFns(tableViewDef?.getFilters());

        if (defaultFilters) {
          arr = arr.filter((v) => defaultFilters.every((fn) => fn(v)));
        }

        if (filters) {
          arr = arr.filter((v) => filters.every((fn) => fn(v)));
        }

        if (tableType) {
          const columnId = sorting[0]?.id;
          const isDesc = sorting[0]?.desc;

          arr = inPlaceSort(arr)?.[isDesc ? 'desc' : 'asc'](
            getContractSortFn(columnId),
          );
        }

        if (searchTerm) {
          arr = new Fuse(arr, {
            keys: ['value.name'],
            threshold: 0.3,
            isCaseSensitive: false,
          })
            .search(removeAccents(searchTerm), { limit: 40 })
            .map((r) => r.item);
        }

        return arr;
      }),
    )
    .with(TableViewType.Invoices, () =>
      store.invoices.toComputedArray((arr) => {
        const defaultFilters = getInvoiceDefaultFilterFns(
          tableViewDef?.getDefaultFilters(),
        );
        const filters = getInvoiceFilterFns(tableViewDef?.getFilters());

        if (defaultFilters) {
          arr = arr.filter((v) => defaultFilters.every((fn) => fn(v)));
        }

        if (filters) {
          arr = arr.filter((v) => filters.every((fn) => fn(v)));
        }

        if (tableType) {
          const columnId = sorting[0]?.id;
          const isDesc = sorting[0]?.desc;

          arr = inPlaceSort(arr)?.[isDesc ? 'desc' : 'asc'](
            getInvoicesSortFn(columnId),
          );
        }

        if (searchTerm) {
          const normalizedSearchTerm = removeAccents(searchTerm);

          arr = new Fuse(arr, {
            keys: ['value.contract.contractName'],
            threshold: 0.3,
            isCaseSensitive: false,
          })
            .search(normalizedSearchTerm, { limit: 40 })
            .map((r) => r.item);
        }

        return arr;
      }),
    )
    .with(TableViewType.Flow, () =>
      store.flows.toComputedArray((arr) => {
        if (tableType !== TableViewType.Flow) return arr;

        const filters = getFlowsFilterFns(tableViewDef?.getFilters());

        if (filters) {
          arr = arr.filter((v) => filters.every((fn) => fn(v)));
        }

        if (tableType) {
          const columnId = sorting[0]?.id;
          const isDesc = sorting[0]?.desc;

          arr = inPlaceSort(arr)?.[isDesc ? 'desc' : 'asc'](
            getFlowsColumnSortFn(columnId),
          );
        }

        if (searchTerm) {
          const normalizedSearchTerm = removeAccents(searchTerm);

          arr = new Fuse(arr, {
            keys: ['value.name'],
            threshold: 0.3,
            isCaseSensitive: false,
          })
            .search(normalizedSearchTerm, { limit: 40 })
            .map((r) => r.item);
        }

        return arr.filter((e) => e.value.status !== 'ARCHIVED');
      }),
    )
    .with(TableViewType.Opportunities, () =>
      store.opportunities.toComputedArray((arr) => {
        if (tableType !== TableViewType.Opportunities) return arr;
        arr = arr.filter((opp) => opp.value.internalType === 'NBO');

        const filters = getOpportunityFilterFns(tableViewDef?.getFilters());

        if (filters) {
          arr = arr.filter((v) => filters.every((fn) => fn(v)));
        }

        if (tableType) {
          const columnId = sorting[0]?.id;
          const isDesc = sorting[0]?.desc;

          arr = inPlaceSort(arr)?.[isDesc ? 'desc' : 'asc'](
            getOpportunitiesSortFn(columnId),
          );
        }

        if (searchTerm) {
          const normalizedSearchTerm = removeAccents(searchTerm);

          arr = new Fuse(arr, {
            keys: ['value.name', 'organization.value.name', 'owner.name'],
            threshold: 0.3,
            isCaseSensitive: false,
          })
            .search(normalizedSearchTerm, { limit: 40 })
            .map((r) => r.item);
        }

        return arr;
      }),
    )
    .otherwise(() => []);
};

function removeAccents(str: string) {
  return str
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '');
}
