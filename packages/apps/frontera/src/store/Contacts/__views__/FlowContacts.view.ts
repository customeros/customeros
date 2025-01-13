import { action, autorun } from 'mobx';
import { inPlaceSort } from 'fast-sort';

import type { Contact } from '../Contact.dto';

import { indexAndSearch } from './util';
import { getContactSortFn } from './sortFns';
import { getContactFilterFns } from './filterFns';
import { ContactsStore } from '../Contacts.store';

// TODO: Cache filtered and sorted results for faster subsequent access
export class FlowContactsView {
  private cachedCombos = new Map<string, string>();

  constructor(private store: ContactsStore) {
    // autorun(() => {
    //   const flows = Array.from(this.store.root.flows.value);
    //   const preset = this.store.root.tableViewDefs.flowContactsPreset;

    //   if (!preset) return;

    //   flows.forEach(([id]) => {
    //     const viewDef = this.store.root.tableViewDefs.getById(preset);
    //     const combo = [
    //       id,
    //       viewDef?.value.defaultFilters,
    //       viewDef?.value.filters,
    //       viewDef?.value.sorting,
    //       JSON.stringify(viewDef?.value.columns),
    //     ].join('-');

    //     if (this.cachedCombos.get(id) === combo) return;

    //     this.cachedCombos.set(id, combo);

    //     this.store.search(preset);
    //   });
    // });

    autorun(() => {
      const flows = Array.from(this.store.root.flows.value);
      const preset = this.store.root.tableViewDefs.flowContactsPreset;

      if (!preset) return;

      // const dataSize = this.store.size;
      const dataVersion = this.store.version;

      flows.forEach(([id]) => {
        const viewDef = this.store.root.tableViewDefs.getById(preset);
        const combo = [
          id,
          viewDef?.value.defaultFilters,
          viewDef?.value.filters,
          viewDef?.value.sorting,
          JSON.stringify(viewDef?.value.columns),
          dataVersion,
          // dataSize,
        ].join('-');

        if (this.cachedCombos.get(id) === combo) return;

        this.cachedCombos.set(id, combo);

        this.update(id);
      });
    });

    // reaction(() => {
    //   const preset = this.store.root.tableViewDefs.flowContactsPreset;

    //   return preset ? this.store.getSearchTermByView(preset) : '';
    // }, this.update);
    // reaction(() => {
    //   const preset = this.store.root.tableViewDefs.flowContactsPreset;

    //   return preset ? this.store.availableCounts.get(preset) : 0;
    // }, this.update);
    // reaction(() => this.store.value.size, this.update);
    // reaction(() => this.store.version, this.update);
    // reaction(() => {
    //   const preset = this.store.root.tableViewDefs.flowContactsPreset;

    //   return this.store.cursors.get(preset!);
    // }, this.update);
    // reaction(
    //   () => {
    //     const preset = this.store.root.tableViewDefs.flowContactsPreset;

    //     if (!preset) return '';

    //     const viewDef = this.store.root.tableViewDefs.getById(preset);

    //     const columns = JSON.stringify(viewDef?.value.columns);

    //     return `${viewDef?.value.filters ?? ''}-${
    //       viewDef?.value.defaultFilters ?? ''
    //     }-${viewDef?.value.sorting}-${columns}-${this.store.size}`;
    //   },
    //   () =>
    //     this.store.search(this.store.root.tableViewDefs.flowContactsPreset!),
    // );
    // reaction(() => {
    //   const preset = this.store.root.tableViewDefs.flowContactsPreset;

    //   if (!preset) return '';

    //   const viewDef = this.store.root.tableViewDefs.getById(preset);
    //   const columns = JSON.stringify(viewDef?.value.columns);

    //   return `${viewDef?.value.filters ?? ''}-${
    //     viewDef?.value.defaultFilters ?? ''
    //   }-${viewDef?.value.sorting}-${columns}`;
    // }, this.update);
  }

  @action
  public update = (flowId: string) => {
    const preset = this.store.root.tableViewDefs.flowContactsPreset;

    if (!preset) return;

    const viewDef = this.store.root.tableViewDefs.getById(preset);

    if (!viewDef) return;

    const defaultFilters = getContactFilterFns(viewDef.getDefaultFilters());
    const activeFilters = getContactFilterFns(viewDef.getFilters(), flowId);
    const sorting = JSON.parse(viewDef.value.sorting);

    this.store.setView(flowId, (data) => {
      const columnId = sorting?.id as string;
      const isDesc = sorting?.desc as boolean;

      const filteredIdsWithSortValues = (data as Contact[]).reduce(
        (acc, curr) => {
          if (!curr || !flowId) return acc;
          if (!curr.flowsIds?.includes(flowId)) return acc;

          if (
            defaultFilters.every((fn) => fn(curr)) &&
            activeFilters.every((fn) => fn(curr))
          ) {
            const sortValue = getContactSortFn(columnId)(curr);

            acc.push({ record: curr, sortValue });
          }

          return acc;
        },
        [] as {
          record: Contact;
          sortValue:
            | number
            | boolean
            | Date
            | null
            | (string | undefined)[]
            | string
            | undefined;
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
