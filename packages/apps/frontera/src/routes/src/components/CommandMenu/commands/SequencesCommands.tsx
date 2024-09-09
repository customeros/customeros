import { observer } from 'mobx-react-lite';

import { Edit03 } from '@ui/media/icons/Edit03';
import { Delete } from '@ui/media/icons/Delete';
import { Archive } from '@ui/media/icons/Archive';
import { useStore } from '@shared/hooks/useStore';
import { Columns03 } from '@ui/media/icons/Columns03.tsx';
import { ArrowBlockUp } from '@ui/media/icons/ArrowBlockUp';
import { Kbd, CommandKbd, CommandItem } from '@ui/overlay/CommandMenu';
import { CommandsContainer } from '@shared/components/CommandMenu/commands/shared';
import { sequenceKeywords } from '@shared/components/CommandMenu/commands/sequences/keywords';

export const SequenceCommands = observer(() => {
  const store = useStore();
  const id = (store.ui.commandMenu.context.ids as string[])?.[0];
  const sequence = store.flowSequences.value.get(id);
  const label = `Sequence - ${sequence?.value.name}`;

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
            store.ui.commandMenu.setType('RenameSequence');
            store.ui.commandMenu.setContext({
              ...store.ui.commandMenu.context,
              property: 'name',
            });
          }}
        >
          Rename Sequence
        </CommandItem>

        <CommandItem
          leftAccessory={<Columns03 />}
          keywords={sequenceKeywords.status_update}
          onSelect={() => {
            store.ui.commandMenu.setType('ChangeSequenceStatus');
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
          Change sequence status
        </CommandItem>
        <CommandItem
          leftAccessory={<Archive />}
          keywords={sequenceKeywords.archive_sequence}
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
          Archive sequence
        </CommandItem>
      </>
    </CommandsContainer>
  );
});
