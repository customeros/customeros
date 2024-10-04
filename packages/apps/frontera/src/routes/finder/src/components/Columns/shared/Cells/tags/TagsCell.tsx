import React from 'react';

import { Tag } from '@graphql/types';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.tsx';

interface ContactCardProps {
  tags: Tag[];
  isHovered?: boolean;
}

export const TagsCell = ({ tags, isHovered }: ContactCardProps) => {
  return (
    <>
      {!tags?.length && <p className='text-gray-400'>No tags set</p>}

      {!!tags?.length && tags.length > 0 && (
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
            <div
              className={'bg-gray-100 rounded-md px-1.5 truncate'}
              style={{
                maxWidth: isHovered ? '80px' : '100px',
              }}
            >
              {tags?.[0].name}
            </div>

            {tags?.length > 1 && !isHovered && (
              <div className='rounded-md w-fit px-1.5 ml-1 text-gray-500 '>
                +{tags?.length - 1}
              </div>
            )}
          </div>
        </Tooltip>
      )}
    </>
  );
};
