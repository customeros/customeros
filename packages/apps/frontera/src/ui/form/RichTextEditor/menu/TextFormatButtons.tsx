import { useActive, useCommands } from '@remirror/react';

import { Bold01 } from '@ui/media/icons/Bold01';
import { Italic01 } from '@ui/media/icons/Italic01';
import { Heading01 } from '@ui/media/icons/Heading01';
import { Strikethrough01 } from '@ui/media/icons/Strikethrough01';

import { ToolbarButton } from './ToolbarButton';

export const TextFormatButtons = () => {
  const { toggleBold, toggleItalic, toggleHeading, toggleStrike, focus } =
    useCommands();
  const active = useActive();

  return (
    <div className='flex gap-2'>
      <ToolbarButton
        label='Bold'
        isActive={active.bold()}
        icon={<Bold01 className='text-inherit' />}
        onClick={() => {
          toggleBold();
          focus();
        }}
      />
      <ToolbarButton
        label='Italic'
        isActive={active.italic()}
        icon={<Italic01 className='text-inherit' />}
        onClick={() => {
          toggleItalic();
          focus();
        }}
      />
      <ToolbarButton
        label='Strikethrough'
        isActive={active.strike()}
        icon={<Strikethrough01 className='text-inherit' />}
        onClick={() => {
          toggleStrike();
          focus();
        }}
      />
      <ToolbarButton
        label='Heading'
        isActive={active.heading()}
        icon={<Heading01 className='text-inherit' />}
        onClick={() => {
          toggleHeading();
          focus();
        }}
      />
    </div>
  );
};
