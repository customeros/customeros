import { useState } from 'react';

import { useKey } from 'rooks';
import { observer } from 'mobx-react-lite';

import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const EditEmail = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const contact = store.contacts.value.get(context.ids?.[0] as string);
  const [email, setEmail] = useState(
    () => contact?.value.emails?.[0]?.email ?? '',
  );
  const emailAdress = contact?.value?.emails?.[0]?.email ?? '';

  const label = `Contact - ${contact?.value.name}`;

  const handleChangeEmail = () => {
    if (!contact) return;

    if (emailAdress?.length === 0) {
      contact?.addEmail();
    } else {
      contact?.update((o) => {
        o.emails[0].email = email;

        return o;
      });
    }

    if (email.length === 0) {
      contact?.removeEmail();
    }

    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.setType('ContactCommands');
  };

  useKey('Enter', handleChangeEmail);

  return (
    <Command>
      <CommandInput
        label={label}
        value={email || ''}
        placeholder='Edit email'
        onValueChange={(value) => setEmail(value)}
      />
      <Command.List>
        <CommandItem
          leftAccessory={<Edit03 />}
          onSelect={handleChangeEmail}
        >{`Rename email to "${email}"`}</CommandItem>
      </Command.List>
    </Command>
  );
});
