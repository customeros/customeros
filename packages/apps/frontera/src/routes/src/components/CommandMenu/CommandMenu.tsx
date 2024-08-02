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
  ChangeTags,
  AssignOwner,
  ChangeStage,
  OpportunityHub,
  ChangeCurrency,
  OrganizationHub,
  ChangeArrEstimate,
  ChangeRelationship,
  UpdateHealthStatus,
  OpportunityCommands,
  OrganizationCommands,
  RenameOpportunityName,
  SetOpportunityNextSteps,
  DeleteConfirmationModal,
  AddContactViaLinkedInUrl,
  RenameOrganizationProperty,
} from './commands';

const Commands: Record<CommandMenuType, ReactElement> = {
  GlobalHub: <GlobalHub />,
  AssignOwner: <AssignOwner />,
  ChangeCurrency: <ChangeCurrency />,
  UpdateHealthStatus: <UpdateHealthStatus />,
  OpportunityHub: <OpportunityHub />,
  ChangeArrEstimate: <ChangeArrEstimate />,
  OpportunityCommands: <OpportunityCommands />,
  OrganizationHub: <OrganizationHub />,
  OrganizationCommands: <OrganizationCommands />,
  ChangeRelationship: <ChangeRelationship />,
  ChangeStage: <ChangeStage />,
  RenameOrganizationProperty: <RenameOrganizationProperty />,
  ChangeTags: <ChangeTags />,
  RenameOpportunityName: <RenameOpportunityName />,
  AddContactViaLinkedInUrl: <AddContactViaLinkedInUrl />,
  SetOpportunityNextSteps: <SetOpportunityNextSteps />,
  DeleteConfirmationModal: <DeleteConfirmationModal />,
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
        {/* z-[5001] is needed to ensure tooltips are not overlapping  - tooltips have zIndex of 5000 - this should be revisited */}
        <ModalOverlay className='z-[5001]'>
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
