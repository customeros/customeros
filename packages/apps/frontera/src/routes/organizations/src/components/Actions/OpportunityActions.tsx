import { useState, useEffect } from 'react';

import { useKeys, useKeyBindings } from 'rooks';
import { ContactStore } from '@store/Contacts/Contact.store';
import { CommandMenuType } from '@store/UI/CommandMenu.store.ts';

import { useModKey } from '@shared/hooks/useModKey';
import { TableInstance } from '@ui/presentation/Table';
import { SharedTableActions } from '@organizations/components/Actions/components/SharedActions.tsx';

interface TableActionsProps {
  focusedId: string | null;
  onOpenCommandK: () => void;
  enableKeyboardShortcuts?: boolean;
  table: TableInstance<ContactStore>;
  handleOpen: (type: CommandMenuType) => void;
}

export const OpportunitiesTableActions = ({
  table,
  enableKeyboardShortcuts,
  onOpenCommandK,
  handleOpen,
  focusedId,
}: TableActionsProps) => {
  const [targetId, setTargetId] = useState<string | null>(null);
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
    ['Shift', 'S'],
    (e) => {
      e.stopPropagation();
      e.preventDefault();
      handleOpen('ChangeStage');
    },
    { when: enableKeyboardShortcuts },
  );
  useKeys(
    ['Shift', 'O'],
    (e) => {
      e.stopPropagation();
      e.preventDefault();
      handleOpen('AssignOwner');
    },
    { when: enableKeyboardShortcuts && (selectCount === 1 || !!focusedId) },
  );

  useKeys(
    ['Shift', 'R'],
    (e) => {
      e.stopPropagation();
      e.preventDefault();
      handleOpen('RenameOpportunityName');
    },
    { when: enableKeyboardShortcuts && (selectCount === 1 || !!focusedId) },
  );

  useModKey(
    'Backspace',
    () => {
      handleOpen('DeleteConfirmationModal');
    },
    { when: enableKeyboardShortcuts },
  );
  useKeyBindings(
    {
      Escape: clearSelection,
    },
    { when: enableKeyboardShortcuts },
  );

  if (!selectCount && !targetId) return null;

  return (
    <SharedTableActions
      table={table}
      handleOpen={handleOpen}
      onOpenCommandK={onOpenCommandK}
      onHide={() => handleOpen('DeleteConfirmationModal')}
    />
  );
};
