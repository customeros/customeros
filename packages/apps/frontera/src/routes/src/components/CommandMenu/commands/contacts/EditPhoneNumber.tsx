import { set } from 'lodash';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Command, CommandInput } from '@ui/overlay/CommandMenu';

export const EditPhoneNumber = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const contact = store.contacts.value.get(context.ids?.[0] as string);
  const phoneNumber = contact?.value.phoneNumbers?.[0]?.rawPhoneNumber ?? '';

  const label = `Contact - ${contact?.value.name}`;

  const handleChangePhoneNumber = (phoneNumber: string) => {
    if (!context.ids?.[0]) return;

    if (!contact) return;

    if (!contact?.value.phoneNumbers?.[0]?.id) {
      contact?.addPhoneNumber();
    } else {
      contact?.updatePhoneNumber();
    }
    contact?.update(
      (value) => {
        set(value, 'phoneNumbers[0].rawPhoneNumber', phoneNumber);

        return value;
      },
      { mutate: false },
    );
  };

  return (
    <Command>
      <CommandInput
        label={label}
        value={phoneNumber}
        placeholder='Edit phone number'
        onValueChange={(value) => handleChangePhoneNumber(value)}
      />
      <Command.List></Command.List>
    </Command>
  );
});
