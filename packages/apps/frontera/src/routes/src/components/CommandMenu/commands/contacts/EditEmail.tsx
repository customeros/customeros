import { set } from 'lodash';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Command, CommandInput } from '@ui/overlay/CommandMenu';

export const EditEmail = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const contact = store.contacts.value.get(context.ids?.[0] as string);
  const emailAdress = contact?.value?.emails?.[0]?.email ?? '';

  const label = `Contact - ${contact?.value.name}`;

  const handleChangeEmail = (email: string) => {
    if (!context.ids?.[0]) return;

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
    contact?.update(
      (value) => {
        set(value, 'emails[0].email', email);

        return value;
      },
      { mutate: false },
    );
  };

  return (
    <Command>
      <CommandInput
        label={label}
        placeholder='Edit email'
        value={emailAdress || ''}
        onValueChange={(value) => handleChangeEmail(value)}
      />
      <Command.List></Command.List>
    </Command>
  );
});
