import { observer } from 'mobx-react-lite';

import { Tag01 } from '@ui/media/icons/Tag01';
import { User01 } from '@ui/media/icons/User01';
import { Copy07 } from '@ui/media/icons/Copy07';
import { Archive } from '@ui/media/icons/Archive';
import { useStore } from '@shared/hooks/useStore';
import { Activity } from '@ui/media/icons/Activity';
import { Columns03 } from '@ui/media/icons/Columns03';
import { CoinsStacked01 } from '@ui/media/icons/CoinsStacked01';
import { OrganizationStage, OrganizationRelationship } from '@graphql/types';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';
import { StageSubMenu } from '@shared/components/CommandMenu/commands/shared';
import { AlignHorizontalCentre02 } from '@ui/media/icons/AlignHorizontalCentre02';
import {
  RelationshipSubMenu,
  UpdateHealthStatusSubMenu,
} from '@shared/components/CommandMenu/commands/organization';

// TODO - uncomment keyboard shortcuts when they are implemented
export const OrganizationBulkCommands = observer(() => {
  const store = useStore();
  const selectedIds = store.ui.commandMenu.context.ids;

  const organizations = selectedIds?.map((e: string) =>
    store.organizations.value.get(e),
  );
  const label = `${selectedIds?.length} organizations`;

  return (
    <Command>
      <CommandInput label={label} placeholder='Type a command or search' />
      <Command.List>
        <CommandItem
          leftAccessory={<Tag01 />}
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

        <CommandItem
          leftAccessory={<Tag01 />}
          onSelect={() => {
            store.organizations.removeTags(selectedIds);
            store.ui.commandMenu.setOpen(false);
          }}
        >
          Remove tags
        </CommandItem>

        <CommandItem
          leftAccessory={<AlignHorizontalCentre02 />}
          onSelect={() => {
            store.ui.commandMenu.setType('ChangeRelationship');
          }}
        >
          Change relationship...
        </CommandItem>

        <RelationshipSubMenu
          selectedIds={selectedIds}
          closeMenu={() => store.ui.commandMenu.setOpen(false)}
          updateRelationship={store.organizations.updateRelationship}
        />

        {organizations?.every(
          (organization) =>
            organization?.value?.relationship ===
            OrganizationRelationship.Prospect,
        ) && (
          <CommandItem
            leftAccessory={<Columns03 />}
            onSelect={() => {
              store.ui.commandMenu.setType('ChangeStage');
            }}
          >
            Change org stage...
          </CommandItem>
        )}

        <StageSubMenu
          selectedIds={selectedIds}
          updateStage={store.organizations.updateStage}
          closeMenu={() => store.ui.commandMenu.setOpen(false)}
        />

        <CommandItem
          leftAccessory={<Archive />}
          keywords={['delete', 'archive']}
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
          leftAccessory={<Copy07 />}
          onSelect={() => {
            const [primaryId, ...restIds] = selectedIds;

            store.organizations.merge(primaryId, restIds);
            store.ui.commandMenu.setOpen(false);
          }}
        >
          Merge
        </CommandItem>

        <CommandItem
          leftAccessory={<Activity />}
          onSelect={() => {
            store.ui.commandMenu.setType('UpdateHealthStatus');
          }}
        >
          Change health status...
        </CommandItem>

        <UpdateHealthStatusSubMenu
          selectedIds={selectedIds}
          updateHealth={store.organizations.updateHealth}
          closeMenu={() => store.ui.commandMenu.setOpen(false)}
        />

        <CommandItem
          leftAccessory={<User01 />}
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
          onSelect={() => {
            store.organizations.updateStage(
              selectedIds,
              OrganizationStage.Engaged,
            );
            store.ui.commandMenu.setOpen(false);
          }}
        >
          Create new opportunity...
        </CommandItem>

        {/*<CommandItem*/}
        {/*  leftAccessory={<CoinsStacked01 />}*/}
        {/*  onSelect={() => {*/}
        {/*    // store.organizations.updateStage(*/}
        {/*    //   [organizations?.id as string],*/}
        {/*    //   OrganizationStage.Engaged,*/}
        {/*    // );*/}
        {/*    store.ui.commandMenu.setOpen(false);*/}
        {/*  }}*/}
        {/*>*/}
        {/*  Create new opportunity...*/}
        {/*</CommandItem>*/}
        {/*<CommandItem*/}
        {/*  leftAccessory={<Trophy01 />}*/}
        {/*  onSelect={() => {*/}
        {/*    store.ui.commandMenu.setType('AssignOwner');*/}
        {/*  }}*/}
        {/*>*/}
        {/*  Change onboarding stage*/}
        {/*</CommandItem>*/}
      </Command.List>
    </Command>
  );
});
