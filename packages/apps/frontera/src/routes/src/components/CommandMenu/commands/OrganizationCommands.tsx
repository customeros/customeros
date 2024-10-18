import { observer } from 'mobx-react-lite';

import { Tag01 } from '@ui/media/icons/Tag01';
import { Edit03 } from '@ui/media/icons/Edit03';
import { User01 } from '@ui/media/icons/User01';
import { Delete } from '@ui/media/icons/Delete';
import { User03 } from '@ui/media/icons/User03';
import { Archive } from '@ui/media/icons/Archive';
import { useStore } from '@shared/hooks/useStore';
import { Globe01 } from '@ui/media/icons/Globe01';
import { Activity } from '@ui/media/icons/Activity';
import { Columns03 } from '@ui/media/icons/Columns03';
import { ArrowBlockUp } from '@ui/media/icons/ArrowBlockUp';
import { CoinsStacked01 } from '@ui/media/icons/CoinsStacked01';
import { Kbd, CommandKbd, CommandItem } from '@ui/overlay/CommandMenu';
import { AlignHorizontalCentre02 } from '@ui/media/icons/AlignHorizontalCentre02';
import {
  InternalType,
  InternalStage,
  OrganizationRelationship,
} from '@graphql/types';
import { organizationKeywords } from '@shared/components/CommandMenu/commands/organization/keywords.ts';
import {
  CommandsContainer,
  StageSubItemGroup,
} from '@shared/components/CommandMenu/commands/shared';
import {
  RelationshipSubItemGroup,
  UpdateHealthStatusSubItemGroup,
} from '@shared/components/CommandMenu/commands/organization';

import { OwnerSubItemGroup } from './shared/OwnerSubItemGroup';
import { AddTagSubItemGroup } from './organization/AddTagSubItemGroup';

