import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';
import { useStore } from '@shared/hooks/useStore';
import { SearchSm } from '@ui/media/icons/SearchSm';
import { Tumbleweed } from '@ui/media/icons/Tumbleweed';
import { InputGroup, LeftElement } from '@ui/form/InputGroup';
import { EntityType } from '@shared/types/__generated__/graphql.types';

import { TagList } from './TagsList.tsx';

const entityTypes = {
  [EntityType.Organization]: { label: 'Company' },
  [EntityType.Contact]: { label: 'Contact' },
};
export const TagsManager = observer(() => {
  const store = useStore();
  const [searchTerm, setSearchTerm] = useState('');

  const filteredTags = store.tags.toComputedArray((arr) => {
    arr = arr.filter((entity) => {
      const name = entity.value.name.toLowerCase().includes(searchTerm || '');

      return name;
    });

    return arr;
  });

  const filteredTagsByType = (entityType: EntityType) => {
    const allTags = store.tags.getByEntityType(entityType);

    if (searchTerm) {
      return allTags.filter((tag) =>
        tag.value.name.toLowerCase().includes(searchTerm),
      );
    }

    return allTags;
  };

  const tags = store.tags.toArray().length;

  return (
    <>
      <div className='px-6 pb-4 max-w-[500px] h-full overflow-y-auto border-r border-gray-200'>
        <div className='flex flex-col'>
          <div className='flex justify-between items-center pt-2 sticky top-0 bg-gray-25'>
            <p className='text-gray-700 font-semibold'>Tags</p>
          </div>
          <p className='mb-4 text-sm'>Manage your workspace tags</p>

          {tags > 0 && (
            <div className='mb-4'>
              <InputGroup className='gap-2'>
                <LeftElement>
                  <SearchSm />
                </LeftElement>
                <Input
                  size='xs'
                  className='w-full'
                  variant='unstyled'
                  placeholder='Search tags...'
                  onKeyDown={(e) => {
                    e.stopPropagation();
                  }}
                  onChange={(e) => {
                    setSearchTerm(e.target.value.toLowerCase());
                    e.stopPropagation();
                  }}
                />
              </InputGroup>
            </div>
          )}

          {filteredTags.length === 0 && searchTerm ? (
            <div className='flex justify-center items-center h-full'>
              <div className='flex flex-col items-center mt-4 gap-2'>
                <Tumbleweed className='w-8 h-8 text-gray-400' />
                <span className='text-sm text-gray-500'>
                  Empty here in
                  <span className='font-semibold'> No Resultsville</span>
                </span>
              </div>
            </div>
          ) : (
            <>
              {Object.entries(entityTypes).map(([type, meta]) => (
                <TagList
                  key={type}
                  title={meta.label}
                  entityType={type as EntityType}
                  tags={filteredTagsByType(type as EntityType)}
                />
              ))}
            </>
          )}
        </div>
      </div>
    </>
  );
});
