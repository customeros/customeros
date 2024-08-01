import type { ReactElement } from 'react';

import { useKey } from 'rooks';
import { observer } from 'mobx-react-lite';
import { CommandMenuType } from '@store/UI/CommandMenu.store';

import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import { OrganizationHub } from '@shared/components/CommandMenu/commands/OrganizationHub.tsx';
import { ChangeTags } from '@shared/components/CommandMenu/commands/organization/ChangeTags.tsx';
import { ChangeStage } from '@shared/components/CommandMenu/commands/organization/ChangeStage.tsx';
import { OrganizationCommands } from '@shared/components/CommandMenu/commands/OrganizationCommands.tsx';
import {
  Modal,
  ModalBody,
  ModalPortal,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal/Modal';
import { ChangeRelationship } from '@shared/components/CommandMenu/commands/organization/ChangeRelationship.tsx';
import { UpdateHealthStatus } from '@shared/components/CommandMenu/commands/organization/UpdateHealthStatus.tsx';
import { AddContactViaLinkedInUrl } from '@shared/components/CommandMenu/commands/organization/AddContactViaLinkedInUrl.tsx';
import { RenameOrganizationProperty } from '@shared/components/CommandMenu/commands/organization/RenameOrganizationProperty.tsx';

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
  AddContactViaLinkedInUrl: <AddContactViaLinkedInUrl />,
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
