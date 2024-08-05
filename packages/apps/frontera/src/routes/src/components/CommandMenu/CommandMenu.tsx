import { useRef } from 'react';
import type { ReactElement } from 'react';
import { useLocation } from 'react-router-dom';

import { useKey } from 'rooks';
import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import { CommandMenuType } from '@store/UI/CommandMenu.store';

import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick.ts';
import {
  Modal,
  ModalBody,
  ModalPortal,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal/Modal';
import { OrganizationBulkCommands } from '@shared/components/CommandMenu/commands/OrganizationBulkCommands.tsx';

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
  ChooseOpportunityStage,
  SetOpportunityNextSteps,
  DeleteConfirmationModal,
  AddContactViaLinkedInUrl,
  RenameOrganizationProperty,
  ChooseOpportunityOrganization,
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
  OrganizationBulkCommands: <OrganizationBulkCommands />,
  ChooseOpportunityStage: <ChooseOpportunityStage />,
  ChooseOpportunityOrganization: <ChooseOpportunityOrganization />,
};

export const CommandMenu = observer(() => {
  const location = useLocation();

  const store = useStore();
  const commandRef = useRef(null);

  useOutsideClick({
    ref: commandRef,
    handler: () => store.ui.commandMenu.setOpen(false),
  });

  useKey('Escape', () => {
    match(location.pathname)
      .with('/prospects', () => {
        store.ui.commandMenu.setOpen(false);
        store.ui.commandMenu.setType('OpportunityHub');
      })
      .otherwise(() => {
        store.ui.commandMenu.setOpen(false);
      });
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
            <ModalContent ref={commandRef}>
              {Commands[store.ui.commandMenu.type ?? 'GlobalHub']}
            </ModalContent>
          </ModalBody>
        </ModalOverlay>
      </ModalPortal>
    </Modal>
  );
});
