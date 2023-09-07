import React from 'react';
import { RichTextEditor } from '@ui/form/RichTextEditor/RichTextEditor';
import { useRemirror } from '@remirror/react';
import { basicEditorExtensions } from '@ui/form/RichTextEditor/extensions';
import { Box, Flex } from '@chakra-ui/react';
import { Button } from '@ui/form/Button';
import { TagSuggestor } from './TagSuggestor';

import { TagsSelect } from './TagSelect';
import Image from 'next/image';
import noteIcon from 'public/images/event-ill-log.png';

export const Logger: React.FC = () => {
  const remirrorProps = useRemirror({
    extensions: basicEditorExtensions,
  });
  return (
    <Flex
      flexDirection='column'
      position='relative'
      className='customeros-logger'
    >
      <Box position='absolute' top={-6} right={-6}>
        <Image src={noteIcon} alt='' height={123} width={174} />
      </Box>

      <RichTextEditor
        {...remirrorProps}
        placeholder='Log conversation you had with a customer'
        formId={''}
        name='content'
        showToolbar={false}
      >
        <TagSuggestor />
      </RichTextEditor>
      <Flex justifyContent='space-between' zIndex={3}>
        <TagsSelect />
        <Button
          className='customeros-remirror-submit-button'
          variant='outline'
          colorScheme='gray'
          fontWeight={600}
          borderRadius='lg'
          pt={1}
          pb={1}
          pl={3}
          pr={3}
          size='sm'
          fontSize='sm'
          // isDisabled={isSending}
          // isLoading={isSending}
          loadingText='Sending'
          // onClick={han}
        >
          Log
        </Button>
      </Flex>
    </Flex>
  );
};
