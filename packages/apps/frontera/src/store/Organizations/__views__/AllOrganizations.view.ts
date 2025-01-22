import { inPlaceSort } from 'fast-sort';
import { action, reaction } from 'mobx';

import type { Organization } from '../Organization.dto';

import { indexAndSearch } from './util';
import { getOrganizationSortFn } from './sortFns';
import { getOrganizationFilterFns } from './filterFns';
import { OrganizationsStore } from '../Organizations.store';

/* TODO: Explore implementing a push based strategy for this functionality.
 * I.e. we could call the update() method imperatively from within the TableViewDef entity
 * when a filter/sort rule is updated. This would avoid the need for the reaction and make
 * the execution more predictable.
 *
 * Subjects that could push this update are:
 *  - TableViewDef.appendDefaultFilter
 *  - TableViewDef.appendFilter
 *  - TableViewDef.setDefaultFilters
 *  - TableViewDef.removeFilter
 *  - TableViewDef.removeFilters
 *  - TableViewDef.toggleFilter
 *  - TableViewDef.setSorting
 *  - TableViewDef.setFilter
 *  - TableViewDef.setPropertyFilter
 *  - TableViewDef.orderColumnsByVisibility?
 *  - TableViewDef.setColumnName
 *  - TableViewDef.setColumnSize
 *
 *  - Organizations.search
 *  - Organizations.retrieve
 */
export class AllOrganizationsView {
  constructor(private store: OrganizationsStore) {
    // when new entites are added, removed or updated(by bumping the store version)
    // -> re-compute the view eagerly
    reaction(() => {
      const preset = this.store.root.tableViewDefs.organizationsPreset;

      return preset ? this.store.availableCounts.get(preset) : 0;
    }, this.update);
    reaction(() => this.store.value.size, this.update);
    reaction(() => this.store.version, this.update);

    // when cursor is updated by loading next chunk (via scrolling the table)
    reaction(() => {
      const preset = this.store.root.tableViewDefs.organizationsPreset;

      return this.store.cursors.get(preset!);
    }, this.update);

    // when the table view def is changed by updating a filter/sort rule
    // -> do an async search to fetch a new list of ids
    reaction(
      () => {
        const preset = this.store.root.tableViewDefs.organizationsPreset;

        if (!preset) return '';

        const viewDef = this.store.root.tableViewDefs.getById(preset);

        const columns = JSON.stringify(viewDef?.value.columns);

        return `${viewDef?.value.filters ?? ''}-${
          viewDef?.value.defaultFilters ?? ''
        }-${viewDef?.value.sorting}-${columns}`;
      },
      () =>
        this.store.search(this.store.root.tableViewDefs.organizationsPreset!),
    );

    // when the table view def is changed by updating a filter/sort rule
    // -> filter + sort the current entities in the store to optimistically update the table view
    reaction(() => {
      const preset = this.store.root.tableViewDefs.organizationsPreset;

      if (!preset) return '';

      const viewDef = this.store.root.tableViewDefs.getById(preset);
      const columns = JSON.stringify(viewDef?.value.columns);

      return `${viewDef?.value.filters ?? ''}-${
        viewDef?.value.defaultFilters ?? ''
      }-${viewDef?.value.sorting}-${columns}`;
    }, this.update);
  }

  @action
  public update = () => {
    const preset = this.store.root.tableViewDefs.organizationsPreset;

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
