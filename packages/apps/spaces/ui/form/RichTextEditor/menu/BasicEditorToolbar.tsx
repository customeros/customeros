import React, { FC } from 'react';

import { useActive, useCommands } from '@remirror/react';
import { Flex, HStack, StackDivider } from '@chakra-ui/react';

import { Button } from '@ui/form/Button';
import { Quote } from '@ui/media/icons/Quote';

import { ListButtons } from './ListButtons';
import { ToolbarButton } from './ToolbarButton';
import { TextFormatButtons } from './TextFormatButtons';

export const BasicEditorToolbar: FC<{
  isSending: boolean;
  onSubmit: () => void;
}> = ({ isSending, onSubmit }) => {
  const { toggleBlockquote, focus } = useCommands();
  const active = useActive();

  return (
    <Flex justifyContent='space-between' alignItems='center' flex={1} minH={8}>
      <HStack
        w='full'
        bg='transparent'
        divider={
          <StackDivider
            m={0}
            borderColor='gray.200'
            borderWidth='1px'
            marginInlineStart={0}
            marginInlineEnd={0}
          />
        }
      >
        <TextFormatButtons />
        <ListButtons />
        <ToolbarButton
          label='Quote'
          onClick={() => {
            toggleBlockquote();
            focus();
          }}
          isActive={active.blockquote()}
          icon={<Quote color='gray.400' />}
        />
      </HStack>
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
        isDisabled={isSending}
        isLoading={isSending}
        loadingText='Sending'
        onClick={onSubmit}
      >
        Send
      </Button>
    </Flex>
  );
};
