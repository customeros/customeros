import type { Store } from '@store/active-record';

import { reaction } from 'mobx';
import { inPlaceSort } from 'fast-sort';

import type { Organization } from '../Organization';

import { getOrganizationSortFn } from './sortFns';
import { getOrganizationFilterFns } from './filterFns';

// TODO: Cache filtered and sorted results for faster subsequent access
export class AllOrganizationsView {
  constructor(private store: Store<Organization>) {
    reaction(() => this.store.size, this.update);
    reaction(() => this.store.version, this.update);
    reaction(() => {
      const preset = this.store.root.tableViewDefs.organizationsPreset;

      if (!preset) return '';

      const viewDef = this.store.root.tableViewDefs.getById(preset);

      return `${viewDef?.value.filters ?? ''}-${
        viewDef?.value.defaultFilters ?? ''
      }-${viewDef?.value.sorting}`;
    }, this.update);
  }

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

      const filteredIdsWithSortValues = data.reduce((acc, curr) => {
        if (!curr) return acc;

        if (
          defaultFilters.every((fn) => fn(curr)) &&
          activeFilters.every((fn) => fn(curr))
        ) {
          const sortValue = getOrganizationSortFn(columnId)(curr);

          acc.push({ id: this.store.options.getId(curr), sortValue });
        }

        return acc;
      }, [] as { id: string; sortValue: string | number | boolean | Date | null | undefined }[]);

      const sortedIds = inPlaceSort(filteredIdsWithSortValues)
        [isDesc ? 'desc' : 'asc']((entry) => entry.sortValue)
        .map((entry) => entry.id);

      const splicedIds = sortedIds.splice(
        this.store.range[0],
        this.store.range[1] + 1,
      );

      return splicedIds;
    });
  };
}
