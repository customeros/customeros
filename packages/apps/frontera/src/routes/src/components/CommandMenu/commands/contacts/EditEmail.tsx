import { useMemo } from 'react';

import { set } from 'lodash';
import { observer } from 'mobx-react-lite';

import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const EditEmail = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const selectedId = store.ui.selectionId;

  const contact = store.contacts.value.get(context.ids?.[0] as string);
  const oldEmail = useMemo(
    () => contact?.value?.emails?.[selectedId ?? 0]?.email,
    [contact?.isLoading, selectedId],
  );
  const emailAdress = contact?.value?.emails?.[selectedId ?? 0]?.email ?? '';

  const label = `Contact - ${contact?.name}`;

  const handleSaveEmail = () => {
    contact?.updateEmail(oldEmail ?? '', selectedId ?? 0);
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.setType('ContactCommands');
  };

  return (
    <Command>
      <CommandInput
        label={label}
        value={emailAdress}
        placeholder='Edit email'
        onKeyDownCapture={(e) => {
          if (e.key === ' ') {
            e.stopPropagation();
          }
        }}
        onValueChange={(newValue) => {
          contact?.update(
            (value) => {
              set(value, ['emails', selectedId ?? 0, 'email'], newValue);

              return value;
            },
            { mutate: false },
          );
        }}
      />
      <Command.List>
        <CommandItem
          leftAccessory={<Edit03 />}
          onSelect={handleSaveEmail}
        >{`Rename email to "${emailAdress}"`}</CommandItem>
      </Command.List>
    </Command>
  );
});
