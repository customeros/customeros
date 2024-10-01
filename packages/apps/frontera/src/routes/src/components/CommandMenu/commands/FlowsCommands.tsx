import { observer } from 'mobx-react-lite';

import { Edit03 } from '@ui/media/icons/Edit03';
import { Delete } from '@ui/media/icons/Delete';
import { Archive } from '@ui/media/icons/Archive';
import { useStore } from '@shared/hooks/useStore';
import { Columns03 } from '@ui/media/icons/Columns03.tsx';
import { ArrowBlockUp } from '@ui/media/icons/ArrowBlockUp';
import { Kbd, CommandKbd, CommandItem } from '@ui/overlay/CommandMenu';
import { CommandsContainer } from '@shared/components/CommandMenu/commands/shared';
import { flowKeywords } from '@shared/components/CommandMenu/commands/flows/keywords.ts';
import { UpdateStatusSubItemGroup } from '@shared/components/CommandMenu/commands/flows/UpdateStatusSubItemGroup.tsx';

export const FlowsCommands = observer(() => {
  const store = useStore();
  const id = (store.ui.commandMenu.context.ids as string[])?.[0];
  const sequence = store.flows.value.get(id);
  const label = `Flow - ${sequence?.value.name}`;

  return (
    <CommandsContainer label={label}>
      <>
        <CommandItem
          leftAccessory={<Edit03 />}
          keywords={[
            'rename',
            'org',
            'Sequence',
            'company',
            'update',
            'edit',
            'change',
          ]}
          rightAccessory={
            <>
              <Kbd>
                <ArrowBlockUp className='size-3' />
              </Kbd>
              <Kbd>R</Kbd>
            </>
          }
          onSelect={() => {
            store.ui.commandMenu.setType('RenameFlow');
            store.ui.commandMenu.setContext({
              ...store.ui.commandMenu.context,
              property: 'name',
            });
          }}
        >
          Rename flow
        </CommandItem>

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
          Change flow status...
        </CommandItem>
        <UpdateStatusSubItemGroup />
        <CommandItem
          leftAccessory={<Archive />}
          keywords={flowKeywords.archive_flow}
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
          Archive flow
        </CommandItem>
      </>
    </CommandsContainer>
  );
});
