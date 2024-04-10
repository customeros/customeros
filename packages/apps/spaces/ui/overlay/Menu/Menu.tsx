import React, { forwardRef } from 'react';

import { twMerge } from 'tailwind-merge';
import * as DropdownMenuPrimitive from '@radix-ui/react-dropdown-menu';

export const Menu = DropdownMenuPrimitive.Root;
export const MenuGroup = DropdownMenuPrimitive.Group;

interface MenuItemProps extends DropdownMenuPrimitive.DropdownMenuItemProps {
  className?: string;
  children: React.ReactNode;
}

export const MenuItem = forwardRef<HTMLDivElement, MenuItemProps>(
  ({ children, className, ...props }, forwardedRef) => {
    return (
      <DropdownMenuPrimitive.DropdownMenuItem
        {...props}
        ref={forwardedRef}
        className={twMerge(
          'flex text-start py-[0.375rem] px-3 outline-none cursor-pointer hover:bg-gray-50 data-[highlighted]:bg-gray-50 data-[highlighted]:text-gray-700 data-[disabled]:text-gray-400 data-[disabled]:pointer-events-none',
          className,
        )}
      >
        {children}
      </DropdownMenuPrimitive.DropdownMenuItem>
    );
  },
);

interface MenuListProps extends DropdownMenuPrimitive.DropdownMenuContentProps {
  className?: string;
  hasArrow?: boolean;
  children: React.ReactNode;
  align?: 'start' | 'end' | 'center';
  side?: 'top' | 'right' | 'bottom' | 'left';
}

export const MenuList = forwardRef<HTMLDivElement, MenuListProps>(
  (
    { children, hasArrow, side = 'right', align = 'end', className, ...props },
    forwardedRef,
  ) => {
    return (
      <DropdownMenuPrimitive.Portal>
        <DropdownMenuPrimitive.Content
          {...props}
          ref={forwardedRef}
          align={align}
          side={side}
          sideOffset={5}
          className={twMerge(
            className,
            'bg-white min-w-[auto] py-2 shadow-lg border rounded-md data-[side=top]:animate-slideDownAndFade data-[side=right]:animate-slideLeftAndFade data-[side=bottom]:animate-slideUpAndFade data-[side=left]:animate-slideRightAndFade z-10',
          )}
        >
          {children}
          {hasArrow && <DropdownMenuPrimitive.Arrow />}
        </DropdownMenuPrimitive.Content>
      </DropdownMenuPrimitive.Portal>
    );
  },
);

export const MenuButton = DropdownMenuPrimitive.Trigger;
MenuButton.defaultProps = {
  className: 'outline-none',
};
