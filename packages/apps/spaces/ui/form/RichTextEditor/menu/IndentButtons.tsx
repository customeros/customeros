import { useCommands } from '@remirror/react';
import { Flex } from '@chakra-ui/react';
import { IconButton } from '@ui/form/IconButton';
import RightIndent from '../../../../components/ui/media/icons/RightIndent';
import LeftIndent from '../../../../components/ui/media/icons/LeftIndent';
import React from 'react';

export const IndentButtons = () => {
  const commands = useCommands();
  return (
    <Flex gap={2}>
      <IconButton
        className='customeros-remirror-button'
        bg='transparent'
        variant='ghost'
        aria-label='Italic'
        onClick={() => {
          commands.decreaseIndent();
        }}
        icon={<RightIndent currentColor='inherit' />}
      />
      <IconButton
        className='customeros-remirror-button'
        bg='transparent'
        variant='ghost'
        aria-label='Strikethrough'
        onClick={() => {
          commands.increaseIndent();
        }}
        icon={<LeftIndent currentColor='inherit' />}
      />
    </Flex>
  );
};
