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
          'flex gap-2 items-center text-start py-[6px] px-[10px] leading-[18px] text-gray-700  rounded-sm outline-none cursor-pointer hover:bg-gray-50 hover:rounded-md ',
          'data-[highlighted]:bg-gray-50 data-[highlighted]:text-gray-700 data-[disabled]:opacity-50 data-[disabled]:cursor-not-allowed hover:data-[disabled]:bg-transparent',
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
    { children, hasArrow, side = 'bottom', align = 'end', className, ...props },
    forwardedRef,
  ) => {
    return (
      <DropdownMenuPrimitive.Content
        {...props}
        side={side}
        align={align}
        sideOffset={5}
        ref={forwardedRef}
        className={twMerge(
          className,
          'bg-white min-w-[auto] py-1.5 px-[6px] shadow-lg border rounded-md data-[side=top]:animate-slideDownAndFade data-[side=right]:animate-slideLeftAndFade data-[side=bottom]:animate-slideUpAndFade data-[side=left]:animate-slideRightAndFade z-10',
        )}
      >
        {children}
        {hasArrow && <DropdownMenuPrimitive.Arrow />}
      </DropdownMenuPrimitive.Content>
    );
  },
);

export const MenuLabel = forwardRef<
  HTMLDivElement,
  DropdownMenuPrimitive.MenuLabelProps
>(({ className, ...props }, ref) => {
  return (
    <DropdownMenuPrimitive.Label
      ref={ref}
      {...props}
      className={twMerge(
        'text-xs text-gray-500 uppercase px-3 pt-[10px] pb-1 focus:outline-none',
        className,
      )}
    />
  );
});

export const MenuButton = DropdownMenuPrimitive.Trigger;
MenuButton.defaultProps = {
  className: 'focus:outline-none',
};
