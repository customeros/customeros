import { useState, useEffect } from 'react';

import { useKeys, useKeyBindings } from 'rooks';
import { ContactStore } from '@store/Contacts/Contact.store';
import { CommandMenuType } from '@store/UI/CommandMenu.store.ts';

import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import { TableInstance } from '@ui/presentation/Table';

import { SharedTableActions } from './components/SharedActions';

interface TableActionsProps {
  focusedId: string | null;
  onOpenCommandK: () => void;
  onHideContacts: () => void;
  enableKeyboardShortcuts?: boolean;
  table: TableInstance<ContactStore>;
  handleOpen: (type: CommandMenuType) => void;
}

export const ContactTableActions = ({
  table,
  enableKeyboardShortcuts,
  onOpenCommandK,
  handleOpen,
  onHideContacts,
  focusedId,
}: TableActionsProps) => {
  const [targetId, setTargetId] = useState<string | null>(null);
  const store = useStore();
  const selection = table.getState().rowSelection;
  const selectedIds = Object.keys(selection);
  const selectCount = selectedIds.length;

  const clearSelection = () => table.resetRowSelection();

  useEffect(() => {
    if (selectCount === 1) {
      setTargetId(selectedIds[0]);
    }

    if (selectCount < 1) {
      setTargetId(null);
    }
  }, [selectCount]);

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
      handleOpen('EditEmail');
    },
    { when: enableKeyboardShortcuts && (selectCount === 1 || !!focusedId) },
  );

  useKeys(
    ['Shift', 'R'],
    (e) => {
      e.stopPropagation();
      e.preventDefault();
      handleOpen('EditName');
    },
    { when: enableKeyboardShortcuts && (selectCount === 1 || !!focusedId) },
  );
  useKeyBindings(
    {
      Escape: clearSelection,
    },
    { when: enableKeyboardShortcuts },
  );

  useKeyBindings(
    {
      Space: (e) => {
        e.stopPropagation();
        e.preventDefault();
        store.ui.setContactPreviewCardOpen(true);
      },
    },
    { when: !!focusedId },
  );

  useModKey(
    'Backspace',
    () => {
      onHideContacts();
    },
    { when: enableKeyboardShortcuts },
  );

  if (!selectCount && !targetId) return null;

  return (
    <SharedTableActions
      table={table}
      onHide={onHideContacts}
      handleOpen={handleOpen}
      onOpenCommandK={onOpenCommandK}
    />
  );
};
