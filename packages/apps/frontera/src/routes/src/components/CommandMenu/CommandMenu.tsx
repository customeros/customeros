import type { ReactElement } from 'react';

import { useKey } from 'rooks';
import { observer } from 'mobx-react-lite';
import { CommandMenuType } from '@store/UI/CommandMenu.store';

import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import {
  Modal,
  ModalBody,
  ModalPortal,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal/Modal';

import {
  GlobalHub,
  AssignOwner,
  OpportunityHub,
  ChangeCurrency,
  ChangeArrEstimate,
  OpportunityCommands,
} from './commands';

const Commands: Record<CommandMenuType, ReactElement> = {
  GlobalHub: <GlobalHub />,
  AssignOwner: <AssignOwner />,
  ChangeCurrency: <ChangeCurrency />,
  OpportunityHub: <OpportunityHub />,
  ChangeArrEstimate: <ChangeArrEstimate />,
  OpportunityCommands: <OpportunityCommands />,
};

export const CommandMenu = observer(() => {
  const store = useStore();

  useKey('Escape', () => {
    store.ui.commandMenu.setOpen(false);
  });
  useModKey('k', () => {
    store.ui.commandMenu.setOpen(true);
  });

  return (
    <Modal
      open={store.ui.commandMenu.isOpen}
      onOpenChange={store.ui.commandMenu.setOpen}
    >
      <ModalPortal>
        <ModalOverlay className='z-[100]'>
          <ModalBody>
            <ModalContent>
              {Commands[store.ui.commandMenu.type ?? 'GlobalHub']}
            </ModalContent>
          </ModalBody>
        </ModalOverlay>
      </ModalPortal>
    </Modal>
  );
});
