import { useMemo } from 'react';

import { observer } from 'mobx-react-lite';

import { Tag01 } from '@ui/media/icons/Tag01';
import { Tag as TagType } from '@graphql/types';
import { Check } from '@ui/media/icons/Check.tsx';
import { useStore } from '@shared/hooks/useStore';
import { useModKey } from '@shared/hooks/useModKey';
import { CommandSubItem } from '@ui/overlay/CommandMenu';

export const AddPersonaTagSubItemGroup = observer(() => {
  const store = useStore();

  const context = store.ui.commandMenu.context;

  const contact = store.contacts.value.get(context.ids?.[0] as string);

  const selectedIds = context.ids;

  const handleSelect = (t: TagType) => {
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

  useModKey('Enter', (e) => {
    e.stopPropagation();
    store.ui.commandMenu.setOpen(false);
  });

  return (
    <>
      {sortedTags?.map((tag) => (
        <CommandSubItem
          key={tag.id}
          icon={<Tag01 />}
          leftLabel='Add tag'
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
