import { useState } from 'react';

import { set } from 'lodash';
import { observer } from 'mobx-react-lite';

import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const EditEmail = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const contact = store.contacts.value.get(context.ids?.[0] as string);
  const [value, setValue] = useState(
    () => contact?.value.emails?.[0]?.email ?? '',
  );
  const emailAdress = contact?.value?.emails?.[0]?.email ?? '';

  const label = `Contact - ${contact?.value.name}`;

  const handleChangeEmail = () => {
    if (!context.ids?.[0]) return;

    if (!contact) return;

    if (emailAdress?.length === 0) {
      contact?.addEmail();
    } else {
      contact?.update((o) => {
        o.emails[0].email = value;

        return o;
      });
    }

    if (value.length === 0) {
      contact?.removeEmail();
    }
    contact?.update(
      (value) => {
        set(value, 'emails[0].email', value);

        return value;
      },
      { mutate: false },
    );
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.setType('ContactCommands');
  };

  return (
    <Command>
      <CommandInput
        label={label}
        placeholder='Edit email'
        value={emailAdress || ''}
        onValueChange={(value) => setValue(value)}
      />
      <Command.List>
        <CommandItem
          leftAccessory={<Edit03 />}
          onSelect={handleChangeEmail}
        >{`Rename email to "${value}"`}</CommandItem>
      </Command.List>
    </Command>
  );
});
