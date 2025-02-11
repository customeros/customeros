import { useKey } from 'rooks';
import { RootStore } from '@store/root';

export const useKeyboardShortcuts = (_id: string, store: RootStore) => {
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
