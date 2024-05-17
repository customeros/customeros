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
        onClick={() => {
          toggleBold();
          focus();
        }}
        isActive={active.bold()}
        icon={<Bold01 className='text-inherit' />}
      />
      <ToolbarButton
        label='Italic'
        onClick={() => {
          console.log('🏷️ ----- : Italic ');
          toggleItalic();
          focus();
        }}
        isActive={active.italic()}
        icon={<Italic01 className='text-inherit' />}
      />
      <ToolbarButton
        label='Strikethrough'
        onClick={() => {
          toggleStrike();
          focus();
        }}
        isActive={active.strike()}
        icon={<Strikethrough01 className='text-inherit' />}
      />
      <ToolbarButton
        label='Heading'
        onClick={() => {
          toggleHeading();
          focus();
        }}
        isActive={active.heading()}
        icon={<Heading01 className='text-inherit' />}
      />
    </div>
  );
};
