import { useActive, useCommands } from '@remirror/react';
import { Flex } from '@chakra-ui/react';
import { IconButton } from '@ui/form/IconButton';
import Bold from '../../../../components/ui/media/icons/Bold';
import Italic from '../../../../components/ui/media/icons/Italic';
import Strikethrough from '../../../../components/ui/media/icons/Strikethrough';
import Underline from '../../../../components/ui/media/icons/Underline';
import React from 'react';

export const TextFormatButtons = () => {
  const { toggleBold, toggleItalic, toggleStrike, toggleUnderline, focus } =
    useCommands();
  const active = useActive();
  return (
    <Flex gap={2}>
      <IconButton
        className='customeros-remirror-button'
        bg='transparent'
        variant='ghost'
        aria-label='Bold'
        onClick={() => {
          toggleBold();
          focus();
        }}
        isActive={active.bold()}
        icon={<Bold currentColor='inherit' />}
      />
      <IconButton
        className='customeros-remirror-button'
        bg='transparent'
        variant='ghost'
        aria-label='Italic'
        onClick={() => {
          toggleItalic();
          focus();
        }}
        isActive={active.italic()}
        icon={<Italic currentColor='inherit' />}
      />
      <IconButton
        className='customeros-remirror-button'
        bg='transparent'
        variant='ghost'
        aria-label='Strikethrough'
        onClick={() => {
          toggleStrike();
          focus();
        }}
        isActive={active.strike()}
        icon={<Strikethrough currentColor='inherit' />}
      />
      <IconButton
        className='customeros-remirror-button'
        bg='transparent'
        variant='ghost'
        aria-label='Underline'
        onClick={() => {
          toggleUnderline();
          focus();
        }}
        isActive={active.underline()}
        icon={<Underline currentColor='inherit' />}
      />
    </Flex>
  );
};
