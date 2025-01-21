import { inPlaceSort } from 'fast-sort';
import { action, reaction } from 'mobx';

import type { Contact } from '../Contact.dto';

import { indexAndSearch } from './util';
import { getContactSortFn } from './sortFns';
import { getContactFilterFns } from './filterFns';
import { ContactsStore } from '../Contacts.store';

export class ContactsView {
  constructor(private store: ContactsStore) {
    reaction(() => {
      const preset = this.store.root.tableViewDefs.contactsPreset;

      return preset ? this.store.getSearchTermByView(preset) : '';
    }, this.update);
    reaction(() => {
      const preset = this.store.root.tableViewDefs.contactsPreset;

      return preset ? this.store.availableCounts.get(preset) : 0;
    }, this.update);
    reaction(() => this.store.value.size, this.update);
    reaction(() => this.store.version, this.update);
    reaction(() => {
      const preset = this.store.root.tableViewDefs.contactsPreset;

      return this.store.cursors.get(preset!);
    }, this.update);
    reaction(
      () => {
        const preset = this.store.root.tableViewDefs.contactsPreset;

        if (!preset) return '';

        const viewDef = this.store.root.tableViewDefs.getById(preset);

        const columns = JSON.stringify(viewDef?.value.columns);

        return `${viewDef?.value.filters ?? ''}-${
          viewDef?.value.defaultFilters ?? ''
        }-${viewDef?.value.sorting}-${columns}`;
      },
      () => this.store.search(this.store.root.tableViewDefs.contactsPreset!),
    );
    reaction(() => {
      const preset = this.store.root.tableViewDefs?.contactsPreset;

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
    const preset = this.store.root.tableViewDefs.contactsPreset;

    if (!preset) return;

    const viewDef = this.store.root.tableViewDefs.getById(preset);

    if (!viewDef) return;

    const defaultFilters = getContactFilterFns(viewDef.getDefaultFilters());
    const activeFilters = getContactFilterFns(viewDef.getFilters());
    const sorting = JSON.parse(viewDef.value.sorting);

    this.store.setView(preset, (data) => {
      const columnId = sorting?.id as string;
      const isDesc = sorting?.desc as boolean;

      const filteredIdsWithSortValues = (data as Contact[]).reduce(
        (acc, curr) => {
          if (!curr) return acc;

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
