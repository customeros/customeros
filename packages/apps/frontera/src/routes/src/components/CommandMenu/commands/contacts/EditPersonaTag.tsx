import { useMemo, useState } from 'react';

import { CommandGroup } from 'cmdk';
import { observer } from 'mobx-react-lite';

import { Tag } from '@graphql/types';
import { DataSource } from '@graphql/types';
import { Plus } from '@ui/media/icons/Plus.tsx';
import { Check } from '@ui/media/icons/Check.tsx';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const EditPersonaTag = observer(() => {
  const store = useStore();
  const [search, setSearch] = useState('');

  const context = store.ui.commandMenu.context;

  const contact = store.contacts.value.get(context.ids?.[0] as string);

  const selectedIds = context.ids;
  const label =
    selectedIds?.length === 1
      ? `Contact - ${contact?.name}`
      : `${selectedIds?.length} contacts`;

  const handleSelect = (t: Tag) => () => {
    if (!context.ids?.[0]) return;

    if (!contact) return;

    if (selectedIds?.length === 1) {
      const foundIndex = contact.value?.tags?.findIndex(
        (e) => e.metadata.id === t.metadata.id,
      );

      if (foundIndex !== -1) {
        contact.value.tags?.splice(foundIndex || 0, 1);
      } else {
        contact.value.tags = contact.value.tags ?? [];
        contact.value.tags.push(t);
      }
      contact.commit();
    } else {
      store.contacts.updateTags(selectedIds, [t]);
    }
  };

  const handleCreateOption = (value: string) => {
    store.tags?.create({ name: value });
    contact?.value.tags?.push({
      id: value,
      name: value,
      metadata: {
        id: value,
        source: DataSource.Openline,
        sourceOfTruth: DataSource.Openline,
        appSource: 'organization',
        created: new Date().toISOString(),
        lastUpdated: new Date().toISOString(),
      },
      appSource: 'organization',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      source: DataSource.Openline,
    });
    contact?.commit();

    setSearch('');
  };

  const newSelectedTags = new Set(
    (contact?.value?.tags ?? []).map((tag) => tag.name),
  );
  const contactTags = useMemo(
    () => new Set((contact?.value?.tags ?? []).map((tag) => tag.name)),
    [],
  );

  const sortedTags = store.tags
    ?.toArray()
    .filter((e) => !!e.value.name)
    .sort((a, b) => {
      const aInOrg = contactTags.has(a.value.name);
      const bInOrg = contactTags.has(b.value.name);

      if (aInOrg && !bInOrg) return -1;
      if (!aInOrg && bInOrg) return 1;

      return 0;
    });
  const filteredTags = sortedTags?.filter((tag) =>
    tag.value.name.toLowerCase().includes(search.toLowerCase()),
  );

  useModKey('Enter', (e) => {
    e.stopPropagation();
    store.ui.commandMenu.setOpen(false);
  });

  return (
    <Command shouldFilter={false} label='Change or add tags...'>
      <CommandInput
        label={label}
        value={search}
        onValueChange={setSearch}
        placeholder='Edit persona tag...'
        onKeyDownCapture={(e) => {
          if (e.key === ' ') {
            e.stopPropagation();
          }

          if (e.metaKey && e.key === 'Enter') {
            e.stopPropagation();
            store.ui.commandMenu.setOpen(false);
          } else {
            handleSelect(search as unknown as Tag);
          }
        }}
      />

      <CommandGroup>
        <Command.List>
          {filteredTags?.map((tag) => (
            <CommandItem
              key={tag.id}
              onSelect={handleSelect(tag.value)}
              rightAccessory={
                newSelectedTags.has(tag.value.name) ? <Check /> : null
              }
            >
              {tag.value.name}
            </CommandItem>
          ))}
          {search && (
            <CommandItem
              leftAccessory={<Plus />}
              onSelect={() => handleCreateOption(search)}
            >
              <span className='text-gray-700 ml-1'>Create new tag:</span>
              <span className='text-gray-500 ml-1'>{search}</span>
            </CommandItem>
          )}
        </Command.List>
      </CommandGroup>
    </Command>
  );
});
