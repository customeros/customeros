import React, { useMemo } from 'react';

import { useKeyBindings } from 'rooks';
import { observer } from 'mobx-react-lite';

import { DataSource } from '@graphql/types';
import { Plus } from '@ui/media/icons/Plus.tsx';
import { Tag as TagType } from '@graphql/types';
import { Check } from '@ui/media/icons/Check.tsx';
import { useStore } from '@shared/hooks/useStore';
import {
  Command,
  CommandItem,
  CommandInput,
  useCommandState,
} from '@ui/overlay/CommandMenu';

export const EditPersonaTag = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;

  const contact = store.contacts.value.get(context.ids?.[0] as string);
  const label = `Contact - ${contact?.value?.name}`;
  const [search, setSearch] = React.useState('');

  const handleSelect = (t: TagType) => () => {
    if (!context.ids?.[0]) return;

    if (!contact) return;

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

    // clear search
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

  return (
    <Command label='Change or add tags...'>
      <CommandInput
        label={label}
        value={search}
        onValueChange={setSearch}
        placeholder='Edit persona tag...'
      />

      <EmptySearch createOption={handleCreateOption} />
      <Command.List>
        {sortedTags?.map((tag) => (
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
      </Command.List>
    </Command>
  );
});

const EmptySearch = ({
  createOption,
}: {
  createOption: (data: string) => void;
}) => {
  const search = useCommandState((state) => state.search);

  useKeyBindings({
    Enter: () => {
      createOption(search);
    },
  });

  return (
    <Command.Empty>
      <div
        tabIndex={0}
        role='button'
        onClick={() => createOption(search)}
        className='mx-5 my-3 p-2 flex flex-1 items-center text-gray-500 text-sm hover:bg-gray-50 rounded cursor-pointer'
      >
        <Plus />
        <span className='text-gray-700 ml-1'>Create new tag:</span>
        <span className='text-gray-500 ml-1'>{search}</span>
      </div>
    </Command.Empty>
  );
};
