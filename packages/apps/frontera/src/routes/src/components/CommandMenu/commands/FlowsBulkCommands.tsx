import { observer } from 'mobx-react-lite';

import { Archive } from '@ui/media/icons/Archive';
import { useStore } from '@shared/hooks/useStore';
import { Delete } from '@ui/media/icons/Delete.tsx';
import { Columns03 } from '@ui/media/icons/Columns03.tsx';
import { ArrowBlockUp } from '@ui/media/icons/ArrowBlockUp.tsx';
import { Kbd, CommandKbd, CommandItem } from '@ui/overlay/CommandMenu';
import { CommandsContainer } from '@shared/components/CommandMenu/commands/shared';
import { flowKeywords } from '@shared/components/CommandMenu/commands/flows/keywords';
import { organizationKeywords } from '@shared/components/CommandMenu/commands/organization';

export const FlowsBulkCommands = observer(() => {
  const store = useStore();
  const selectedIds = store.ui.commandMenu.context.ids;

  const label = `${selectedIds?.length} sequences`;

  return (
    <CommandsContainer label={label}>
      <>
        <CommandItem
          leftAccessory={<Columns03 />}
          keywords={flowKeywords.status_update}
          onSelect={() => {
            store.ui.commandMenu.setType('ChangeFlowStatus');
          }}
          rightAccessory={
            <>
              <Kbd>
                <ArrowBlockUp className='size-3' />
              </Kbd>
              <Kbd>S</Kbd>
            </>
          }
        >
          Change sequence status...
        </CommandItem>
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
      </>
    </CommandsContainer>
  );
});
