import { inPlaceSort } from 'fast-sort';
import { toJS, action, autorun } from 'mobx';

import type { Organization } from '../Organization.dto';

import { indexAndSearch } from './util';
import { getOrganizationSortFn } from './sortFns';
import { getOrganizationFilterFns } from './filterFns';
import { OrganizationsStore } from '../Organizations.store';

export class TeamViews {
  private cachedCombos = new Map<string, string>();
  private cachedSearchCombos = new Map<string, string>();

  constructor(private store: OrganizationsStore) {
    autorun(() => {
      const teamViewDefs = this.store.root.tableViewDefs.teamPresets;

      teamViewDefs.forEach((viewDef) => {
        const preset = viewDef?.value.id;
        const combo = [
          viewDef.value?.defaultFilters,
          viewDef.value?.filters,
          viewDef.value?.sorting,
          JSON.stringify(toJS(viewDef.value.columns)),
        ].join('-');

        if (this.cachedSearchCombos.get(preset) === combo) return;

        this.cachedSearchCombos.set(preset, combo);

        this.store.search(preset);
      });
    });

    autorun(() => {
      const teamViewDefs = this.store.root.tableViewDefs.teamPresets;

      const dataSize = this.store.value.size;
      const dataVersion = this.store.version;

      teamViewDefs.forEach((viewDef) => {
        const preset = viewDef?.value.id;
        const searchTerm = this.store.getSearchTermByView(preset);
        const combo = [
          viewDef.value?.defaultFilters,
          viewDef.value?.filters,
          viewDef.value?.sorting,
          searchTerm,
          JSON.stringify(toJS(viewDef.value.columns)),
          dataSize,
          dataVersion,
        ].join('-');

        if (this.cachedCombos.get(preset) === combo) return;

        this.cachedCombos.set(preset, combo);

        this.update(preset, searchTerm);
      });
    });
  }

  @action
  public update = (preset: string, searchTerm?: string) => {
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

      if (searchTerm) {
        sorted = indexAndSearch(sorted, searchTerm);
      }

      return sorted;
    });
  };
}
