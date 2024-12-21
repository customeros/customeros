import { useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import { AddEmailCase } from '@domain/usecases/command-menu/add-email.usecase';

import { Mail02 } from '@ui/media/icons/Mail02';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

const addNewEmailCase = new AddEmailCase();

export const AddEmail = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const contact = store.contacts.value.get(context.ids?.[0] as string);

  const label = `Contact - ${contact?.name}`;

  useEffect(() => {
    if (contact) {
      addNewEmailCase.setEntity(contact);
    }
  }, [contact?.id]);

  return (
    <Command shouldFilter={false}>
      <CommandInput
        label={label}
        placeholder='Add new email'
        value={addNewEmailCase.inputValue}
        onValueChange={(value) => {
          addNewEmailCase.setInputValue(value);
        }}
        onKeyDownCapture={(e) => {
          if (e.key === '') {
            e.stopPropagation();
          }
        }}
      />
      <Command.List>
        <CommandItem
          leftAccessory={<Mail02 />}
          onSelect={() => {
            addNewEmailCase.submit();
            store.ui.commandMenu.setOpen(false);
            store.ui.commandMenu.setType('ContactCommands');
          }}
        >
          {`Add new email ${addNewEmailCase.inputValue}`}
        </CommandItem>
      </Command.List>
    </Command>
  );
});
