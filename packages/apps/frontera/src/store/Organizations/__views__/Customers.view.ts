import { reaction } from 'mobx';
import { inPlaceSort } from 'fast-sort';

import type { Organization } from '../Organization.dto';

import { indexAndSearch } from './util';
import { getOrganizationSortFn } from './sortFns';
import { getOrganizationFilterFns } from './filterFns';
import { OrganizationsStore } from '../Organizations.store';

// TODO: Cache filtered and sorted results for faster subsequent access
export class CustomersView {
  constructor(private store: OrganizationsStore) {
    reaction(() => {
      const preset = this.store.root.tableViewDefs.customersPreset;

      return preset ? this.store.getSearchTermByView(preset) : '';
    }, this.update);
    reaction(() => {
      const preset = this.store.root.tableViewDefs.customersPreset;

      return preset ? this.store.availableCounts.get(preset) : 0;
    }, this.update);
    reaction(() => this.store.value.size, this.update);
    reaction(() => this.store.version, this.update);
    reaction(() => {
      const preset = this.store.root.tableViewDefs.customersPreset;

      return this.store.cursors.get(preset!);
    }, this.update);
    reaction(
      () => {
        const preset = this.store.root.tableViewDefs.customersPreset;

        if (!preset) return '';

        const viewDef = this.store.root.tableViewDefs.getById(preset);

        const columns = JSON.stringify(viewDef?.value.columns);

        return `${viewDef?.value.filters ?? ''}-${
          viewDef?.value.defaultFilters ?? ''
        }-${viewDef?.value.sorting}-${columns}`;
      },
      () => this.store.search(this.store.root.tableViewDefs.customersPreset!),
    );
    reaction(() => {
      const preset = this.store.root.tableViewDefs.customersPreset;

      if (!preset) return '';

      const viewDef = this.store.root.tableViewDefs.getById(preset);
      const columns = JSON.stringify(viewDef?.value.columns);

      return `${viewDef?.value.filters ?? ''}-${
        viewDef?.value.defaultFilters ?? ''
      }-${viewDef?.value.sorting}-${columns}`;
    }, this.update);
  }

  public update = () => {
    const preset = this.store.root.tableViewDefs.customersPreset;

    if (!preset) return;

    const viewDef = this.store.root.tableViewDefs.getById(preset);

    if (!viewDef) return;

    const defaultFilters = getOrganizationFilterFns(
      viewDef.getDefaultFilters(),
    );
    const activeFilters = getOrganizationFilterFns(viewDef.getFilters());
    const sorting = JSON.parse(viewDef.value.sorting);

    this.store.setView(preset, (data) => {
      const columnId = sorting?.id as string;
      const isDesc = sorting?.desc as boolean;

      const filteredIdsWithSortValues = (data as Organization[]).reduce(
        (acc, curr) => {
          if (!curr) return acc;

          if (
            defaultFilters.every((fn) => fn(curr)) &&
            activeFilters.every((fn) => fn(curr))
          ) {
            const sortValue = getOrganizationSortFn(columnId)(curr);

            acc.push({ record: curr, sortValue });
          }

          return acc;
        },
        [] as {
          record: Organization;
          sortValue: string | number | boolean | Date | null | undefined;
        }[],
      );

      let sorted = inPlaceSort(filteredIdsWithSortValues)
        [isDesc ? 'desc' : 'asc']((entry) => entry.sortValue)
        .map((entry) => entry.record);

      const searchTerm = this.store.getSearchTermByView(preset);

      if (searchTerm) {
        sorted = indexAndSearch(sorted, searchTerm);
      }

      return sorted;
    });
  };
}
