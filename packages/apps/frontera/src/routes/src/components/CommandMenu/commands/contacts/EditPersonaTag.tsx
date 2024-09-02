import { useMemo, useState } from 'react';

import { CommandGroup } from 'cmdk';
import { observer } from 'mobx-react-lite';

import { DataSource } from '@graphql/types';
import { Plus } from '@ui/media/icons/Plus.tsx';
import { Tag as TagType } from '@graphql/types';
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
      ? `Contact - ${contact?.value.name}`
      : `${selectedIds?.length} contacts`;

  const handleSelect = (t: TagType) => () => {
    if (!context.ids?.[0]) return;

    if (!contact) return;

    if (selectedIds?.length === 1) {
      contact?.update((o) => {
        const existingIndex = o.tags?.find((e) => e.name === t.name);

        if (existingIndex) {
          const newTags = o.tags?.filter((e) => e.name !== t.name);

          o.tags = newTags;
        }

        if (!existingIndex) {
          o.tags = [...(o.tags ?? []), t];
        }

        return o;
      });
    } else {
      selectedIds.forEach((id) => {
        const contact = store.contacts.value.get(id);

        if (contact) {
          contact.update((o) => {
            const existingIndex = o.tags?.find((e) => e.name === t.name);

            if (existingIndex) {
              const newTags = o.tags?.filter((e) => e.name !== t.name);

              o.tags = newTags;
            }

            if (!existingIndex) {
              o.tags = [...(o.tags ?? []), t];
            }

            return o;
          });
        }
      });
    }
  };

  const handleCreateOption = (value: string) => {
    store.tags?.create({ name: value });
    contact?.update((c) => {
      c.tags = [
        ...(c.tags || []),
        {
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
        },
      ];

      return c;
    });

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

  useModKey('Enter', () => {
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
            store.ui.commandMenu.setOpen(false);
          } else {
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            handleSelect(search as any);
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
