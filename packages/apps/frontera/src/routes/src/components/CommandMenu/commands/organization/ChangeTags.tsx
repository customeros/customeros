import { useMemo, useState } from 'react';

import { match } from 'ts-pattern';
import { CommandGroup } from 'cmdk';
import unionBy from 'lodash/unionBy';
import { observer } from 'mobx-react-lite';
import { TagDatum } from '@store/Tags/Tag.store';
import { Organization } from '@store/Organizations/Organization.dto';

import { EntityType } from '@graphql/types';
import { Plus } from '@ui/media/icons/Plus.tsx';
import { Check } from '@ui/media/icons/Check.tsx';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import { Command, CommandItem, CommandInput } from '@ui/overlay/CommandMenu';

export const ChangeTags = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;
  const [newTags, setNewTags] = useState<TagDatum[]>([]);

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

  const handleSelect = (t?: TagDatum | null) => {
    if (!t) return;
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
    if (store.tags.toArray().find((e) => e.value.name === value)) return;
    store.tags?.create(
      { name: value, entityType: EntityType.Organization },
      {
        onSucces: (id) => {
          setNewTags((prev) => [store.tags.getById(id)!.value, ...prev]);

          match(context.entity)
            .with('Organization', () => {
              const organization = entity as Organization;

              organization.draft();

              organization?.value.tags?.push({
                name: value,
                metadata: {
                  id,
                },
                entityType: EntityType.Organization,
              });

              organization.commit();
            })
            .with('Organizations', () => {
              store.organizations.updateTags(context.ids as string[], [
                {
                  name: value,
                  entityType: EntityType.Organization,
                  metadata: {
                    id: value,
                  },
                },
              ]);
            });

          setSearch('');
        },
      },
    );
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

  const selectedTags = useMemo(
    () =>
      match(context.entity)
        .with('Organization', () =>
          ((entity as Organization)?.value?.tags ?? []).slice().reverse(),
        )
        .with('Organizations', () => {
          const mappedTags = (entity as Organization[])
            .map((e) => e.value?.tags)
            .reverse()
            .flat()
            .filter((e) => Boolean(e));

          return mappedTags;
        })
        .otherwise(() => []),
    [],
  );

  const allTags = store.tags
    .getByEntityType(EntityType.Organization)
    .map((s) => s.value);

  const uniqueTags = unionBy(newTags, selectedTags, allTags, 'metadata.id');

  useModKey('Enter', () => {
    store.ui.commandMenu.setOpen(false);
  });

  const filteredTags = uniqueTags?.filter((tag) => {
    return tag?.name.toLowerCase().includes(search.toLowerCase());
  });

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
          }
        }}
      />
      <CommandGroup>
        <Command.List>
          {filteredTags?.map((tag) => (
            <CommandItem
              key={tag?.metadata.id}
              rightAccessory={newSelectedTags.has(tag?.name) ? <Check /> : null}
              onSelect={() => {
                handleSelect(tag);
              }}
              onKeyDown={(e) => {
                if (e.metaKey && e.key === 'Enter') {
                  e.stopPropagation();
                  e.preventDefault();
                  store.ui.commandMenu.setOpen(false);

                  return;
                }
              }}
            >
              {tag?.name}
            </CommandItem>
          ))}
          {search && (
            <CommandItem
              leftAccessory={<Plus />}
              onSelect={() => {
                handleCreateOption(search);
              }}
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
