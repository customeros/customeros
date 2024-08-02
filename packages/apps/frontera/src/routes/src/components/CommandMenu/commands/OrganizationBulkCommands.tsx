import { observer } from 'mobx-react-lite';

import { User01 } from '@ui/media/icons/User01';
import { Archive } from '@ui/media/icons/Archive';
import { useStore } from '@shared/hooks/useStore';
import { Copy07 } from '@ui/media/icons/Copy07.tsx';
import { Columns03 } from '@ui/media/icons/Columns03';
import { OrganizationRelationship } from '@graphql/types';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

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
        {/*<CommandItem*/}
        {/*  leftAccessory={<Tag01 />}*/}
        {/*  onSelect={() => {*/}
        {/*    store.ui.commandMenu.setType('ChangeTags');*/}
        {/*  }}*/}
        {/*  // rightAccessory={*/}
        {/*  //   <>*/}
        {/*  //     <Kbd className='px-1.5'>*/}
        {/*  //       <ArrowBlockUp className='size-3' />*/}
        {/*  //     </Kbd>*/}
        {/*  //     <Kbd className='px-1.5'>T</Kbd>*/}
        {/*  //   </>*/}
        {/*  // }*/}
        {/*>*/}
        {/*  Change or add tags...*/}
        {/*</CommandItem>*/}

        {/*{organizations?.some(*/}
        {/*  (organization) => !!organization?.value?.tags?.length,*/}
        {/*) && (*/}
        {/*  <CommandItem*/}
        {/*    leftAccessory={<Tag01 />}*/}
        {/*    onSelect={() => {*/}
        {/*      // todo*/}
        {/*      store.ui.commandMenu.setOpen(false);*/}
        {/*    }}*/}
        {/*  >*/}
        {/*    Remove tags*/}
        {/*  </CommandItem>*/}
        {/*)}*/}

        {/*<CommandItem*/}
        {/*  leftAccessory={<AlignHorizontalCentre02 />}*/}
        {/*  onSelect={() => {*/}
        {/*    store.ui.commandMenu.setType('ChangeRelationship');*/}
        {/*  }}*/}
        {/*>*/}
        {/*  Change relationship...*/}
        {/*</CommandItem>*/}

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

        <CommandItem
          leftAccessory={<Archive />}
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

        {/*<CommandItem*/}
        {/*  leftAccessory={<Activity />}*/}
        {/*  onSelect={() => {*/}
        {/*    store.ui.commandMenu.setType('UpdateHealthStatus');*/}
        {/*  }}*/}
        {/*>*/}
        {/*  Change health status...*/}
        {/*</CommandItem>*/}

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
