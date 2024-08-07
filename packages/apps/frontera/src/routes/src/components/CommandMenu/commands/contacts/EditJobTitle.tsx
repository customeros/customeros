import { useState, useEffect } from 'react';

import { set } from 'lodash';
import { CommandList } from 'cmdk';
import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Command, CommandInput } from '@ui/overlay/CommandMenu';

export const EditJobTitle = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const [name, setName] = useState('');
  const selectedIds = context.ids;
  const contact = store.contacts.value.get(context.ids?.[0] as string);

  const contactJobTitle =
    contact?.value?.jobRoles?.map((jobRole) => jobRole.jobTitle)?.[0] ?? '';

  const label =
    selectedIds?.length === 1
      ? `Contact - ${contact?.value.name}`
      : `${selectedIds?.length} contacts`;

  const handleChangeJobTitle = (jobTitle: string) => {
    if (!context.ids?.[0]) return;

    if (!contact) return;

    if (selectedIds?.length === 1) {
      contact?.update((value) => {
        set(value, 'jobRoles[0].jobTitle', jobTitle);

        return value;
      });
    } else {
      selectedIds.forEach((id) => {
        const contact = store.contacts.value.get(id);

        if (contact) {
          contact.update((value) => {
            value.jobRoles[0].jobTitle = name;

            return value;
          });
        }
      });
    }
  };

  useEffect(() => {
    if (name.length !== 0) {
      handleChangeJobTitle(name);
    }
  }, [name]);

  return (
    <Command label={label}>
      <CommandInput
        label={label}
        placeholder='Edit job title'
        value={selectedIds.length === 1 ? contactJobTitle : name}
        onValueChange={
          selectedIds.length === 1 ? handleChangeJobTitle : setName
        }
      />
      <CommandList></CommandList>
    </Command>
  );
});
