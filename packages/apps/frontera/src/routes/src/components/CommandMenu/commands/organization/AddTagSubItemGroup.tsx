import React, { useMemo } from 'react';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import { OrganizationStore } from '@store/Organizations/Organization.store';

import { Tag01 } from '@ui/media/icons/Tag01';
import { Tag as TagType } from '@graphql/types';
import { Check } from '@ui/media/icons/Check.tsx';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import { CommandSubItem } from '@ui/overlay/CommandMenu';

export const AddTagSubItemGroup = observer(() => {
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

  return (
    <>
      {sortedTags?.map((tag) => (
        <CommandSubItem
          key={tag.id}
          icon={<Tag01 />}
          leftLabel='Change tag'
          rightLabel={tag.value.name}
          rightAccessory={
            newSelectedTags.has(tag.value.name) ? <Check /> : null
          }
          onSelectAction={() => {
            handleSelect(tag.value);
            store.ui.commandMenu.setOpen(false);
          }}
        />
      ))}
    </>
  );
});