// TODO - uncomment keyboard shortcuts when they are implemented
export const OrganizationCommands = observer(() => {
  const store = useStore();
  const selectedIds = store.ui.commandMenu.context.ids;
  const id = (store.ui.commandMenu.context.ids as string[])?.[0];
  const organization = store.organizations.value.get(id);
  const label = `Organization - ${organization?.value.name}`;

  return (
    <CommandsContainer label={label}>
      <>
        <CommandItem
          leftAccessory={<User03 />}
          rightAccessory={<Kbd>C</Kbd>}
          keywords={organizationKeywords.add_contact}
          onSelect={() => {
            store.ui.commandMenu.setType('AddContactViaLinkedInUrl');
          }}
        >
          Add contact via LinkedIn
        </CommandItem>

        <CommandItem
          leftAccessory={<Tag01 />}
          keywords={organizationKeywords.change_or_add_tags}
          onSelect={() => {
            store.ui.commandMenu.setType('ChangeTags');
          }}
          rightAccessory={
            <>
              <Kbd>
                <ArrowBlockUp className='text-inherit size-3' />
              </Kbd>
              <Kbd>T</Kbd>
            </>
          }
        >
          Change or add tags...
        </CommandItem>
        <AddTagSubItemGroup />

        {!!organization?.value?.tags?.length && (
          <CommandItem
            leftAccessory={<Tag01 />}
            keywords={['change', 'add', 'tags', 'update', 'edit']}
            onSelect={() => {
              const tagCount = organization?.value?.tags?.length ?? 0;

              for (let i = 0; i < tagCount; i++) {
                organization?.value?.tags?.pop();
                organization?.commit();
              }
              store.ui.toastSuccess(
                'All tags were removed',
                'tags-remove-success',
              );

              store.ui.commandMenu.setOpen(false);
            }}
          >
            Remove tags
          </CommandItem>
        )}

        <CommandItem
          leftAccessory={<Edit03 />}
          rightAccessory={
            <>
              <Kbd>
                <ArrowBlockUp className='size-3' />
              </Kbd>
              <Kbd>R</Kbd>
            </>
          }
          keywords={[
            'rename',
            'org',
            'organization',
            'company',
            'update',
            'edit',
            'change',
          ]}
          onSelect={() => {
            store.ui.commandMenu.setType('RenameOrganizationProperty');
            store.ui.commandMenu.setContext({
              ...store.ui.commandMenu.context,
              property: 'name',
            });
          }}
        >
          Rename organization
        </CommandItem>

        <CommandItem
          leftAccessory={<Globe01 />}
          keywords={[
            'edit',
            'website',
            'change',
            'domain',
            'link',
            'url',
            'web address',
          ]}
          onSelect={() => {
            store.ui.commandMenu.setType('RenameOrganizationProperty');
            store.ui.commandMenu.setContext({
              ...store.ui.commandMenu.context,
              property: 'website',
            });
          }}
        >
          Edit website
        </CommandItem>
        <CommandItem
          leftAccessory={<AlignHorizontalCentre02 />}
          onSelect={() => {
            store.ui.commandMenu.setType('ChangeRelationship');
          }}
          keywords={[
            'change',
            'relationship',
            'status',
            'update',
            'edit',
            'customer',
            'prospect',
            'former customer',
            'unqualified',
          ]}
        >
          Change relationship...
        </CommandItem>
        <RelationshipSubItemGroup
          selectedIds={selectedIds}
          closeMenu={() => store.ui.commandMenu.setOpen(false)}
          updateRelationship={store.organizations.updateRelationship}
        />

        {organization?.value?.relationship ===
          OrganizationRelationship.Prospect && (
          <CommandItem
            leftAccessory={<Columns03 />}
            keywords={organizationKeywords.change_org_stage}
            onSelect={() => {
              store.ui.commandMenu.setType('ChangeStage');
            }}
          >
            Change org stage...
          </CommandItem>
        )}
        <StageSubItemGroup
          selectedIds={selectedIds}
          updateStage={store.organizations.updateStage}
          closeMenu={() => store.ui.commandMenu.setOpen(false)}
        />

        <CommandItem
          leftAccessory={<Archive />}
          keywords={organizationKeywords.archive_org}
          onSelect={() => {
            store.ui.commandMenu.setType('DeleteConfirmationModal');
          }}
          rightAccessory={
            <>
              <CommandKbd />
              <Kbd>
                <Delete className='text-inherit size-3' />
              </Kbd>
            </>
          }
        >
          Archive org
        </CommandItem>

        <CommandItem
          leftAccessory={<Activity />}
          keywords={organizationKeywords.change_health_status}
          onSelect={() => {
            store.ui.commandMenu.setType('UpdateHealthStatus');
          }}
        >
          Health status...
        </CommandItem>
        <UpdateHealthStatusSubItemGroup
          selectedIds={selectedIds}
          updateHealth={store.organizations.updateHealth}
          closeMenu={() => store.ui.commandMenu.setOpen(false)}
        />

        <CommandItem
          leftAccessory={<User01 />}
          keywords={organizationKeywords.assign_owner}
          onSelect={() => {
            store.ui.commandMenu.setType('AssignOwner');
          }}
          rightAccessory={
            <>
              <Kbd>
                <ArrowBlockUp className='text-inherit size-3' />
              </Kbd>
              <Kbd>O</Kbd>
            </>
          }
        >
          Assign owner...
        </CommandItem>
        <OwnerSubItemGroup />

        <CommandItem
          rightAccessory={<Kbd>O</Kbd>}
          leftAccessory={<CoinsStacked01 />}
          keywords={organizationKeywords.create_new_opportunity}
          onSelect={() => {
            store.opportunities.create({
              organization: organization?.value,
              name: `${organization?.value.name}'s opportunity`,
              internalType: InternalType.Nbo,
              externalStage: String(
                store.settings.tenant.value?.opportunityStages[0].value,
              ),
              internalStage: InternalStage.Open,
            });
            store.ui.commandMenu.setOpen(false);
          }}
        >
          Create new opportunity
        </CommandItem>

        {/*<CommandItem*/}
        {/*  leftAccessory={<Trophy01 />}*/}
        {/*  onSelect={() => {*/}
        {/*    store.ui.commandMenu.setType('AssignOwner');*/}
        {/*  }}*/}
        {/*>*/}
        {/*  Change onboarding stage*/}
        {/*</CommandItem>*/}
      </>
    </CommandsContainer>
  );
});
