import { observer } from 'mobx-react-lite';

import { Tag01 } from '@ui/media/icons/Tag01';
import { useStore } from '@shared/hooks/useStore';
import { Kbd, CommandItem } from '@ui/overlay/CommandMenu';
import { ArrowBlockUp } from '@ui/media/icons/ArrowBlockUp.tsx';
import { CommandsContainer } from '@shared/components/CommandMenu/commands/shared';

export const ContactHub = observer(() => {
  const store = useStore();

  const context = store.ui.commandMenu.context;
  const contact = store.contacts.value.get(context.ids?.[0] as string);

  const label = `Contact - ${contact?.value.name}`;

  return (
    <CommandsContainer label={label}>
      <CommandItem
        leftAccessory={<Tag01 />}
        onSelect={() => {
          store.ui.commandMenu.setType('EditPersonaTag');
        }}
        rightAccessory={
          <>
            <Kbd>
              <ArrowBlockUp className='text-inherit size-3' />
            </Kbd>
            <Kbd>T</Kbd>
          </>
        }
      >
        Edit persona tag...
      </CommandItem>
    </CommandsContainer>
  );
});
