import { useMemo } from 'react';

import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import { TagDatum } from '@store/Tags/Tag.store';
import { Organization } from '@store/Organizations/Organization.dto';

import { EntityType } from '@graphql/types';
import { Tag01 } from '@ui/media/icons/Tag01';
import { Check } from '@ui/media/icons/Check.tsx';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import { CommandSubItem } from '@ui/overlay/CommandMenu';

export const AddTagSubItemGroup = observer(() => {
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

  const handleSelect = (t: TagDatum) => () => {
    if (!context.ids?.[0]) return;

    if (!entity) return;

    match(context.entity)
      .with('Organization', () => {
        const organization = entity as Organization;

        const foundIndex = organization.value.tags?.findIndex(
          (e) => e.name === t.name,
        );

        if (foundIndex && foundIndex > -1) {
          organization.value.tags?.splice(foundIndex, 1);
        } else {
          organization?.value?.tags?.push(t);
        }

        organization.commit();
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
