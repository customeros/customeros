import React, { useMemo } from 'react';

import { match } from 'ts-pattern';
import { CommandGroup } from 'cmdk';
import { observer } from 'mobx-react-lite';
import { OrganizationStore } from '@store/Organizations/Organization.store';

import { Plus } from '@ui/media/icons/Plus.tsx';
import { Check } from '@ui/media/icons/Check.tsx';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import { DataSource, Tag as TagType } from '@graphql/types';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const ChangeTags = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;

  const entity = match(context.entity)
    .returnType<OrganizationStore | OrganizationStore[] | undefined>()
    .with('Organization', () =>
      store.organizations.value.get(context.ids?.[0] as string),
    )
    .with(
      'Organizations',
      () =>
        context.ids?.map((e: string) =>
          store.organizations.value.get(e),
        ) as OrganizationStore[],
    )
    .otherwise(() => undefined);
  const label = match(context.entity)
    .with(
      'Organization',
      () => `Organization - ${(entity as OrganizationStore)?.value?.name}`,
    )
    .with('Organizations', () => `${context.ids?.length} organizations`)
    .otherwise(() => '');

  const [search, setSearch] = React.useState('');

  const handleSelect = (t: TagType) => () => {
    if (!context.ids?.[0]) return;

    if (!entity) return;

    match(context.entity)
      .with('Organization', () => {
        (entity as OrganizationStore)?.update((o) => {
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
      })
      .with('Organizations', () => {
        store.organizations.updateTags(context.ids as string[], [t]);
      });
  };

  const handleCreateOption = (value: string) => {
    if (store.tags.toArray().find((e) => e.value.name === value)) return;
    store.tags?.create({ name: value });

    match(context.entity)
      .with('Organization', () => {
        (entity as OrganizationStore)?.update((org) => {
          org.tags = [
            ...(org.tags || []),
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

          return org;
        });
      })
      .with('Organizations', () => {
        store.organizations.updateTags(context.ids as string[], [
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
        ]);
      });

    // clear search
    setSearch('');
  };

  const newSelectedTags = match(context.entity)
    .with(
      'Organization',
      () =>
        new Set(
          ((entity as OrganizationStore)?.value?.tags ?? []).map(
            (tag) => tag?.name,
          ),
        ),
    )
    .with('Organizations', () => {
      const mappedTags = (entity as OrganizationStore[])
        .map((e) => e.value?.tags)
        .flat()
        .filter((e) => Boolean(e));

      return new Set((mappedTags ?? []).map((tag) => tag?.name));
    })
    .otherwise(() => new Set([]));

  const orgTags = useMemo(() => {
    return match(context.entity)
      .with(
        'Organization',
        () =>
          new Set(
            ((entity as OrganizationStore)?.value?.tags ?? []).map(
              (tag) => tag?.name,
            ),
          ),
      )
      .with('Organizations', () => {
        const mappedTags = (entity as OrganizationStore[])
          .map((e) => e.value?.tags)
          .flat()
          .filter((e) => Boolean(e));

        return new Set(mappedTags.map((tag) => tag?.name));
      })
      .otherwise(() => new Set([]));
  }, []);

  const sortedTags = store.tags
    ?.toArray()
    .filter((e) => !!e.value.name)
    .sort((a, b) => {
      const aInOrg = orgTags.has(a.value.name);
      const bInOrg = orgTags.has(b.value.name);

      if (aInOrg && !bInOrg) return -1;
      if (!aInOrg && bInOrg) return 1;

      return 0;
    });

  useModKey('Enter', () => {
    store.ui.commandMenu.setOpen(false);
  });

  const filteredTags = sortedTags?.filter((tag) =>
    tag.value.name.toLowerCase().includes(search.toLowerCase()),
  );

  return (
    <Command
      shouldFilter={false}
      label='Change or add tags...'
      onKeyDown={(e) => {
        e.stopPropagation();
      }}
    >
      <CommandInput
        label={label}
        value={search}
        onValueChange={setSearch}
        placeholder='Change or add tags...'
        onKeyDownCapture={(e) => {
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
