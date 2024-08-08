import { observer } from 'mobx-react-lite';

import { Edit03 } from '@ui/media/icons/Edit03';
import { User01 } from '@ui/media/icons/User01';
import { Tag01 } from '@ui/media/icons/Tag01.tsx';
import { Archive } from '@ui/media/icons/Archive';
import { useStore } from '@shared/hooks/useStore';
import { Globe01 } from '@ui/media/icons/Globe01.tsx';
import { Columns03 } from '@ui/media/icons/Columns03';
import { CommandItem } from '@ui/overlay/CommandMenu';
import { Activity } from '@ui/media/icons/Activity.tsx';
import { CoinsStacked01 } from '@ui/media/icons/CoinsStacked01.tsx';
import { OrganizationStage, OrganizationRelationship } from '@graphql/types';
import { AlignHorizontalCentre02 } from '@ui/media/icons/AlignHorizontalCentre02';
import { organizationKeywords } from '@shared/components/CommandMenu/commands/organization/keywords.ts';
import {
  CommandsContainer,
  StageSubItemGroup,
} from '@shared/components/CommandMenu/commands/shared';
import {
  RelationshipSubItemGroup,
  UpdateHealthStatusSubItemGroup,
} from '@shared/components/CommandMenu/commands/organization';

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
          leftAccessory={<User01 />}
          keywords={organizationKeywords.add_contact}
          // rightAccessory={<Kbd className='px-1.5'>C</Kbd>}
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
          // rightAccessory={
          //   <>
          //     <Kbd className='px-1.5'>
          //       <ArrowBlockUp className='size-3' />
          //     </Kbd>
          //     <Kbd className='px-1.5'>T</Kbd>
          //   </>
          // }
        >
          Change or add tags...
        </CommandItem>

        {!!organization?.value?.tags?.length && (
          <CommandItem
            leftAccessory={<Tag01 />}
            keywords={organizationKeywords.change_or_add_tags}
            onSelect={() => {
              organization?.removeAllTagsFromOrganization();
              store.ui.commandMenu.setOpen(false);
            }}
          >
            Remove tags
          </CommandItem>
        )}

        <CommandItem
          leftAccessory={<Edit03 />}
          keywords={organizationKeywords.rename_org}
          // rightAccessory={
          //   <>
          //     <Kbd className='px-1.5'>
          //       <ArrowBlockUp className='size-3' />
          //     </Kbd>
          //     <Kbd className='px-1.5'>R</Kbd>
          //   </>
          // }
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
          keywords={organizationKeywords.edit_website}
          onSelect={() => {
            store.ui.commandMenu.setType('RenameOrganizationProperty');
            store.ui.commandMenu.setContext({
              ...store.ui.commandMenu.context,
              property: 'website',
            });
          }}
        >
          Edit website...
        </CommandItem>
        <CommandItem
          leftAccessory={<AlignHorizontalCentre02 />}
          keywords={organizationKeywords.change_relationship}
          onSelect={() => {
            store.ui.commandMenu.setType('ChangeRelationship');
          }}
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
          // rightAccessory={
          //   <>
          //     <Kbd className='px-1.5'>
          //       <CommandIcon className='size-3' />
          //     </Kbd>
          //     <Kbd className='px-1.5'>
          //       <Delete className='size-3' />
          //     </Kbd>
          //   </>
          // }
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
          Change health status...
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
          // rightAccessory={
          //   <>
          //     <Kbd className='px-1.5'>
          //       <ArrowBlockUp className='size-3' />
          //     </Kbd>
          //     <Kbd className='px-1.5'>O</Kbd>
          //   </>
          // }
        >
          Assign owner...
        </CommandItem>

        <CommandItem
          leftAccessory={<CoinsStacked01 />}
          keywords={organizationKeywords.create_new_opportunity}
          onSelect={() => {
            store.organizations.updateStage(
              [organization?.id as string],
              OrganizationStage.Engaged,
            );
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
