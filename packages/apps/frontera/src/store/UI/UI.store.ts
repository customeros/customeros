import { makeAutoObservable } from 'mobx';

import { toastError, toastSuccess } from '@ui/presentation/Toast';

export class UIStore {
  isSearching: string | null = null;
  isEditingTableCell: boolean = false;

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

  setIsEditingTableCell(value: boolean) {
    this.isEditingTableCell = value;
  }
}
