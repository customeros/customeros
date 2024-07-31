import { useActive, useCommands } from '@remirror/react';

import { ListNumbered } from '@ui/media/icons/ListNumbered';
import { ListBulleted } from '@ui/media/icons/ListBulleted';

import { ToolbarButton } from './ToolbarButton';

export const ListButtons = () => {
  const { toggleOrderedList, toggleBulletList, focus } = useCommands();
  const active = useActive();

  return (
    <div className='flex gap-2'>
      <ToolbarButton
        label='Numbered list'
        isActive={active.orderedList()}
        icon={<ListNumbered className='text-inherit' />}
        onClick={() => {
          toggleOrderedList();
          focus();
        }}
      />
      <ToolbarButton
        label='Bulleted list'
        isActive={active.bulletList()}
        icon={<ListBulleted className='text-inherit' />}
        onClick={() => {
          toggleBulletList();
          focus();
        }}
      />
    </div>
  );
};
