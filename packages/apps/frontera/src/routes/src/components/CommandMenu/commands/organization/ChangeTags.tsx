import React from 'react';

import { useKeyBindings } from 'rooks';
import { observer } from 'mobx-react-lite';

import { Check } from '@ui/media/icons/Check.tsx';
import { useStore } from '@shared/hooks/useStore';
import { Tag, TagLabel } from '@ui/presentation/Tag';
import { Tags } from '@organization/components/Tabs';
import { SelectOption } from '@shared/types/SelectOptions.ts';
import { DataSource, OpportunityRenewalLikelihood } from '@graphql/types';
import {
  Command,
  CommandItem,
  CommandInput,
  useCommandState,
} from '@ui/overlay/CommandMenu';

export const ChangeTags = observer(() => {
  const store = useStore();
  const context = store.ui.commandMenu.context;

  const organization = store.organizations.value.get(context.id as string);
  const label = `Organization - ${organization?.value?.name}`;

  const handleSelect = (t) => () => {
    if (!context.id) return;

    if (!organization) return;

    console.log('🏷️ ----- t: ', t);
    organization?.update((o) => {
      o.tags = [t] ?? [];

      return o;
    });
  };

  const handleCreateOption = (value: string) => {
    store.tags?.create(undefined, {
      onSucces: (serverId) => {
        store.tags?.value.get(serverId)?.update((tag) => {
          tag.name = value;

          return tag;
        });
      },
    });

    organization?.update((org) => {
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
  };

  const tags = (organization?.value?.tags ?? []).filter((d) => !!d?.name);

  return (
    <Command label='Change Stage'>
      <div className='flex p-5 gap-1'>
        {tags?.map((t) => (
          <Tag key={t.id} variant='subtle'>
            <TagLabel>{t.name}</TagLabel>
          </Tag>
        ))}
        <Command.Input placeholder='Change tags...' />
      </div>
      <EmptySearch createOption={handleCreateOption} />
      <Command.List>
        {store.tags
          ?.toArray()
          .filter((e) => !!e.value.name)
          .map((tag) => (
            <CommandItem
              key={tag.id}
              onSelect={handleSelect(tag.value)}
              onCreateOption={handleCreateOption}
              rightAccessory={
                tags.find((e) => e.name === tag.value.name) ? <Check /> : null
              }
            >
              {tag.value.name}
            </CommandItem>
          ))}
      </Command.List>
    </Command>
  );
});

const EmptySearch = ({ createOption }: { createOption: any }) => {
  const search = useCommandState((state) => state.search);

  useKeyBindings({
    Enter: () => {
      createOption(search);
    },
  });

  return <Command.Empty>{`Press enter to create "${search}".`}</Command.Empty>;
};
