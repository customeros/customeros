import { useRef } from 'react';
import type { ReactElement } from 'react';
import { useLocation } from 'react-router-dom';

import { useKey } from 'rooks';
import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import { CommandMenuType } from '@store/UI/CommandMenu.store';

import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick';
import { DuplicateView } from '@shared/components/CommandMenu/commands/tableViewDef/DuplicateView';
import { OpportunityBulkCommands } from '@shared/components/CommandMenu/commands/OpportunityBulkCommands';
import {
  Modal,
  ModalBody,
  ModalPortal,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal/Modal';
import { ChangeBulkArrEstimate } from '@shared/components/CommandMenu/commands/opportunity/ChangeBulkArrEstimate';

import {
  EditName,
  GlobalHub,
  EditEmail,
  ChangeTags,
  ContactHub,
  AssignOwner,
  ChangeStage,
  EditJobTitle,
  EditTimeZone,
  OpportunityHub,
  ChangeCurrency,
  EditPersonaTag,
  OrganizationHub,
  EditPhoneNumber,
  ContactCommands,
  ChangeArrEstimate,
  ChangeRelationship,
  UpdateHealthStatus,
  AddNewOrganization,
  RenameTableViewDef,
  OpportunityCommands,
  ChangeOrAddJobRoles,
  ContactBulkCommands,
  OrganizationCommands,
  RenameOpportunityName,
  ChooseOpportunityStage,
  MergeConfirmationModal,
  SetOpportunityNextSteps,
  DeleteConfirmationModal,
  AddContactViaLinkedInUrl,
  OrganizationBulkCommands,
  RenameOrganizationProperty,
  ChooseOpportunityOrganization,
  ContactEmailVerificationInfoModal,
} from './commands';

//can we keep this in a nice order ? Thanks
const Commands: Record<CommandMenuType, ReactElement> = {
  EditName: <EditName />,
  GlobalHub: <GlobalHub />,
  EditEmail: <EditEmail />,
  ChangeTags: <ChangeTags />,
  ContactHub: <ContactHub />,
  AssignOwner: <AssignOwner />,
  ChangeStage: <ChangeStage />,
  EditTimeZone: <EditTimeZone />,
  EditJobTitle: <EditJobTitle />,
  ChangeCurrency: <ChangeCurrency />,
  OpportunityHub: <OpportunityHub />,
  EditPersonaTag: <EditPersonaTag />,
  OrganizationHub: <OrganizationHub />,
  ContactCommands: <ContactCommands />,
  EditPhoneNumber: <EditPhoneNumber />,
  ChangeArrEstimate: <ChangeArrEstimate />,
  RenameTableViewDef: <RenameTableViewDef />,
  AddNewOrganization: <AddNewOrganization />,
  ChangeRelationship: <ChangeRelationship />,
  UpdateHealthStatus: <UpdateHealthStatus />,
  ChangeOrAddJobRoles: <ChangeOrAddJobRoles />,
  ContactBulkCommands: <ContactBulkCommands />,
  OpportunityCommands: <OpportunityCommands />,
  OrganizationCommands: <OrganizationCommands />,
  RenameOpportunityName: <RenameOpportunityName />,
  MergeConfirmationModal: <MergeConfirmationModal />,
  ChooseOpportunityStage: <ChooseOpportunityStage />,
  SetOpportunityNextSteps: <SetOpportunityNextSteps />,
  DeleteConfirmationModal: <DeleteConfirmationModal />,
  OrganizationBulkCommands: <OrganizationBulkCommands />,
  AddContactViaLinkedInUrl: <AddContactViaLinkedInUrl />,
  RenameOrganizationProperty: <RenameOrganizationProperty />,
  ChooseOpportunityOrganization: <ChooseOpportunityOrganization />,
  ContactEmailVerificationInfoModal: <ContactEmailVerificationInfoModal />,
  DuplicateView: <DuplicateView />,
  OpportunityBulkCommands: <OpportunityBulkCommands />,
  ChangeBulkArrEstimate: <ChangeBulkArrEstimate />,
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
            <ModalContent
              ref={commandRef}
              //We added this onKeyDown here to not refactor the command componenet and is more easy to extend in the future if needed to add more keybindings
              //the escape command ( is a default keybinding to close the modal) is not needed in the command component
              onKeyDown={(e) => {
                if (e.key !== 'Escape' && e.key !== 'Enter') {
                  e.stopPropagation();
                }
              }}
            >
              {Commands[store.ui.commandMenu.type ?? 'GlobalHub']}
            </ModalContent>
          </ModalBody>
        </ModalOverlay>
      </ModalPortal>
    </Modal>
  );
});
