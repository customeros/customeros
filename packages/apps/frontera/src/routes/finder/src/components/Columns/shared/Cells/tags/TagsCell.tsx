import { Tag } from '@graphql/types';

interface ContactCardProps {
  tags: Tag[];
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
                  <div
                    className={'bg-gray-100 rounded-md px-1.5 truncate mr-1'}
                    style={{
                      maxWidth: isHovered ? '80px' : '100px',
                    }}
                  >
                    {e.name}
                  </div>
                </div>
              );
            })}
          </>
        )}
      </div>
    </>
  );
};
