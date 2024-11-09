import { useState } from 'react';

import { set } from 'lodash';
import { observer } from 'mobx-react-lite';

import { Check } from '@ui/media/icons/Check';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';
import { timezoneOptions } from '@organization/components/Tabs/panels/PeoplePanel/util';

export const EditTimeZone = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const [search, setSearch] = useState('');

  const contact = store.contacts.value.get(context.ids?.[0] as string);

  const selectedIds = context.ids;
  const label =
    selectedIds?.length === 1
      ? `Contact - ${contact?.name}`
      : `${selectedIds?.length} contacts`;

  const contactTimeZone = timezoneOptions.find(
    (v) => v.value === contact?.value?.timezone,
  );

  const handleChangeTimeZone = (timezone: string) => {
    if (!context.ids?.[0]) return;

    if (!contact) return;

    if (selectedIds?.length === 1) {
      set(contact.value, 'timezone', timezone);
      contact.commit();
    } else {
      selectedIds.forEach((id) => {
        const contact = store.contacts.value.get(id);

        if (contact) {
          set(contact.value, 'timezone', timezone);
          contact.commit();
        }
      });
    }
    store.ui.commandMenu.setOpen(false);
  };

  return (
    <Command>
      <CommandInput
        label={label}
        onValueChange={setSearch}
        placeholder='Choose a timezone...'
        value={search ?? contactTimeZone?.label}
        onKeyDownCapture={(e) => {
          if (e.key === ' ') {
            e.stopPropagation();
          }
        }}
      />
      <Command.List>
        {timezoneOptions
          .filter((v) => v.label.toLowerCase().includes(search.toLowerCase()))
          .map((v) => (
            <CommandItem
              key={v.value}
              onSelect={() => {
                handleChangeTimeZone(v.value);
              }}
              rightAccessory={
                v.value === contact?.value?.timezone ? <Check /> : undefined
              }
            >
              {v.label}
            </CommandItem>
          ))}
      </Command.List>
    </Command>
  );
});
