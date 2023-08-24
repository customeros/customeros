import React, { useCallback } from 'react';
import { htmlToProsemirrorNode } from 'remirror';
import {
  ItalicExtension,
  BoldExtension,
  StrikeExtension,
  UnderlineExtension,
  OrderedListExtension,
  BulletListExtension,
  FontSizeExtension,
  FontFamilyExtension,
  BlockquoteExtension,
  HeadingExtension,
  NodeFormattingExtension,
} from 'remirror/extensions';
import {
  Remirror,
  ThemeProvider,
  Toolbar,
  useRemirror,
  useCommands,
  useActive,
} from '@remirror/react';
import Bold from '../../../components/ui/media/icons/Bold';
import Italic from '../../../components/ui/media/icons/Italic';
import Strikethrough from '../../../components/ui/media/icons/Strikethrough';
import Underline from '../../../components/ui/media/icons/Underline';
import { IconButton } from '@ui/form/IconButton';
import { Center, Divider, Flex, HStack, StackDivider } from '@chakra-ui/react';
import OrderedList from '../../../components/ui/media/icons/OrderedList';
import UnorderedList from '../../../components/ui/media/icons/UnorderedList';
import Quote from '../../../components/ui/media/icons/Quote';
import LeftIndent from '../../../components/ui/media/icons/LeftIndent';
import RightIndent from '../../../components/ui/media/icons/RightIndent';
import { Button } from '@ui/form/Button';
const extensions = () => [
  new ItalicExtension(),
  new BoldExtension(),
  new StrikeExtension(),
  new UnderlineExtension(),
  new OrderedListExtension(),
  new BulletListExtension(),
  new FontSizeExtension(),
  new FontFamilyExtension(),
  new BlockquoteExtension(),
  new HeadingExtension(),
  new NodeFormattingExtension(),
];

const IndentButtons = () => {
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

const ListButtons = () => {
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
        icon={<OrderedList currentColor='inherit' />}
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
        icon={<UnorderedList currentColor='inherit' />}
      />
    </Flex>
  );
};

const TextFormatButtons = () => {
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
export const Menu = () => {
  const { toggleBlockquote, focus } = useCommands();
  const active = useActive();

  return (
    <Flex>
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
        variant='outline'
        fontWeight={600}
        borderRadius='lg'
        pt={0}
        pb={0}
        pl={3}
        pr={3}
        size='sm'
        fontSize='sm'
        loadingText='Sending'
      >
        Send
      </Button>
    </Flex>
  );
};

export const SocialEditor = () => {
  const { manager, state, onChange } = useRemirror({
    extensions: extensions,
    content: '<p>Text in <i>italic</i></p>',
    stringHandler: htmlToProsemirrorNode,
  });
  return (
    <ThemeProvider>
      <Remirror
        manager={manager}
        autoFocus
        onChange={onChange}
        initialContent={state}
        autoRender='end'
      >
        <Toolbar>
          <Menu />
        </Toolbar>
      </Remirror>
    </ThemeProvider>
  );
};
