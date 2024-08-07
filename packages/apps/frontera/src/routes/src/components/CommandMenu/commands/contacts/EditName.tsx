import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Command, CommandInput } from '@ui/overlay/CommandMenu';

export const EditName = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const contact = store.contacts.value.get(context.ids?.[0] as string);
  const contactName = contact?.value?.name ?? '';

  const label = `Contact - ${contact?.value.name}`;

  const handleChangeName = (name: string) => {
    if (!context.ids?.[0]) return;

    if (!contact) return;
    contact?.update((o) => {
      o.name = name;

      return o;
    });
  };

  return (
    <Command>
      <CommandInput
        label={label}
        placeholder='Edit name'
        value={contactName || ''}
        onValueChange={(value) => handleChangeName(value)}
      />
      <Command.List></Command.List>
    </Command>
  );
});
