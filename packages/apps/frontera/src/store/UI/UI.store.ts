import { makeAutoObservable } from 'mobx';

import { toastError, toastSuccess } from '@ui/presentation/Toast';

export class UIStore {
  isSearching: string | null = null;
  isFilteringTable: boolean = false;
  isEditingTableCell: boolean = false;
  dirtyEditor: string | null = null;
  activeConfirmation: string | null = null;

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

  confirmAction(id: string) {
    this.activeConfirmation = id;
  }

  clearConfirmAction() {
    this.activeConfirmation = null;
  }
}
