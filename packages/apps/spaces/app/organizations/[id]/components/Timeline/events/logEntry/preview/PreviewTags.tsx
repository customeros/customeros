import React from 'react';
import { Text } from '@ui/typography/Text';

export const PreviewTags: React.FC<{
  isAuthor: boolean;
  tags: Array<{ name: string; id: string }>;
  formId: string;
}> = ({ isAuthor, tags = [], formId }) => {
  return (
    <>
      {/*{!isAuthor && (*/}
      <Text fontSize='sm' fontWeight='medium'>
        {tags.map(({ name }) => `#${name}`).join(' ')}
      </Text>
      {/*)}*/}

      {/*{isAuthor && (*/}
      {/*  <Text*/}
      {/*    fontSize='sm'*/}
      {/*    fontWeight='medium'*/}
      {/*    sx={{*/}
      {/*      '--tag-select-font-size': `14px`,*/}
      {/*    }}*/}
      {/*  >*/}
      {/*    <TagsSelect formId={formId} name='tags' tags={data?.tags} />*/}
      {/*  </Text>*/}
      {/*)}*/}
    </>
  );
};
