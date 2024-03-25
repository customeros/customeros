import React from 'react';

import { twMerge } from 'tailwind-merge';
import * as DropdownMenuPrimitive from '@radix-ui/react-dropdown-menu';

export const Menu = DropdownMenuPrimitive.Root;
export const MenuGroup = DropdownMenuPrimitive.Group;

interface MenuItemProps extends DropdownMenuPrimitive.DropdownMenuItemProps {
  className?: string;
  children: React.ReactNode;
}

export const MenuItem = React.forwardRef<HTMLDivElement, MenuItemProps>(
  ({ children, className, ...props }, forwardedRef) => {
    return (
      <DropdownMenuPrimitive.DropdownMenuItem
        {...props}
        ref={forwardedRef}
        className={twMerge(
          'flex text-start  py-[0.375rem] px-3 outline-none cursor-pointer hover:bg-gray-200 ',
          className,
        )}
      >
        <span className=''>{children}</span>
      </DropdownMenuPrimitive.DropdownMenuItem>
    );
  },
);

interface MenuListProps extends DropdownMenuPrimitive.DropdownMenuContentProps {
  className?: string;
  hasArrow?: boolean;
  children: React.ReactNode;
}

export const MenuList = React.forwardRef<HTMLDivElement, MenuListProps>(
  (
    { children, hasArrow, align = 'end', className, ...props },
    forwardedRef,
  ) => {
    return (
      <DropdownMenuPrimitive.Portal>
        <DropdownMenuPrimitive.Content
          {...props}
          ref={forwardedRef}
          align={align}
          className={twMerge(
            className,
            'bg-white min-w-56 py-2 border-b-[1px] shadow-xs outline-offset-[2px] outline-[2px] rounded-md shadow-[0 1px 2px 0 rgba(0,0,0,0.05)] data-[side=top]:animate-slideDownAndFade data-[side=right]:animate-slideLeftAndFade data-[side=bottom]:animate-slideUpAndFade data-[side=left]:animate-slideRightAndFade ',
          )}
        >
          {children}
          {hasArrow && <DropdownMenuPrimitive.Arrow />}
        </DropdownMenuPrimitive.Content>
      </DropdownMenuPrimitive.Portal>
    );
  },
);

interface MenuButtonProps
  extends DropdownMenuPrimitive.DropdownMenuTriggerProps {
  className?: string;
  children: React.ReactNode;
}
export const MenuButton = ({
  className,
  children,
  ...props
}: MenuButtonProps) => {
  return (
    <DropdownMenuPrimitive.Trigger
      {...props}
      className={twMerge('outline-none', className)}
    >
      {children}
    </DropdownMenuPrimitive.Trigger>
  );
};
