import { action, autorun } from 'mobx';
import { inPlaceSort } from 'fast-sort';

import type { Contact } from '../Contact.dto';

import { indexAndSearch } from './util';
import { getContactSortFn } from './sortFns';
import { getContactFilterFns } from './filterFns';
import { ContactsStore } from '../Contacts.store';

export class FlowContactsView {
  private cachedCombos = new Map<string, string>();
  private cachedSearchCombos = new Map<string, string>();

  constructor(private store: ContactsStore) {
    autorun(() => {
      if (!this.store.root.flowParticipants?.value) return;

      const flowParticipantIds = Array.from(
        this.store.root.flowParticipants?.value.values(),
      ).reduce((acc, curr) => {
        if (!this.store.value.has(curr?.value?.entityId)) {
          acc.push(curr?.value?.entityId);
        }

        return acc;
      }, [] as string[]);

      if (flowParticipantIds.length) {
        this.store.retrieve(flowParticipantIds);
      }
    });

    autorun(() => {
      const flowContactsViewDefs =
        this.store.root.tableViewDefs.flowContactsPresets;

      const dataSize = this.store.value.size;
      const dataVersion = this.store.version;

      flowContactsViewDefs.forEach((viewDef) => {
        const preset = viewDef?.value.id;
        const searchTerm = this.store.getSearchTermByView(preset);
        const combo = [
          viewDef.value?.defaultFilters,
          viewDef.value?.filters,
          viewDef.value?.sorting,
          searchTerm,
          JSON.stringify(viewDef.value.columns),
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

      if (searchTerm) {
        sorted = indexAndSearch(sorted, searchTerm);
      }

      return sorted;
    });
  };

  // TODO: Continue exploring this approach to reconcile the entity in the view
  _reconcile(entityId: string, preset: string) {
    const entity = this.store.getById(entityId)!;

    const viewDef = this.store.root.tableViewDefs.getById(preset);

    if (!viewDef) return;

    const defaultFilters = getContactFilterFns(viewDef.getDefaultFilters());
    const activeFilters = getContactFilterFns(viewDef.getFilters());
    const sorting = JSON.parse(viewDef.value.sorting);
    const columnId = sorting?.id as string;
    const isDesc = sorting?.desc as boolean;

    if (
      defaultFilters.every((fn) => fn(entity)) &&
      activeFilters.every((fn) => fn(entity))
    ) {
      const currentView = this.store.views.get(preset);

      if (currentView) {
        currentView.unshift(entity);

        const sorted = inPlaceSort(
          currentView.map((e) => ({
            record: e,
            sortValue: getContactSortFn(columnId)(e),
          })),
        )
          [isDesc ? 'desc' : 'asc']((entry) => entry.sortValue)
          .map((entry) => entry.record);

        this.store.setView(preset, (_) => sorted);
      }
    }
  }
}
