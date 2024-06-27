import React, { useRef, useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';

import { Tag, DataSource } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.tsx';
import { SelectOption } from '@shared/types/SelectOptions';
import { Tags } from '@organization/components/Tabs/shared/';
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick';

interface ContactCardProps {
  id: string;
}

export const TagsCell = observer(({ id }: ContactCardProps) => {
  const store = useStore();
  const [isEdit, setIsEdit] = useState(false);
  const contactStore = store.contacts.value.get(id);
  const ref = useRef(null);
  useOutsideClick({
    ref: ref,
    handler: (e) => {
      // @ts-expect-error e.target.id can be undefined
      if (e?.target?.id.includes('react-select')) {
        e.preventDefault();
        e.stopPropagation();

        return;
      }
      setIsEdit(false);
    },
  });

  useEffect(() => {
    store.ui.setIsEditingTableCell(isEdit);
  }, [isEdit]);
  const handleCreateOption = (value: string) => {
    store.tags?.create(undefined, {
      onSucces: (serverId) => {
        store.tags?.value.get(serverId)?.update((tag) => {
          tag.name = value;

          return tag;
        });
      },
    });

    contactStore?.update((contact) => {
      contact.tags = [
        ...(contact.tags || []),
        {
          id: value,
          name: value,
          appSource: 'organization',
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
          source: DataSource.Openline,
        },
      ];

      return contact;
    });
  };

  const tags = contactStore?.value?.tags;

  return (
    <div onDoubleClick={() => setIsEdit(true)} ref={ref}>
      {!isEdit && !contactStore?.value?.tags?.length && (
        <p className='text-gray-400'>Unknown</p>
      )}

      {!isEdit && tags?.length && (
        <Tooltip
          label={
            tags.length > 1
              ? tags
                  .slice(1, tags.length)
                  ?.map((e) => e.name)
                  .join(', ')
              : ''
          }
        >
          <div className='flex w-fit'>
            <div className='bg-gray-100 rounded-md w-fit px-1.5 py-0.5'>
              {tags?.[0].name}
            </div>
            {tags?.length > 1 && (
              <div className='rounded-md w-fit px-1.5 py-0.5 ml-1 text-gray-500'>
                +{tags?.length - 1}
              </div>
            )}
          </div>
        </Tooltip>
      )}
      {isEdit && (
        <Tags
          icon={null}
          placeholder='Persona'
          onChange={(e) => {
            contactStore?.update((c) => {
              c.tags = e.map(
                (option: SelectOption) =>
                  ({
                    id: option.value,
                    name: option.label,
                  } as Tag),
              );

              return c;
            });
          }}
          value={
            contactStore?.value?.tags?.map((t) => ({
              label: t.name,
              value: t.id,
            })) ?? []
          }
          onCreateOption={handleCreateOption}
        />
      )}
    </div>
  );
});
