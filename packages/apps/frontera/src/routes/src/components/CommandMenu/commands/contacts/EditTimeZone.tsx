import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { Check } from '@ui/media/icons/Check';
import { useStore } from '@shared/hooks/useStore';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';
import { timezoneOptions } from '@organization/components/Tabs/panels/PeoplePanel/util';

export const EditTimeZone = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const [search, setSearch] = useState('');

  const label = `Contact - ${
    store.contacts.value.get(context.ids?.[0] as string)?.value.name
  }`;

  const contact = store.contacts.value.get(context.ids?.[0] as string);

  const contactTimeZone = timezoneOptions.find(
    (v) => v.value === contact?.value?.timezone,
  );

  const handleChangeTimeZone = (timezone: string) => {
    if (!context.ids?.[0]) return;

    if (!contact) return;
    contact?.update((o) => {
      o.timezone = timezone;

      return o;
    });
    store.ui.commandMenu.setOpen(false);
    store.ui.commandMenu.setType('ContactHub');
  };

  return (
    <Command>
      <CommandInput
        label={label}
        onValueChange={setSearch}
        placeholder='Choose a timezone...'
        value={search ?? contactTimeZone?.label}
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
