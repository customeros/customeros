import { observer } from 'mobx-react-lite';

import { Copy07 } from '@ui/media/icons/Copy07';
import { Archive } from '@ui/media/icons/Archive';
import { useStore } from '@shared/hooks/useStore';
import { Delete } from '@ui/media/icons/Delete.tsx';
import { Kbd, CommandKbd, CommandItem } from '@ui/overlay/CommandMenu';
import { CommandsContainer } from '@shared/components/CommandMenu/commands/shared';
import { organizationKeywords } from '@shared/components/CommandMenu/commands/organization';

export const SequencesBulkCommands = observer(() => {
  const store = useStore();
  const selectedIds = store.ui.commandMenu.context.ids;

  const label = `${selectedIds?.length} sequenced`;

  return (
    <CommandsContainer label={label}>
      <>
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
                <Delete className='size-3' />
              </Kbd>
            </>
          }
        >
          Archive sequences
        </CommandItem>

        <CommandItem
          leftAccessory={<Copy07 />}
          onSelect={() => {
            store.ui.commandMenu.setType('MergeConfirmationModal');
          }}
        >
          Merge
        </CommandItem>
      </>
    </CommandsContainer>
  );
});
