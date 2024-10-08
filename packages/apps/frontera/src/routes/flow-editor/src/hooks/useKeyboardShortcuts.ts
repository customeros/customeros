import { useKey, useKeys } from 'rooks';
import { RootStore } from '@store/root';

export const useKeyboardShortcuts = (id: string, store: RootStore) => {
  useKeys(['Shift', 'S'], (e) => {
    e.stopPropagation();
    e.preventDefault();

    store.ui.commandMenu.setContext({
      ids: [id || ''],
      entity: 'Flow',
    });
    store.ui.commandMenu.setType('ChangeFlowStatus');
    store.ui.commandMenu.setOpen(true);
  });

  useKeys(['Shift', 'R'], (e) => {
    e.stopPropagation();
    e.preventDefault();
    store.ui.commandMenu.setContext({
      ids: [id || ''],
      entity: 'Flow',
      property: 'name',
    });
    store.ui.commandMenu.setType('RenameFlow');
    store.ui.commandMenu.setOpen(true);
  });

  useKey(
    'Escape',
    () => {
      store.ui.flowCommandMenu.setOpen(false);
    },
    {
      when: store.ui.flowCommandMenu.isOpen,
    },
  );
};
