import { autorun } from 'mobx';
import { inPlaceSort } from 'fast-sort';

import type { Organization } from '../Organization.dto';

import { getOrganizationSortFn } from './sortFns';
import { getOrganizationFilterFns } from './filterFns';
import { OrganizationsStore } from '../Organizations.store';

// TODO: Cache filtered and sorted results for faster subsequent access
export class CustomView {
  constructor(private store: OrganizationsStore) {
    autorun(() => {
      const p = this.store.root.tableViewDefs.customPresets;

      p.forEach((v) => {
        const id = v.value.id;

        this.update(id);
      });
    });
  }

  public update = (preset: string) => {
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

      const sorted = inPlaceSort(filteredIdsWithSortValues)
        [isDesc ? 'desc' : 'asc']((entry) => entry.sortValue)
        .map((entry) => entry.record);

      // const splicedIds = sortedIds.splice(
      //   this.store.range[0],
      //   this.store.range[1] + 1,
      // );

      return sorted;
    });
  };
}
