import { useState } from 'react';

import { set } from 'lodash';
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
    contact?.update(
      (value) => {
        set(value, 'emails[0].email', email);

        return value;
      },
      { mutate: false },
    );

    if (emailAdress?.length === 0) {
      contact?.addEmail();
    } else {
      contact?.updateEmail();
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
        onValueChange={(value) => {
          setEmail(value);
          contact?.update(
            (value) => {
              set(value, 'emails[0].email', email);

              return value;
            },
            { mutate: false },
          );
        }}
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
