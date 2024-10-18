import { useState } from 'react';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const RenameOrganizationProperty = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const entity = store.organizations.value.get(context.ids?.[0] as string);
  const label = `Organization - ${entity?.value?.name}`;
  const property = context.property as 'name' | 'website';
  const defaultValue = match({ property })
    .with({ property: 'name' }, () => entity?.value?.name ?? '')
    .with({ property: 'website' }, () => entity?.value?.website ?? '')
    .otherwise(() => '');

  const [name, setName] = useState(() => defaultValue ?? '');

  const handleClose = () => {
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.setType('OrganizationCommands');
  };

  const handleSelect = () => {
    if (!context.ids?.[0]) return;
    const property = context.property as 'name' | 'website';

    if (!entity || !property) return;

    if (!name.trim()?.length) {
      handleClose();

      return;
    }

    entity.value[property] = name;
    entity.commit();

    handleClose();
  };

  const placeholder = match({ property })
    .with({ property: 'name' }, () => 'Rename organization...')
    .with({ property: 'website' }, () => 'Edit website')
    .otherwise(() => '');

  return (
    <Command label={`Rename ${context.property}`}>
      <CommandInput
        value={name}
        label={label}
        placeholder={placeholder}
        onValueChange={(value) => setName(value)}
      />
      <Command.List>
        <CommandItem
          onSelect={handleSelect}
        >{`Rename ${context.property} to "${name}"`}</CommandItem>
      </Command.List>
    </Command>
  );
});
