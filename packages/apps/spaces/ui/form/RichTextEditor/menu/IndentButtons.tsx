import { useCommands } from '@remirror/react';
import { Flex } from '@chakra-ui/react';
import { RightIndent } from '@ui/media/icons/RightIndent';
import { LeftIndent } from '@ui/media/icons/LeftIndent';
import React from 'react';
import { ToolbarButton } from './ToolbarButton';

export const IndentButtons = () => {
  const commands = useCommands();
  return (
    <Flex gap={2}>
      <ToolbarButton
        label='Indent'
        onClick={() => {
          commands.decreaseIndent();
        }}
        icon={<RightIndent color='inherit' />}
      />
      <ToolbarButton
        label='Outdent'
        onClick={() => {
          commands.increaseIndent();
        }}
        icon={<LeftIndent color='inherit' />}
      />
    </Flex>
  );
};
