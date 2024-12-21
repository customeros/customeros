import { useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import { AddLinkedin } from '@domain/usecases/command-menu/add-linkedin.usecase';

import { useStore } from '@shared/hooks/useStore';
import { Linkedin } from '@ui/media/icons/Linkedin';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

const addNewLinkedinUrl = new AddLinkedin();

export const AddLinkedinUrl = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const contact = store.contacts.value.get(context.ids?.[0] as string);

  const label = `Contact - ${contact?.name}`;

  useEffect(() => {
    if (contact) {
      addNewLinkedinUrl.setEntity(contact);
    }
  }, [contact?.id]);

  return (
    <Command shouldFilter={false}>
      <CommandInput
        label={label}
        placeholder='Add LinkedIn'
        value={addNewLinkedinUrl.inputValue}
        onValueChange={(value) => {
          addNewLinkedinUrl.setInputValue(value);
        }}
        onKeyDownCapture={(e) => {
          if (e.key === '') {
            e.stopPropagation();
          }
        }}
      />
      <Command.List>
        <CommandItem
          leftAccessory={<Linkedin />}
          onSelect={() => {
            addNewLinkedinUrl.setLinkedInUrl();
            store.ui.commandMenu.setOpen(false);
            store.ui.commandMenu.setType('ContactCommands');
          }}
        >
          {`Add LinkedIn url: ${addNewLinkedinUrl.inputValue}`}
        </CommandItem>
      </Command.List>
    </Command>
  );
});
