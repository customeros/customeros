import { useState } from 'react';

import { set } from 'lodash';
import { observer } from 'mobx-react-lite';

import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const EditPhoneNumber = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const contact = store.contacts.value.get(context.ids?.[0] as string);
  const phoneNumber = contact?.value.phoneNumbers?.[0]?.rawPhoneNumber ?? '';

  const [number, setPhone] = useState(() => phoneNumber);

  const label = `Contact - ${contact?.name}`;

  const handleChangePhoneNumber = () => {
    if (!contact) return;

    set(contact.value, 'phoneNumbers[0].rawPhoneNumber', number);
    contact.commit();

    if (!contact?.value.phoneNumbers?.[0]?.id) {
      contact?.addPhoneNumber();
    } else {
      contact?.updatePhoneNumber();
    }

    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.setType('ContactCommands');
  };

  return (
    <Command>
      <CommandInput
        label={label}
        value={number}
        placeholder='Edit phone number'
        onValueChange={(value) => setPhone(value)}
        onKeyDownCapture={(e) => {
          if (e.key === ' ') {
            e.stopPropagation();
          }
        }}
      />
      <Command.List>
        <CommandItem
          leftAccessory={<Edit03 />}
          onSelect={handleChangePhoneNumber}
        >
          {`Edit phone number to "${number}"`}
        </CommandItem>
      </Command.List>
    </Command>
  );
});
