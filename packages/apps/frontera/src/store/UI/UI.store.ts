import { makeAutoObservable } from 'mobx';

import { toastError, toastSuccess } from '@ui/presentation/Toast';

import { CommandMenuStore } from './CommandMenu.store';

export class UIStore {
  searchCount: number = 0;
  filteredTable: Array<unknown> = [];
  isSearching: string | null = null;
  isFilteringTable: boolean = false;
  isFilteringICP: boolean = false;
  isEditingTableCell: boolean = false;
  dirtyEditor: string | null = null;
  activeConfirmation: string | null = null;
  movedIcpOrganization: number = 0;
  commandMenu = new CommandMenuStore();
  private activeConfirmationCallback: () => void = () => {};

  constructor() {
    makeAutoObservable(this);
  }

  toastSuccess(text: string, id: string) {
    // redundant call to toastSuccess - should be refactored
    toastSuccess(text, id);
  }

  toastError(text: string, id: string) {
    // redundant call to toastError - should be refactored
    toastError(text, id);
  }

  setIsSearching(value: string | null) {
    this.isSearching = value;
  }

  setIsFilteringTable(value: boolean) {
    this.isFilteringTable = value;
  }

  setIsEditingTableCell(value: boolean) {
    this.isEditingTableCell = value;
  }

  setDirtyEditor(value: string | null) {
    this.dirtyEditor = value;
  }

  clearDirtyEditor() {
    this.dirtyEditor = null;
  }

  confirmAction(id: string, callback?: () => void) {
    this.activeConfirmation = id;
    callback && (this.activeConfirmationCallback = callback);
  }

  clearConfirmAction() {
    this.activeConfirmation = null;
    this.activeConfirmationCallback?.();
  }

  setIsFilteringICP(value: boolean) {
    this.isFilteringICP = value;
  }

  setSearchCount(value: number) {
    this.searchCount = value;
  }

  setFilteredTable(data: Array<unknown>) {
    this.filteredTable = data;
  }

  setMovedIcpOrganization(value: number) {
    this.movedIcpOrganization = value;
  }
}
