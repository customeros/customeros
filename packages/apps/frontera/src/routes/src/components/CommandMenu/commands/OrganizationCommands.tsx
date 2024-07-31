import { observer } from 'mobx-react-lite';

import { Edit03 } from '@ui/media/icons/Edit03';
import { User01 } from '@ui/media/icons/User01';
import { Tag01 } from '@ui/media/icons/Tag01.tsx';
import { Archive } from '@ui/media/icons/Archive';
import { useStore } from '@shared/hooks/useStore';
import { Delete } from '@ui/media/icons/Delete.tsx';
import { Globe01 } from '@ui/media/icons/Globe01.tsx';
import { Columns03 } from '@ui/media/icons/Columns03';
import { Trophy01 } from '@ui/media/icons/Trophy01.tsx';
import { Activity } from '@ui/media/icons/Activity.tsx';
import { Calculator } from '@ui/media/icons/Calculator';
import { ArrowBlockUp } from '@ui/media/icons/ArrowBlockUp.tsx';
import { CoinsStacked01 } from '@ui/media/icons/CoinsStacked01.tsx';
import { Command as CommandIcon } from '@ui/media/icons/Command.tsx';
import { CurrencyDollarCircle } from '@ui/media/icons/CurrencyDollarCircle';
import { AlignHorizontalCentre02 } from '@ui/media/icons/AlignHorizontalCentre02.tsx';
import {
  Kbd,
  Command,
  CommandItem,
  CommandInput,
} from '@ui/overlay/CommandMenu';

export const OrganizationCommands = observer(() => {
  const store = useStore();
  const organization = store.organizations.value.get(
    store.ui.commandMenu.context.id as string,
  );
  const label = `Organization - ${organization?.value.name}`;

  return (
    <Command>
      <CommandInput label={label} placeholder='Type a command or search' />
      <Command.List>
        <CommandItem
          onSelect={() => {}}
          leftAccessory={<User01 />}
          rightAccessory={<Kbd className='px-1.5'>P</Kbd>}
        >
          Add contact
        </CommandItem>

        <CommandItem
          leftAccessory={<Tag01 />}
          onSelect={() => {
            store.ui.commandMenu.setType('ChangeTags');
          }}
          rightAccessory={
            <>
              <Kbd className='px-1.5'>
                <ArrowBlockUp className='size-3' />
              </Kbd>
              <Kbd className='px-1.5'>T</Kbd>
            </>
          }
        >
          Change or add tags
        </CommandItem>
        <CommandItem
          leftAccessory={<Tag01 />}
          onSelect={() => {
            organization?.update((value) => {
              value.tags = [];
              return value;
            });
          }}
        >
          Remove tags
        </CommandItem>
        <CommandItem
          leftAccessory={<Edit03 />}
          rightAccessory={
            <>
              <Kbd className='px-1.5'>
                <ArrowBlockUp className='size-3' />
              </Kbd>
              <Kbd className='px-1.5'>R</Kbd>
            </>
          }
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
          onSelect={() => {
            store.ui.commandMenu.setType('ChangeRelationship');
          }}
        >
          Change relationship...
        </CommandItem>
        <CommandItem
          leftAccessory={<Columns03 />}
          onSelect={() => {
            store.ui.commandMenu.setType('ChangeStage');
          }}
        >
          Change org stage...
        </CommandItem>
        <CommandItem
          leftAccessory={<Archive />}
          onSelect={() => {
            store.organizations.hide([organization?.id as string]);
            store.ui.commandMenu.setOpen(false);
          }}
          rightAccessory={
            <>
              <Kbd className='px-1.5'>
                <CommandIcon className='size-3' />
              </Kbd>
              <Kbd className='px-1.5'>
                <Delete className='size-3' />
              </Kbd>
            </>
          }
        >
          Archive org
        </CommandItem>

        <CommandItem
          leftAccessory={<Activity />}
          onSelect={() => {
            store.ui.commandMenu.setType('UpdateHealthStatus');
          }}
        >
          Change health status...
        </CommandItem>

        <CommandItem
          leftAccessory={<User01 />}
          onSelect={() => {
            store.ui.commandMenu.setType('AssignOwner');
          }}
          rightAccessory={
            <>
              <Kbd className='px-1.5'>
                <ArrowBlockUp className='size-3' />
              </Kbd>
              <Kbd className='px-1.5'>O</Kbd>
            </>
          }
        >
          Assign owner...
        </CommandItem>

        {/*<CommandItem*/}
        {/*  leftAccessory={<CoinsStacked01 />}*/}
        {/*  rightAccessory={<Kbd className='px-1.5'>O</Kbd>}*/}
        {/*  onSelect={() => {*/}
        {/*    store.ui.commandMenu.setType('AssignOwner');*/}
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
