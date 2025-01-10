import { observer } from 'mobx-react-lite';

import { Tag as TagType } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { Tag, TagLabel } from '@ui/presentation/Tag';

interface ContactCardProps {
  tags: TagType[];
  isHovered?: boolean;
}

export const TagsCell = ({ tags, isHovered }: ContactCardProps) => {
  return (
    <>
      {!tags?.length && <p className='text-gray-400 truncate'>No tags set</p>}

      <div className='flex '>
        {!!tags?.length && tags.length > 0 && (
          <>
            {tags.map((e) => {
              return (
                <div key={e.metadata.id} className='flex w-fit'>
                  <TagDisplay id={e.metadata.id} isHovered={isHovered} />
                </div>
              );
            })}
          </>
        )}
      </div>
    </>
  );
};

export const TagDisplay = observer(
  ({ id, isHovered }: { id: string; isHovered?: boolean }) => {
    const store = useStore();
    const tag = store.tags.getById(id);

    return (
      <>
        <Tag
          variant='subtle'
          className={'mr-1'}
          style={{
            maxWidth: isHovered ? '80px' : '100px',
          }}
          colorScheme={
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            (tag?.value?.colorCode as unknown as any) ?? 'grayModern'
          }
        >
          <TagLabel className={'truncate'}>{tag?.tagName}</TagLabel>
        </Tag>
      </>
    );
  },
);
