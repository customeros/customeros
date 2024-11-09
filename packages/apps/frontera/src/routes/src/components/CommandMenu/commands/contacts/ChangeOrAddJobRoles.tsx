import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { Check } from '@ui/media/icons/Check';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
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
  const [search, setSearch] = useState('');

  const context = store.ui.commandMenu.context;

  const contact = store.contacts.value.get(context.ids?.[0] as string);
  const selectedIds = context.ids;

  const label =
    selectedIds?.length === 1
      ? `Contact - ${contact?.name}`
      : `${selectedIds?.length} contacts`;

  const handleSelect = (opt: SelectOption[]) => {
    if (!context.ids?.[0] || !contact) return;

    if (selectedIds?.length === 1) {
      contact.value.jobRoles.map((jobRole) => {
        jobRole.description = opt.map((v) => v.value).join(',');
      });
    } else {
      selectedIds.forEach((id) => {
        const contact = store.contacts.value.get(id);

        if (contact) {
          contact.value.jobRoles.map((jobRole) => {
            jobRole.description = opt.map((v) => v.value).join(',');
          });
        }
      });
    }
    contact.commit();
  };

  useModKey('Enter', () => {
    store.ui.commandMenu.setOpen(false);
  });

  return (
    <Command label='Edit job roles...'>
      <CommandInput
        label={label}
        value={search}
        onValueChange={setSearch}
        placeholder='Edit job roles...'
        onKeyDownCapture={(e) => {
          if (e.key === ' ') {
            e.stopPropagation();
          }

          if (e.metaKey && e.key === 'Enter') {
            store.ui.commandMenu.setOpen(false);
          } else {
            handleSelect([{ value: search, label: search }]);
          }
        }}
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
