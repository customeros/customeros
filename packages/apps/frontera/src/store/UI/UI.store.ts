import type { Channel } from 'phoenix';
import type { RootStore } from '@store/root';

import { Persister } from '@store/persister';
import { Transport } from '@store/transport';
import { when, makeAutoObservable } from 'mobx';
import { SystemSyncPacket } from '@store/types';
import { FlowStepCommandMenuStore } from '@store/UI/FlowStepCommandMenu.store.ts';

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
  contactPreviewCardOpen: boolean = false;
  movedIcpOrganization: number = 0;
  focusRow: number | string | null = null;
  commandMenu = new CommandMenuStore();
  selectionId: number | null = null;
  flowCommandMenu = new FlowStepCommandMenuStore();
  isSystemNotificationOpen = false;
  private channel?: Channel;
  private activeConfirmationCallback: () => void = () => {};

  constructor(private root: RootStore, private transport: Transport) {
    makeAutoObservable(this);
    this.toastSuccess = this.toastSuccess.bind(this);
    this.purgeLocalData = this.purgeLocalData.bind(this);

    when(
      () => !!this.root?.session?.value?.tenant && !this.root.demoMode,
      async () => {
        try {
          await this.initChannelConnection();
        } catch (e) {
          console.error(e);
        }
      },
    );
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

  setContactPreviewCardOpen(value: boolean) {
    this.contactPreviewCardOpen = value;
  }

  setFocusRow(value: number | string | null) {
    this.focusRow = value;
  }

  setSelectionId(value: number | null) {
    this.selectionId = value;
  }

  purgeLocalData() {
    Persister.attemptPurge({ force: true });
    this.toastSuccess(
      'Re-sync done. Refreshing in order to continue...',
      're-sync',
    );
    setTimeout(() => window.location.reload(), 2000);
  }

  showSystemNotification() {
    this.isSystemNotificationOpen = true;
  }

  private async initChannelConnection() {
    try {
      const connection = await this.transport.join('System', 'all', 0, true);

      if (!connection) return;

      this.channel = connection.channel;
      this.subscribe();
    } catch (e) {
      console.error(e);
    }
  }

  private subscribe() {
    if (!this.channel || this.root.demoMode) return;

    this.channel.on('sync_group_packet', (packet: SystemSyncPacket) => {
      if (packet.action === 'NEW_VERSION_AVAILABLE') {
        this.showSystemNotification();
      }
    });
  }
}
