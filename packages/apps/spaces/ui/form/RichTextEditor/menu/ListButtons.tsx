import { useActive, useCommands } from '@remirror/react';
import { Flex } from '@chakra-ui/react';
import { IconButton } from '@ui/form/IconButton';
import OrderedList from '../../../../components/ui/media/icons/OrderedList';
import UnorderedList from '../../../../components/ui/media/icons/UnorderedList';
import React from 'react';

export const ListButtons = () => {
  const { toggleOrderedList, toggleBulletList, focus } = useCommands();
  const active = useActive();
  return (
    <Flex gap={2}>
      <IconButton
        className='customeros-remirror-button'
        bg='transparent'
        variant='ghost'
        aria-label='Strikethrough'
        onClick={() => {
          toggleOrderedList();
          focus();
        }}
        isActive={active.orderedList()}
        icon={<OrderedList color='inherit' />}
      />
      <IconButton
        className='customeros-remirror-button'
        bg='transparent'
        variant='ghost'
        aria-label='Underline'
        onClick={() => {
          toggleBulletList();
          focus();
        }}
        isActive={active.bulletList()}
        icon={<UnorderedList color='inherit' />}
      />
    </Flex>
  );
};
