import { useMemo } from 'react';

import { observer } from 'mobx-react-lite';
import { EditContactNameUseCase } from '@domain/usecases/command-menu/edit-contact-name.usecase';

import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const EditName = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const contact = store.contacts.value.get(context.ids?.[0] as string);

  const label = `Contact - ${contact?.name}`;

  const contactNameUseCase = useMemo(
    () => new EditContactNameUseCase(contact?.id || ''),
    [contact?.id],
  );

  const handleClose = () => {
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.setType('ContactCommands');
  };

  const handleChangeName = () => {
    contactNameUseCase.execute();
    handleClose();
  };

  return (
    <Command>
      <CommandInput
        label={label}
        value={contact?.name}
        placeholder='Edit name'
        onValueChange={(value) => contact?.setName(value)}
        onKeyDownCapture={(e) => {
          if (e.key === ' ') {
            e.stopPropagation();
          }
        }}
      />
      <Command.List>
        <CommandItem
          leftAccessory={<Edit03 />}
          onSelect={() => handleChangeName()}
        >{`Rename name to "${contact?.name}"`}</CommandItem>
      </Command.List>
    </Command>
  );
});
