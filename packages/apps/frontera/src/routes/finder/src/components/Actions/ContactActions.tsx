import { useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import { useKeys, useKeyBindings } from 'rooks';
import { Contact } from '@store/Contacts/Contact.dto';
import { CommandMenuType } from '@store/UI/CommandMenu.store.ts';
import { EditEmailCase } from '@domain/usecases/command-menu/edit-email.usecase';

import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import { TableInstance } from '@ui/presentation/Table';

import { SharedTableActions } from './components/SharedActions';

interface TableActionsProps {
  selection: string[];
  focusedId?: string | null;
  isCommandMenuOpen: boolean;
  table: TableInstance<Contact>;
  enableKeyboardShortcuts?: boolean;
}
const editEmailUseCase = EditEmailCase.getInstance();

export const ContactTableActions = observer(
  ({
    table,
    enableKeyboardShortcuts,
    focusedId,
    selection,
  }: TableActionsProps) => {
    const [targetId, setTargetId] = useState<string | null>(null);
    const store = useStore();
    const selectCount = selection?.length;

    const clearSelection = () => table.resetRowSelection();

    const onOpenCommandK = () => {
      if (selection?.length === 1) {
        store.ui.commandMenu.setType('ContactCommands');
        store.ui.commandMenu.setContext({
          entity: 'Contact',
          ids: selection,
        });
        store.ui.commandMenu.setOpen(true);
      } else {
        store.ui.commandMenu.setType('ContactBulkCommands');
        store.ui.commandMenu.setContext({
          entity: 'Contact',
          ids: selection,
        });
        store.ui.commandMenu.setOpen(true);
      }
    };

    const handleOpen = (type: CommandMenuType, property?: string) => {
      if (selection?.length >= 1) {
        store.ui.commandMenu.setContext({
          ids: selection,
          entity: 'Contact',
          property: property,
        });
      } else {
        store.ui.commandMenu.setContext({
          ids: [focusedId || ''],
          entity: 'Contact',
          property: property,
        });
      }

      if (selection?.length === 1) {
        store.ui.commandMenu.setContext({
          ids: selection,
          entity: 'Contact',
          property: property,
        });
      }

      store.ui.commandMenu.setType(type);
      store.ui.commandMenu.setOpen(true);
    };

    const onHideContacts = () => {
      store.ui.commandMenu.setType('DeleteConfirmationModal');
      store.ui.commandMenu.setOpen(true);
      store.ui.commandMenu.setCallback(() => {
        clearSelection();
      });
    };

    useEffect(() => {
      if (selectCount === 1) {
        setTargetId(selection[0]);
      }

      if (selectCount < 1) {
        setTargetId(null);
      }
    }, [selectCount, focusedId]);

    useKeys(
      ['Shift', 'T'],
      (e) => {
        e.stopPropagation();
        e.preventDefault();
        handleOpen('EditPersonaTag');
      },
      { when: enableKeyboardShortcuts },
    );
    useKeys(
      ['Shift', 'E'],
      (e) => {
        e.stopPropagation();
        e.preventDefault();
        editEmailUseCase.setEmail(
          table
            .getRow(focusedId || '')
            .original.value.emails?.find((e) => e.primary)?.email || '',
        );
        handleOpen('EditEmail', 'email');
      },
      { when: enableKeyboardShortcuts && (selectCount === 1 || !!focusedId) },
    );

    useModKey(
      'Backspace',
      () => {
        onHideContacts();
      },
      { when: enableKeyboardShortcuts },
    );

    useKeys(
      ['Shift', 'R'],
      (e) => {
        e.stopPropagation();
        e.preventDefault();
        handleOpen('EditName', 'name');
      },
      { when: enableKeyboardShortcuts && (selectCount === 1 || !!focusedId) },
    );

    useKeys(
      ['Shift', 'A'],
      (e) => {
        e.stopPropagation();
        e.preventDefault();
        handleOpen('EditContactFlow');
      },
      { when: enableKeyboardShortcuts },
    );

    useKeyBindings(
      {
        Escape: clearSelection,
      },
      { when: enableKeyboardShortcuts },
    );

    useKeyBindings({
      c: (e) => {
        e.stopPropagation();
        e.preventDefault();
        handleOpen('AddContactsBulk');
      },
    });

    useKeyBindings(
      {
        Space: (e) => {
          e.stopPropagation();
          e.preventDefault();

          if (store.ui.showPreviewCard && focusedId) {
            if (focusedId === store.ui.focusRow) {
              store.ui.setShowPreviewCard(false);

              return;
            }

            store.ui.setFocusRow(focusedId);

            return;
          }

          store.ui.setShowPreviewCard(true);
        },
      },
      { when: !!focusedId },
    );

    if (!selectCount && !targetId) return null;

    return (
      <SharedTableActions
        table={table}
        onHide={onHideContacts}
        handleOpen={handleOpen}
        selectCount={selectCount}
        onOpenCommandK={onOpenCommandK}
      />
    );
  },
);
