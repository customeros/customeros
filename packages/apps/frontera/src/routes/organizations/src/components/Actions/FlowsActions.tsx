import { useState, useEffect } from 'react';

import { useKeys, useKeyBindings } from 'rooks';
import { ContactStore } from '@store/Contacts/Contact.store';
import { CommandMenuType } from '@store/UI/CommandMenu.store.ts';

import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import { TableInstance } from '@ui/presentation/Table';
import { SharedTableActions } from '@organizations/components/Actions/components/SharedActions.tsx';

interface TableActionsProps {
  selection: string[];
  focusedId?: string | null;
  enableKeyboardShortcuts?: boolean;
  table: TableInstance<ContactStore>;
}

export const FlowSequencesTableActions = ({
  table,
  enableKeyboardShortcuts,
  selection,
  focusedId,
}: TableActionsProps) => {
  const store = useStore();

  const [targetId, setTargetId] = useState<string | null>(null);
  const selectCount = selection?.length;
  const clearSelection = () => table.resetRowSelection();

  const onOpenCommandK = () => {
    if (selection.length === 1) {
      store.ui.commandMenu.setType('FlowCommands');
      store.ui.commandMenu.setContext({
        entity: 'Flow',
        ids: selection,
      });
      store.ui.commandMenu.setOpen(true);
    } else {
      store.ui.commandMenu.setType('FlowsBulkCommands');
      store.ui.commandMenu.setContext({
        entity: 'Flows',
        ids: selection,
      });
      store.ui.commandMenu.setOpen(true);
    }
  };

  const handleOpen = (type: CommandMenuType, property?: string) => {
    if (selection?.length > 1) {
      store.ui.commandMenu.setContext({
        ids: selection,
        entity: 'Flows',
        property: property,
      });
    } else {
      store.ui.commandMenu.setContext({
        ids: [focusedId || ''],
        entity: 'Flow',
        property: property,
      });
    }

    store.ui.commandMenu.setType(type);
    store.ui.commandMenu.setOpen(true);
  };

  useEffect(() => {
    if (selectCount === 1) {
      setTargetId(selection[0]);
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
      handleOpen('ChangeFlowStatus');
    },
    { when: enableKeyboardShortcuts },
  );

  useKeys(
    ['Shift', 'R'],
    (e) => {
      e.stopPropagation();
      e.preventDefault();
      handleOpen('RenameFlow', 'name');
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
      selectCount={selectCount}
      onOpenCommandK={onOpenCommandK}
      onHide={() => handleOpen('DeleteConfirmationModal')}
    />
  );
};
