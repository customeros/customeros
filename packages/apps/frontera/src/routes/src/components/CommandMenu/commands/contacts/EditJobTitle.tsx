import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { Edit03 } from '@ui/media/icons/Edit03';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const EditJobTitle = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const selectedIds = context.ids;
  const contact = store.contacts.value.get(context.ids?.[0] as string);

  const contactJobTitle =
    contact?.value?.jobRoles?.map((jobRole) => jobRole.jobTitle)?.[0] ?? '';
  const [name, setName] = useState(() => contactJobTitle);

  const label =
    selectedIds?.length === 1
      ? `Contact - ${contact?.name}`
      : `${selectedIds?.length} contacts`;

  const handleChangeJobTitle = () => {
    if (!contact) return;

    if (selectedIds?.length === 1) {
      contact.value.jobRoles.map((jobRole) => {
        jobRole.jobTitle = name;
      });
    } else {
      selectedIds.forEach((id) => {
        const contact = store.contacts.value.get(id);

        if (contact) {
          contact.value.jobRoles.map((jobRole) => {
            jobRole.jobTitle = name;
          });
        }
      });
    }
    contact.commit();
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.setType('ContactCommands');
  };

  return (
    <Command label={label}>
      <CommandInput
        label={label}
        value={name || ''}
        placeholder='Edit job title'
        onValueChange={(value) => setName(value)}
        onKeyDownCapture={(e) => {
          if (e.key === ' ') {
            e.stopPropagation();
          }
        }}
      />
      <Command.List>
        <CommandItem
          leftAccessory={<Edit03 />}
          onSelect={handleChangeJobTitle}
        >{`Rename job title to "${name}"`}</CommandItem>
      </Command.List>
    </Command>
  );
});
