import { useActive, useCommands } from '@remirror/react';
import { Flex } from '@chakra-ui/react';
import { Bold01 } from '@ui/media/icons/Bold01';
import { Italic01 } from '@ui/media/icons/Italic01';
import { Strikethrough01 } from '@ui/media/icons/Strikethrough01';
import { ToolbarButton } from './ToolbarButton';
import { Heading01 } from '@ui/media/icons/Heading01';

export const TextFormatButtons = () => {
  const { toggleBold, toggleItalic, toggleHeading, toggleStrike, focus } =
    useCommands();
  const active = useActive();
  return (
    <Flex gap={2}>
      <ToolbarButton
        label='Bold'
        onClick={() => {
          toggleBold();
          focus();
        }}
        isActive={active.bold()}
        icon={<Bold01 color='gray.400' />}
      />
      <ToolbarButton
        label='Italic'
        onClick={() => {
          toggleItalic();
          focus();
        }}
        isActive={active.italic()}
        icon={<Italic01 color='gray.400' />}
      />
      <ToolbarButton
        label='Strikethrough'
        onClick={() => {
          toggleStrike();
          focus();
        }}
        isActive={active.strike()}
        icon={<Strikethrough01 color='gray.400' />}
      />
      <ToolbarButton
        label='Heading'
        onClick={() => {
          toggleHeading();
          focus();
        }}
        isActive={active.heading()}
        icon={<Heading01 color='gray.400' />}
      />
    </Flex>
  );
};
