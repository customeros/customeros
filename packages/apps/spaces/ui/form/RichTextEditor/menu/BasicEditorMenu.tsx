import React, { FC } from 'react';
import { useActive, useCommands } from '@remirror/react';
import { Flex, HStack, StackDivider } from '@chakra-ui/react';
import { IconButton } from '@ui/form/IconButton';
import Quote from '../../../../components/ui/media/icons/Quote';
import { TextFormatButtons } from './TextFormatButtons';
import { ListButtons } from './ListButtons';
import { IndentButtons } from './IndentButtons';
import { Button } from '@ui/form/Button';

export const BasicEditorMenu: FC<{
  isSending: boolean;
  onSubmit: () => void;
}> = ({ isSending, onSubmit }) => {
  const { toggleBlockquote, focus } = useCommands();
  const active = useActive();

  return (
    <Flex justifyContent='space-between' flex={1} minH={8}>
      <HStack
        w='full'
        bg='transparent'
        divider={
          <StackDivider
            m={0}
            borderColor='gray.200'
            marginInlineStart={0}
            marginInlineEnd={0}
          />
        }
      >
        <TextFormatButtons />
        <ListButtons />
        <IndentButtons />
        <IconButton
          className='customeros-remirror-button'
          bg='transparent'
          variant='ghost'
          aria-label='Quote'
          onClick={() => {
            toggleBlockquote();
            focus();
          }}
          isActive={active.blockquote()}
          icon={<Quote currentColor='inherit' />}
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
