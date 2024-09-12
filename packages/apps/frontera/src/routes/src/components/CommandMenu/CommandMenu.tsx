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
import { SequenceHub } from '@shared/components/CommandMenu/commands/SequenceHub.tsx';
import { SequenceCommands } from '@shared/components/CommandMenu/commands/SequencesCommands.tsx';
import { DuplicateView } from '@shared/components/CommandMenu/commands/tableViewDef/DuplicateView';
import { RenameSequence } from '@shared/components/CommandMenu/commands/sequences/RenameSequence.tsx';
import { SequencesBulkCommands } from '@shared/components/CommandMenu/commands/SequencesBulkCommands.tsx';
import { OpportunityBulkCommands } from '@shared/components/CommandMenu/commands/OpportunityBulkCommands';
import {
  Modal,
  ModalBody,
  ModalPortal,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal/Modal';
import { ConfirmBulkSequenceEdit } from '@shared/components/CommandMenu/commands/contacts/ConfirmBulkSequenceEdit.tsx';

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
  CreateNewSequence,
  ChangeRelationship,
  UpdateHealthStatus,
  AddNewOrganization,
  RenameTableViewDef,
  OpportunityCommands,
  ChangeOrAddJobRoles,
  ContactBulkCommands,
  EditContactSequence,
  OrganizationCommands,
  ChangeSequenceStatus,
  RenameOpportunityName,
  ChangeBulkArrEstimate,
  ChooseOpportunityStage,
  MergeConfirmationModal,
  SetOpportunityNextSteps,
  DeleteConfirmationModal,
  AddContactViaLinkedInUrl,
  OrganizationBulkCommands,
  UnlinkContactFromSequence,
  RenameOrganizationProperty,
  ChooseOpportunityOrganization,
  ContactEmailVerificationInfoModal,
} from './commands';

//can we keep this in a nice order ? Thanks
const Commands: Record<CommandMenuType, ReactElement> = {
  // Shared
  EditName: <EditName />,
  GlobalHub: <GlobalHub />,
  ChangeTags: <ChangeTags />,
  DuplicateView: <DuplicateView />,
  ChangeStage: <ChangeStage />,
  AssignOwner: <AssignOwner />,
  ChangeCurrency: <ChangeCurrency />,
  ChangeArrEstimate: <ChangeArrEstimate />,
  ChangeRelationship: <ChangeRelationship />,
  UpdateHealthStatus: <UpdateHealthStatus />,
  DeleteConfirmationModal: <DeleteConfirmationModal />,

  // Contact
  ContactHub: <ContactHub />,
  ContactBulkCommands: <ContactBulkCommands />,
  ContactCommands: <ContactCommands />,
  EditEmail: <EditEmail />,
  EditTimeZone: <EditTimeZone />,
  EditJobTitle: <EditJobTitle />,
  EditPersonaTag: <EditPersonaTag />,
  EditPhoneNumber: <EditPhoneNumber />,
  ChangeOrAddJobRoles: <ChangeOrAddJobRoles />,
  ContactEmailVerificationInfoModal: <ContactEmailVerificationInfoModal />,
  UnlinkContactFromSequence: <UnlinkContactFromSequence />,
  ConfirmBulkSequenceEdit: <ConfirmBulkSequenceEdit />,

  // Opportunity
  OpportunityHub: <OpportunityHub />,
  OpportunityBulkCommands: <OpportunityBulkCommands />,
  OpportunityCommands: <OpportunityCommands />,
  RenameOpportunityName: <RenameOpportunityName />,
  ChooseOpportunityStage: <ChooseOpportunityStage />,
  SetOpportunityNextSteps: <SetOpportunityNextSteps />,
  ChooseOpportunityOrganization: <ChooseOpportunityOrganization />,

  // Organization
  OrganizationHub: <OrganizationHub />,
  OrganizationBulkCommands: <OrganizationBulkCommands />,
  OrganizationCommands: <OrganizationCommands />,
  ChangeBulkArrEstimate: <ChangeBulkArrEstimate />,
  MergeConfirmationModal: <MergeConfirmationModal />,
  AddNewOrganization: <AddNewOrganization />,
  AddContactViaLinkedInUrl: <AddContactViaLinkedInUrl />,
  RenameOrganizationProperty: <RenameOrganizationProperty />,

  // Sequence
  SequenceHub: <SequenceHub />,
  SequencesBulkCommands: <SequencesBulkCommands />,
  SequenceCommands: <SequenceCommands />,
  CreateNewSequence: <CreateNewSequence />,
  RenameSequence: <RenameSequence />,
  ChangeSequenceStatus: <ChangeSequenceStatus />,
  EditContactSequence: <EditContactSequence />,

  //TableViewDef
  RenameTableViewDef: <RenameTableViewDef />,
};

export const CommandMenu = observer(() => {
  const location = useLocation();

  const store = useStore();
  const commandRef = useRef(null);

  useOutsideClick({
    ref: commandRef,
    handler: () => store.ui.commandMenu.setOpen(false),
  });

  useKey('Escape', (e) => {
    e.stopPropagation();
    match(location.pathname)
      .with('/prospects', () => {
        store.ui.commandMenu.setOpen(false);
        store.ui.commandMenu.setType('OpportunityHub');
      })

      .otherwise(() => {
        store.ui.commandMenu.setOpen(false);
      });
  });

  useModKey('k', (e) => {
    e.stopPropagation();
    store.ui.commandMenu.setOpen(true);
  });
  useModKey('Enter', (e) => {
    e.stopPropagation();
  });

  return (
    <Modal
      open={store.ui.commandMenu.isOpen}
      onOpenChange={store.ui.commandMenu.setOpen}
    >
      <ModalPortal>
        {/* z-[5001] is needed to ensure tooltips are not overlapping  - tooltips have zIndex of 5000 - this should be revisited */}
        <ModalOverlay
          className='z-[5001]'
          // Prevent event propagation to parent elements, except for 'Escape' key
          onKeyDown={(e) => e.key !== 'Escape' && e.stopPropagation()}
        >
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
