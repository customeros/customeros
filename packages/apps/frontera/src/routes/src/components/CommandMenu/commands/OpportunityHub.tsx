import { observer } from 'mobx-react-lite';

// import { Delete } from '@ui/media/icons/Delete';
import { useStore } from '@shared/hooks/useStore';
import { PlusCircle } from '@ui/media/icons/PlusCircle';
// import { Command as CommandIcon } from '@ui/media/icons/Command';
import {
  // Kbd,
  CommandItem,
} from '@ui/overlay/CommandMenu';
import { CommandsContainer } from '@shared/components/CommandMenu/commands/shared';

export const OpportunityHub = observer(() => {
  const store = useStore();

  return (
    <CommandsContainer label='Opportunities'>
      <CommandItem
        leftAccessory={<PlusCircle />}
        onSelect={() => {
          store.ui.commandMenu.setType('ChooseOpportunityStage');
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
        Add new opportunity...
      </CommandItem>
    </CommandsContainer>
  );
});
