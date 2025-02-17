import { useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import { EditEmailCase } from '@domain/usecases/command-menu/edit-email.usecase';

import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

const editEmailUseCase = EditEmailCase.getInstance();

export const EditEmail = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;

  const contact = store.contacts.value.get(context.ids?.[0] as string);

  const label = `Contact - ${contact?.name}`;

  useEffect(() => {
    if (contact) {
      editEmailUseCase.setEntity(contact);
    }
  }, [contact?.id]);

  if (!contact) return;

  useEffect(() => {
    if (contact.id) {
      editEmailUseCase.setOldEmail(editEmailUseCase.email);
    }
  }, [contact?.id]);

  return (
    <Command shouldFilter={false}>
      <CommandInput
        label={label}
        placeholder={'Edit email'}
        value={editEmailUseCase.email}
        onValueChange={(newValue) => {
          editEmailUseCase.setEmail(newValue);
        }}
        onKeyDownCapture={(e) => {
          if (e.key === ' ') {
            e.stopPropagation();
          }
        }}
      />

      <Command.List>
        <CommandItem
          leftAccessory={<Edit03 />}
          onSelect={() => {
            editEmailUseCase.submit();
            store.ui.commandMenu.setOpen(false);
            store.ui.commandMenu.setType('ContactCommands');
          }}
        >
          {`Rename email to "${editEmailUseCase.email}"`}
        </CommandItem>
      </Command.List>
    </Command>
  );
});
