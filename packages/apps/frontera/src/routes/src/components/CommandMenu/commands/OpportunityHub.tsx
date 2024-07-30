import { observer } from 'mobx-react-lite';

import { Delete } from '@ui/media/icons/Delete';
import { PlusCircle } from '@ui/media/icons/PlusCircle';
import { Command as CommandIcon } from '@ui/media/icons/Command';
import {
  Kbd,
  Command,
  CommandItem,
  CommandInput,
} from '@ui/overlay/CommandMenu';

export const OpportunityHub = observer(() => {
  return (
    <Command>
      <CommandInput
        label='Opportunities Hub'
        placeholder='Type a command or search'
      />
      <Command.List>
        <CommandItem
          onSelect={() => {}}
          leftAccessory={<PlusCircle />}
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
          Add new opportunity...
        </CommandItem>
      </Command.List>
    </Command>
  );
});
