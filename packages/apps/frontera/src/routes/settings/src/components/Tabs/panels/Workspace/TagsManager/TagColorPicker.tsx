import { useState } from 'react';

import { cn } from '@ui/utils/cn.ts';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

export const TagColorPicker = ({
  colorCode,
  onSelect,
}: {
  colorCode: string;
  onSelect: (color: string) => void;
}) => {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <Menu onOpenChange={(data) => setIsOpen(data)}>
      <MenuButton onBlur={(e) => e.stopPropagation()}>
        <div
          className={cn(`size-2 rounded-full ml-3 bg-grayModern-500 `, {
            [`bg-${colorCode}-400`]: colorCode,
            [`ring-1 ring-offset-1 ring-${colorCode}-400`]: isOpen,
          })}
        />
      </MenuButton>
      <MenuList
        side='bottom'
        align='start'
        onBlur={(e) => e.stopPropagation()}
        onKeyDown={(e) => e.stopPropagation()}
        className={'w-fit flex gap-2 p-[6px] bg-grayModern-700'}
      >
        {Object.entries(colorMap).map(([color, bg]) => (
          <MenuItem
            key={color}
            onClick={() => {
              onSelect(color);
            }}
            className='p-0 hover:bg-transparent hover:focus:bg-transparent cursor-pointer'
          >
            <div
              className={cn(`size-4 rounded-full  `, bg, {
                [`ring-1 ring-offset-2 ring-offset-grayModern-700`]:
                  colorCode === color,
              })}
            />
          </MenuItem>
        ))}
      </MenuList>
    </Menu>
  );
};
const colorMap = {
  grayModern: 'bg-grayModern-400 ring-grayModern-400',
  error: 'bg-error-400 ring-error-400',
  warning: 'bg-warning-400 ring-warning-400',
  success: 'bg-success-400 ring-success-400',
  grayWarm: 'bg-grayWarm-400 ring-grayWarm-400',
  moss: 'bg-moss-400 ring-moss-400',
  blueLight: 'bg-blueLight-400 ring-blueLight-400',
  indigo: 'bg-indigo-400 ring-indigo-400',
  violet: 'bg-violet-400 ring-violet-400',
  pink: 'bg-pink-400 ring-pink-400',
};
