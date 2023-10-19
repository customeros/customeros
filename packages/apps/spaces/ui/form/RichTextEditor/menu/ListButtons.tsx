import { useActive, useCommands } from '@remirror/react';
import { Flex } from '@chakra-ui/react';
import { ListNumbered } from '@ui/media/icons/ListNumbered';
import { ListBulleted } from '@ui/media/icons/ListBulleted';
import React from 'react';
import { ToolbarButton } from './ToolbarButton';

export const ListButtons = () => {
  const { toggleOrderedList, toggleBulletList, focus } = useCommands();
  const active = useActive();
  return (
    <Flex gap={2}>
      <ToolbarButton
        label='Numbered list'
        onClick={() => {
          toggleOrderedList();
          focus();
        }}
        isActive={active.orderedList()}
        icon={<ListNumbered color='gray.400' />}
      />
      <ToolbarButton
        label='Bulleted list'
        onClick={() => {
          toggleBulletList();
          focus();
        }}
        isActive={active.bulletList()}
        icon={<ListBulleted color='gray.400' />}
      />
    </Flex>
  );
};
