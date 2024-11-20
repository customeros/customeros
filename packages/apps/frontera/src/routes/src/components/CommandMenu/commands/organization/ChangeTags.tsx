import { useMemo, useState } from 'react';

import { match } from 'ts-pattern';
import { CommandGroup } from 'cmdk';
import { observer } from 'mobx-react-lite';
import { Organization } from '@store/Organizations/Organization.dto';

import { Plus } from '@ui/media/icons/Plus.tsx';
import { Check } from '@ui/media/icons/Check.tsx';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import { DataSource, EntityType, Tag as TagType } from '@graphql/types';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const ChangeTags = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;

  const entity = match(context.entity)
    .returnType<Organization | Organization[] | undefined>()
    .with('Organization', () =>
      store.organizations.getById(context.ids?.[0] as string),
    )
    .with(
      'Organizations',
      () =>
        context.ids?.map((e: string) =>
          store.organizations.getById(e),
        ) as Organization[],
    )
    .otherwise(() => undefined);
  const label = match(context.entity)
    .with(
      'Organization',
      () => `Organization - ${(entity as Organization)?.value?.name}`,
    )
    .with('Organizations', () => `${context.ids?.length} organizations`)
    .otherwise(() => '');

  const [search, setSearch] = useState('');

  const handleSelect = (t: TagType) => {
    if (!context.ids?.[0]) return;

    if (!entity) return;

    match(context.entity)
      .with('Organization', () => {
        const organization = entity as Organization;

        const foundIndex = organization.value.tags?.findIndex(
          (e) => e.metadata.id === t.metadata.id,
        );

        organization.draft();

        if (typeof foundIndex !== 'undefined' && foundIndex > -1) {
          organization.value.tags?.splice(foundIndex, 1);
        } else {
          if (!Array.isArray(organization?.value?.tags)) {
            organization.value.tags = [];
          }
          organization?.value?.tags?.push(t);
        }

        organization.commit();
      })
      .with('Organizations', () => {
        store.organizations.updateTags(context.ids as string[], [t]);
      });
  };

  const handleCreateOption = (value: string) => {
    if (
      store.tags
        .getByEntityType(EntityType.Organization)
        .find((e) => e.value.name === value)
    )
      return;
    store.tags?.create({ name: value });

    match(context.entity)
      .with('Organization', () => {
        const organization = entity as Organization;

        organization?.value.tags?.push({
          id: value,
          name: value,
          appSource: 'organization',
          entityType: EntityType.Organization,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
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
            entityType: EntityType.Organization,
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
          ((entity as Organization)?.value?.tags ?? []).map((tag) => tag?.name),
        ),
    )
    .with('Organizations', () => {
      const mappedTags = (entity as Organization[])
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
            ((entity as Organization)?.value?.tags ?? []).map(
              (tag) => tag?.name,
            ),
          ),
      )
      .with('Organizations', () => {
        const mappedTags = (entity as Organization[])
          .map((e) => e.value?.tags)
          .flat()
          .filter((e) => Boolean(e));

        return new Set(mappedTags.map((tag) => tag?.name));
      })
      .otherwise(() => new Set([]));
  }, []);

  const sortedTags = store.tags
    ?.getByEntityType(EntityType.Organization)
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
    <Command shouldFilter={false} label='Change or add tags...'>
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
              onSelect={() => {
                handleSelect(tag.value);
              }}
              rightAccessory={
                newSelectedTags.has(tag.value.name) ? <Check /> : null
              }
              onKeyDown={(e) => {
                if (e.metaKey && e.key === 'Enter') {
                  e.stopPropagation();
                  e.preventDefault();
                  store.ui.commandMenu.setOpen(false);

                  return;
                }
              }}
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
