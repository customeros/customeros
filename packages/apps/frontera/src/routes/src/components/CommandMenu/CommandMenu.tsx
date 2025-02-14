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

import { AddEmail } from './commands/contacts/AddEmail';
import { AddLinkedinUrl } from './commands/contacts/AddLinkedin';
import { EditLatestOrgActive } from './commands/contacts/EditLatestOrgActive';
import {
  FlowHub,
  EditSku,
  EditName,
  StopFlow,
  GlobalHub,
  EditEmail,
  StartFlow,
  AddNewSku,
  ChangeTags,
  ContactHub,
  RenameFlow,
  AssignOwner,
  ChangeStage,
  CreateAgent,
  RenameAgent,
  EditJobTitle,
  AddNewDomain,
  RemoveDomain,
  ArchiveAgent,
  CreateNewFlow,
  FlowsCommands,
  DuplicateFlow,
  AgentCommands,
  AgentsCommands,
  OpportunityHub,
  ChangeCurrency,
  EditPersonaTag,
  DuplicateAgent,
  OrganizationHub,
  ContactCommands,
  EditContactFlow,
  AddContactsBulk,
  SwitchWorkspace,
  ChangeFlowStatus,
  AddSingleContact,
  ChangeArrEstimate,
  FlowsBulkCommands,
  AddContactsToFlow,
  ChangeRelationship,
  UpdateHealthStatus,
  AddNewOrganization,
  RenameTableViewDef,
  OpportunityCommands,
  ContactBulkCommands,
  ConfirmBulkFlowEdit,
  EditCompanyLinkedin,
  OrganizationCommands,
  ActiveFlowUpdateInfo,
  RenameOpportunityName,
  ChangeBulkArrEstimate,
  UnlinkContactFromFlow,
  ConfirmSingleFlowEdit,
  FlowValidationMessage,
  ChooseOpportunityStage,
  MergeConfirmationModal,
  SetOpportunityNextSteps,
  DeleteConfirmationModal,
  GetBrowserExtensionLink,
  OrganizationBulkCommands,
  InstallLinkedInExtension,
  ConfirmEmailContentChanges,
  ChooseOpportunityOrganization,
  ContactEmailVerificationInfoModal,
} from './commands';

//can we keep this in a nice order ? Thanks
const Commands: Record<CommandMenuType, ReactElement> = {
  // Shared
  EditName: <EditName />,
  GlobalHub: <GlobalHub />,
  SwitchWorkspace: <SwitchWorkspace />,
  ChangeTags: <ChangeTags />,
  DuplicateView: <DuplicateView />,
  ChangeStage: <ChangeStage />,
  AssignOwner: <AssignOwner />,
  ChangeCurrency: <ChangeCurrency />,
  ChangeArrEstimate: <ChangeArrEstimate />,
  ChangeRelationship: <ChangeRelationship />,
  UpdateHealthStatus: <UpdateHealthStatus />,
  DeleteConfirmationModal: <DeleteConfirmationModal />,
  AddLinkedin: <AddLinkedinUrl />,

  // Contact
  ContactHub: <ContactHub />,
  ContactBulkCommands: <ContactBulkCommands />,
  ContactCommands: <ContactCommands />,
  EditEmail: <EditEmail />,
  EditJobTitle: <EditJobTitle />,
  EditPersonaTag: <EditPersonaTag />,
  ContactEmailVerificationInfoModal: <ContactEmailVerificationInfoModal />,
  UnlinkContactFromFlow: <UnlinkContactFromFlow />,
  ConfirmBulkFlowEdit: <ConfirmBulkFlowEdit />,
  ConfirmSingleFlowEdit: <ConfirmSingleFlowEdit />,
  AddContactsBulk: <AddContactsBulk />,
  AddEmail: <AddEmail />,
  EditLatestOrgActive: <EditLatestOrgActive />,
  AddSingleContact: <AddSingleContact />,

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
  AddNewDomain: <AddNewDomain />,
  RemoveDomain: <RemoveDomain />,
  OrganizationCommands: <OrganizationCommands />,
  ChangeBulkArrEstimate: <ChangeBulkArrEstimate />,
  MergeConfirmationModal: <MergeConfirmationModal />,
  EditCompanyLinkedin: <EditCompanyLinkedin />,
  AddNewOrganization: <AddNewOrganization />,

  // Flows
  FlowHub: <FlowHub />,
  FlowsBulkCommands: <FlowsBulkCommands />,
  FlowCommands: <FlowsCommands />,
  StartFlow: <StartFlow />,
  StopFlow: <StopFlow />,
  CreateNewFlow: <CreateNewFlow />,
  DuplicateFlow: <DuplicateFlow />,
  RenameFlow: <RenameFlow />,
  ChangeFlowStatus: <ChangeFlowStatus />,
  EditContactFlow: <EditContactFlow />,
  AddContactsToFlow: <AddContactsToFlow />,
  FlowValidationMessage: <FlowValidationMessage />,
  ActiveFlowUpdateInfo: <ActiveFlowUpdateInfo />,
  GetBrowserExtensionLink: <GetBrowserExtensionLink />,
  ConfirmEmailContentChanges: <ConfirmEmailContentChanges />,
  InstallLinkedInExtension: <InstallLinkedInExtension />,

  //TableViewDef
  RenameTableViewDef: <RenameTableViewDef />,

  //SKU
  EditSku: <EditSku />,
  AddNewSku: <AddNewSku />,

  //Agent
  CreateAgent: <CreateAgent />,
  AgentsCommands: <AgentsCommands />,
  AgentCommands: <AgentCommands />,
  DuplicateAgent: <DuplicateAgent />,
  RenameAgent: <RenameAgent />,
  ArchiveAgent: <ArchiveAgent />,
};

export const CommandMenu = observer(() => {
  const location = useLocation();

  const store = useStore();
  const commandRef = useRef(null);

  useOutsideClick({
    ref: commandRef,
    handler: () => {
      store.ui.commandMenu.setOpen(false);
    },
  });

  useKey('Escape', (e) => {
    e.stopPropagation();
    match(location.pathname)
      .with('/prospects', () => {
        store.ui.commandMenu.setOpen(false);
        store.ui.commandMenu.setType('OpportunityHub');
      })
      .with('/agents', () => {
        if (store.ui.commandMenu.type === 'AgentsCommands') {
          store.ui.commandMenu.setOpen(false);
        } else {
          store.ui.commandMenu.setType('AgentsCommands');
          store.ui.commandMenu.setOpen(false);
        }
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
      // onOpenChange={store.ui.commandMenu.setOpen}
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
