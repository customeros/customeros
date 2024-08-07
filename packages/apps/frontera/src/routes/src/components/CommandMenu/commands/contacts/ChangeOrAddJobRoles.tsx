import { useState } from 'react';

import set from 'lodash/set';
import { observer } from 'mobx-react-lite';

import { Check } from '@ui/media/icons/Check';
import { useStore } from '@shared/hooks/useStore';
import { SelectOption } from '@shared/types/SelectOptions';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

const roleOptions = [
  { value: 'Decision Maker', label: 'Decision Maker' },
  { value: 'Influencer', label: 'Influencer' },
  { value: 'User', label: 'User' },
  { value: 'Stakeholder', label: 'Stakeholder' },
  { value: 'Gatekeeper', label: 'Gatekeeper' },
  { value: 'Champion', label: 'Champion' },
  { value: 'Data Owner', label: 'Data Owner' },
];

export const ChangeOrAddJobRoles = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;

  const contact = store.contacts.value.get(context.ids?.[0] as string);
  const label = `Contact - ${contact?.value?.name}`;
  const [search, setSearch] = useState('');

  const handleSelect = (opt: SelectOption[]) => {
    if (!context.ids?.[0] || !contact) return;

    contact.update((value) => {
      const selectedValues = opt.map((v) => v.value).join(',');

      set(value, 'jobRoles[0].description', selectedValues);

      return value;
    });
  };

  return (
    <Command label='Change or add job roles...'>
      <CommandInput
        label={label}
        value={search}
        onValueChange={setSearch}
        placeholder='Change or add job roles...'
      />

      <Command.List>
        {roleOptions.map((role, idx) => {
          const selectedDescriptions =
            contact?.value?.jobRoles?.[0]?.description?.split(',') || [];

          const isSelected = selectedDescriptions.includes(role.value);

          return (
            <CommandItem
              key={idx}
              rightAccessory={isSelected ? <Check /> : undefined}
              onSelect={() => {
                const newSelections = isSelected
                  ? selectedDescriptions.filter((desc) => desc !== role.value)
                  : [...selectedDescriptions, role.value];

                const newOptions = roleOptions.filter((r) =>
                  newSelections.includes(r.value),
                );

                handleSelect(newOptions);
              }}
            >
              {role.label}
            </CommandItem>
          );
        })}
      </Command.List>
    </Command>
  );
});
