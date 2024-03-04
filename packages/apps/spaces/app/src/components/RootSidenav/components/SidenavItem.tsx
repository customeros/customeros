'use client';
import { ReactElement, MouseEventHandler } from 'react';

import { Button } from '@ui/form/Button/Button';

interface SidenavItemProps {
  href?: string;
  label: string;
  isActive?: boolean;
  onClick?: () => void;
  icon: ((isActive: boolean) => ReactElement) | ReactElement;
}

export const SidenavItem = ({
  label,
  icon,
  onClick,
  isActive,
}: SidenavItemProps) => {
  const handleClick: MouseEventHandler = (e) => {
    e.preventDefault();
    onClick?.();
  };

  const className = isActive
    ? `w-full justify-start font-semibold`
    : `w-full justify-start font-normal`;

  return (
    <Button
      size='md'
      variant='tertiary'
      isDestructive
      // px='3'
      // w='full'
      // size='md'
      // variant='ghost'
      // fontSize='sm'
      // textDecoration='none'
      // fontWeight={isActive ? 'semibold' : 'regular'}
      // justifyContent='flex-start'
      // borderRadius='md'
      // bg={isActive ? 'gray.100' : 'transparent'}
      // color={isActive ? 'gray.700' : 'gray.500'}
      onClick={handleClick}
      // leftIcon={typeof icon === 'function' ? icon(!!isActive) : icon}
      // _focus={{
      //   boxShadow: 'sidenavItemFocus',
      // }}
      className={className}
    >
      {label}
    </Button>
  );
};
