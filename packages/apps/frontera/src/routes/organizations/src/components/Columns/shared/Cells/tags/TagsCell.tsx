import React from 'react';

import { Tag } from '@graphql/types';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.tsx';
import { SelectOption } from '@shared/types/SelectOptions';
import { Tags } from '@organization/components/Tabs/shared/';

interface ContactCardProps {
  tags: Tag[];
  isEdit: boolean;
  setIsEdit: (isEdit: boolean) => void;
  onCreateOption: (value: string) => void;
  onChange: (value: SelectOption<string>[]) => void;
}

export const TagsCell = ({
  isEdit,
  tags,
  onCreateOption,
  onChange,
}: ContactCardProps) => {
  return (
    <>
      {!isEdit && !tags?.length && <p className='text-gray-400'>Unknown</p>}

      {!isEdit && !!tags?.length && tags.length > 0 && (
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
          hideBorder
          icon={null}
          menuPortalTarget={document.body}
          placeholder='Persona'
          onChange={(e) => {
            onChange(e);
          }}
          value={
            tags?.map((t) => ({
              label: t.name,
              value: t.id,
            })) ?? []
          }
          onCreateOption={onCreateOption}
        />
      )}
    </>
  );
};
