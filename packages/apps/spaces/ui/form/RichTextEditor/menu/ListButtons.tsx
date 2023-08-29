import { useActive, useCommands } from '@remirror/react';
import { Flex } from '@chakra-ui/react';
import OrderedList from '@ui/media/icons/OrderedList';
import UnorderedList from '@ui/media/icons/UnorderedList';
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
        icon={<OrderedList color='inherit' />}
      />
      <ToolbarButton
        label='Bullet list'
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
